# API Gateway Outputs
output "api_gateway_url" {
  description = "URL of the API Gateway"
  value       = module.api_gateway.api_gateway_url
}

output "api_gateway_id" {
  description = "ID of the API Gateway"
  value       = module.api_gateway.api_gateway_id
}

# Lambda Outputs
output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = module.lambda.lambda_function_name
}

output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = module.lambda.lambda_arn
}

# MongoDB Atlas Outputs
output "mongodb_connection_string" {
  description = "MongoDB Atlas connection string"
  value       = module.mongodb_atlas.connection_string
  sensitive   = true
}

output "mongodb_cluster_name" {
  description = "MongoDB Atlas cluster name"
  value       = module.mongodb_atlas.cluster_name
}

# Application Endpoints
output "customer_api_endpoints" {
  description = "Customer API endpoints"
  value = {
    base_url       = module.api_gateway.api_gateway_url
    auth           = "${module.api_gateway.api_gateway_url}/auth"
    customers      = "${module.api_gateway.api_gateway_url}/customers"
    customer_by_id = "${module.api_gateway.api_gateway_url}/customers/{id}"
  }
}