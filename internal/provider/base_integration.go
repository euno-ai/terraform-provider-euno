package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BaseIntegrationResourceModel contains the common fields for all integration resources
type BaseIntegrationResourceModel struct {
	ID                        types.Int64  `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Active                    types.Bool   `tfsdk:"active"`
	Schedule                  *ScheduleModel `tfsdk:"schedule"`
	TriggerSecret             types.String `tfsdk:"trigger_secret"`
	TriggerURL                types.String `tfsdk:"trigger_url"`
	InvalidationStrategy      *InvalidationStrategyModel `tfsdk:"invalidation_strategy"`
	PendingCredentialsLookupKey types.String `tfsdk:"pending_credentials_lookup_key"`
	LastUpdatedAt             types.String `tfsdk:"last_updated_at"`
	CreatedAt                 types.String `tfsdk:"created_at"`
}

// ScheduleModel describes the schedule configuration
type ScheduleModel struct {
	TimeZone     types.String `tfsdk:"time_zone"`
	RepeatOn     types.List   `tfsdk:"repeat_on"`
	RepeatTime   types.String `tfsdk:"repeat_time"`
	RepeatPeriod types.Int64  `tfsdk:"repeat_period"`
}


// InvalidationStrategyModel describes the invalidation strategy configuration
type InvalidationStrategyModel struct {
	RevisionID types.Int64 `tfsdk:"revision_id"`
	TTLDays    types.Int64 `tfsdk:"ttl_days"`
}

// BaseIntegrationResource provides common functionality for all integration resources
type BaseIntegrationResource struct {
	client *EunoClient
}

// Configure adds the provider configured client to the resource.
func (r *BaseIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*EunoClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *EunoClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}


// getCommonBlocks returns the common blocks for pull integration resources
func getCommonBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"schedule": schema.SingleNestedBlock{
			MarkdownDescription: "The schedule configuration for the integration",
			Attributes: map[string]schema.Attribute{
				"time_zone": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The time zone for the schedule",
				},
				"repeat_on": schema.ListAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					MarkdownDescription: "The days of the week to repeat on",
				},
				"repeat_time": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The time to repeat at (HH:MM:SS format)",
				},
				"repeat_period": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The period in hours to repeat",
				},
			},
		},
		"invalidation_strategy": schema.SingleNestedBlock{
			MarkdownDescription: "The invalidation strategy configuration",
			Attributes: map[string]schema.Attribute{
				"revision_id": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The revision ID for invalidation",
				},
				"ttl_days": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "The TTL in days for invalidation",
				},
			},
		},
	}
}

// getCommonBlocksForPull returns the common blocks for pull integration resources
func getCommonBlocksForPull() map[string]schema.Block {
	return getCommonBlocks()
}

// getCommonBlocksForPush returns the common blocks for push integration resources (no schedule)
func getCommonBlocksForPush() map[string]schema.Block {
	return map[string]schema.Block{
		"invalidation_strategy": schema.SingleNestedBlock{
			MarkdownDescription: "The invalidation strategy configuration",
			Attributes: map[string]schema.Attribute{
				"revision_id": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The revision ID for invalidation",
				},
				"ttl_days": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "The TTL in days for invalidation",
				},
			},
		},
	}
}

// convertScheduleToAPI converts Terraform schedule to API format
func convertScheduleToAPI(schedule *ScheduleModel) *IntegrationSchedule {
	if schedule == nil {
		return nil
	}

	apiSchedule := &IntegrationSchedule{
		TimeZone: schedule.TimeZone.ValueString(),
	}

	if !schedule.RepeatOn.IsNull() {
		repeatOn := make([]string, 0, len(schedule.RepeatOn.Elements()))
		for _, elem := range schedule.RepeatOn.Elements() {
			if strElem, ok := elem.(types.String); ok {
				repeatOn = append(repeatOn, strElem.ValueString())
			}
		}
		apiSchedule.RepeatOn = repeatOn
	}

	if !schedule.RepeatTime.IsNull() {
		apiSchedule.RepeatTime = schedule.RepeatTime.ValueString()
	}

	if !schedule.RepeatPeriod.IsNull() {
		period := int(schedule.RepeatPeriod.ValueInt64())
		apiSchedule.RepeatPeriod = &period
	}

	return apiSchedule
}

// convertScheduleFromAPI converts API schedule to Terraform format
func convertScheduleFromAPI(apiSchedule *IntegrationSchedule) *ScheduleModel {
	if apiSchedule == nil {
		return nil
	}

	schedule := &ScheduleModel{
		TimeZone: types.StringValue(apiSchedule.TimeZone),
	}

	if apiSchedule.RepeatOn != nil {
		repeatOn := make([]attr.Value, len(apiSchedule.RepeatOn))
		for i, day := range apiSchedule.RepeatOn {
			repeatOn[i] = types.StringValue(day)
		}
		schedule.RepeatOn = types.ListValueMust(types.StringType, repeatOn)
	}

	if apiSchedule.RepeatTime != "" {
		schedule.RepeatTime = types.StringValue(apiSchedule.RepeatTime)
	}

	if apiSchedule.RepeatPeriod != nil && *apiSchedule.RepeatPeriod != 0 {
		schedule.RepeatPeriod = types.Int64Value(int64(*apiSchedule.RepeatPeriod))
	}

	return schedule
}

// convertInvalidationStrategyToAPI converts Terraform invalidation strategy to API format
func convertInvalidationStrategyToAPI(strategy *InvalidationStrategyModel) *InvalidationStrategy {
	if strategy == nil {
		return nil
	}

	apiStrategy := &InvalidationStrategy{
		TTLDays: int(strategy.TTLDays.ValueInt64()),
	}

	if !strategy.RevisionID.IsNull() {
		revisionID := int(strategy.RevisionID.ValueInt64())
		apiStrategy.RevisionID = &revisionID
	}

	return apiStrategy
}

// convertInvalidationStrategyFromAPI converts API invalidation strategy to Terraform format
func convertInvalidationStrategyFromAPI(apiStrategy *InvalidationStrategy) *InvalidationStrategyModel {
	if apiStrategy == nil {
		return nil
	}

	strategy := &InvalidationStrategyModel{
		TTLDays: types.Int64Value(int64(apiStrategy.TTLDays)),
	}

	if apiStrategy.RevisionID != nil {
		strategy.RevisionID = types.Int64Value(int64(*apiStrategy.RevisionID))
	}

	return strategy
}


// ImportState imports the resource from the API.
func (r *BaseIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is just the integration ID since account_id is in the provider
	integrationID := req.ID
	if integrationID == "" {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be the integration ID")
		return
	}

	// Convert string ID to Int64
	id, err := strconv.ParseInt(integrationID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to parse integration ID %q: %s", integrationID, err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(id))...)
}

// getCommonAttributes returns the common attributes for all integration resources
func getCommonAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "The ID of the integration",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name of the integration",
		},
		"active": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Whether the integration is active",
		},
		"trigger_secret": schema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			MarkdownDescription: "The secret key for triggering the integration",
		},
		"trigger_url": schema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			MarkdownDescription: "The URL for triggering the integration",
		},
		"pending_credentials_lookup_key": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The pending credentials lookup key",
		},
		"last_updated_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The last updated timestamp",
		},
		"created_at": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The creation timestamp",
		},
	}
}

