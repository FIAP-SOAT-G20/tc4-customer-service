# API Gateway REST API
resource "aws_api_gateway_rest_api" "customer_service" {
  name        = "${var.project_name}-${var.environment}-api"
  description = "Customer Service API Gateway for ${var.environment}"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  tags = var.common_tags
}

# API Gateway Deployment
resource "aws_api_gateway_deployment" "customer_service" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.auth.id,
      aws_api_gateway_resource.customers.id,
      aws_api_gateway_resource.customer_id.id,
      aws_api_gateway_method.auth_post.id,
      aws_api_gateway_method.customers_get.id,
      aws_api_gateway_method.customers_post.id,
      aws_api_gateway_method.customer_get.id,
      aws_api_gateway_method.customer_put.id,
      aws_api_gateway_method.customer_delete.id,
      aws_api_gateway_integration.auth_post.id,
      aws_api_gateway_integration.customers_get.id,
      aws_api_gateway_integration.customers_post.id,
      aws_api_gateway_integration.customer_get.id,
      aws_api_gateway_integration.customer_put.id,
      aws_api_gateway_integration.customer_delete.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

# API Gateway Stage
resource "aws_api_gateway_stage" "customer_service" {
  deployment_id = aws_api_gateway_deployment.customer_service.id
  rest_api_id   = aws_api_gateway_rest_api.customer_service.id
  stage_name    = var.environment

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway.arn
    format = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      caller         = "$context.identity.caller"
      user           = "$context.identity.user"
      requestTime    = "$context.requestTime"
      httpMethod     = "$context.httpMethod"
      resourcePath   = "$context.resourcePath"
      status         = "$context.status"
      protocol       = "$context.protocol"
      responseLength = "$context.responseLength"
    })
  }

  tags = var.common_tags
}

# CloudWatch Log Group for API Gateway
resource "aws_cloudwatch_log_group" "api_gateway" {
  name              = "API-Gateway-Execution-Logs_${aws_api_gateway_rest_api.customer_service.id}/${var.environment}"
  retention_in_days = var.log_retention_days

  tags = var.common_tags
}

# Lambda Permission is managed in the Lambda module

# API Gateway Resources and Methods

# /auth resource
resource "aws_api_gateway_resource" "auth" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  parent_id   = aws_api_gateway_rest_api.customer_service.root_resource_id
  path_part   = "auth"
}

# POST /auth
resource "aws_api_gateway_method" "auth_post" {
  rest_api_id   = aws_api_gateway_rest_api.customer_service.id
  resource_id   = aws_api_gateway_resource.auth.id
  http_method   = "POST"
  authorization = "NONE"

  request_validator_id = aws_api_gateway_request_validator.customer_service.id
}

resource "aws_api_gateway_integration" "auth_post" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.auth.id
  http_method = aws_api_gateway_method.auth_post.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# /customers resource
resource "aws_api_gateway_resource" "customers" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  parent_id   = aws_api_gateway_rest_api.customer_service.root_resource_id
  path_part   = "customers"
}

# GET /customers (list and get by CPF)
resource "aws_api_gateway_method" "customers_get" {
  rest_api_id   = aws_api_gateway_rest_api.customer_service.id
  resource_id   = aws_api_gateway_resource.customers.id
  http_method   = "GET"
  authorization = "NONE"

  request_parameters = {
    "method.request.querystring.cpf"   = false
    "method.request.querystring.page"  = false
    "method.request.querystring.limit" = false
  }
}

resource "aws_api_gateway_integration" "customers_get" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.customers.id
  http_method = aws_api_gateway_method.customers_get.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# POST /customers (create)
resource "aws_api_gateway_method" "customers_post" {
  rest_api_id   = aws_api_gateway_rest_api.customer_service.id
  resource_id   = aws_api_gateway_resource.customers.id
  http_method   = "POST"
  authorization = "NONE"

  request_validator_id = aws_api_gateway_request_validator.customer_service.id
}

resource "aws_api_gateway_integration" "customers_post" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.customers.id
  http_method = aws_api_gateway_method.customers_post.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# /customers/{id} resource
resource "aws_api_gateway_resource" "customer_id" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  parent_id   = aws_api_gateway_resource.customers.id
  path_part   = "{id}"
}

# GET /customers/{id}
resource "aws_api_gateway_method" "customer_get" {
  rest_api_id   = aws_api_gateway_rest_api.customer_service.id
  resource_id   = aws_api_gateway_resource.customer_id.id
  http_method   = "GET"
  authorization = "NONE"

  request_parameters = {
    "method.request.path.id" = true
  }
}

resource "aws_api_gateway_integration" "customer_get" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.customer_id.id
  http_method = aws_api_gateway_method.customer_get.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# PUT /customers/{id}
resource "aws_api_gateway_method" "customer_put" {
  rest_api_id   = aws_api_gateway_rest_api.customer_service.id
  resource_id   = aws_api_gateway_resource.customer_id.id
  http_method   = "PUT"
  authorization = "NONE"

  request_parameters = {
    "method.request.path.id" = true
  }

  request_validator_id = aws_api_gateway_request_validator.customer_service.id
}

resource "aws_api_gateway_integration" "customer_put" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.customer_id.id
  http_method = aws_api_gateway_method.customer_put.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# DELETE /customers/{id}
resource "aws_api_gateway_method" "customer_delete" {
  rest_api_id   = aws_api_gateway_rest_api.customer_service.id
  resource_id   = aws_api_gateway_resource.customer_id.id
  http_method   = "DELETE"
  authorization = "NONE"

  request_parameters = {
    "method.request.path.id" = true
  }
}

resource "aws_api_gateway_integration" "customer_delete" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.customer_id.id
  http_method = aws_api_gateway_method.customer_delete.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# Request Validator
resource "aws_api_gateway_request_validator" "customer_service" {
  name                        = "${var.project_name}-${var.environment}-validator"
  rest_api_id                 = aws_api_gateway_rest_api.customer_service.id
  validate_request_body       = true
  validate_request_parameters = true
}

# CORS Configuration for all methods
resource "aws_api_gateway_method" "customers_options" {
  rest_api_id   = aws_api_gateway_rest_api.customer_service.id
  resource_id   = aws_api_gateway_resource.customers.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "customers_options" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.customers.id
  http_method = aws_api_gateway_method.customers_options.http_method
  type        = "MOCK"

  request_templates = {
    "application/json" = "{\"statusCode\": 200}"
  }
}

resource "aws_api_gateway_method_response" "customers_options" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.customers.id
  http_method = aws_api_gateway_method.customers_options.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
}

resource "aws_api_gateway_integration_response" "customers_options" {
  rest_api_id = aws_api_gateway_rest_api.customer_service.id
  resource_id = aws_api_gateway_resource.customers.id
  http_method = aws_api_gateway_method.customers_options.http_method
  status_code = aws_api_gateway_method_response.customers_options.status_code

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'GET,POST,OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }
}