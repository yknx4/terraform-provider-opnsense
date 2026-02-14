terraform {
  required_providers {
    opnsense = {
      source = "yknx4/opnsense"
    }
  }
}

provider "opnsense" {
  url        = var.opnsense_url
  api_key    = var.opnsense_api_key
  api_secret = var.opnsense_api_secret
  insecure   = var.opnsense_insecure
}

# Example firewall filter rule - Allow HTTPS
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

# Example firewall filter rule - Allow SSH from specific network
resource "opnsense_firewall_filter" "allow_ssh_from_trusted" {
  enabled     = true
  description = "Allow SSH from trusted network"
  interface   = "wan"
  action      = "pass"
  direction   = "in"
  protocol    = "TCP"
  source      = "10.0.0.0/8"
  destination = "192.168.1.100"
  dest_port   = "22"
  log         = true
}

# Example firewall filter rule - Block all other inbound
resource "opnsense_firewall_filter" "block_inbound_default" {
  enabled     = true
  description = "Block all other inbound traffic"
  interface   = "wan"
  action      = "block"
  direction   = "in"
  protocol    = "any"
  source      = "any"
  destination = "any"
  log         = true
}

# Example static route
resource "opnsense_route" "example_route" {
  network     = "10.9.0.0/24"
  gateway     = "LAN_DHCP"
  description = "Example route to internal network"
  disabled    = false
}

# Example static route - VPN network
resource "opnsense_route" "vpn_route" {
  network     = "172.16.0.0/16"
  gateway     = "VPN_GATEWAY"
  description = "Route to VPN network"
  disabled    = false
}

