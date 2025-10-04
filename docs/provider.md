# Euno Provider

The Euno provider is used to interact with Euno's data integration platform.

## Configuration

### Provider Arguments

The provider accepts the following configuration arguments:

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `account_id` | Your Euno account ID. | `number` | n/a | *yes* |
| `endpoint` | The Euno API endpoint. URL including scheme. | `string` | `"https://api.euno.ai"` | no |
| `timeout` | API request timeout in seconds. | `number` | `30` | no |
| `retry_max` | Maximum number of retries for failed requests. | `number` | `3` | no |

### Example Configuration

```hcl
provider "euno" {
  account_id = 123
  endpoint   = "https://api.euno.ai"
  timeout    = 60
  retry_max  = 5
}
```

### Environment Variables

You can also configure the provider using environment variables:

| Variable | Description | Equivalent |
|----------|-------------|------------|
| `EUNO_ACCOUNT_ID` | Your Euno account ID | `account_id` |
| `EUNO_ENDPOINT` | The Euno API endpoint | `endpoint` |
| `EUNO_TIMEOUT` | API request timeout | `timeout` |
| `EUNO_RETRY_MAX` | Maximum number of retries | `retry_max` |

#### Example with Environment Variables

```bash
export EUNO_ACCOUNT_ID=123
export EUNO_ENDPOINT=https://api.euno.ai
```

```hcl
provider "euno" {}
```

## Provider Features

- **Automated API Key Management**: Each integration handles its own authentication credentials
- **Push vs Pull Integrations**: Supports both scheduled pull integrations and webhook push integrations
- **Flexible Scheduling**: Configure recurring data synchronization with cron expressions
- **Data Validation Strategy**: Control how data is validated and invalidated
- **Rate Limiting**: Built-in rate limiting to respect API quotas

## Integration Types

### Pull Integrations

Pull integrations are scheduled to automatically fetch data from external sources:

- **Fivetran**: Automated data pipeline synchronization
- **Snowflake**: Query-based data loading and transformation
- **Hex**: Notebook-based data processing integration

Features:
- Configurable schedules (hourly, daily, weekly, custom cron)
- Retry mechanisms for failed runs
- Data validation and invalidation strategies

### Push Integrations

Push integrations receive data via webhooks triggered by external systems:

- **DBT Core**: Integration with DBT transformation projects

Features:
- Webhook URL and secret key generation
- Real-time data processing triggers
- Secure webhook authentication

## Rate Limiting

The provider includes built-in rate limiting to ensure compliance with API quotas:

- Maximum 3 concurrent API requests
- Automatic retry with exponential backoff
- Request queuing for burst protection

## Error Handling

The provider handles various error scenarios:

- **Authentication Errors**: Clear guidance on credential configuration
- **Rate Limit Exceeded**: Automatic retry with backoff
- **Validation Errors**: Detailed field-level error messages
- **Network Issues**: Retry logic for transient failures

## Data Sources

Currently, the provider does not include data sources. All resources are managed resources that create Euno integrations.