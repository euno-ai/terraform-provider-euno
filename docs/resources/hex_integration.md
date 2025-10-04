# euno_hex_integration

Manages a Hex integration in Euno. This resource creates a scheduled data synchronization integration that crawls Hex notebooks and processes their output.

~> **Note:** This is a pull-type integration that runs on a schedule. Use [`euno_dbt_core_integration`](dbt_core_integration.md) for push-type integrations.

## Example Usage

```hcl
resource "euno_hex_integration" "main" {
  name   = "hex-daily-analytics"
  active = true

  schedule {
    frequency       = "daily"
    cron_expression = "0 8 * * * *"  # 8 AM daily
  }

  invalidation_strategy {
    revision_id = null
    ttl_days    = 7
  }

  configuration {
    api_token                = "your-hex-api-token"
    project_id              = "your-hex-project-id"
    namespace_id            = "your-hex-namespace-id"
    exclude_deleted_notebooks = true
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
| `configuration` | Hex-specific configuration. | `object` | n/a | *yes* |

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
| `revision_id` | Specific revision ID to validate against. | `string` | Required | *yes* |
| `ttl_days` | Number of days after which data is considered stale and gets invalidated. | `number` | n/a | *yes* |

#### Configuration Block

The `configuration` block supports the following:

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `api_token` | Your Hex API token for authentication. | `string` | n/a | *yes* |
| `project_id` | The Hex project ID to crawl for notebooks. | `string` | n/a | *yes* |
| `namespace_id` | The Hex namespace ID where the project is located. | `string` | n/a | *yes* |
| `exclude_deleted_notebooks` | Whether to exclude deleted notebooks from crawling. | `bool` | `true` | no |

### Computed Attributes (Read-Only)

| Name | Description | Type |
|------|-------------|------|
| `id` | The unique ID of the integration as assigned by Euno. | `number` |
| `last_updated_at` | Timestamp of the last update to this integration. | `string` |
| `created_at` | Timestamp when this integration was created. | `string` |

## Import

Hex integrations can be imported using the integration ID:

```bash
terraform import euno_hex_integration.main 123
```

Where `123` is the integration ID returned by the Euno API.

## Getting Hex Credentials

### API Token

1. Log in to your Hex account
2. Navigate to Account Settings → API Tokens
3. Create a new API token or use an existing one
4. Ensure the token has appropriate permissions for your project

### Project and Namespace IDs

You can find these values in your Hex account:

1. **Project ID**: Found in the project URL: `https://app.hex.tech/projects/YOUR_PROJECT_ID`
2. **Namespace ID**: Available in project settings or the project URL structure

## Scheduling Examples

### Daily Notebook Processing

```hcl
schedule {
  frequency       = "daily"
  cron_expression = "0 8 * * * *"  # 8 AM daily
}
```

### Hourly Monitoring

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
  cron_expression = "0 0 9 * * MON-FRI"  # 9 AM weekdays
}
```

## Configuration Examples

### Basic Configuration

```hcl
configuration {
  api_token                = var.hex_api_token
  project_id              = "proj_abc123"
  namespace_id            = "ns_def456"
  exclude_deleted_notebooks = true
}
```

### Using Variables for Security

```hcl
variable "hex_api_token" {
  description = "Hex API token"
  type        = string
  sensitive   = true
}

resource "euno_hex_integration" "main" {
  name   = "hex-analytics"
  active = true

  schedule {
    frequency = "daily"
  }

  invalidation_strategy {
    revision_id = "daily-run"
    ttl_days    = 7
  }

  configuration {
    api_token                = var.hex_api_token
    project_id              = "proj_abc123"
    namespace_id            = "ns_def456"
    exclude_deleted_notebooks = true
  }
}
```

## Validation Strategy

Hex integrations require a specific `revision_id` for validation:

```hcl
invalidation_strategy {
  revision_id = "notebook-revision-v1.2"  # Specific Hex notebook revision
  ttl_days    = 7                         # Refresh data weekly
}
```

The `revision_id` should correspond to a specific Hex notebook revision that you want to validate against. This ensures consistency in data processing.

## Security Best Practices

1. **API Token Management**: Store API tokens securely using Terraform variables with `sensitive = true`

2. **Access Control**: Ensure your Hex API token has the minimum necessary permissions
   - Read access to the specific project
   - Permission to access notebook metadata and outputs

3. **Validation Strategy**: Use meaningful revision IDs that correspond to your notebook versions for reliable validation

## Troubleshooting

### Common Issues

**Authentication Errors**
- Verify your Hex API token is valid and has appropriate permissions
- Check that the token hasn't expired
- Ensure the token can access the specified project and namespace

**Project Access Issues**
- Verify the `project_id` and `namespace_id` are correct
- Check that the API token has access to the specified project
- Ensure the project exists and is accessible

**Notebook Processing Errors**
- Check Hex notebook execution logs for any runtime errors
- Verify notebook dependencies are properly configured
- Ensure notebooks don't have infinite loops or infinite loops or resource-intensive operations

**Validation Failures**
- Verify the `revision_id` corresponds to a valid notebook revision
- Check that the notebook revision produces consistent output
- Ensure TTL settings align with your data freshness requirements