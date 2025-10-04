terraform {
  required_providers {
    euno = {
      source  = "euno-ai/euno"
      version = "0.0.1"
    }
  }
}

provider "euno" {
  account_id = 123
  server_url = "https://api.euno.ai"
  api_key    = "test-api-key"
}

# Test Fivetran integration
resource "euno_fivetran_integration" "test" {
  name   = "ci-test-fivetran"
  active = true

  configuration {
    api_key    = "placeholder-key"
    api_secret = "placeholder-secret"
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