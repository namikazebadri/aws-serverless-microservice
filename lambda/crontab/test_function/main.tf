locals {
  project         = var.PROJECT
  module          = "crontab"
  function        = "test_function"
  module_function = "${local.module}/${local.function}"
  src_path        = "./lambda/${local.module_function}"
  binary_path     = "./bin/${local.module_function}/bootstrap"
  archive_path    = "./bin/${local.module_function}/${local.function}.zip"
}

resource "null_resource" "function_binary" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "export PATH=/home/runner/go/bin:$PATH && go version && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ${local.src_path}"
  }
}

data "archive_file" "function_archive" {
  depends_on = [null_resource.function_binary]

  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

resource "aws_lambda_function" "crontab_test_function" {
  function_name = "${local.project}-${var.ENV}_${local.module}_${local.function}"
  description   = "Lambda function."
  role          = var.LAMBDA_ROLE
  handler       = local.function
  memory_size   = 128

  filename         = local.archive_path
  source_code_hash = data.archive_file.function_archive.output_base64sha256

  runtime = "provided.al2"

  timeout = 10

  environment {
    variables = {
      ENV = var.ENV

      SECRET_MANAGER_NAME = var.SECRET_MANAGER_NAME
    }
  }
}

resource "aws_cloudwatch_event_rule" "crontab_test_function_trigger" {
  name                = "${local.project}-${var.ENV}_${local.module}_${local.function}"
  schedule_expression = "cron(0 * * * ? *)"
}

resource "aws_cloudwatch_event_target" "crontab_test_function_trigger_target" {
  rule      = aws_cloudwatch_event_rule.crontab_test_function_trigger.name
  target_id = "LambdaTarget"
  arn       = aws_lambda_function.crontab_test_function.arn
  input = jsonencode({
    source = "AWS EventBridge"
  })
}

resource "aws_lambda_permission" "allow_event_bridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.crontab_test_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.crontab_test_function_trigger.arn
}