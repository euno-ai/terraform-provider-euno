package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"euno": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the environment, or API keys, etc.
}

func TestAccEunoProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFivetranIntegrationConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("euno_fivetran_integration.test", "name", "test-fivetran"),
					resource.TestCheckResourceAttr("euno_fivetran_integration.test", "active", "true"),
				),
			},
		},
	})
}

func testAccFivetranIntegrationConfig() string {
	return `
provider "euno" {
  account_id = 123
  server_url = "http://localhost:8000"
  api_key    = "test-api-key"
}

resource "euno_fivetran_integration" "test" {
  name   = "test-fivetran"
  active = true
  
  configuration {
    api_key    = "test-key"
    api_secret = "test-secret"
    base_url   = "https://api.fivetran.com/v1"
  }
  
  schedule {
    time_zone   = "America/Los_Angeles"
    repeat_on   = ["Mon", "Wed", "Fri"]
    repeat_time = "10:00:00"
  }
  
  invalidation_strategy {
    ttl_days = 7
  }
}
`
}
