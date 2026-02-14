#!/bin/bash
set -e

# sdbx - Deploy to Production Environment

echo "🚀 Deploying sdbx to Production..."
echo ""

# Navigate to prod environment
cd "$(dirname "$0")/../terraform/environments/prod"

# Check if terraform.tfvars exists
if [ ! -f terraform.tfvars ]; then
    echo "📝 Creating terraform.tfvars from example..."
    cp terraform.tfvars.example terraform.tfvars
    echo "  ✓ Created terraform.tfvars"
    echo ""
    echo "⚠️  Please review terraform.tfvars and customize if needed."
    echo "   Press Enter to continue or Ctrl+C to exit..."
    read
fi

# Initialize Terraform
echo "🔧 Initializing Terraform..."
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
BACKEND_BUCKET="sdbx-terraform-state-${ACCOUNT_ID}"
echo "  Backend bucket: ${BACKEND_BUCKET}"
terraform init -backend-config="bucket=${BACKEND_BUCKET}"
echo ""

# Ensure IP hash salt exists in Parameter Store
echo "🔑 Checking IP hash salt in Parameter Store..."
PARAM_NAME="/sdbx/prod/ip-hash-salt"
if aws ssm get-parameter --name "${PARAM_NAME}" --query "Parameter.Name" --output text 2>/dev/null; then
    echo "  ✓ Salt found in Parameter Store"
else
    echo "  Salt not found — initializing automatically..."
    "$(dirname "$0")/init-ip-hash-salt.sh" sdbx prod
fi
echo ""

# Validate configuration
echo "✅ Validating Terraform configuration..."
terraform validate
echo ""

# Format check
echo "📐 Checking Terraform formatting..."
terraform fmt -check -recursive || {
    echo "  ⚠️  Formatting issues found. Running terraform fmt..."
    terraform fmt -recursive
}
echo ""

# Plan deployment
echo "📋 Planning deployment..."
terraform plan -out=tfplan
echo ""

# Prompt for confirmation
echo "⚠️  Ready to deploy PRODUCTION infrastructure to AWS."
echo ""
echo "This will create:"
echo "  - 2 S3 buckets (encrypted files + static frontend)"
echo "  - 1 DynamoDB table"
echo "  - 5 Lambda functions"
echo "  - 1 API Gateway"
echo "  - 1 CloudFront distribution"
echo "  - CloudWatch alarms and dashboard"
echo ""
echo "Estimated monthly cost: ~\$5-20 (depending on usage)"
echo ""
read -p "Do you want to proceed? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ Deployment cancelled"
    rm -f tfplan
    exit 0
fi

# Apply deployment
echo ""
echo "🚀 Deploying infrastructure..."
terraform apply tfplan
rm -f tfplan

echo ""
echo "✅ Production deployment complete!"
echo ""

# Show outputs
echo "📊 Infrastructure outputs:"
terraform output
echo ""

echo "Next steps:"
echo "  1. Deploy frontend: cd ../../../frontend && aws s3 sync . s3://\$(terraform -chdir=../terraform/environments/prod output -raw static_bucket_name)/"
echo "  2. Test API endpoint: curl \$(terraform -chdir=../terraform/environments/prod output -raw api_endpoint)/upload/init"
echo "  3. Open CloudFront URL: \$(terraform -chdir=../terraform/environments/prod output -raw cloudfront_domain)"
echo ""
