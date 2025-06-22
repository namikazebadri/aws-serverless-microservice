locals {
  project         = var.PROJECT
  module          = "https"
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
    command = "PATH=$PATH:/home/runner/go/bin && go version && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ${local.src_path}"
    interpreter = ["/bin/bash", "-c"]
  }
}

data "archive_file" "function_archive" {
  depends_on = [null_resource.function_binary]

  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

resource "aws_lambda_function" "https_test_function" {
  function_name = "${local.project}-${var.ENV}_${local.module}_${local.function}"
  description   = "Lambda for ${local.module} module."
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

resource "aws_lambda_function_url" "https_test_function_url" {
  function_name      = aws_lambda_function.https_test_function.function_name
  authorization_type = "AWS_IAM"
}

resource "aws_lambda_permission" "https_test_function_url_permission" {
  statement_id  = "https_test_function_UrlAllowIAM"
  action        = "lambda:InvokeFunctionUrl"
  function_name = aws_lambda_function.https_test_function.function_name
  principal     = "*"
  function_url_auth_type = "AWS_IAM"
}