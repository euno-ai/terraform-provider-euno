package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure SnowflakeIntegrationResource satisfies various resource interfaces.
var _ resource.Resource = &SnowflakeIntegrationResource{}
var _ resource.ResourceWithImportState = &SnowflakeIntegrationResource{}

// SnowflakeIntegrationResourceModel describes the Snowflake integration resource data model.
type SnowflakeIntegrationResourceModel struct {
	BaseIntegrationResourceModel
	Configuration SnowflakeConfigurationModel `tfsdk:"configuration"`
}

// SnowflakeConfigurationModel describes the Snowflake-specific configuration
type SnowflakeConfigurationModel struct {
	Host                                      types.String  `tfsdk:"host"`
	User                                      types.String  `tfsdk:"user"`
	Password                                  types.String  `tfsdk:"password"`
	PrivateKey                                types.String  `tfsdk:"private_key"`
	Role                                      types.String  `tfsdk:"role"`
	Warehouse                                 types.String  `tfsdk:"warehouse"`
	Database                                  types.String  `tfsdk:"database"`
	TableToUseForQueryHistory                 types.String  `tfsdk:"table_to_use_for_query_history"`
	AdditionalWhereClauseForQueryHistoryQuery types.String  `tfsdk:"additional_where_clause_for_query_history_query"`
	OverridePlatformURI                       types.String  `tfsdk:"override_platform_uri"`
	OverrideBaseURI                           types.String  `tfsdk:"override_base_uri"`
	ExtractViews                              types.Bool    `tfsdk:"extract_views"`
	ExtractTables                             types.Bool    `tfsdk:"extract_tables"`
	ExtractTableauUsage                       types.Bool    `tfsdk:"extract_tableau_usage"`
	ExtractDailyUsage                         types.Bool    `tfsdk:"extract_daily_usage"`
	ExtractDailyDMLSummary                    types.Bool    `tfsdk:"extract_daily_dml_summary"`
	ExtractMaterializedViewsRefreshHistory    types.Bool    `tfsdk:"extract_materialized_views_refresh_history"`
	ExtractHexUsage                           types.Bool    `tfsdk:"extract_hex_usage"`
	ExtractHexLineage                         types.Bool    `tfsdk:"extract_hex_lineage"`
	ExtractHexLineageLookbackDays             types.Int64   `tfsdk:"extract_hex_lineage_lookback_days"`
	CostPerCredit                             types.Float64 `tfsdk:"cost_per_credit"`
	StorageCostPerTB                          types.Float64 `tfsdk:"storage_cost_per_tb"`
	ObserveWarehouses                         types.Bool    `tfsdk:"observe_warehouses"`
	UseSnowflakeDatabase                      types.Bool    `tfsdk:"use_snowflake_database"`
	ExtractLineageFromQueryHistory            types.Bool    `tfsdk:"extract_lineage_from_query_history"`
	LineageLookbackDays                       types.Int64   `tfsdk:"lineage_lookback_days"`
	ObserveInboundShares                      types.Bool    `tfsdk:"observe_inbound_shares"`
}

// SnowflakeIntegrationResource defines the Snowflake integration resource implementation.
type SnowflakeIntegrationResource struct {
	BaseIntegrationResource
}

// NewSnowflakeIntegrationResource is a helper function to simplify the provider server and testing implementation.
func NewSnowflakeIntegrationResource() resource.Resource {
	return &SnowflakeIntegrationResource{}
}

// Metadata returns the resource type name.
func (r *SnowflakeIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snowflake_integration"
}

