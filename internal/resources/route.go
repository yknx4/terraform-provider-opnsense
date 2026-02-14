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

var _ resource.Resource = &RouteResource{}
var _ resource.ResourceWithImportState = &RouteResource{}

type RouteResource struct {
	client *client.Client
}

type RouteResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Network     types.String `tfsdk:"network"`
	Gateway     types.String `tfsdk:"gateway"`
	Description types.String `tfsdk:"description"`
	Disabled    types.Bool   `tfsdk:"disabled"`
}

func NewRouteResource() resource.Resource {
	return &RouteResource{}
}

func (r *RouteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

func (r *RouteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an OPNsense static route.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the route.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network": schema.StringAttribute{
				Description: "The destination network in CIDR notation (e.g., '10.0.0.0/24').",
				Required:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "The gateway to use for this route.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the route.",
				Optional:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Whether the route is disabled.",
				Optional:    true,
			},
		},
	}
}

func (r *RouteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RouteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request body
	requestBody := map[string]interface{}{
		"route": map[string]interface{}{
			"network":  plan.Network.ValueString(),
			"gateway":  plan.Gateway.ValueString(),
			"descr":    plan.Description.ValueString(),
			"disabled": plan.Disabled.ValueBool(),
		},
	}

	// Make API call to create the route
	respBody, err := r.client.Post("/api/routes/routes/addroute", requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating route",
			"Could not create route: "+err.Error(),
		)
		return
	}

	// Parse response to get the route ID
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
			"Error extracting route ID",
			"Could not extract route ID from API response",
		)
		return
	}

	// Apply the configuration
	_, err = r.client.Post("/api/routes/routes/reconfigure", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying route configuration",
			"Could not apply route configuration: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RouteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the route from the API
	path := fmt.Sprintf("/api/routes/routes/getroute/%s", state.ID.ValueString())
	respBody, err := r.client.Get(path)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading route",
			"Could not read route: "+err.Error(),
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
	if route, ok := result["route"].(map[string]interface{}); ok {
		if network, ok := route["network"].(string); ok {
			state.Network = types.StringValue(network)
		}
		if gateway, ok := route["gateway"].(string); ok {
			state.Gateway = types.StringValue(gateway)
		}
		if description, ok := route["descr"].(string); ok {
			state.Description = types.StringValue(description)
		}
		if disabled, ok := route["disabled"].(string); ok {
			state.Disabled = types.BoolValue(disabled == "1")
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RouteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request body
	requestBody := map[string]interface{}{
		"route": map[string]interface{}{
			"network":  plan.Network.ValueString(),
			"gateway":  plan.Gateway.ValueString(),
			"descr":    plan.Description.ValueString(),
			"disabled": plan.Disabled.ValueBool(),
		},
	}

	// Make API call to update the route
	path := fmt.Sprintf("/api/routes/routes/setroute/%s", plan.ID.ValueString())
	_, err := r.client.Post(path, requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating route",
			"Could not update route: "+err.Error(),
		)
		return
	}

	// Apply the configuration
	_, err = r.client.Post("/api/routes/routes/reconfigure", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying route configuration",
			"Could not apply route configuration: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RouteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the route via API
	path := fmt.Sprintf("/api/routes/routes/delroute/%s", state.ID.ValueString())
	_, err := r.client.Post(path, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting route",
			"Could not delete route: "+err.Error(),
		)
		return
	}

	// Apply the configuration
	_, err = r.client.Post("/api/routes/routes/reconfigure", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying route configuration",
			"Could not apply route configuration: "+err.Error(),
		)
		return
	}
}

func (r *RouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
