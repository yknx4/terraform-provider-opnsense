package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yknx4/terraform-provider-opnsense/internal/client"
)

func main() {
	// Parse command-line flags
	url := flag.String("url", os.Getenv("OPNSENSE_URL"), "OPNsense URL")
	apiKey := flag.String("api-key", os.Getenv("OPNSENSE_API_KEY"), "API Key")
	apiSecret := flag.String("api-secret", os.Getenv("OPNSENSE_API_SECRET"), "API Secret")
	insecure := flag.Bool("insecure", os.Getenv("OPNSENSE_INSECURE") == "true", "Skip TLS verification")
	output := flag.String("output", "import.tf", "Output file for import configuration")

	flag.Parse()

	// Validate required flags
	if *url == "" || *apiKey == "" || *apiSecret == "" {
		log.Fatal("Error: url, api-key, and api-secret are required (via flags or environment variables)")
	}

	// Create API client
	c := client.NewClient(*url, *apiKey, *apiSecret, *insecure)

	fmt.Println("Discovering OPNsense resources...")

	// Discover resources
	var importBlocks []string

	// Discover firewall filter rules
	firewallRules, err := discoverFirewallRules(c)
	if err != nil {
		log.Printf("Warning: Failed to discover firewall rules: %v", err)
	} else {
		importBlocks = append(importBlocks, firewallRules...)
		fmt.Printf("Found %d firewall filter rules\n", len(firewallRules))
	}

	// Discover routes
	routes, err := discoverRoutes(c)
	if err != nil {
		log.Printf("Warning: Failed to discover routes: %v", err)
	} else {
		importBlocks = append(importBlocks, routes...)
		fmt.Printf("Found %d routes\n", len(routes))
	}

	// Generate import.tf file
	if len(importBlocks) == 0 {
		fmt.Println("No resources found to import")
		return
	}

	// Create output file
	f, err := os.Create(*output)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	// Write header comment
	f.WriteString("# Generated import configuration for OPNsense resources\n")
	f.WriteString("# Run 'terraform plan -generate-config-out=generated.tf' to generate resource configurations\n\n")

	// Write import blocks
	for _, block := range importBlocks {
		f.WriteString(block)
		f.WriteString("\n")
	}

	fmt.Printf("\nImport configuration written to %s\n", *output)
	fmt.Println("Next steps:")
	fmt.Println("  1. Review the generated import.tf file")
	fmt.Println("  2. Run: terraform plan -generate-config-out=generated.tf")
	fmt.Println("  3. Review and edit generated.tf as needed")
	fmt.Println("  4. Run: terraform apply")
}

func discoverFirewallRules(c *client.Client) ([]string, error) {
	resp, err := c.Get("/api/firewall/filter/searchRule")
	if err != nil {
		return nil, fmt.Errorf("failed to search firewall rules: %w", err)
	}

	var result struct {
		Rows []struct {
			UUID        string `json:"uuid"`
			Description string `json:"description"`
		} `json:"rows"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse firewall rules response: %w", err)
	}

	var blocks []string
	for _, rule := range result.Rows {
		resourceName := sanitizeResourceName(rule.Description)
		if resourceName == "" {
			resourceName = rule.UUID[:8]
		}

		block := fmt.Sprintf(`import {
  to = opnsense_firewall_filter.%s
  id = "%s"
}`, resourceName, rule.UUID)
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func discoverRoutes(c *client.Client) ([]string, error) {
	resp, err := c.Get("/api/routes/routes/searchroute")
	if err != nil {
		return nil, fmt.Errorf("failed to search routes: %w", err)
	}

	var result struct {
		Rows []struct {
			UUID    string `json:"uuid"`
			Descr   string `json:"descr"`
			Network string `json:"network"`
		} `json:"rows"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse routes response: %w", err)
	}

	var blocks []string
	for _, route := range result.Rows {
		resourceName := sanitizeResourceName(route.Descr)
		if resourceName == "" {
			// Use network as fallback
			resourceName = sanitizeResourceName(route.Network)
		}
		if resourceName == "" {
			resourceName = route.UUID[:8]
		}

		block := fmt.Sprintf(`import {
  to = opnsense_route.%s
  id = "%s"
}`, resourceName, route.UUID)
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// sanitizeResourceName converts a description into a valid Terraform resource name
func sanitizeResourceName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace spaces and special characters with underscores
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, name)

	// Remove leading/trailing underscores
	name = strings.Trim(name, "_")

	// Collapse multiple underscores
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}

	return name
}
