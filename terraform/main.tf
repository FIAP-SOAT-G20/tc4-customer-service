terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.15"
    }
  }

  backend "s3" {
    bucket  = "tc4-customer-terraform-state"
    key     = "customer-service/terraform.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = var.common_tags
  }
}

provider "mongodbatlas" {
  public_key  = var.mongodb_atlas_public_key
  private_key = var.mongodb_atlas_private_key
}

# MongoDB Atlas Cluster
module "mongodb_atlas" {
  source = "./modules/mongodb_atlas"

  project_name        = var.project_name
  environment         = var.environment
  mongodb_org_id      = var.mongodb_atlas_org_id
  mongodb_version     = var.mongodb_version
  cluster_tier        = var.mongodb_cluster_tier
  region              = var.mongodb_region
  mongodb_username    = var.mongodb_username
  mongodb_password    = var.mongodb_password
  allowed_cidr_blocks = var.mongodb_allowed_cidr_blocks
  common_tags         = var.common_tags
}

# Lambda Function
module "lambda" {
  source = "./modules/lambda"

  project_name    = var.project_name
  environment     = var.environment
  lambda_zip_path = var.lambda_zip_path
  lambda_handler  = var.lambda_handler
  lambda_runtime  = var.lambda_runtime
  lambda_memory   = var.lambda_memory
  lambda_timeout  = var.lambda_timeout
  common_tags = var.common_tags

  # Environment variables for Lambda
  environment_variables = {
    ENVIRONMENT      = var.environment
    MONGODB_URI      = module.mongodb_atlas.connection_string
    MONGODB_DATABASE = "customer_service"
    JWT_SECRET       = var.jwt_secret
    LOG_LEVEL        = var.log_level
  }

  depends_on = [module.mongodb_atlas]
}

# API Gateway
module "api_gateway" {
  source = "./modules/api_gateway"

  project_name      = var.project_name
  environment       = var.environment
  lambda_arn        = module.lambda.lambda_arn
  lambda_invoke_arn = module.lambda.lambda_invoke_arn
  common_tags       = var.common_tags

  depends_on = [module.lambda]
}

# Lambda Permission for API Gateway (managed separately to avoid circular dependency)
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${module.api_gateway.api_gateway_execution_arn}/*/*"

  depends_on = [module.lambda, module.api_gateway]
}