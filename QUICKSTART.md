# Quick Start Guide

Get up and running with the OPNsense Terraform provider in minutes.

## Prerequisites

- OPNsense firewall instance
- API credentials (key + secret)
- Terraform >= 1.5
- Go >= 1.21 (for building from source)

## 5-Minute Setup

### 1. Clone and Build

```bash
git clone https://github.com/yknx4/terraform-provider-opnsense.git
cd terraform-provider-opnsense
make deps
make build
make import-tool
```

### 2. Get OPNsense API Credentials

In your OPNsense web interface:
1. Go to **System → Access → Users**
2. Create or edit a user
3. Click **Generate API Key**
4. Save the key and secret

### 3. Set Environment Variables

```bash
export OPNSENSE_URL="https://your-opnsense.example.com"
export OPNSENSE_API_KEY="your-api-key"
export OPNSENSE_API_SECRET="your-api-secret"
```

### 4. Create Your First Resource

Create a file `main.tf`:

```hcl
terraform {
  required_providers {
    opnsense = {
      source = "yknx4/opnsense"
    }
  }
}

provider "opnsense" {
  # Uses environment variables
}

resource "opnsense_firewall_filter" "allow_https" {
  enabled     = true
  description = "Allow HTTPS"
  interface   = "wan"
  action      = "pass"
  direction   = "in"
  protocol    = "TCP"
  dest_port   = "443"
  log         = true
}
```

### 5. Apply

```bash
terraform init
terraform plan
terraform apply
```

## Import Existing Resources

Already have OPNsense resources? Import them:

```bash
# Discover resources
./bin/opnsense-import

# Generate Terraform configuration
terraform plan -generate-config-out=generated.tf

# Review generated.tf and apply
terraform apply
```

## Common Tasks

### Add a Firewall Rule

```hcl
resource "opnsense_firewall_filter" "allow_ssh" {
  enabled     = true
  description = "Allow SSH"
  interface   = "wan"
  action      = "pass"
  direction   = "in"
  protocol    = "TCP"
  source      = "10.0.0.0/8"
  dest_port   = "22"
}
```

### Add a Static Route

```hcl
resource "opnsense_route" "vpn_network" {
  network     = "172.16.0.0/16"
  gateway     = "VPN_GW"
  description = "VPN Network Route"
}
```

### Import a Specific Resource

If you know the UUID:

```hcl
import {
  to = opnsense_firewall_filter.my_rule
  id = "abc123-def456-ghi789"
}
```

Then run:
```bash
terraform plan -generate-config-out=my_rule.tf
```

## Troubleshooting

### "Connection refused"
- Check OPNSENSE_URL is correct
- Ensure OPNsense is accessible from your machine

### "Authentication failed"
- Verify API key and secret
- Check user permissions in OPNsense

### "Certificate verification failed"
- Use self-signed cert? Set `insecure = true` in provider config
- Or add `-insecure` flag to import tool

### "Resource already exists"
- You may have already imported it
- Check `terraform state list`
- Use `terraform import` if needed

## Next Steps

- Read the [full documentation](docs/README.md)
- Check out [examples](examples/)
- Learn about [importing resources](docs/IMPORT_GUIDE.md)
- See [contributing guide](CONTRIBUTING.md) to add features

## Testing Your Setup

Run the included tests:

```bash
make test
```

Expected output:
```
=== RUN   TestNewClient
--- PASS: TestNewClient
...
PASS
```

## Getting Help

- Check documentation in [docs/](docs/)
- Review [examples/](examples/)
- Open an issue on GitHub
- Read the [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)

## Quick Reference

| Command | Description |
|---------|-------------|
| `make build` | Build provider |
| `make import-tool` | Build import tool |
| `make test` | Run tests |
| `make install` | Install provider locally |
| `make clean` | Clean build artifacts |
| `terraform init` | Initialize Terraform |
| `terraform plan` | Preview changes |
| `terraform apply` | Apply changes |
| `./bin/opnsense-import` | Discover resources |
