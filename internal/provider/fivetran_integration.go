package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure FivetranIntegrationResource satisfies various resource interfaces.
var _ resource.Resource = &FivetranIntegrationResource{}
var _ resource.ResourceWithImportState = &FivetranIntegrationResource{}

// FivetranIntegrationResourceModel describes the Fivetran integration resource data model.
type FivetranIntegrationResourceModel struct {
	BaseIntegrationResourceModel
	Configuration FivetranConfigurationModel `tfsdk:"configuration"`
}

// FivetranConfigurationModel describes the Fivetran-specific configuration
type FivetranConfigurationModel struct {
	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
	BaseURL   types.String `tfsdk:"base_url"`
}

// FivetranIntegrationResource defines the Fivetran integration resource implementation.
type FivetranIntegrationResource struct {
	BaseIntegrationResource
}

// NewFivetranIntegrationResource is a helper function to simplify the provider server and testing implementation.
func NewFivetranIntegrationResource() resource.Resource {
	return &FivetranIntegrationResource{}
}

// Metadata returns the resource type name.
func (r *FivetranIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fivetran_integration"
}

// Schema defines the schema for the resource.
func (r *FivetranIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Euno Fivetran Integration resource",

		Attributes: getCommonAttributes(),
		Blocks:     getCommonBlocks(),
	}

	// Add Fivetran-specific configuration block
	resp.Schema.Blocks["configuration"] = schema.SingleNestedBlock{
		MarkdownDescription: "Fivetran-specific configuration",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Fivetran API key",
			},
			"api_secret": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Fivetran API secret",
			},
			"base_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Fivetran API base URL (defaults to https://api.fivetran.com/v1)",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *FivetranIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FivetranIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration to API format
	configMap := map[string]interface{}{
		"api_key":    data.Configuration.APIKey.ValueString(),
		"api_secret": data.Configuration.APISecret.ValueString(),
	}

	if !data.Configuration.BaseURL.IsNull() {
		configMap["base_url"] = data.Configuration.BaseURL.ValueString()
	}

	// Convert Terraform data to API format
	integration := IntegrationIn{
		IntegrationType: "fivetran",
		Name:            data.Name.ValueString(),
		Active:          data.Active.ValueBool(),
		Configuration:   configMap,
		Schedule:        convertScheduleToAPI(data.Schedule),
		InvalidationStrategy: convertInvalidationStrategyToAPI(data.InvalidationStrategy),
	}

	// Create the integration
	result, err := r.client.CreateIntegration(ctx, integration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Fivetran integration, got error: %s", err))
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
		if apiKey, ok := result.Configuration["api_key"].(string); ok {
			data.Configuration.APIKey = types.StringValue(apiKey)
		}
		if apiSecret, ok := result.Configuration["api_secret"].(string); ok {
			data.Configuration.APISecret = types.StringValue(apiSecret)
		}
		if baseURL, ok := result.Configuration["base_url"].(string); ok {
			data.Configuration.BaseURL = types.StringValue(baseURL)
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
func (r *FivetranIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FivetranIntegrationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the integration from the API
	result, err := r.client.GetIntegration(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Fivetran integration, got error: %s", err))
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
		if apiKey, ok := result.Configuration["api_key"].(string); ok {
			data.Configuration.APIKey = types.StringValue(apiKey)
		}
		if apiSecret, ok := result.Configuration["api_secret"].(string); ok {
			data.Configuration.APISecret = types.StringValue(apiSecret)
		}
		if baseURL, ok := result.Configuration["base_url"].(string); ok {
			data.Configuration.BaseURL = types.StringValue(baseURL)
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
func (r *FivetranIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FivetranIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration to API format
	configMap := map[string]interface{}{
		"api_key":    data.Configuration.APIKey.ValueString(),
		"api_secret": data.Configuration.APISecret.ValueString(),
	}

	if !data.Configuration.BaseURL.IsNull() {
		configMap["base_url"] = data.Configuration.BaseURL.ValueString()
	}

	// Convert Terraform data to API format
	integration := IntegrationIn{
		IntegrationType: "fivetran",
		Name:            data.Name.ValueString(),
		Active:          data.Active.ValueBool(),
		Configuration:   configMap,
		Schedule:        convertScheduleToAPI(data.Schedule),
		InvalidationStrategy: convertInvalidationStrategyToAPI(data.InvalidationStrategy),
	}

	// Update the integration
	result, err := r.client.UpdateIntegration(ctx, int(data.ID.ValueInt64()), integration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Fivetran integration, got error: %s", err))
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
		if apiKey, ok := result.Configuration["api_key"].(string); ok {
			data.Configuration.APIKey = types.StringValue(apiKey)
		}
		if apiSecret, ok := result.Configuration["api_secret"].(string); ok {
			data.Configuration.APISecret = types.StringValue(apiSecret)
		}
		if baseURL, ok := result.Configuration["base_url"].(string); ok {
			data.Configuration.BaseURL = types.StringValue(baseURL)
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
func (r *FivetranIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FivetranIntegrationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the integration
	err := r.client.DeleteIntegration(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Fivetran integration, got error: %s", err))
		return
	}
}

// ImportState imports the resource from the API.
func (r *FivetranIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.BaseIntegrationResource.ImportState(ctx, req, resp)
}
