# Documentation

Welcome to the Terraform Provider for Euno documentation. This provider allows you to manage your Euno integrations declaratively using Terraform.

## Getting Started

- [Provider Configuration](provider.md) - Configure the Euno provider
- [Installation Guide](../README.md#installation) - Install and configure the provider

## Resources

### Pull Integrations
Pull integrations run on a schedule and fetch data from external sources:

- **[Fivetran Integration](resources/fivetran_integration.md)** - Automated data pipeline synchronization from Fivetran
- **[Snowflake Integration](resources/snowflake_integration.md)** - Execute SQL queries against Snowflake databases
- **[Hex Integration](resources/hex_integration.md)** - Crawl and process Hex notebook outputs

### Push Integrations  
Push integrations receive data via webhooks triggered by external systems:

- **[DBT Core Integration](resources/dbt_core_integration.md)** - Webhook integration for DBT project runs

## Examples

Working examples for each integration type are available in the [examples/](../examples/) directory:

- [Fivetran Example](../examples/fivetran/)
- [Snowflake Example](../examples/snowflake/)  
- [Hex Example](../examples/hex/)
- [DBT Core Example](../examples/dbt_core/)

## Integration Types Overview

### Pull Integrations
Pull integrations are scheduled to run automatically and fetch data from external sources. They support:
- Configurable scheduling (hourly, daily, weekly, custom cron)
- Retry mechanisms for failed runs
- Data validation and invalidation strategies
- Authentication credential management

### Push Integrations
Push integrations receive data via webhooks triggered by external systems. They provide:
- Webhook URL and secret key generation
- Real-time data processing triggers
- Secure webhook authentication
- Data validation strategies

## Support

- 📖 [Provider Documentation](https://registry.terraform.io/providers/euno-ai/euno/latest/docs)
- 🐛 [Report Issues](https://github.com/euno-ai/terraform-provider-euno/issues)
- 💬 [Discussions](https://github.com/euno-ai/terraform-provider-euno/discussions)