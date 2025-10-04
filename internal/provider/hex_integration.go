package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure HexIntegrationResource satisfies various resource interfaces.
var _ resource.Resource = &HexIntegrationResource{}
var _ resource.ResourceWithImportState = &HexIntegrationResource{}

// HexIntegrationResourceModel describes the Hex integration resource data model.
type HexIntegrationResourceModel struct {
	BaseIntegrationResourceModel
	Configuration HexConfigurationModel `tfsdk:"configuration"`
}

// HexConfigurationModel describes the Hex-specific configuration
type HexConfigurationModel struct {
	APIToken     types.String `tfsdk:"api_token"`
	BaseURL      types.String `tfsdk:"base_url"`
	WorkspaceID  types.String `tfsdk:"workspace_id"`
	WorkspaceName types.String `tfsdk:"workspace_name"`
}

// HexIntegrationResource defines the Hex integration resource implementation.
type HexIntegrationResource struct {
	BaseIntegrationResource
}

// NewHexIntegrationResource is a helper function to simplify the provider server and testing implementation.
func NewHexIntegrationResource() resource.Resource {
	return &HexIntegrationResource{}
}

// Metadata returns the resource type name.
func (r *HexIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hex_integration"
}

// Schema defines the schema for the resource.
func (r *HexIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Euno Hex Integration resource",

		Attributes: getCommonAttributes(),
		Blocks:     getCommonBlocks(),
	}

	// Add Hex-specific configuration block
	resp.Schema.Blocks["configuration"] = schema.SingleNestedBlock{
		MarkdownDescription: "Hex-specific configuration",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Hex API token",
			},
			"base_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Hex API base URL (defaults to https://app.hex.tech/api/v1)",
			},
			"workspace_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Hex workspace ID. Used to create links to projects in the workspace and generate URIs",
			},
			"workspace_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Hex workspace name (defaults to hex_workspace)",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *HexIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HexIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration to API format
	configMap := map[string]interface{}{
		"api_token":    data.Configuration.APIToken.ValueString(),
		"workspace_id": data.Configuration.WorkspaceID.ValueString(),
	}

	if !data.Configuration.BaseURL.IsNull() {
		configMap["base_url"] = data.Configuration.BaseURL.ValueString()
	}
	if !data.Configuration.WorkspaceName.IsNull() {
		configMap["workspace_name"] = data.Configuration.WorkspaceName.ValueString()
	}

	// Convert Terraform data to API format
	integration := IntegrationIn{
		IntegrationType: "hex",
		Name:            data.Name.ValueString(),
		Active:          data.Active.ValueBool(),
		Configuration:   configMap,
		Schedule:        convertScheduleToAPI(data.Schedule),
		InvalidationStrategy: convertInvalidationStrategyToAPI(data.InvalidationStrategy),
	}

	// Create the integration
	result, err := r.client.CreateIntegration(ctx, integration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Hex integration, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.ID = types.Int64Value(int64(result.ID))
	data.Name = types.StringValue(result.Name)
	if result.Active != nil {
		data.Active = types.BoolValue(*result.Active)
	}
	data.CreatedAt = types.StringValue(result.CreatedAt)
	data.LastUpdatedAt = types.StringValue(result.LastUpdatedAt)

	// Convert configuration back to Terraform format
	if result.Configuration != nil {
		if apiToken, ok := result.Configuration["api_token"].(string); ok {
			data.Configuration.APIToken = types.StringValue(apiToken)
		}
		if baseURL, ok := result.Configuration["base_url"].(string); ok {
			data.Configuration.BaseURL = types.StringValue(baseURL)
		}
		if workspaceID, ok := result.Configuration["workspace_id"].(string); ok {
			data.Configuration.WorkspaceID = types.StringValue(workspaceID)
		}
		if workspaceName, ok := result.Configuration["workspace_name"].(string); ok {
			data.Configuration.WorkspaceName = types.StringValue(workspaceName)
		}
	}

	// Convert schedule and invalidation strategy back
	data.Schedule = convertScheduleFromAPI(result.Schedule)
	data.InvalidationStrategy = convertInvalidationStrategyFromAPI(result.InvalidationStrategy)

	if result.PendingCredentialsLookupKey != nil {
		data.PendingCredentialsLookupKey = types.StringValue(*result.PendingCredentialsLookupKey)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *HexIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HexIntegrationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the integration from the API
	result, err := r.client.GetIntegration(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Hex integration, got error: %s", err))
		return
	}

	// Map the response back to the model
	data.ID = types.Int64Value(int64(result.ID))
	data.Name = types.StringValue(result.Name)
	if result.Active != nil {
		data.Active = types.BoolValue(*result.Active)
	}
	data.CreatedAt = types.StringValue(result.CreatedAt)
	data.LastUpdatedAt = types.StringValue(result.LastUpdatedAt)

	// Convert configuration back to Terraform format
	if result.Configuration != nil {
		if apiToken, ok := result.Configuration["api_token"].(string); ok {
			data.Configuration.APIToken = types.StringValue(apiToken)
		}
		if baseURL, ok := result.Configuration["base_url"].(string); ok {
			data.Configuration.BaseURL = types.StringValue(baseURL)
		}
		if workspaceID, ok := result.Configuration["workspace_id"].(string); ok {
			data.Configuration.WorkspaceID = types.StringValue(workspaceID)
		}
		if workspaceName, ok := result.Configuration["workspace_name"].(string); ok {
			data.Configuration.WorkspaceName = types.StringValue(workspaceName)
		}
	}

	// Convert schedule and invalidation strategy back
	data.Schedule = convertScheduleFromAPI(result.Schedule)
	data.InvalidationStrategy = convertInvalidationStrategyFromAPI(result.InvalidationStrategy)

	if result.PendingCredentialsLookupKey != nil {
		data.PendingCredentialsLookupKey = types.StringValue(*result.PendingCredentialsLookupKey)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *HexIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HexIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration to API format
	configMap := map[string]interface{}{
		"api_token":    data.Configuration.APIToken.ValueString(),
		"workspace_id": data.Configuration.WorkspaceID.ValueString(),
	}

	if !data.Configuration.BaseURL.IsNull() {
		configMap["base_url"] = data.Configuration.BaseURL.ValueString()
	}
	if !data.Configuration.WorkspaceName.IsNull() {
		configMap["workspace_name"] = data.Configuration.WorkspaceName.ValueString()
	}

	// Convert Terraform data to API format
	integration := IntegrationIn{
		IntegrationType: "hex",
		Name:            data.Name.ValueString(),
		Active:          data.Active.ValueBool(),
		Configuration:   configMap,
		Schedule:        convertScheduleToAPI(data.Schedule),
		InvalidationStrategy: convertInvalidationStrategyToAPI(data.InvalidationStrategy),
	}

	// Update the integration
	result, err := r.client.UpdateIntegration(ctx, int(data.ID.ValueInt64()), integration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Hex integration, got error: %s", err))
		return
	}

	// Update the model with the response data
	data.ID = types.Int64Value(int64(result.ID))
	data.Name = types.StringValue(result.Name)
	if result.Active != nil {
		data.Active = types.BoolValue(*result.Active)
	}
	data.CreatedAt = types.StringValue(result.CreatedAt)
	data.LastUpdatedAt = types.StringValue(result.LastUpdatedAt)

	// Convert configuration back to Terraform format
	if result.Configuration != nil {
		if apiToken, ok := result.Configuration["api_token"].(string); ok {
			data.Configuration.APIToken = types.StringValue(apiToken)
		}
		if baseURL, ok := result.Configuration["base_url"].(string); ok {
			data.Configuration.BaseURL = types.StringValue(baseURL)
		}
		if workspaceID, ok := result.Configuration["workspace_id"].(string); ok {
			data.Configuration.WorkspaceID = types.StringValue(workspaceID)
		}
		if workspaceName, ok := result.Configuration["workspace_name"].(string); ok {
			data.Configuration.WorkspaceName = types.StringValue(workspaceName)
		}
	}

	// Convert schedule and invalidation strategy back
	data.Schedule = convertScheduleFromAPI(result.Schedule)
	data.InvalidationStrategy = convertInvalidationStrategyFromAPI(result.InvalidationStrategy)

	if result.PendingCredentialsLookupKey != nil {
		data.PendingCredentialsLookupKey = types.StringValue(*result.PendingCredentialsLookupKey)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *HexIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HexIntegrationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the integration
	err := r.client.DeleteIntegration(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Hex integration, got error: %s", err))
		return
	}
}

// ImportState imports the resource from the API.
func (r *HexIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.BaseIntegrationResource.ImportState(ctx, req, resp)
}
