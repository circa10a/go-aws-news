terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

# Download the Lambda binary from GitHub Releases
locals {
  release_url = "https://github.com/circa10a/go-aws-news/releases/download/v${var.app_version}/go-aws-news_${var.app_version}_linux_amd64.tar.gz"
  build_dir   = "${path.module}/.build"
  lambda_zip  = var.lambda_zip_path != null ? var.lambda_zip_path : "${local.build_dir}/lambda.zip"
  fetch_zip   = var.lambda_zip_path == null
}

resource "terraform_data" "download_release" {
  count = local.fetch_zip ? 1 : 0

  triggers_replace = var.app_version

  provisioner "local-exec" {
    command = <<-EOT
      set -e
      mkdir -p "${local.build_dir}"
      echo "Downloading go-aws-news v${var.app_version}..."
      curl -fsSL "${local.release_url}" -o "${local.build_dir}/release.tar.gz"
      tar -xzf "${local.build_dir}/release.tar.gz" -C "${local.build_dir}"
      mv "${local.build_dir}/awsnews" "${local.build_dir}/bootstrap"
      chmod +x "${local.build_dir}/bootstrap"
      cd "${local.build_dir}" && zip -j lambda.zip bootstrap
      rm -f "${local.build_dir}/release.tar.gz" "${local.build_dir}/bootstrap"
      echo "Lambda zip ready at ${local.build_dir}/lambda.zip"
    EOT
  }
}

# Read the zip so we can compute a content hash for Lambda updates
data "local_file" "lambda_zip" {
  filename   = local.lambda_zip
  depends_on = [terraform_data.download_release]
}


# SSM Parameter for provider config
resource "aws_ssm_parameter" "config" {
  name        = var.ssm_parameter_name
  description = "go-aws-news provider configuration"
  type        = "SecureString"
  value       = var.config_yaml

  tags = var.tags
}


# IAM Role for Lambda
data "aws_iam_policy_document" "lambda_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "${var.function_name}-role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  tags               = var.tags
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "ssm_read" {
  statement {
    actions   = ["ssm:GetParameter"]
    resources = [aws_ssm_parameter.config.arn]
  }
}

resource "aws_iam_role_policy" "ssm_read" {
  name   = "${var.function_name}-ssm-read"
  role   = aws_iam_role.lambda.id
  policy = data.aws_iam_policy_document.ssm_read.json
}

# Lambda Function
resource "aws_lambda_function" "this" {
  function_name    = var.function_name
  role             = aws_iam_role.lambda.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["x86_64"]
  filename         = local.lambda_zip
  source_code_hash = data.local_file.lambda_zip.content_base64sha256
  timeout          = var.lambda_timeout
  memory_size      = var.lambda_memory_size

  environment {
    variables = {
      GO_AWS_NEWS_CONFIG_NAME = var.ssm_parameter_name
    }
  }

  depends_on = [
    aws_cloudwatch_log_group.lambda,
    terraform_data.download_release,
  ]

  tags = var.tags
}


# EventBridge Schedule
resource "aws_cloudwatch_event_rule" "schedule" {
  name                = "${var.function_name}-schedule"
  description         = "Triggers go-aws-news on a schedule"
  schedule_expression = var.schedule_expression
  tags                = var.tags
}

resource "aws_cloudwatch_event_target" "lambda" {
  rule = aws_cloudwatch_event_rule.schedule.name
  arn  = aws_lambda_function.this.arn
}

resource "aws_lambda_permission" "eventbridge" {
  statement_id  = "AllowEventBridgeInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.schedule.arn
}

# CloudWatch Log Group (explicit so Terraform manages retention/cleanup)
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = var.log_retention_days
  tags              = var.tags
}
