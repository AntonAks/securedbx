# One-time bootstrap: lets GitHub Actions deploy to AWS without any stored keys.
# Apply locally once per account:
#   cd terraform/github-oidc && terraform init && terraform apply
# Then set the output role ARN as repo variable AWS_DEPLOY_ROLE_ARN:
#   gh variable set AWS_DEPLOY_ROLE_ARN --body "$(terraform output -raw deploy_role_arn)"
#
# State is local on purpose — this is applied once from a trusted machine.

terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "eu-central-1"
}

resource "aws_iam_openid_connect_provider" "github" {
  url            = "https://token.actions.githubusercontent.com"
  client_id_list = ["sts.amazonaws.com"]
  # NOT a secret: public SHA-1 fingerprint of the DigiCert root CA behind
  # token.actions.githubusercontent.com — the same well-known constant from the
  # AWS/GitHub OIDC docs for everyone. Since 2023 AWS validates GitHub's cert
  # chain itself, but the API still requires the field to be present.
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]
}

data "aws_iam_policy_document" "trust" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }
    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values = [
        "repo:AntonAks/securedbx:ref:refs/heads/main",
        "repo:AntonAks/securedbx:ref:refs/heads/dev",
        "repo:AntonAks/securedbx:ref:refs/tags/v*", # prod releases are v* tags
      ]
    }
  }
}

resource "aws_iam_role" "deploy" {
  name               = "sdbx-github-deploy"
  assume_role_policy = data.aws_iam_policy_document.trust.json
}

# ponytail: AdministratorAccess — terraform here manages IAM/S3/Lambda/CloudFront/
# DynamoDB/SSM, a least-privilege policy is a moving target for a solo project.
# Tighten to a scoped policy if this account ever holds anything beyond sdbx.
resource "aws_iam_role_policy_attachment" "deploy_admin" {
  role       = aws_iam_role.deploy.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}

output "deploy_role_arn" {
  value = aws_iam_role.deploy.arn
}
