# OPNsense Terraform Provider Examples

This directory contains example Terraform configurations for the OPNsense provider.

## Prerequisites

1. OPNsense instance with API access enabled
2. API key and secret generated in OPNsense
3. Terraform >= 1.5 installed

## Setup

1. Copy `terraform.tfvars.example` to `terraform.tfvars`
2. Fill in your OPNsense credentials
3. Run `terraform init`
4. Run `terraform plan` to see what will be created
5. Run `terraform apply` to apply the configuration

## Examples

### main.tf
Basic provider configuration and example resources including:
- Firewall filter rules
- Static routes

### Using Environment Variables

Instead of using `terraform.tfvars`, you can set environment variables:

```bash
export OPNSENSE_URL="https://opnsense.example.com"
export OPNSENSE_API_KEY="your-api-key"
export OPNSENSE_API_SECRET="your-api-secret"
export OPNSENSE_INSECURE="false"

terraform apply
```

## Security Note

Never commit your `terraform.tfvars` file or any file containing credentials to version control!
