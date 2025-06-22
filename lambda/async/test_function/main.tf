locals {
  project         = var.PROJECT
  module          = "async"
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
    command = <<EOT
      set -e

      sudo apt-get update -y
      sudo apt-get install -y golang

      # Cek versi go
      go version

      # Build binary pakai Go dari folder itu
      GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ${local.src_path}

      echo "✅ Build berhasil."
    EOT
    interpreter = ["/bin/bash", "-c"]
  }
}

data "archive_file" "function_archive" {
  depends_on = [null_resource.function_binary]

  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

resource "aws_sqs_queue" "async_test_function_dead_letter_queue" {
  name                        = "${local.project}-${var.ENV}_${local.module}_${local.function}-DLQ"
  delay_seconds               = 2
  max_message_size            = 262144
  message_retention_seconds   = 60
  visibility_timeout_seconds  = 60
}

resource "aws_sqs_queue" "async_test_function_queue" {
  name                      = "${local.project}-${var.ENV}_${local.module}_${local.function}"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  redrive_policy            = "{\"deadLetterTargetArn\":\"${aws_sqs_queue.async_test_function_dead_letter_queue.arn}\",\"maxReceiveCount\":1}"
}

resource "aws_lambda_function" "async_test_function" {
  function_name = "${local.project}-${var.ENV}_${local.module}_${local.function}"
  description   = "Lambda for ${local.module} module."
  role          = var.LAMBDA_ROLE
  handler       = local.function
  memory_size   = 128

  filename         = local.archive_path
  source_code_hash = data.archive_file.function_archive.output_base64sha256
  reserved_concurrent_executions = 50

  runtime = "provided.al2"

  timeout = 10

  environment {
    variables = {
      ENV = var.ENV

      SECRET_MANAGER_NAME = var.SECRET_MANAGER_NAME
    }
  }
}

resource "aws_lambda_event_source_mapping" "event_source_mapping" {
  batch_size        = 1
  event_source_arn  = aws_sqs_queue.async_test_function_queue.arn
  enabled           = true
  function_name     = aws_lambda_function.async_test_function.arn
}