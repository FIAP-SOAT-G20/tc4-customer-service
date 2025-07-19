variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "lambda_arn" {
  description = "ARN of the Lambda function"
  type        = string
}

variable "lambda_invoke_arn" {
  description = "Invoke ARN of the Lambda function"
  type        = string
}

variable "log_retention_days" {
  description = "Number of days to retain CloudWatch logs"
  type        = number
  default     = 14
}

variable "cors_allow_origin" {
  description = "CORS allow origin"
  type        = string
  default     = "*"
}

variable "cors_allow_headers" {
  description = "CORS allow headers"
  type        = string
  default     = "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"
}

variable "cors_allow_methods" {
  description = "CORS allow methods"
  type        = string
  default     = "GET,POST,PUT,DELETE,OPTIONS"
}

variable "common_tags" {
  description = "Common tags to apply to all resources"
  type = map(string)
  default = {}
}