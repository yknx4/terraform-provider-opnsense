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
  insecure   = false
}

# Example firewall filter rule
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

# Example static route
resource "opnsense_route" "example_route" {
  network     = "10.9.0.0/24"
  gateway     = "LAN_DHCP"
  description = "Example route to internal network"
  disabled    = false
}
