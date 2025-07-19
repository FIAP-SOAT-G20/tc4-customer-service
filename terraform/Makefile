.DEFAULT_GOAL := help

# Variables
ENV ?= dev
TERRAFORM_DIR = .
VAR_FILE = environments/$(ENV)/terraform.tfvars

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
BLUE = \033[0;34m
NC = \033[0m # No Color

.PHONY: help
help: ## 📜 Show this help message
	@echo "Terraform deployment commands for Customer Service"
	@echo ""
	@echo "Usage: make <command> [ENV=dev|staging|prod]"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

.PHONY: validate-env
validate-env: ## ✅ Validate environment
	@if [ ! -f "$(VAR_FILE)" ]; then \
		echo "$(RED)❌ Environment file not found: $(VAR_FILE)$(NC)"; \
		echo "$(BLUE)💡 Copy from: environments/$(ENV)/terraform.tfvars.example$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Environment $(ENV) is valid$(NC)"

.PHONY: check-lambda
check-lambda: ## 🔍 Check if Lambda package exists
	@if [ ! -f "../dist/function.zip" ]; then \
		echo "$(RED)❌ Lambda package not found: ../dist/function.zip$(NC)"; \
		echo "$(BLUE)💡 Run: make package (in project root)$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Lambda package exists$(NC)"

.PHONY: init
init: ## 🚀 Initialize Terraform
	@echo "$(BLUE)🔧 Initializing Terraform...$(NC)"
	@terraform init
	@echo "$(GREEN)✅ Terraform initialized$(NC)"

.PHONY: workspace
workspace: init ## 🏷️  Create or select workspace
	@echo "$(BLUE)🏷️  Setting up workspace: $(ENV)$(NC)"
	@terraform workspace select $(ENV) 2>/dev/null || terraform workspace new $(ENV)
	@echo "$(GREEN)✅ Workspace $(ENV) ready$(NC)"

.PHONY: validate
validate: workspace ## ✅ Validate Terraform configuration
	@echo "$(BLUE)🔍 Validating Terraform configuration...$(NC)"
	@terraform validate
	@echo "$(GREEN)✅ Configuration is valid$(NC)"

.PHONY: plan
plan: validate-env check-lambda validate ## 📋 Plan Terraform deployment
	@echo "$(BLUE)📋 Planning deployment for $(ENV)...$(NC)"
	@terraform plan -var-file="$(VAR_FILE)" -out="terraform-$(ENV).tfplan"
	@echo "$(GREEN)✅ Plan created: terraform-$(ENV).tfplan$(NC)"

.PHONY: apply
apply: plan ## 🚀 Apply Terraform deployment
	@echo "$(BLUE)🚀 Applying deployment for $(ENV)...$(NC)"
	@terraform apply "terraform-$(ENV).tfplan"
	@rm -f "terraform-$(ENV).tfplan"
	@echo "$(GREEN)✅ Deployment completed for $(ENV)$(NC)"

.PHONY: apply-auto
apply-auto: validate-env check-lambda validate ## 🚀 Apply without interactive approval
	@echo "$(YELLOW)⚠️  Auto-applying deployment for $(ENV) (no confirmation)...$(NC)"
	@terraform apply -var-file="$(VAR_FILE)" -auto-approve
	@echo "$(GREEN)✅ Auto-deployment completed for $(ENV)$(NC)"

.PHONY: destroy
destroy: validate-env workspace ## 💥 Destroy infrastructure
	@echo "$(RED)💥 Destroying infrastructure for $(ENV)...$(NC)"
	@echo "$(YELLOW)⚠️  This will delete all resources!$(NC)"
	@terraform destroy -var-file="$(VAR_FILE)"
	@echo "$(RED)💥 Infrastructure destroyed for $(ENV)$(NC)"

.PHONY: output
output: workspace ## 📤 Show Terraform outputs
	@echo "$(BLUE)📤 Terraform outputs for $(ENV):$(NC)"
	@terraform output

.PHONY: show
show: workspace ## 👁️  Show current state
	@echo "$(BLUE)👁️  Current state for $(ENV):$(NC)"
	@terraform show

.PHONY: fmt
fmt: ## 🎨 Format Terraform files
	@echo "$(BLUE)🎨 Formatting Terraform files...$(NC)"
	@terraform fmt -recursive
	@echo "$(GREEN)✅ Files formatted$(NC)"

.PHONY: clean
clean: ## 🧹 Clean Terraform files
	@echo "$(BLUE)🧹 Cleaning Terraform files...$(NC)"
	@rm -rf .terraform/
	@rm -f *.tfplan
	@rm -f *.tfstate*
	@echo "$(GREEN)✅ Cleaned$(NC)"

.PHONY: refresh
refresh: validate-env workspace ## 🔄 Refresh Terraform state
	@echo "$(BLUE)🔄 Refreshing state for $(ENV)...$(NC)"
	@terraform refresh -var-file="$(VAR_FILE)"
	@echo "$(GREEN)✅ State refreshed$(NC)"

.PHONY: import
import: validate-env workspace ## 📥 Import existing resource
	@echo "$(BLUE)📥 Import resource for $(ENV):$(NC)"
	@echo "Usage: make import RESOURCE=aws_lambda_function.example ID=function-name"
	@if [ -z "$(RESOURCE)" ] || [ -z "$(ID)" ]; then \
		echo "$(RED)❌ RESOURCE and ID are required$(NC)"; \
		exit 1; \
	fi
	@terraform import -var-file="$(VAR_FILE)" $(RESOURCE) $(ID)

.PHONY: deploy-dev
deploy-dev: ## 🚀 Deploy to development
	@$(MAKE) apply ENV=dev

.PHONY: deploy-staging
deploy-staging: ## 🚀 Deploy to staging
	@$(MAKE) apply ENV=staging

.PHONY: deploy-prod
deploy-prod: ## 🚀 Deploy to production
	@$(MAKE) apply ENV=prod

.PHONY: status
status: ## 📊 Show deployment status
	@echo "$(BLUE)📊 Deployment Status$(NC)"
	@echo "Current workspace: $$(terraform workspace show)"
	@echo "Environment: $(ENV)"
	@echo "Config file: $(VAR_FILE)"
	@if [ -f "$(VAR_FILE)" ]; then \
		echo "$(GREEN)✅ Config exists$(NC)"; \
	else \
		echo "$(RED)❌ Config missing$(NC)"; \
	fi
	@if [ -f "../dist/function.zip" ]; then \
		echo "$(GREEN)✅ Lambda package ready$(NC)"; \
	else \
		echo "$(RED)❌ Lambda package missing$(NC)"; \
	fi

.PHONY: endpoints
endpoints: workspace ## 🌐 Show API endpoints
	@echo "$(BLUE)🌐 API Endpoints for $(ENV):$(NC)"
	@terraform output -json customer_api_endpoints | jq -r 'to_entries[] | "\(.key): \(.value)"'

.PHONY: logs
logs: workspace ## 📝 Show recent Lambda logs
	@echo "$(BLUE)📝 Recent Lambda logs for $(ENV):$(NC)"
	@FUNCTION_NAME=$$(terraform output -raw lambda_function_name) && \
	aws logs tail "/aws/lambda/$$FUNCTION_NAME" --follow

.PHONY: test-api
test-api: workspace ## 🧪 Test deployed API
	@echo "$(BLUE)🧪 Testing deployed API for $(ENV):$(NC)"
	@API_URL=$$(terraform output -raw api_gateway_url) && \
	echo "Testing: $$API_URL/customers" && \
	curl -s "$$API_URL/customers" | jq .