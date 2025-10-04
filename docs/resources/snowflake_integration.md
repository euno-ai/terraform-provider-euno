# euno_snowflake_integration

Manages a Snowflake integration in Euno. This resource creates a scheduled data synchronization integration that executes SQL queries against your Snowflake database.

~> **Note:** This is a pull-type integration that runs on a schedule. Use [`euno_dbt_core_integration`](dbt_core_integration.md) for push-type integrations.

## Example Usage

```hcl
resource "euno_snowflake_integration" "main" {
  name   = "snowflake-hourly-analytics"
  active = true

  schedule {
    frequency       = "daily"
    cron_expression = "0 6 * * * *"  # 6 AM daily
  }

  invalidation_strategy {
    revision_id = null
    ttl_days    = 30
  }

  configuration {
    sql   = "SELECT user_id, timestamp, event_type FROM analytics.events WHERE hour(timestamp) >= '${hour}' - 24::INT"
    credential_type = "key_pair"
    role  = "ANALYST_ROLE"
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
| `configuration` | Snowflake-specific configuration. | `object` | n/a | *yes* |

#### Schedule Block

The `schedule` block supports the following:

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `frequency` | Frequency of execution. Must be one of: `hourly`, `daily`, `weekly`, `custom`. | `string` | n/a | *yes* |
| `cron_expression` | Custom cron expression (required when `frequency` is `custom`). | `string` | `null` | no |
| `repeat_on` | Days to repeat (required when `frequency` is `weekly`). Days: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, .. | `list(string)` | `null` | no |

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
| `sql` | The SQL query to execute. Must be a valid Snowflake SQL query. | `string` | n/a | *yes* |
| `credential_type` | Type of authentication credential. Must be one of: `password`, `key_pair`. | `string` | n/a | *yes* |
| `role` | Snowflake role to use for the connection. | `string` | n/a | *yes* |
| `username` | Snowflake username (required for `password` credential type). | `string` | `""` | no |
| `password` | Snowflake password (required for `password` credential type). | `string` | `""` | no |
| `private_key_path` | Path to private key file (required for `key_pair` credential type). | `string` | `""` | no |
| `private_key_passphrase` | Private key passphrase (required for `key_pair` credential type). | `string` | `""` | no |

### Computed Attributes (Read-Only)

| Name | Description | Type |
|------|-------------|------|
| `id` | The unique ID of the integration as assigned by Euno. | `number` |
| `last_updated_at` | Timestamp of the last update to this integration. | `string` |
| `created_at` | Timestamp when this integration was created. | `string` |

## Import

Snowflake integrations can be imported using the integration ID:

```bash
terraform import euno_snowflake_integration.main 123
```

Where `123` is the integration ID returned by the Euno API.

## Authentication Methods

### Password Authentication

Use username and password for authentication:

```hcl
configuration {
  sql            = "SELECT * FROM my_table"
  credential_type = "password"
  role           = "ANALYST_ROLE"
  username       = "myuser"
  password       = "mypassword"
}
```

### Key Pair Authentication

Use private key authentication (recommended for production):

```hcl
configuration {
  sql                    = "SELECT * FROM my_table"
  credential_type       = "key_pair"
  role                  = "ANALYST_ROLE"
  private_key_path      = "/path/to/private_key.pem"
  private_key_passphrase = "key_password"
}
```

### Using Existing Credentials

If credentials are already stored in Euno, you can omit the credential fields:

```hcl
configuration {
  sql   = "SELECT * FROM my_table"
  credential_type = "key_pair"
  role  = "ANALYST_ROLE"
}
```

## Scheduling Examples

### Daily Analytics Query

```hcl
schedule {
  frequency       = "daily"
  cron_expression = "0 6 * * * *"  # 6 AM daily
}
```

### Hourly Data Refresh

```hcl
schedule {
  frequency       = "hourly"
  cron_expression = "0 0 * * * *"
  repeat_on      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
}
```

### Custom Schedule

```hcl
schedule {
  frequency       = "custom"
  cron_expression = "0 */4 * * *"  # Every 4 hours
}
```

## Security Best Practices

1. **Use Key Pair Authentication**: Prefer key pair authentication over password authentication for better security.

2. **Principle of Least Privilege**: Assign the minimum necessary Snowflake role permissions for your integration's SQL queries.

3. **Secure Credential Storage**: Store credentials in Terraform variables or a secure secret management system:
   ```hcl
   variable "snowflake_private_key" {
     description = "Snowflake private key"
     type        = string
     sensitive   = true
   }
   ```

4. **SQL Injection Prevention**: Use parameterized queries and validate input data.

## Troubleshooting

### Common Issues

**Authentication Failures**
- Verify Snowflake credentials are correct
- Check that the role has necessary permissions
- Ensure the account identifier and warehouse are accessible

**SQL Syntax Errors**
- Validate SQL syntax against Snowflake documentation
- Test queries in Snowflake console before using in integration
- Ensure all referenced objects exist and are accessible

**Permission Errors**
- Verify the Snowflake role has SELECT permissions on referenced tables/views
- Check warehouse usage permissions
- Ensure database and schema access permissions

**Network Connectivity**
- Verify Euno can reach your Snowflake instance
- Check firewall rules and network policies
- Ensure proper Snowflake network policies are configured