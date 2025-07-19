variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "mongodb_org_id" {
  description = "MongoDB Atlas organization ID"
  type        = string
  default     = null
}

variable "mongodb_version" {
  description = "MongoDB version"
  type        = string
  default     = "7.0"
}

variable "cluster_tier" {
  description = "MongoDB Atlas cluster tier"
  type        = string
  default     = "M0"
}

variable "region" {
  description = "MongoDB Atlas region"
  type        = string
  default     = "US_EAST_1"
}

variable "mongodb_username" {
  description = "MongoDB database username"
  type        = string
}

variable "mongodb_password" {
  description = "MongoDB database password"
  type        = string
  sensitive   = true
}

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "customer_service"
}

variable "allowed_cidr_blocks" {
  description = "CIDR blocks allowed to access MongoDB"
  type = list(string)
  default = ["0.0.0.0/0"]
}

variable "common_tags" {
  description = "Common tags to apply to all resources"
  type = map(string)
  default = {}
}