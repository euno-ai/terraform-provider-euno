# euno_dbt_core_integration

Manages a DBT Core integration in Euno. This resource creates a push-type integration that receives data via webhooks triggered by DBT project runs.

~> **Note:** This is a push-type integration that provides webhook endpoints for triggering. Use [`euno_fivetran_integration`](fivetran_integration.md), [`euno_snowflake_integration`](snowflake_integration.md), or [`euno_hex_integration`](hex_integration.md) for pull-type integrations.

## Example Usage

```hcl
resource "euno_dbt_core_integration" "main" {
  name   = "analytics-dbt-project"
  active = true

  invalidation_strategy {
    revision_id = null
    ttl_days    = 30
  }

  configuration {
    build_target                              = "prod"
    repository_url                           = "https://github.com/mycompany/dbt-project"
    repository_branch                        = "main"
    dbt_project_root_directory_in_repository = "/"
    allow_resources_with_no_catalog_entry     = false
    override_uri_prefix                      = "dbt.myproject"
    schemas_aliases = {
      "analytics.base" = "prod.base"
      "analytics.marts" = "prod.marts"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

### Top-Level Arguments

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `name` | The name of the integration. Must be unique within your account. | `string` | n/a | *yes* |
| `active` | Whether the integration is active. | `bool` | `true` | no |
| `invalidation_strategy` | Configuration for data validation and invalidation. | `object` | n/a | *yes* |
| `configuration` | DBT Core-specific configuration. | `object` | n/a | *yes* |

### Computed Attributes (Read-Only)

| Name | Description | Type |
|------|-------------|------|
| `id` | The unique ID of the integration as assigned by Euno. | `number` |
| `secret_key` | A secret key for authenticating webhook requests to this integration. | `string` |
| `webhook_url` | The webhook URL where DBT Core sends data. | `string` |
| `last_updated_at` | Timestamp of the last update to this integration. | `string` |
| `created_at` | Timestamp when this integration was created. | `string` |

~> **Important:** `secret_key` and `webhook_url` are sensitive computed attributes that contain authentication credentials and the webhook endpoint URL.

#### Invalidation Strategy Block

The `invalidation_strategy` block supports the following:

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `revision_id` | Specific revision ID to validate against. Set to `null` to validate against the latest revision. | `string` | `null` | no |
| `ttl_days` | Number of days after which data is considered stale and gets invalidated. | `number` | n/a | *yes* |

#### Configuration Block

The `configuration` block supports the following:

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `build_target` | The dbt target to build. Must correspond to a target defined in your `profiles.yml`. | `string` | n/a | *yes* |
| `repository_url` | The URL of the git repository where the dbt project is stored. | `string` | `""` | no |
| `repository_branch` | The branch of the git repository where the dbt project is stored. | `string` | `""` | no |
| `dbt_project_root_directory_in_repository` | The subdirectory within the git repository where the dbt project is stored. | `string` | `"/"` | no |
| `repository_revision` | The revision of the git repository where the dbt project is stored. | `string` | `""` | no |
| `schemas_aliases` | A map of schema aliases for database.schema combinations. | `map(string)` | `{}` | no |
| `allow_resources_with_no_catalog_entry` | Whether to allow dbt resources with no corresponding catalog entry to be ingested. | `bool` | `false` | no |
| `override_uri_prefix` | The prefix to override the URI of the resources. If not set, uses "dbt.<dbt_project_name>". | `string` | `""` | no |
| `stage_build_target` | The stage dbt target to build (for pre-production validation). | `string` | `""` | no |

## Import

DBT Core integrations can be imported using the integration ID:

```bash
terraform import euno_dbt_core_integration.main 123
```

Where `123` is the integration ID returned by the Euno API.

## Configuration Examples

### Basic Configuration

```hcl
configuration {
  build_target = "prod"
}
```

### Repository Configuration

```hcl
configuration {
  build_target                              = "prod"
  repository_url                           = "https://github.com/mycompany/dbt-project"
  repository_branch                        = "main"
  dbt_project_root_directory_in_repository = "/"
  repository_revision                      = "main"
}
```

### Schema Aliases

Map schemas from your dbt project to different database/schema combinations:

```hcl
configuration {
  build_target = "prod"
  schemas_aliases = {
    "analytics.base"    = "prod.base""
    "analytics.marts"   = "prod.marts"
    "staging.raw_data"  = "prod.staging"
  }
}
```

### Complete Configuration

```hcl
configuration {
  build_target                              = "prod"
  stage_build_target                       = "staging"
  repository_url                           = "https://github.com/mycompany/dbt-project"
  repository_branch                        = "main"
  dbt_project_root_directory_in_repository = "/"
  repository_revision                      = "main"
  allow_resources_with_no_catalog_entry     = false
  override_uri_prefix                      = "dbt.myproject"
  schemas_aliases = {
    "analytics.base" = "prod.base"
    "analytics.marts" = "prod.marts"
  }
}
```

## Webhook Integration

Once created, the DBT Core integration provides:

### Secret Key (`secret_key`)
A secure token for authenticating webhook requests. Include this in your DBT Core hooks configuration.

### Webhook URL (`webhook_url`)
The endpoint URL where DBT Core should send data after successful runs.

### Using Webhook in DBT Core

Configure your `dbt_project.yml` to trigger the webhook:

```yaml
on-run-end: "{{ post_hook('curl -X POST {{ env(\"EUNO_WEBHOOK_URL\") }} -H \"Authorization: Bearer {{ env(\"EUNO_SECRET_KEY\") }}\"') }}"
```

Or use dbt Cloud's webhook configuration with the provided URL and secret key.

## Security Best Practices

1. **Secret Key Management**: Store the `secret_key` securely and never commit it to version control:
   ```hcl
   output "dbt_secret_key" {
     value = euno_dbt_core_integration.main.secret_key
     sensitive = true
   }
   ```

2. **Schema Aliases**: Use meaningful schema aliases to organize your data consistently across environments.

3. **Access Control**: Ensure your DBT Core deployment can reach the webhook URL and has the secret key.

4. **Validation Strategy**: Set appropriate TTL values based on how frequently your DBT models refresh data.

## Troubleshooting

### Common Issues

**Webhook Authentication Failures**
- Verify the `secret_key` is correctly included in webhook requests
- Check that the `Authorization` header format is correct
- Ensure the webhook request contains valid JSON data

**DBT Build Target Issues**
- Verify the `build_target` corresponds to an existing target in your `profiles.yml`
- Ensure the target's database connection is properly configured
- Check that DBT can connect to the specified database/schema

**Schema Alias Conflicts**
- Verify schema aliases don't conflict with existing Euno catalog entries
- Check that the target database/schema combinations exist
- Ensure aliases follow the "database.schema" pattern

**Repository Access Issues**
- Verify repository URLs are accessible
- Check that branch and revision references exist
- Ensure proper permissions for accessing the repository