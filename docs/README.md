# Terraform Provider for OPNsense

A Terraform provider for managing OPNsense firewall resources via the OPNsense API.

## Features

- **Provider Configuration**: Authenticate with OPNsense using API keys
- **Firewall Filter Resources**: Manage firewall rules
- **Route Resources**: Manage static routes
- **Import Tool**: Discover and import existing OPNsense resources

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.5
- [Go](https://golang.org/doc/install) >= 1.21 (for building from source)
- OPNsense instance with API access enabled

## Building the Provider

Clone the repository and build the provider:

```bash
git clone https://github.com/yknx4/terraform-provider-opnsense.git
cd terraform-provider-opnsense
make build
```

To install the provider locally for testing:

```bash
make install
```

## Using the Provider

### Provider Configuration

```hcl
terraform {
  required_providers {
    opnsense = {
      source = "yknx4/opnsense"
    }
  }
}

provider "opnsense" {
  url        = "https://opnsense.example.com"
  api_key    = "your-api-key"
  api_secret = "your-api-secret"
  insecure   = false  # Set to true to skip TLS verification
}
```

You can also use environment variables:
- `OPNSENSE_URL`
- `OPNSENSE_API_KEY`
- `OPNSENSE_API_SECRET`
- `OPNSENSE_INSECURE`

### Resources

#### Firewall Filter

Manages firewall filter rules.

```hcl
resource "opnsense_firewall_filter" "allow_https" {
  enabled     = true
  description = "Allow inbound HTTPS traffic"
  interface   = "wan"
  action      = "pass"
  direction   = "in"
  protocol    = "TCP"
  source      = "any"
  destination = "192.168.1.100"
  dest_port   = "443"
  log         = true
}
```

**Arguments:**
- `enabled` (Required) - Whether the rule is enabled
- `description` (Optional) - Description of the rule
- `interface` (Required) - Interface to apply the rule (e.g., "wan", "lan")
- `action` (Required) - Action to take ("pass", "block", "reject")
- `direction` (Required) - Traffic direction ("in", "out")
- `protocol` (Optional) - Protocol ("TCP", "UDP", "ICMP", "any")
- `source` (Optional) - Source address or network
- `destination` (Optional) - Destination address or network
- `source_port` (Optional) - Source port or range
- `dest_port` (Optional) - Destination port or range
- `log` (Optional) - Whether to log matches

#### Route

Manages static routes.

```hcl
resource "opnsense_route" "example" {
  network     = "10.9.0.0/24"
  gateway     = "LAN_DHCP"
  description = "Route to internal network"
  disabled    = false
}
```

**Arguments:**
- `network` (Required) - Destination network in CIDR notation
- `gateway` (Required) - Gateway to use for this route
- `description` (Optional) - Description of the route
- `disabled` (Optional) - Whether the route is disabled

## Import Tool

The import tool discovers existing OPNsense resources and generates an `import.tf` file for importing them into Terraform.

### Building the Import Tool

```bash
make import-tool
```

### Using the Import Tool

```bash
./bin/opnsense-import \
  -url https://opnsense.example.com \
  -api-key YOUR_API_KEY \
  -api-secret YOUR_API_SECRET \
  -output import.tf
```

Or using environment variables:

```bash
export OPNSENSE_URL=https://opnsense.example.com
export OPNSENSE_API_KEY=your-api-key
export OPNSENSE_API_SECRET=your-api-secret
./bin/opnsense-import
```

The tool will:
1. Connect to your OPNsense instance
2. Discover all firewall rules and routes
3. Generate an `import.tf` file with import blocks

### Importing Resources

After generating the import file:

```bash
# Generate Terraform configuration for discovered resources
terraform plan -generate-config-out=generated.tf

# Review the generated configuration
cat generated.tf

# Apply the import
terraform apply
```

## Development

### Running Tests

```bash
make test
```

### Code Formatting

```bash
make fmt
```

### Building

```bash
make build
```

### Cleaning Build Artifacts

```bash
make clean
```

## API Reference

This provider uses the OPNsense API as documented at:
https://docs.opnsense.org/development/how-tos/api.html

### Setting Up API Access in OPNsense

1. Log in to your OPNsense web interface
2. Navigate to System → Access → Users
3. Create or edit a user
4. Generate an API key and secret
5. Ensure the user has appropriate permissions for the resources you want to manage

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See the [LICENSE](LICENSE) file for details.
