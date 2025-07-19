output "project_id" {
  description = "MongoDB Atlas project ID"
  value       = mongodbatlas_project.customer_service.id
}

output "cluster_name" {
  description = "MongoDB Atlas cluster name"
  value       = mongodbatlas_cluster.customer_service.name
}

output "cluster_id" {
  description = "MongoDB Atlas cluster ID"
  value       = mongodbatlas_cluster.customer_service.id
}

output "connection_string" {
  description = "MongoDB Atlas connection string with credentials"
  value       = "mongodb+srv://${var.mongodb_username}:${var.mongodb_password}@tc4-customer-service-de.tntn79t.mongodb.net/${var.database_name}?retryWrites=true&w=majority"
  sensitive   = true
}