terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }
  cloud {

    organization = "bct-research"

    workspaces {
      name = "aws-serverless-microservice"
    }
  }
  required_version = ">= 1.3.7"
}

provider "aws" {
  region  = var.AWS_REGION
  access_key = var.AWS_SECRET_KEY
  secret_key = var.AWS_ACCESS_KEY
}

// ==== HTTPS MODULES ====

module "https_test_function" {
  source = "./lambda/https/test_function"
 
  ENV = var.ENV
  PROJECT = var.PROJECT
  SECRET_MANAGER_NAME = var.SECRET_MANAGER_NAME
  LAMBDA_ROLE = var.LAMBDA_ROLE 
}

// ==== ASYNC MODULES ====

module "async_test_function" {
  source = "./lambda/async/test_function"
 
  ENV = var.ENV
  PROJECT = var.PROJECT
  SECRET_MANAGER_NAME = var.SECRET_MANAGER_NAME
  LAMBDA_ROLE = var.LAMBDA_ROLE 
}

// ==== CRONTAB MODULES ====

module "crontab_test_function" {
  source = "./lambda/crontab/test_function"
 
  ENV = var.ENV
  PROJECT = var.PROJECT
  SECRET_MANAGER_NAME = var.SECRET_MANAGER_NAME
  LAMBDA_ROLE = var.LAMBDA_ROLE 
}

output "status" {
  value = "Success"
}