provider "aws" {
  region = "us-east-1" # Change to your preferred region
}

module "go_aws_news" {
  source = "github.com/circa10a/go-aws-news//terraform"

  # Required: your provider config (enable at least one provider)
  config_yaml = file("${path.module}/config.yaml")

  # Optional: pin to a specific release version
  # app_version = "1.7.3"

  # Optional: override the schedule (default: 2PM UTC / 8AM CST daily)
  # schedule_expression = "cron(0 14 * * ? *)"

  # Optional: customize naming
  # function_name      = "go-aws-news"
  # ssm_parameter_name = "go-aws-news-config"

  # Optional: Lambda tuning
  # lambda_timeout     = 60
  # lambda_memory_size = 128

  # Optional: log retention
  # log_retention_days = 14

  tags = {
    Project   = "go-aws-news"
    ManagedBy = "terraform"
  }
}

output "lambda_function_name" {
  value = module.go_aws_news.lambda_function_name
}

output "lambda_function_arn" {
  value = module.go_aws_news.lambda_function_arn
}
