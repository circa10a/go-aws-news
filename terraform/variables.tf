variable "app_version" {
  description = "go-aws-news release version to deploy (from GitHub Releases)"
  type        = string
  default     = "1.7.3"
}

variable "lambda_zip_path" {
  description = "Optional path to a pre-built Lambda zip. If null, the module downloads the release from GitHub."
  type        = string
  default     = null
}

variable "function_name" {
  description = "Name of the Lambda function"
  type        = string
  default     = "go-aws-news"
}

variable "config_yaml" {
  description = "Contents of your provider config.yaml file"
  type        = string
  sensitive   = true
}

variable "ssm_parameter_name" {
  description = "Name of the SSM Parameter Store parameter for config"
  type        = string
  default     = "go-aws-news-config"
}

variable "schedule_expression" {
  description = "EventBridge schedule expression (cron or rate)"
  type        = string
  default     = "cron(0 14 * * ? *)" # 2PM UTC / 8AM CST daily
}

variable "lambda_timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 60
}

variable "lambda_memory_size" {
  description = "Lambda memory in MB"
  type        = number
  default     = 128
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 14
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
