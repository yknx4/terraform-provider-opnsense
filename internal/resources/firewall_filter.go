package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yknx4/terraform-provider-opnsense/internal/client"
)

var _ resource.Resource = &FirewallFilterResource{}
var _ resource.ResourceWithImportState = &FirewallFilterResource{}

type FirewallFilterResource struct {
	client *client.Client
}

type FirewallFilterResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Description types.String `tfsdk:"description"`
	Interface   types.String `tfsdk:"interface"`
	Action      types.String `tfsdk:"action"`
	Direction   types.String `tfsdk:"direction"`
	Protocol    types.String `tfsdk:"protocol"`
	Source      types.String `tfsdk:"source"`
	Destination types.String `tfsdk:"destination"`
	SourcePort  types.String `tfsdk:"source_port"`
	DestPort    types.String `tfsdk:"dest_port"`
	Log         types.Bool   `tfsdk:"log"`
}

func NewFirewallFilterResource() resource.Resource {
	return &FirewallFilterResource{}
}

func (r *FirewallFilterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_filter"
}

func (r *FirewallFilterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an OPNsense firewall filter rule.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the firewall rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the firewall rule is enabled.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the firewall rule.",
				Optional:    true,
			},
			"interface": schema.StringAttribute{
				Description: "The interface on which the rule applies (e.g., 'wan', 'lan').",
				Required:    true,
			},
			"action": schema.StringAttribute{
				Description: "The action to take (e.g., 'pass', 'block', 'reject').",
				Required:    true,
			},
			"direction": schema.StringAttribute{
				Description: "The direction of traffic (e.g., 'in', 'out').",
				Required:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "The protocol (e.g., 'TCP', 'UDP', 'ICMP', 'any').",
				Optional:    true,
			},
			"source": schema.StringAttribute{
				Description: "The source address or network.",
				Optional:    true,
			},
			"destination": schema.StringAttribute{
				Description: "The destination address or network.",
				Optional:    true,
			},
			"source_port": schema.StringAttribute{
				Description: "The source port or port range.",
				Optional:    true,
			},
			"dest_port": schema.StringAttribute{
				Description: "The destination port or port range.",
				Optional:    true,
			},
			"log": schema.BoolAttribute{
				Description: "Whether to log matches for this rule.",
				Optional:    true,
			},
		},
	}
}

func (r *FirewallFilterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *FirewallFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallFilterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request body
	requestBody := map[string]interface{}{
		"rule": map[string]interface{}{
			"enabled":     plan.Enabled.ValueBool(),
			"description": plan.Description.ValueString(),
			"interface":   plan.Interface.ValueString(),
			"action":      plan.Action.ValueString(),
			"direction":   plan.Direction.ValueString(),
			"protocol":    plan.Protocol.ValueString(),
			"source":      plan.Source.ValueString(),
			"destination": plan.Destination.ValueString(),
			"source_port": plan.SourcePort.ValueString(),
			"dest_port":   plan.DestPort.ValueString(),
			"log":         plan.Log.ValueBool(),
		},
	}

	// Make API call to create the rule
	respBody, err := r.client.Post("/api/firewall/filter/addRule", requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating firewall filter",
			"Could not create firewall filter: "+err.Error(),
		)
		return
	}

	// Parse response to get the rule ID
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing response",
			"Could not parse API response: "+err.Error(),
		)
		return
	}

	// Extract UUID from response
	if uuid, ok := result["uuid"].(string); ok {
		plan.ID = types.StringValue(uuid)
	} else {
		resp.Diagnostics.AddError(
			"Error extracting rule ID",
			"Could not extract rule ID from API response",
		)
		return
	}

	// Apply the configuration
	_, err = r.client.Post("/api/firewall/filter/apply", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying firewall configuration",
			"Could not apply firewall configuration: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallFilterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the rule from the API
	path := fmt.Sprintf("/api/firewall/filter/getRule/%s", state.ID.ValueString())
	respBody, err := r.client.Get(path)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading firewall filter",
			"Could not read firewall filter: "+err.Error(),
		)
		return
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing response",
			"Could not parse API response: "+err.Error(),
		)
		return
	}

	// Update state with values from API
	if rule, ok := result["rule"].(map[string]interface{}); ok {
		if enabled, ok := rule["enabled"].(string); ok {
			state.Enabled = types.BoolValue(enabled == "1")
		}
		if description, ok := rule["description"].(string); ok {
			state.Description = types.StringValue(description)
		}
		if iface, ok := rule["interface"].(string); ok {
			state.Interface = types.StringValue(iface)
		}
		if action, ok := rule["action"].(string); ok {
			state.Action = types.StringValue(action)
		}
		if direction, ok := rule["direction"].(string); ok {
			state.Direction = types.StringValue(direction)
		}
		if protocol, ok := rule["protocol"].(string); ok {
			state.Protocol = types.StringValue(protocol)
		}
		if source, ok := rule["source"].(string); ok {
			state.Source = types.StringValue(source)
		}
		if destination, ok := rule["destination"].(string); ok {
			state.Destination = types.StringValue(destination)
		}
		if sourcePort, ok := rule["source_port"].(string); ok {
			state.SourcePort = types.StringValue(sourcePort)
		}
		if destPort, ok := rule["dest_port"].(string); ok {
			state.DestPort = types.StringValue(destPort)
		}
		if log, ok := rule["log"].(string); ok {
			state.Log = types.BoolValue(log == "1")
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallFilterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request body
	requestBody := map[string]interface{}{
		"rule": map[string]interface{}{
			"enabled":     plan.Enabled.ValueBool(),
			"description": plan.Description.ValueString(),
			"interface":   plan.Interface.ValueString(),
			"action":      plan.Action.ValueString(),
			"direction":   plan.Direction.ValueString(),
			"protocol":    plan.Protocol.ValueString(),
			"source":      plan.Source.ValueString(),
			"destination": plan.Destination.ValueString(),
			"source_port": plan.SourcePort.ValueString(),
			"dest_port":   plan.DestPort.ValueString(),
			"log":         plan.Log.ValueBool(),
		},
	}

	// Make API call to update the rule
	path := fmt.Sprintf("/api/firewall/filter/setRule/%s", plan.ID.ValueString())
	_, err := r.client.Post(path, requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating firewall filter",
			"Could not update firewall filter: "+err.Error(),
		)
		return
	}

	// Apply the configuration
	_, err = r.client.Post("/api/firewall/filter/apply", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying firewall configuration",
			"Could not apply firewall configuration: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallFilterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the rule via API
	path := fmt.Sprintf("/api/firewall/filter/delRule/%s", state.ID.ValueString())
	_, err := r.client.Post(path, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting firewall filter",
			"Could not delete firewall filter: "+err.Error(),
		)
		return
	}

	// Apply the configuration
	_, err = r.client.Post("/api/firewall/filter/apply", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying firewall configuration",
			"Could not apply firewall configuration: "+err.Error(),
		)
		return
	}
}

func (r *FirewallFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
