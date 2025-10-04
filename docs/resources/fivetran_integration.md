# euno_fivetran_integration

Manages a Fivetran integration in Euno. This resource creates a scheduled data synchronization integration that pulls data from Fivetran connectors.

~> **Note:** This is a pull-type integration that runs on a schedule. Use [`euno_dbt_core_integration`](dbt_core_integration.md) for push-type integrations.

## Example Usage

```hcl
resource "euno_fivetran_integration" "main" {
  name   = "fivetran-hourly-sync"
  active = true

  schedule {
    frequency       = "hourly"
    cron_expression = "0 0 * * * *"
    repeat_on       = ["monday", "tuesday", "wednesday", "thursday", "friday"]
  }

  invalidation_strategy {
    revision_id = null
    ttl_days    = 7
  }

  configuration {
    auto_sync_enabled = true
    connector         = "bigquery"  # or "snowflake", "postgres", etc.
    destination_schema_prefix = "fivetran_"
    connector_id     = "sync_your_connector_id"
    api_key          = "your-fivetran-api-key"
    api_secret       = "your-fivetran-api-secret"  # optional if already set
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
| `schedule` | Configuration for scheduled execution. | `object` | n/a | *yes* |
| `invalidation_strategy` | Configuration for data validation and invalidation. | `object` | n/a | *yes* |
| `configuration` | Fivetran-specific configuration. | `object` | n/a | *yes* |

#### Schedule Block

The `schedule` block supports the following:

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `frequency` | Frequency of execution. Must be one of: `hourly`, `daily`, `weekly`, `custom`. | `string` | n/a | *yes* |
| `cron_expression` | Custom cron expression (required when `frequency` is `custom`). | `string` | `null` | no |
| `repeat_on` | Days to repeat (required when `frequency` is `weekly`). Days: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. | `list(string)` | `null` | no |

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
| `auto_sync_enabled` | Whether automatic synchronization is enabled. | `bool` | n/a | *yes* |
| `connector` | The destination type. Must be one of: `bigquery`, `snowflake`, `postgres`. | `string` | n/a | *yes* |
| `destination_schema_prefix` | Prefix for destination schema names. | `string` | n/a | *yes* |
| `connector_id` | The ID of the Fivetran connector to sync. | `string` | n/a | *yes* |
| `api_key` | Your Fivetran API key. | `string` | n/a | *yes* |
| `api_secret` | Your Fivetran API secret. Leave empty to use existing stored secret. | `string` | `""` | no |
| `transform` | Whether Fivetran transformations should be used. | `bool` | `false` | no |
| `day_of_the_week` | Day of the week for the sync schedule. | `number` | n/a | *yes* |
| `hour_of_the_day` | Hour of the day for the sync schedule (0-23). | `number` | n/a | *yes* |
| `version_id` | Specific version ID for the connector. | `string` | `""` | no |

### Computed Attributes (Read-Only)

| Name | Description | Type |
|------|-------------|------|
| `id` | The unique ID of the integration as assigned by Euno. | `number` |
| `last_updated_at` | Timestamp of the last update to this integration. | `string` |
| `created_at` | Timestamp when this integration was created. | `string` |

## Import

Fivetran integrations can be imported using the integration ID:

```bash
terraform import euno_fivetran_integration.main 123
```

Where `123` is the integration ID returned by the Euno API.

## Common Scheduling Patterns

### Hourly Synchronization (Business Hours)

```hcl
schedule {
  frequency       = "hourly"
  cron_expression = "0 0 * * * *"
  repeat_on      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
}
```

### Daily Synchronization

```hcl
schedule {
  frequency       = "daily"
  cron_expression = "0 2 * * * *"  # 2 AM daily
}
```

### Weekly Synchronization

```hcl
schedule {
  frequency = "weekly"
  repeat_on = ["sunday"]
}
```

### Custom Cron Expression

```hcl
schedule {
  frequency       = "custom"
  cron_expression = "0 15 * * 1-5"  # 3 PM weekdays only
}
```

## Best Practices

1. **Unique Names**: Ensure integration names are unique within your account to avoid conflicts.

2. **Authentication**: Store API keys and secrets securely. Consider using Terraform variables or a secrets management system.

3. **Scheduling**: Balance data freshness with resource usage. Hourly syncs ensure up-to-date data but consume more resources.

4. **Validation Strategy**: Set appropriate TTL values based on your data freshness requirements. Longer TTLs reduce validation overhead.

## Troubleshooting

### Common Issues

**Authentication Errors**
- Verify your Fivetran API credentials are correct
- Ensure the API key has necessary permissions for the connectors you want to sync

**Scheduling Issues**
- Check that `repeat_on` days are valid when using `weekly` frequency
- Ensure cron expressions are valid when using `custom` frequency

**Configuration Errors**
- Verify `connector_id` matches an existing Fivetran connector in your account
- Ensure `connector` type matches your target destination