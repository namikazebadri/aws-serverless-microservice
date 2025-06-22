locals {
  project         = var.PROJECT
  module          = "{{MODULE_NAME}}"
  function        = "{{FUNCTION_NAME}}"
  module_function = "${local.module}/${local.function}"
  binary_path     = "./bin/${local.module_function}/bootstrap"
  archive_path    = "./bin/${local.module_function}/${local.function}.zip"
}

data "archive_file" "function_archive" {
  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

resource "aws_lambda_function" "MODULE_NAME_FUNCTION_NAME" {
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

resource "aws_lambda_function_url" "MODULE_NAME_FUNCTION_NAME_url" {
  function_name      = aws_lambda_function.MODULE_NAME_FUNCTION_NAME.function_name
  authorization_type = "AWS_IAM"
}

resource "aws_lambda_permission" "MODULE_NAME_FUNCTION_NAME_url_permission" {
  statement_id  = "MODULE_NAME_FUNCTION_NAME_UrlAllowIAM"
  action        = "lambda:InvokeFunctionUrl"
  function_name = aws_lambda_function.MODULE_NAME_FUNCTION_NAME.function_name
  principal     = "*"
  function_url_auth_type = "AWS_IAM"
}