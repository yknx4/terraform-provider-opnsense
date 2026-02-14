# OPNsense Import Tool Guide

The OPNsense import tool helps you discover existing resources in your OPNsense instance and generate Terraform configuration to import them.

## Overview

The import tool:
1. Connects to your OPNsense API
2. Discovers existing resources (firewall rules, routes, etc.)
3. Generates an `import.tf` file with import blocks
4. Provides instructions for completing the import process

## Installation

Build the import tool:

```bash
make import-tool
```

This creates the binary at `./bin/opnsense-import`.

## Usage

### Basic Usage

```bash
./bin/opnsense-import \
  -url https://opnsense.example.com \
  -api-key YOUR_API_KEY \
  -api-secret YOUR_API_SECRET
```

### Using Environment Variables

You can also use environment variables instead of command-line flags:

```bash
export OPNSENSE_URL=https://opnsense.example.com
export OPNSENSE_API_KEY=your-api-key
export OPNSENSE_API_SECRET=your-api-secret
./bin/opnsense-import
```

### Command-Line Options

- `-url` - OPNsense URL (can use `OPNSENSE_URL` env var)
- `-api-key` - API key (can use `OPNSENSE_API_KEY` env var)
- `-api-secret` - API secret (can use `OPNSENSE_API_SECRET` env var)
- `-insecure` - Skip TLS verification (can use `OPNSENSE_INSECURE=true` env var)
- `-output` - Output file path (default: `import.tf`)

## Complete Import Workflow

### Step 1: Run the Import Tool

```bash
./bin/opnsense-import
```

Output:
```
Discovering OPNsense resources...
Found 5 firewall filter rules
Found 3 routes

Import configuration written to import.tf
Next steps:
  1. Review the generated import.tf file
  2. Run: terraform plan -generate-config-out=generated.tf
  3. Review and edit generated.tf as needed
  4. Run: terraform apply
```

### Step 2: Review Generated import.tf

The tool creates an `import.tf` file that looks like:

```hcl
# Generated import configuration for OPNsense resources
# Run 'terraform plan -generate-config-out=generated.tf' to generate resource configurations

import {
  to = opnsense_firewall_filter.allow_https_traffic
  id = "abc123-def456-ghi789"
}

import {
  to = opnsense_route.internal_network
  id = "xyz789-uvw456-rst123"
}
```

### Step 3: Generate Resource Configuration

Terraform 1.5+ can automatically generate resource configurations:

```bash
terraform plan -generate-config-out=generated.tf
```

This creates a `generated.tf` file with full resource configurations:

```hcl
resource "opnsense_firewall_filter" "allow_https_traffic" {
  enabled     = true
  description = "Allow HTTPS traffic"
  interface   = "wan"
  action      = "pass"
  # ... other attributes
}
```

### Step 4: Review and Customize

1. Review the `generated.tf` file
2. Customize resource names and descriptions as needed
3. Organize resources logically (you can split into multiple files)
4. Remove any resources you don't want to manage with Terraform

### Step 5: Apply the Configuration

```bash
terraform plan
terraform apply
```

Terraform will import the resources into its state without making any changes.

## Resource Naming

The import tool automatically generates Terraform resource names based on:

1. The resource's description (if available)
2. The resource's network/identifier (for routes)
3. The resource's UUID (as a fallback)

Names are sanitized to be valid Terraform identifiers:
- Converted to lowercase
- Spaces and special characters replaced with underscores
- Multiple consecutive underscores collapsed to one

Examples:
- "Allow HTTPS Traffic" → `allow_https_traffic`
- "Route-to-VPN" → `route_to_vpn`
- "10.0.0.0/8" → `10_0_0_0_8`

## Supported Resources

The import tool currently discovers:

- **Firewall Filter Rules** (`opnsense_firewall_filter`)
- **Static Routes** (`opnsense_route`)

More resources will be added in future versions.

## Troubleshooting

### Authentication Errors

If you get authentication errors:
1. Verify your API credentials are correct
2. Check that the API user has sufficient permissions
3. Ensure API access is enabled in OPNsense

### TLS Certificate Errors

If you get TLS certificate errors:
1. Use the `-insecure` flag (not recommended for production)
2. Or, ensure your OPNsense instance has a valid SSL certificate

### No Resources Found

If no resources are found:
1. Verify you have resources configured in OPNsense
2. Check that the API user has permissions to view those resources
3. Check the OPNsense API documentation for the correct endpoints

## Advanced Usage

### Custom Output File

Save the import configuration to a custom location:

```bash
./bin/opnsense-import -output /path/to/my-imports.tf
```

### Import Only Specific Resources

If you only want to import specific resources, you can:
1. Generate the full `import.tf` file
2. Manually edit it to remove unwanted import blocks
3. Run `terraform plan -generate-config-out=generated.tf`

### Incremental Imports

You can run the import tool multiple times:
1. After adding new resources to OPNsense
2. Generate a new import file
3. Merge the new import blocks with your existing configuration

## Best Practices

1. **Always review** the generated configuration before applying
2. **Use version control** to track changes to your Terraform configuration
3. **Test in a non-production environment** first
4. **Keep sensitive data secure** - never commit API credentials to version control
5. **Document customizations** you make to the generated configuration

## Example Session

```bash
# Set up environment
export OPNSENSE_URL=https://opnsense.example.com
export OPNSENSE_API_KEY=my-api-key
export OPNSENSE_API_SECRET=my-api-secret

# Run import tool
./bin/opnsense-import

# Review import.tf
cat import.tf

# Generate configuration
terraform plan -generate-config-out=generated.tf

# Review generated configuration
cat generated.tf

# Customize if needed
vi generated.tf

# Apply
terraform plan
terraform apply
```

## Getting Help

If you encounter issues or have questions:
1. Check this guide
2. Review the main documentation in `docs/README.md`
3. Open an issue on GitHub
