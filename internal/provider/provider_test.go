package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-framework/testing"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"euno": providerserver.Serve(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the environment, or API keys, etc.
}

func TestAccEunoProvider(t *testing.T) {
	testing.Test(t, testing.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []testing.TestStep{
			// Create and Read testing
			{
				Config: testAccFivetranIntegrationConfig(),
				Check: testing.ComposeAggregateTestCheckFunc(
					testing.TestCheckResourceAttr("euno_fivetran_integration.test", "name", "test-fivetran"),
					testing.TestCheckResourceAttr("euno_fivetran_integration.test", "active", "true"),
				),
			},
		},
	})
}

func testAccFivetranIntegrationConfig() string {
	return `
provider "euno" {
  account_id = 123
  endpoint   = "http://localhost:8000"
}

resource "euno_fivetran_integration" "test" {
  name   = "test-fivetran"
  active = true
  
  schedule {
    frequency       = "daily"
    cron_expression = "0 6 * * * *"
  }
  
  invalidation_strategy {
    revision_id = null
    ttl_days    = 7
  }
  
  configuration {
    auto_sync_enabled = true
    connector         = "bigquery"
    destination_schema_prefix = "test_"
    connector_id     = "test-connector-id"
    api_key          = "test-key"
    api_secret       = "test-secret"
    transform        = false
    day_of_the_week  = 1
    hour_of_the_day  = 6
    version_id       = ""
  }
}
`
}
