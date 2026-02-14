# Terraform Provider for OPNsense

A Terraform provider for managing OPNsense firewall resources via the OPNsense API.

## Features

- **Provider Configuration**: Authenticate with OPNsense using API keys
- **Firewall Filter Resources**: Manage firewall rules
- **Route Resources**: Manage static routes  
- **Import Tool**: Discover and import existing OPNsense resources into Terraform
- **Tests**: Comprehensive unit and acceptance tests

## Quick Start

```bash
# Build the provider and import tool
make build
make import-tool

# Run tests
make test

# Install locally for development
make install
```

## Documentation

Detailed documentation is available in the [docs](./docs) directory.

### Basic Usage

```hcl
provider "opnsense" {
  url        = "https://opnsense.example.com"
  api_key    = var.opnsense_api_key
  api_secret = var.opnsense_api_secret
}

resource "opnsense_firewall_filter" "allow_https" {
  enabled     = true
  description = "Allow HTTPS"
  interface   = "wan"
  action      = "pass"
  direction   = "in"
  protocol    = "TCP"
  dest_port   = "443"
}
```

### Import Existing Resources

Use the included import tool to discover and import existing OPNsense resources:

```bash
./bin/opnsense-import -url https://opnsense.example.com \
  -api-key YOUR_KEY -api-secret YOUR_SECRET
terraform plan -generate-config-out=generated.tf
```

## Requirements

- Terraform >= 1.5
- Go >= 1.21 (for building)
- OPNsense instance with API access

## License

See [LICENSE](LICENSE) file.