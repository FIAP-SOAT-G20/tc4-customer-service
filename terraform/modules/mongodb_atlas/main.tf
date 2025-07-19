terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.15"
    }
  }
}

# MongoDB Atlas Project
resource "mongodbatlas_project" "customer_service" {
  name   = "${var.project_name}-${var.environment}"
  org_id = var.mongodb_org_id

  tags = var.common_tags
}

# MongoDB Atlas Cluster
resource "mongodbatlas_cluster" "customer_service" {
  project_id = mongodbatlas_project.customer_service.id
  name = "${var.project_name}-${var.environment}-cluster"

  # Cluster Configuration
  cluster_type = "REPLICASET"

  # MongoDB Version
  mongo_db_major_version = var.mongodb_version

  # Provider Settings
  provider_name               = "TENANT"
  backing_provider_name       = "AWS"
  provider_region_name        = var.region
  provider_instance_size_name = var.cluster_tier
}

# Database User
resource "mongodbatlas_database_user" "customer_service" {
  username           = var.mongodb_username
  password           = var.mongodb_password
  project_id         = mongodbatlas_project.customer_service.id
  auth_database_name = "admin"

  roles {
    role_name     = "readWrite"
    database_name = var.database_name
  }

  labels {
    key   = "environment"
    value = var.environment
  }
}

# IP Access List (Network Access)
resource "mongodbatlas_project_ip_access_list" "customer_service" {
  count = length(var.allowed_cidr_blocks)
  project_id = mongodbatlas_project.customer_service.id
  cidr_block = var.allowed_cidr_blocks[count.index]
  comment    = "CIDR block for ${var.environment} environment"
}

# Data source to get connection string
data "mongodbatlas_cluster" "customer_service" {
  project_id = mongodbatlas_project.customer_service.id
  name       = mongodbatlas_cluster.customer_service.name

  depends_on = [
    mongodbatlas_cluster.customer_service,
    mongodbatlas_database_user.customer_service,
    mongodbatlas_project_ip_access_list.customer_service
  ]
}