package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yknx4/terraform-provider-opnsense/internal/client"
	"github.com/yknx4/terraform-provider-opnsense/internal/resources"
)

var _ provider.Provider = &opnsenseProvider{}

type opnsenseProvider struct {
	version string
}

type opnsenseProviderModel struct {
	URL       types.String `tfsdk:"url"`
	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
	Insecure  types.Bool   `tfsdk:"insecure"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &opnsenseProvider{
			version: version,
		}
	}
}

func (p *opnsenseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opnsense"
	resp.Version = p.version
}

func (p *opnsenseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing OPNsense firewall resources via the OPNsense API.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "The URL of the OPNsense instance (e.g., https://opnsense.example.com). Can also be set via OPNSENSE_URL environment variable.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The API key for OPNsense authentication. Can also be set via OPNSENSE_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"api_secret": schema.StringAttribute{
				Description: "The API secret for OPNsense authentication. Can also be set via OPNSENSE_API_SECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Whether to skip TLS certificate verification. Defaults to false. Can also be set via OPNSENSE_INSECURE environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *opnsenseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config opnsenseProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Use environment variables as defaults if not set in config
	url := os.Getenv("OPNSENSE_URL")
	apiKey := os.Getenv("OPNSENSE_API_KEY")
	apiSecret := os.Getenv("OPNSENSE_API_SECRET")
	insecure := os.Getenv("OPNSENSE_INSECURE") == "true"

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}
	if !config.APISecret.IsNull() {
		apiSecret = config.APISecret.ValueString()
	}
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}

	// Validate that required fields are set
	if url == "" {
		resp.Diagnostics.AddError(
			"Missing OPNsense URL",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the OPNsense URL. "+
				"Set the url value in the configuration or use the OPNSENSE_URL environment variable.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the API key. "+
				"Set the api_key value in the configuration or use the OPNSENSE_API_KEY environment variable.",
		)
	}

	if apiSecret == "" {
		resp.Diagnostics.AddError(
			"Missing API Secret",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the API secret. "+
				"Set the api_secret value in the configuration or use the OPNSENSE_API_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the client
	c := client.NewClient(url, apiKey, apiSecret, insecure)

	// Make the client available to resources and data sources
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *opnsenseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewFirewallFilterResource,
		resources.NewRouteResource,
	}
}

func (p *opnsenseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
