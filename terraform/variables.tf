# General Variables
variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "tc4-customer-service"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  validation {
    condition = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "common_tags" {
  description = "Common tags to apply to all resources"
  type = map(string)
  default = {
    Project   = "TC4-Customer-Service"
    Terraform = "true"
  }
}

# MongoDB Atlas Variables
variable "mongodb_atlas_public_key" {
  description = "MongoDB Atlas public key"
  type        = string
  sensitive   = true
}

variable "mongodb_atlas_private_key" {
  description = "MongoDB Atlas private key"
  type        = string
  sensitive   = true
}

variable "mongodb_atlas_org_id" {
  description = "MongoDB Atlas organization ID"
  type        = string
}

variable "mongodb_version" {
  description = "MongoDB version"
  type        = string
  default     = "7.0"
}

variable "mongodb_cluster_tier" {
  description = "MongoDB Atlas cluster tier"
  type        = string
  default     = "M0" # Free tier
}

variable "mongodb_region" {
  description = "MongoDB Atlas region"
  type        = string
  default     = "US_EAST_1"
}

variable "mongodb_username" {
  description = "MongoDB database username"
  type        = string
  default     = "customer_service"
}

variable "mongodb_password" {
  description = "MongoDB database password"
  type        = string
  sensitive   = true
}

variable "mongodb_allowed_cidr_blocks" {
  description = "CIDR blocks allowed to access MongoDB"
  type = list(string)
  default = ["0.0.0.0/0"] # Allow all for now, restrict in production
}

# Lambda Variables
variable "lambda_zip_path" {
  description = "Path to the Lambda deployment zip file"
  type        = string
  default     = "../dist/function.zip"
}

variable "lambda_handler" {
  description = "Lambda handler"
  type        = string
  default     = "bootstrap"
}

variable "lambda_runtime" {
  description = "Lambda runtime"
  type        = string
  default     = "provided.al2023"
}

variable "lambda_memory" {
  description = "Lambda memory in MB"
  type        = number
  default     = 512
}

variable "lambda_timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 30
}

# Application Variables
variable "jwt_secret" {
  description = "JWT secret key"
  type        = string
  sensitive   = true
}

variable "log_level" {
  description = "Log level"
  type        = string
  default     = "INFO"
  validation {
    condition = contains(["DEBUG", "INFO", "WARN", "ERROR"], var.log_level)
    error_message = "Log level must be one of: DEBUG, INFO, WARN, ERROR."
  }
}