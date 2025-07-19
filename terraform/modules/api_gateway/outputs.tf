output "api_gateway_id" {
  description = "ID of the API Gateway"
  value       = aws_api_gateway_rest_api.customer_service.id
}

output "api_gateway_arn" {
  description = "ARN of the API Gateway"
  value       = aws_api_gateway_rest_api.customer_service.arn
}

output "api_gateway_execution_arn" {
  description = "Execution ARN of the API Gateway"
  value       = aws_api_gateway_rest_api.customer_service.execution_arn
}

output "api_gateway_url" {
  description = "URL of the API Gateway"
  value       = aws_api_gateway_stage.customer_service.invoke_url
}

output "api_gateway_stage_name" {
  description = "Name of the API Gateway stage"
  value       = aws_api_gateway_stage.customer_service.stage_name
}

output "api_gateway_deployment_id" {
  description = "ID of the API Gateway deployment"
  value       = aws_api_gateway_deployment.customer_service.id
}