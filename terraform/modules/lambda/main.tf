# Use existing LabRole (AWS Academy)
data "aws_iam_role" "lab_role" {
  name = "LabRole"
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.project_name}-${var.environment}-customer-service"
  retention_in_days = var.log_retention_days

  tags = var.common_tags
}

# Lambda Function
resource "aws_lambda_function" "customer_service" {
  filename      = var.lambda_zip_path
  function_name = "${var.project_name}-${var.environment}-customer-service"
  role          = data.aws_iam_role.lab_role.arn
  handler       = var.lambda_handler
  runtime       = var.lambda_runtime
  memory_size   = var.lambda_memory
  timeout       = var.lambda_timeout
  source_code_hash = filebase64sha256(var.lambda_zip_path)

  environment {
    variables = var.environment_variables
  }

  depends_on = [
    aws_cloudwatch_log_group.lambda_logs,
  ]

  tags = var.common_tags
}
