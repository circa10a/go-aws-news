mock_provider "aws" {
  override_data {
    target = data.aws_iam_policy_document.lambda_assume
    values = {
      json = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"lambda.amazonaws.com\"},\"Action\":\"sts:AssumeRole\"}]}"
    }
  }
  override_data {
    target = data.aws_iam_policy_document.ssm_read
    values = {
      json = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"ssm:GetParameter\"],\"Resource\":[\"*\"]}]}"
    }
  }
}

variables {
  config_yaml = <<-YAML
    providers:
      slack:
        enabled: true
        webhookURL: "https://hooks.slack.com/services/T00/B00/xxx"
        iconURL: "https://example.com/icon.png"
  YAML

  lambda_zip_path = "tests/testdata/lambda.zip"
  function_name   = "go-aws-news-test"
  app_version     = "1.7.3"

  tags = {
    Environment = "test"
  }
}

run "lambda_function_config" {
  command = plan

  assert {
    condition     = aws_lambda_function.this.function_name == "go-aws-news-test"
    error_message = "Lambda function name should match var.function_name"
  }

  assert {
    condition     = aws_lambda_function.this.runtime == "provided.al2023"
    error_message = "Lambda runtime should be provided.al2023"
  }

  assert {
    condition     = aws_lambda_function.this.handler == "bootstrap"
    error_message = "Lambda handler should be bootstrap"
  }

  assert {
    condition     = aws_lambda_function.this.timeout == 60
    error_message = "Lambda timeout should default to 60"
  }

  assert {
    condition     = aws_lambda_function.this.memory_size == 128
    error_message = "Lambda memory should default to 128"
  }

  assert {
    condition     = aws_lambda_function.this.architectures == tolist(["x86_64"])
    error_message = "Lambda architecture should be x86_64"
  }

  assert {
    condition     = aws_lambda_function.this.environment[0].variables["GO_AWS_NEWS_CONFIG_NAME"] == "go-aws-news-config"
    error_message = "Lambda should pass SSM parameter name as env var"
  }
}

run "iam_role_config" {
  command = plan

  assert {
    condition     = aws_iam_role.lambda.name == "go-aws-news-test-role"
    error_message = "IAM role name should be {function_name}-role"
  }

  assert {
    condition     = aws_iam_role_policy_attachment.lambda_basic.policy_arn == "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
    error_message = "Lambda basic execution policy should be attached"
  }

  assert {
    condition     = aws_iam_role_policy.ssm_read.name == "go-aws-news-test-ssm-read"
    error_message = "SSM read policy name should be {function_name}-ssm-read"
  }
}

run "ssm_parameter_config" {
  command = plan

  assert {
    condition     = aws_ssm_parameter.config.name == "go-aws-news-config"
    error_message = "SSM parameter name should match var.ssm_parameter_name"
  }

  assert {
    condition     = aws_ssm_parameter.config.type == "SecureString"
    error_message = "SSM parameter should be SecureString"
  }

  assert {
    condition     = aws_ssm_parameter.config.description == "go-aws-news provider configuration"
    error_message = "SSM parameter description should be set"
  }
}

run "eventbridge_schedule_config" {
  command = plan

  assert {
    condition     = aws_cloudwatch_event_rule.schedule.name == "go-aws-news-test-schedule"
    error_message = "EventBridge rule name should be {function_name}-schedule"
  }

  assert {
    condition     = aws_cloudwatch_event_rule.schedule.schedule_expression == "cron(0 14 * * ? *)"
    error_message = "Default schedule should be cron(0 14 * * ? *)"
  }

  assert {
    condition     = aws_lambda_permission.eventbridge.action == "lambda:InvokeFunction"
    error_message = "Lambda permission should allow InvokeFunction"
  }

  assert {
    condition     = aws_lambda_permission.eventbridge.principal == "events.amazonaws.com"
    error_message = "Lambda permission principal should be events.amazonaws.com"
  }
}

run "cloudwatch_log_group_config" {
  command = plan

  assert {
    condition     = aws_cloudwatch_log_group.lambda.name == "/aws/lambda/go-aws-news-test"
    error_message = "Log group name should be /aws/lambda/{function_name}"
  }

  assert {
    condition     = aws_cloudwatch_log_group.lambda.retention_in_days == 14
    error_message = "Log retention should default to 14 days"
  }
}

run "tags_propagation" {
  command = plan

  assert {
    condition     = aws_lambda_function.this.tags["Environment"] == "test"
    error_message = "Tags should propagate to Lambda function"
  }

  assert {
    condition     = aws_iam_role.lambda.tags["Environment"] == "test"
    error_message = "Tags should propagate to IAM role"
  }

  assert {
    condition     = aws_ssm_parameter.config.tags["Environment"] == "test"
    error_message = "Tags should propagate to SSM parameter"
  }

  assert {
    condition     = aws_cloudwatch_event_rule.schedule.tags["Environment"] == "test"
    error_message = "Tags should propagate to EventBridge rule"
  }

  assert {
    condition     = aws_cloudwatch_log_group.lambda.tags["Environment"] == "test"
    error_message = "Tags should propagate to CloudWatch log group"
  }
}
