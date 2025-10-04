package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure EunoProvider satisfies various provider interfaces.
var _ provider.Provider = &EunoProvider{}

// EunoProvider defines the provider implementation.
type EunoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "" when the provider is built
	// and ran locally.
	version string
}

// EunoProviderModel describes the provider data model.
type EunoProviderModel struct {
	ServerURL types.String `tfsdk:"server_url"`
	APIKey    types.String `tfsdk:"api_key"`
	AccountID types.Int64  `tfsdk:"account_id"`
}

// Metadata returns the provider type name.
func (p *EunoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "euno"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *EunoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the Euno API server",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for authenticating with the Euno API",
				Required:            true,
				Sensitive:           true,
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "The account ID for Euno integrations",
				Required:            true,
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *EunoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config EunoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := NewEunoClient(config.ServerURL.ValueString(), config.APIKey.ValueString(), int(config.AccountID.ValueInt64()))
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines the resources implemented in the provider.
func (p *EunoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFivetranIntegrationResource,
		NewSnowflakeIntegrationResource,
		NewHexIntegrationResource,
		NewDbtCoreIntegrationResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *EunoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// NewExampleDataSource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &EunoProvider{
			version: version,
		}
	}
}
