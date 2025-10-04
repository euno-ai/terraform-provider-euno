package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure DbtCoreIntegrationResource satisfies various resource interfaces.
var _ resource.Resource = &DbtCoreIntegrationResource{}
var _ resource.ResourceWithImportState = &DbtCoreIntegrationResource{}

// DbtCoreIntegrationResourceModel describes the DBT Core integration resource data model.
type DbtCoreIntegrationResourceModel struct {
	ID                        types.Int64  `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Active                    types.Bool   `tfsdk:"active"`
	TriggerSecret             types.String `tfsdk:"trigger_secret"`
	TriggerURL                types.String `tfsdk:"trigger_url"`
	InvalidationStrategy      *InvalidationStrategyModel `tfsdk:"invalidation_strategy"`
	PendingCredentialsLookupKey types.String `tfsdk:"pending_credentials_lookup_key"`
	LastUpdatedAt             types.String `tfsdk:"last_updated_at"`
	CreatedAt                 types.String `tfsdk:"created_at"`
	Configuration             DbtCoreConfigurationModel `tfsdk:"configuration"`
}

// DbtCoreConfigurationModel describes the DBT Core-specific configuration
type DbtCoreConfigurationModel struct {
	SchemasAliases                                  types.Map    `tfsdk:"schemas_aliases"`
	RepositoryURL                                   types.String `tfsdk:"repository_url"`
	BuildTarget                                     types.String `tfsdk:"build_target"`
	StageBuildTarget                                types.String `tfsdk:"stage_build_target"`
	RepositoryBranch                                types.String `tfsdk:"repository_branch"`
	DbtProjectRootDirectoryInRepository            types.String `tfsdk:"dbt_project_root_directory_in_repository"`
	RepositoryRevision                              types.String `tfsdk:"repository_revision"`
	AllowResourcesWithNoCatalogEntry               types.Bool  `tfsdk:"allow_resources_with_no_catalog_entry"`
	OverrideURIPrefix                               types.String `tfsdk:"override_uri_prefix"`
}

// DbtCoreIntegrationResource defines the DBT Core integration resource implementation.
type DbtCoreIntegrationResource struct {
	BaseIntegrationResource
}

// NewDbtCoreIntegrationResource is a helper function to simplify the provider server and testing implementation.
func NewDbtCoreIntegrationResource() resource.Resource {
	return &DbtCoreIntegrationResource{}
}

// Metadata returns the resource type name.
func (r *DbtCoreIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbt_core_integration"
}

// Schema defines the schema for the resource.
func (r *DbtCoreIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Get common attributes (which now include trigger_secret and trigger_url)
	attrs := getCommonAttributes()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Euno DBT Core Integration resource (push integration)",
		Attributes:           attrs,
		Blocks:              getCommonBlocksForPush(),
	}

	// Add DBT Core-specific configuration block
	resp.Schema.Blocks["configuration"] = schema.SingleNestedBlock{
		MarkdownDescription: "DBT Core-specific configuration",
		Attributes: map[string]schema.Attribute{
			"schemas_aliases": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A dictionary of schema aliases, where the keys and values have the template db.schema. Euno will ingest dbt resources (nodes and sources) to the database and schema stated in the manifest file, unless the database.schema combination appears in this mapping.",
			},
			"repository_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The URL of the git repository where the dbt project is stored",
			},
			"build_target": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The dht target to build",
			},
			"stage_build_target": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The stage dbt target to build",
			},
			"repository_branch": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The branch of the git repository where the dbt project is stored",
			},
			"dbt_project_root_directory_in_repository": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The subdirectory within the git repository where the dbt project is stored (defaults to '/')",
			},
			"repository_revision": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The revision of the git repository where the dbt project is stored",
			},
			"allow_resources_with_no_catalog_entry": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether to allow dbt resources with no corresponding catalog entry to be ingested (defaults to false)",
			},
			"override_uri_prefix": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The prefix to override the URI of the resources. If not set, we use 'dbt'.<dbt project name>",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *DbtCoreIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DbtCoreIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize trigger attributes for push integration
	data.TriggerSecret = types.StringNull()
	data.TriggerURL = types.StringNull()

	// Convert configuration to API format
	configMap := map[string]interface{}{
		"build_target": data.Configuration.BuildTarget.ValueString(),
	}

	// Add optional fields
	if !data.Configuration.SchemasAliases.IsNull() {
		schemasAliases := make(map[string]string)
		for key, value := range data.Configuration.SchemasAliases.Elements() {
			if strValue, ok := value.(types.String); ok {
				schemasAliases[key] = strValue.ValueString()
			}
		}
		configMap["schemas_aliases"] = schemasAliases
	}
	if !data.Configuration.RepositoryURL.IsNull() {
		configMap["repository_url"] = data.Configuration.RepositoryURL.ValueString()
	}
	if !data.Configuration.StageBuildTarget.IsNull() {
		configMap["stage_build_target"] = data.Configuration.StageBuildTarget.ValueString()
	}
	if !data.Configuration.RepositoryBranch.IsNull() {
		configMap["repository_branch"] = data.Configuration.RepositoryBranch.ValueString()
	}
	if !data.Configuration.DbtProjectRootDirectoryInRepository.IsNull() {
		configMap["dbt_project_root_directory_in_repository"] = data.Configuration.DbtProjectRootDirectoryInRepository.ValueString()
	}
	if !data.Configuration.RepositoryRevision.IsNull() {
		configMap["repository_revision"] = data.Configuration.RepositoryRevision.ValueString()
	}
	if !data.Configuration.AllowResourcesWithNoCatalogEntry.IsNull() {
		configMap["allow_resources_with_no_catalog_entry"] = data.Configuration.AllowResourcesWithNoCatalogEntry.ValueBool()
	}
	if !data.Configuration.OverrideURIPrefix.IsNull() {
		configMap["override_uri_prefix"] = data.Configuration.OverrideURIPrefix.ValueString()
	}

	// Convert Terraform data to API format
	// Note: For push integrations, we don't include schedule
	integration := IntegrationIn{
		IntegrationType: "dbt_core",
		Name:            data.Name.ValueString(),
		Active:          data.Active.ValueBool(),
		Configuration:   configMap,
		InvalidationStrategy: convertInvalidationStrategyToAPI(data.InvalidationStrategy),
	}

	// Create the integration
	result, err := r.client.CreateIntegration(ctx, integration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create DBT Core integration, got error: %s", err))
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

	// Convert trigger back to Terraform format for push integrations
	// Convert trigger values back to Terraform attributes for push integrations
	if result.TriggerSecret != nil {
		data.TriggerSecret = types.StringValue(*result.TriggerSecret)
	} else {
		data.TriggerSecret = types.StringNull()
	}
	if result.TriggerURL != nil {
		data.TriggerURL = types.StringValue(*result.TriggerURL)
	} else {
		data.TriggerURL = types.StringNull()
	}

	// Convert configuration back to Terraform format
	if result.Configuration != nil {
		if buildTarget, ok := result.Configuration["build_target"].(string); ok {
			data.Configuration.BuildTarget = types.StringValue(buildTarget)
		}
		if schemasAliases, ok := result.Configuration["schemas_aliases"].(map[string]interface{}); ok {
			schemasAliasesMap := make(map[string]attr.Value)
			for key, value := range schemasAliases {
				if strValue, ok := value.(string); ok {
					schemasAliasesMap[key] = types.StringValue(strValue)
				}
			}
			data.Configuration.SchemasAliases = types.MapValueMust(types.StringType, schemasAliasesMap)
		}
		if repositoryURL, ok := result.Configuration["repository_url"].(string); ok {
			data.Configuration.RepositoryURL = types.StringValue(repositoryURL)
		}
		if stageBuildTarget, ok := result.Configuration["stage_build_target"].(string); ok {
			data.Configuration.StageBuildTarget = types.StringValue(stageBuildTarget)
		}
		if repositoryBranch, ok := result.Configuration["repository_branch"].(string); ok {
			data.Configuration.RepositoryBranch = types.StringValue(repositoryBranch)
		}
		if dbtProjectRootDirectoryInRepository, ok := result.Configuration["dbt_project_root_directory_in_repository"].(string); ok {
			data.Configuration.DbtProjectRootDirectoryInRepository = types.StringValue(dbtProjectRootDirectoryInRepository)
		}
		if repositoryRevision, ok := result.Configuration["repository_revision"].(string); ok {
			data.Configuration.RepositoryRevision = types.StringValue(repositoryRevision)
		}
		if allowResourcesWithNoCatalogEntry, ok := result.Configuration["allow_resources_with_no_catalog_entry"].(bool); ok {
			data.Configuration.AllowResourcesWithNoCatalogEntry = types.BoolValue(allowResourcesWithNoCatalogEntry)
		}
		if overrideURIPrefix, ok := result.Configuration["override_uri_prefix"].(string); ok {
			data.Configuration.OverrideURIPrefix = types.StringValue(overrideURIPrefix)
		}
	}

	// Convert invalidation strategy back
	data.InvalidationStrategy = convertInvalidationStrategyFromAPI(result.InvalidationStrategy)

	if result.PendingCredentialsLookupKey != nil {
		data.PendingCredentialsLookupKey = types.StringValue(*result.PendingCredentialsLookupKey)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *DbtCoreIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DbtCoreIntegrationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the integration from the API
	result, err := r.client.GetIntegration(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read DBT Core integration, got error: %s", err))
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

	// Convert trigger back to Terraform format for push integrations
	// Convert trigger values back to Terraform attributes for push integrations
	if result.TriggerSecret != nil {
		data.TriggerSecret = types.StringValue(*result.TriggerSecret)
	} else {
		data.TriggerSecret = types.StringNull()
	}
	if result.TriggerURL != nil {
		data.TriggerURL = types.StringValue(*result.TriggerURL)
	} else {
		data.TriggerURL = types.StringNull()
	}

	// Convert configuration back to Terraform format (same as Create)
	if result.Configuration != nil {
		if buildTarget, ok := result.Configuration["build_target"].(string); ok {
			data.Configuration.BuildTarget = types.StringValue(buildTarget)
		}
		// Add other fields as needed...
	}

	// Convert invalidation strategy back
	data.InvalidationStrategy = convertInvalidationStrategyFromAPI(result.InvalidationStrategy)

	if result.PendingCredentialsLookupKey != nil {
		data.PendingCredentialsLookupKey = types.StringValue(*result.PendingCredentialsLookupKey)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *DbtCoreIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DbtCoreIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration to API format (same as Create)
	configMap := map[string]interface{}{
		"build_target": data.Configuration.BuildTarget.ValueString(),
	}

	// Add optional fields (same logic as Create)
	if !data.Configuration.SchemasAliases.IsNull() {
		schemasAliases := make(map[string]string)
		for key, value := range data.Configuration.SchemasAliases.Elements() {
			if strValue, ok := value.(types.String); ok {
				schemasAliases[key] = strValue.ValueString()
			}
		}
		configMap["schemas_aliases"] = schemasAliases
	}
	// Add other optional fields...

	// Convert Terraform data to API format
	// Note: For push integrations, we don't include schedule
	integration := IntegrationIn{
		IntegrationType: "dbt_core",
		Name:            data.Name.ValueString(),
		Active:          data.Active.ValueBool(),
		Configuration:   configMap,
		InvalidationStrategy: convertInvalidationStrategyToAPI(data.InvalidationStrategy),
	}

	// Update the integration
	result, err := r.client.UpdateIntegration(ctx, int(data.ID.ValueInt64()), integration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update DBT Core integration, got error: %s", err))
		return
	}

	// Update the model with the response data (same as Create)
	data.ID = types.Int64Value(int64(result.ID))
	data.Name = types.StringValue(result.Name)
	if result.Active != nil {
		data.Active = types.BoolValue(*result.Active)
	}
	data.CreatedAt = types.StringValue(result.CreatedAt)
	data.LastUpdatedAt = types.StringValue(result.LastUpdatedAt)

	// Convert trigger back to Terraform format for push integrations
	// Convert trigger values back to Terraform attributes for push integrations
	if result.TriggerSecret != nil {
		data.TriggerSecret = types.StringValue(*result.TriggerSecret)
	} else {
		data.TriggerSecret = types.StringNull()
	}
	if result.TriggerURL != nil {
		data.TriggerURL = types.StringValue(*result.TriggerURL)
	} else {
		data.TriggerURL = types.StringNull()
	}

	// Convert configuration back to Terraform format (same as Create)
	if result.Configuration != nil {
		if buildTarget, ok := result.Configuration["build_target"].(string); ok {
			data.Configuration.BuildTarget = types.StringValue(buildTarget)
		}
		// Add other fields as needed...
	}

	// Convert invalidation strategy back
	data.InvalidationStrategy = convertInvalidationStrategyFromAPI(result.InvalidationStrategy)

	if result.PendingCredentialsLookupKey != nil {
		data.PendingCredentialsLookupKey = types.StringValue(*result.PendingCredentialsLookupKey)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *DbtCoreIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DbtCoreIntegrationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the integration
	err := r.client.DeleteIntegration(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete DBT Core integration, got error: %s", err))
		return
	}
}

// ImportState imports the resource from the API.
func (r *DbtCoreIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.BaseIntegrationResource.ImportState(ctx, req, resp)
}
