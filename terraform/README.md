# go-aws-news Terraform Module <!-- omit in toc -->

Deploy **go-aws-news** to AWS Lambda with a single `terraform apply` — no need to clone the repo or install Go.

- [Quick Start](#quick-start)
- [What Gets Created](#what-gets-created)
- [Variables](#variables)
- [Outputs](#outputs)
- [Updating Config](#updating-config)
- [Upgrading Versions](#upgrading-versions)
- [Tear Down](#tear-down)
- [Advanced: Custom Lambda Zip](#advanced-custom-lambda-zip)

## Quick Start

**Prerequisites:** [Terraform >= 1.0](https://developer.hashicorp.com/terraform/install) and AWS credentials configured (`aws configure` or environment variables).

**1.** Create a new directory and add two files:

```shell
mkdir go-aws-news-deploy && cd go-aws-news-deploy
```

**2.** Create `config.yaml` — enable at least one provider:

```yaml
providers:
  slack:
    enabled: true
    webhookURL: "https://hooks.slack.com/services/T.../B.../xxx"
    iconURL: "https://i.imgur.com/vfdpgpG.png"
```

> See the [provider docs](../providers/) for all available providers and their settings.

**3.** Create `main.tf`:

```hcl
provider "aws" {
  region = "us-east-1"
}

module "go_aws_news" {
  source = "github.com/circa10a/go-aws-news//terraform"

  config_yaml = file("${path.module}/config.yaml")
}
```

**4.** Deploy:

```shell
terraform init
terraform apply
```

That's it. The module downloads the release binary from GitHub, creates the Lambda function, IAM role, SSM config parameter, EventBridge schedule, and CloudWatch log group.

**5.** Verify:

```shell
aws lambda invoke --function-name go-aws-news /dev/stdout
```

## What Gets Created

| Resource | Purpose |
|---|---|
| **Lambda Function** | Runs go-aws-news on the `provided.al2023` runtime |
| **IAM Role + Policies** | Allows Lambda to write logs and read SSM config |
| **SSM Parameter** (SecureString) | Stores your provider config.yaml securely |
| **EventBridge Rule** | Triggers the Lambda on a cron schedule |
| **CloudWatch Log Group** | Captures Lambda execution logs with configurable retention |

## Variables

| Name | Description | Default |
|---|---|---|
| `config_yaml` | **(required)** Contents of your provider config.yaml | — |
| `app_version` | Release version to deploy | `1.7.3` |
| `function_name` | Lambda function name | `go-aws-news` |
| `ssm_parameter_name` | SSM parameter name for config | `go-aws-news-config` |
| `schedule_expression` | EventBridge cron/rate expression | `cron(0 14 * * ? *)` (2PM UTC) |
| `lambda_timeout` | Lambda timeout (seconds) | `60` |
| `lambda_memory_size` | Lambda memory (MB) | `128` |
| `log_retention_days` | CloudWatch log retention (days) | `14` |
| `lambda_zip_path` | Override: path to a pre-built Lambda zip | `null` (downloads from GitHub) |
| `tags` | Tags applied to all resources | `{}` |

## Outputs

| Name | Description |
|---|---|
| `lambda_function_name` | Name of the deployed Lambda function |
| `lambda_function_arn` | ARN of the deployed Lambda function |
| `lambda_role_arn` | ARN of the Lambda execution role |
| `ssm_parameter_name` | SSM parameter name holding the config |
| `schedule_rule_arn` | ARN of the EventBridge schedule rule |
| `log_group_name` | CloudWatch log group name |

## Updating Config

Edit `config.yaml` and re-apply:

```shell
terraform apply
```

Terraform detects the config change and updates the SSM parameter.

## Upgrading Versions

Bump `app_version` to pull a newer release:

```hcl
module "go_aws_news" {
  source      = "github.com/circa10a/go-aws-news//terraform"
  app_version = "1.8.0"
  config_yaml = file("${path.module}/config.yaml")
}
```

```shell
terraform apply
```

## Tear Down

```shell
terraform destroy
```

Removes all resources created by the module.

## Advanced: Custom Lambda Zip

If you want to build from source instead of downloading a release:

```shell
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w' -o bootstrap .
zip lambda.zip bootstrap
```

Then point the module at it:

```hcl
module "go_aws_news" {
  source         = "github.com/circa10a/go-aws-news//terraform"
  lambda_zip_path = "${path.module}/lambda.zip"
  config_yaml     = file("${path.module}/config.yaml")
}
```
