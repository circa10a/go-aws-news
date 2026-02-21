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
  config_yaml     = "providers: {}"
  lambda_zip_path = "tests/testdata/lambda.zip"
}

run "custom_function_name" {
  command = plan

  variables {
    function_name = "my-custom-news"
  }

  assert {
    condition     = aws_lambda_function.this.function_name == "my-custom-news"
    error_message = "Custom function name should be respected"
  }

  assert {
    condition     = aws_iam_role.lambda.name == "my-custom-news-role"
    error_message = "IAM role name should derive from custom function name"
  }

  assert {
    condition     = aws_cloudwatch_event_rule.schedule.name == "my-custom-news-schedule"
    error_message = "EventBridge rule name should derive from custom function name"
  }

  assert {
    condition     = aws_cloudwatch_log_group.lambda.name == "/aws/lambda/my-custom-news"
    error_message = "Log group should derive from custom function name"
  }
}

run "custom_ssm_parameter_name" {
  command = plan

  variables {
    ssm_parameter_name = "my-app-config"
  }

  assert {
    condition     = aws_ssm_parameter.config.name == "my-app-config"
    error_message = "Custom SSM parameter name should be respected"
  }

  assert {
    condition     = aws_lambda_function.this.environment[0].variables["GO_AWS_NEWS_CONFIG_NAME"] == "my-app-config"
    error_message = "Lambda env var should reference custom SSM parameter name"
  }
}

run "custom_schedule" {
  command = plan

  variables {
    schedule_expression = "rate(1 hour)"
  }

  assert {
    condition     = aws_cloudwatch_event_rule.schedule.schedule_expression == "rate(1 hour)"
    error_message = "Custom schedule expression should be respected"
  }
}

run "custom_lambda_sizing" {
  command = plan

  variables {
    lambda_timeout     = 120
    lambda_memory_size = 256
  }

  assert {
    condition     = aws_lambda_function.this.timeout == 120
    error_message = "Custom timeout should be respected"
  }

  assert {
    condition     = aws_lambda_function.this.memory_size == 256
    error_message = "Custom memory size should be respected"
  }
}

run "custom_log_retention" {
  command = plan

  variables {
    log_retention_days = 30
  }

  assert {
    condition     = aws_cloudwatch_log_group.lambda.retention_in_days == 30
    error_message = "Custom log retention should be respected"
  }
}

run "custom_lambda_zip_skips_download" {
  command = plan

  variables {
    lambda_zip_path = "tests/testdata/lambda.zip"
  }

  assert {
    condition     = length(terraform_data.download_release) == 0
    error_message = "Download should be skipped when lambda_zip_path is provided"
  }
}
