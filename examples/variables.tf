variable "opnsense_url" {
  description = "The URL of the OPNsense instance"
  type        = string
}

variable "opnsense_api_key" {
  description = "The API key for OPNsense"
  type        = string
  sensitive   = true
}

variable "opnsense_api_secret" {
  description = "The API secret for OPNsense"
  type        = string
  sensitive   = true
}

variable "opnsense_insecure" {
  description = "Whether to skip TLS certificate verification"
  type        = bool
  default     = false
}
