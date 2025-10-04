package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"euno": providerserver.NewProtocol6Server(New("test")()),
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
				Config: testAccFivetranSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("euno_fivetran_source.test", "name", "test-fivetran"),
					resource.TestCheckResourceAttr("euno_fivetran_source.test", "active", "true"),
				),
			},
		},
	})
}

func testAccFivetranSourceConfig() string {
	return `
provider "euno" {
  server_url = "http://localhost:8000"
  api_key    = "test-api-key"
}

resource "euno_fivetran_source" "test" {
  account_id = 1
  name       = "test-fivetran"
  active     = true
  
  configuration = {
    api_key    = "test-key"
    api_secret = "test-secret"
  }
  
  schedule {
    time_zone   = "America/Los_Angeles"
    repeat_on   = ["Mon", "Wed", "Fri"]
    repeat_time = "10:00:00"
  }
}
`
}
