#!/bin/bash

# Initialize Terraform
echo "Initializing Terraform..."
terraform init

# Check if the terraform.tfvars file exists
if [ ! -f "terraform.tfvars" ]; then
    echo "terraform.tfvars file not found"
    exit 1
fi

# Apply Terraform configuration with terraform.tfvars
echo "Applying Terraform configuration..."
terraform apply -var-file="terraform.tfvars"