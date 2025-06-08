#!/bin/bash

# Initialize Terraform
echo "Initializing Terraform..."
terraform init

# Apply Terraform configuration with terraform.tfvars
echo "Applying Terraform configuration..."
terraform apply -var-file="terraform.tfvars" -auto-approve