// Schema defines the schema for the resource.
func (r *SnowflakeIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Euno Snowflake Integration resource",

		Attributes: getCommonAttributes(),
		Blocks:     getCommonBlocks(),
	}

	// Add Snowflake-specific configuration block
	resp.Schema.Blocks["configuration"] = schema.SingleNestedBlock{
		MarkdownDescription: "Snowflake-specific configuration",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Snowflake host",
			},
			"user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Snowflake user",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Snowflake password (deprecated, use private_key instead)",
			},
			"private_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Snowflake private key for key-pair authentication",
			},
			"role": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Snowflake role",
			},
			"warehouse": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Snowflake warehouse",
			},
			"database": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Snowflake database",
			},
			"table_to_use_for_query_history": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Table to use for query history (defaults to snowflake.account_usage.query_history)",
			},
			"additional_where_clause_for_query_history_query": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Additional WHERE clause to add to the query history query",
			},
			"override_platform_uri": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "String to use for the URI. If not provided, the host will be used",
			},
			"override_base_uri": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "String to use for the base URI. If not provided, the host will be used",
			},
			"extract_views": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract views (defaults to true)",
			},
			"extract_tables": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract tables (defaults to true)",
			},
			"extract_tableau_usage": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract Tableau usage (defaults to true)",
			},
			"extract_daily_usage": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract daily usage (defaults to true)",
			},
			"extract_daily_dml_summary": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract daily DML summary (defaults to true)",
			},
			"extract_materialized_views_refresh_history": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract materialized views refresh history (defaults to false)",
			},
			"extract_hex_usage": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract Hex usage (defaults to false)",
			},
			"extract_hex_lineage": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract Hex lineage (defaults to false)",
			},
			"extract_hex_lineage_lookback_days": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Number of days to look back for Hex lineage (defaults to 7)",
			},
			"cost_per_credit": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Cost per credit in dollars (defaults to 3.0)",
			},
			"storage_cost_per_tb": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Storage cost per TB in dollars (defaults to 23)",
			},
			"observe_warehouses": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether to observe warehouse information (defaults to false)",
			},
			"use_snowflake_database": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Use Snowflake system database to poll views (defaults to false)",
			},
			"extract_lineage_from_query_history": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Extract lineage from query history (defaults to true)",
			},
			"lineage_lookback_days": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Number of days to look back for lineage (defaults to 7)",
			},
			"observe_inbound_shares": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Observe Inbound Snowflake Shares (defaults to true)",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *SnowflakeIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SnowflakeIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration to API format
	configMap := map[string]interface{}{
		"host": data.Configuration.Host.ValueString(),
		"user": data.Configuration.User.ValueString(),
	}

	// Add optional fields
	if !data.Configuration.Password.IsNull() {
		configMap["password"] = data.Configuration.Password.ValueString()
	}
	if !data.Configuration.PrivateKey.IsNull() {
		configMap["private_key"] = data.Configuration.PrivateKey.ValueString()
	}
	if !data.Configuration.Role.IsNull() {
		configMap["role"] = data.Configuration.Role.ValueString()
	}
	if !data.Configuration.Warehouse.IsNull() {
		configMap["warehouse"] = data.Configuration.Warehouse.ValueString()
	}
	if !data.Configuration.Database.IsNull() {
		configMap["database"] = data.Configuration.Database.ValueString()
	}
	if !data.Configuration.TableToUseForQueryHistory.IsNull() {
		configMap["table_to_use_for_query_history"] = data.Configuration.TableToUseForQueryHistory.ValueString()
	}
	if !data.Configuration.AdditionalWhereClauseForQueryHistoryQuery.IsNull() {
		configMap["additional_where_clause_for_query_history_query"] = data.Configuration.AdditionalWhereClauseForQueryHistoryQuery.ValueString()
	}
	if !data.Configuration.OverridePlatformURI.IsNull() {
		configMap["override_platform_uri"] = data.Configuration.OverridePlatformURI.ValueString()
	}
	if !data.Configuration.OverrideBaseURI.IsNull() {
		configMap["override_base_uri"] = data.Configuration.OverrideBaseURI.ValueString()
	}
	if !data.Configuration.ExtractViews.IsNull() {
		configMap["extract_views"] = data.Configuration.ExtractViews.ValueBool()
	}
	if !data.Configuration.ExtractTables.IsNull() {
		configMap["extract_tables"] = data.Configuration.ExtractTables.ValueBool()
	}
	if !data.Configuration.ExtractTableauUsage.IsNull() {
		configMap["extract_tableau_usage"] = data.Configuration.ExtractTableauUsage.ValueBool()
	}
	if !data.Configuration.ExtractDailyUsage.IsNull() {
		configMap["extract_daily_usage"] = data.Configuration.ExtractDailyUsage.ValueBool()
	}
	if !data.Configuration.ExtractDailyDMLSummary.IsNull() {
		configMap["extract_daily_dml_summary"] = data.Configuration.ExtractDailyDMLSummary.ValueBool()
	}
	if !data.Configuration.ExtractMaterializedViewsRefreshHistory.IsNull() {
		configMap["extract_materialized_views_refresh_history"] = data.Configuration.ExtractMaterializedViewsRefreshHistory.ValueBool()
	}
	if !data.Configuration.ExtractHexUsage.IsNull() {
		configMap["extract_hex_usage"] = data.Configuration.ExtractHexUsage.ValueBool()
	}
	if !data.Configuration.ExtractHexLineage.IsNull() {
		configMap["extract_hex_lineage"] = data.Configuration.ExtractHexLineage.ValueBool()
	}
	if !data.Configuration.ExtractHexLineageLookbackDays.IsNull() {
		configMap["extract_hex_lineage_lookback_days"] = data.Configuration.ExtractHexLineageLookbackDays.ValueInt64()
	}
	if !data.Configuration.CostPerCredit.IsNull() {
		configMap["cost_per_credit"] = data.Configuration.CostPerCredit.ValueFloat64()
	}
	if !data.Configuration.StorageCostPerTB.IsNull() {
		configMap["storage_cost_per_tb"] = data.Configuration.StorageCostPerTB.ValueFloat64()
	}
	if !data.Configuration.ObserveWarehouses.IsNull() {
		configMap["observe_warehouses"] = data.Configuration.ObserveWarehouses.ValueBool()
	}
	if !data.Configuration.UseSnowflakeDatabase.IsNull() {
		configMap["use_snowflake_database"] = data.Configuration.UseSnowflakeDatabase.ValueBool()
	}
	if !data.Configuration.ExtractLineageFromQueryHistory.IsNull() {
		configMap["extract_lineage_from_query_history"] = data.Configuration.ExtractLineageFromQueryHistory.ValueBool()
	}
	if !data.Configuration.LineageLookbackDays.IsNull() {
		configMap["lineage_lookback_days"] = data.Configuration.LineageLookbackDays.ValueInt64()
	}
	if !data.Configuration.ObserveInboundShares.IsNull() {
		configMap["observe_inbound_shares"] = data.Configuration.ObserveInboundShares.ValueBool()
	}

	// Convert Terraform data to API format
	integration := IntegrationIn{
		IntegrationType:      "snowflake",
		Name:                 data.Name.ValueString(),
		Active:               data.Active.ValueBool(),
		Configuration:        configMap,
		Schedule:             convertScheduleToAPI(data.Schedule),
		InvalidationStrategy: convertInvalidationStrategyToAPI(data.InvalidationStrategy),
	}

	// Create the integration
	result, err := r.client.CreateIntegration(ctx, integration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Snowflake integration, got error: %s", err))
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
		if host, ok := result.Configuration["host"].(string); ok {
			data.Configuration.Host = types.StringValue(host)
		}
		if user, ok := result.Configuration["user"].(string); ok {
			data.Configuration.User = types.StringValue(user)
		}
		if password, ok := result.Configuration["password"].(string); ok {
			data.Configuration.Password = types.StringValue(password)
		}
		if privateKey, ok := result.Configuration["private_key"].(string); ok {
			data.Configuration.PrivateKey = types.StringValue(privateKey)
		}
		if role, ok := result.Configuration["role"].(string); ok {
			data.Configuration.Role = types.StringValue(role)
		}
		if warehouse, ok := result.Configuration["warehouse"].(string); ok {
			data.Configuration.Warehouse = types.StringValue(warehouse)
		}
		if database, ok := result.Configuration["database"].(string); ok {
			data.Configuration.Database = types.StringValue(database)
		}
		// Add other fields as needed...
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
func (r *SnowflakeIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SnowflakeIntegrationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the integration from the API
	result, err := r.client.GetIntegration(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Snowflake integration, got error: %s", err))
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
		if host, ok := result.Configuration["host"].(string); ok {
			data.Configuration.Host = types.StringValue(host)
		}
		if user, ok := result.Configuration["user"].(string); ok {
			data.Configuration.User = types.StringValue(user)
		}
		// Add other fields as needed...
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
func (r *SnowflakeIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SnowflakeIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration to API format (same as Create)
	configMap := map[string]interface{}{
		"host": data.Configuration.Host.ValueString(),
		"user": data.Configuration.User.ValueString(),
	}

	// Add optional fields (same logic as Create)
	if !data.Configuration.Password.IsNull() {
		configMap["password"] = data.Configuration.Password.ValueString()
	}
	if !data.Configuration.PrivateKey.IsNull() {
		configMap["private_key"] = data.Configuration.PrivateKey.ValueString()
	}
	// Add other optional fields...

	// Convert Terraform data to API format
	integration := IntegrationIn{
		IntegrationType:      "snowflake",
		Name:                 data.Name.ValueString(),
		Active:               data.Active.ValueBool(),
		Configuration:        configMap,
		Schedule:             convertScheduleToAPI(data.Schedule),
		InvalidationStrategy: convertInvalidationStrategyToAPI(data.InvalidationStrategy),
	}

	// Update the integration
	result, err := r.client.UpdateIntegration(ctx, int(data.ID.ValueInt64()), integration)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Snowflake integration, got error: %s", err))
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

	// Convert configuration back to Terraform format (same as Create)
	if result.Configuration != nil {
		if host, ok := result.Configuration["host"].(string); ok {
			data.Configuration.Host = types.StringValue(host)
		}
		if user, ok := result.Configuration["user"].(string); ok {
			data.Configuration.User = types.StringValue(user)
		}
		// Add other fields as needed...
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
func (r *SnowflakeIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SnowflakeIntegrationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the integration
	err := r.client.DeleteIntegration(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Snowflake integration, got error: %s", err))
		return
	}
}

// ImportState imports the resource from the API.
func (r *SnowflakeIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.BaseIntegrationResource.ImportState(ctx, req, resp)
}
