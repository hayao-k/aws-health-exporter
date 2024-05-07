package health

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/health/types"
)

func setEventFilter(input *health.DescribeEventsForOrganizationInput,
	service, eventStatus, eventCategory, region, accountId,
	startTime, endTime, lastUpdatedTime string,
) {
	if service != "" || eventStatus != "" || eventCategory != "" ||
		region != "" || accountId != "" || startTime != "" ||
		endTime != "" || lastUpdatedTime != "" {
		input.Filter = &types.OrganizationEventFilter{}
		if service != "" {
			input.Filter.Services = []string{service}
		}
		if eventStatus != "" {
			input.Filter.EventStatusCodes = []types.EventStatusCode{
				types.EventStatusCode(eventStatus),
			}
		}
		if eventCategory != "" {
			input.Filter.EventTypeCategories = []types.EventTypeCategory{
				types.EventTypeCategory(eventCategory),
			}
		}
		if region != "" {
			input.Filter.Regions = []string{region}
		}
		if accountId != "" {
			input.Filter.AwsAccountIds = []string{accountId}
		}
		if startTime != "" {
			input.Filter.StartTime = parseTimeRange(startTime)
		}
		if endTime != "" {
			input.Filter.EndTime = parseTimeRange(endTime)
		}
		if lastUpdatedTime != "" {
			input.Filter.LastUpdatedTime = parseTimeRange(lastUpdatedTime)
		}
	}
}

func parseTimeRange(timeRangeStr string) *types.DateTimeRange {
	if timeRangeStr == "" {
		return nil
	}
	var from, to string
	trimmed := strings.Trim(timeRangeStr, "{}")
	parts := strings.Split(trimmed, ",")
	for _, part := range parts {
		if strings.HasPrefix(part, "from=") {
			from = strings.TrimPrefix(part, "from=")
		} else if strings.HasPrefix(part, "to=") {
			to = strings.TrimPrefix(part, "to=")
		}
	}
	timeRange := &types.DateTimeRange{}
	if from != "" {
		fromTime, err := time.Parse(time.RFC3339, from)
		if err == nil {
			timeRange.From = aws.Time(fromTime)
		}
	}
	if to != "" {
		toTime, err := time.Parse(time.RFC3339, to)
		if err == nil {
			timeRange.To = aws.Time(toTime)
		}
	}

	return timeRange
}

func DescribeEventsForOrganizationInput(
	service, eventStatus, eventCategory, region, accountId,
	startTime, endTime, lastUpdatedTime string,
) *health.DescribeEventsForOrganizationInput {
	input := &health.DescribeEventsForOrganizationInput{
		MaxResults: aws.Int32(100),
	}
	setEventFilter(input, service, eventStatus, eventCategory, region, accountId,
		startTime, endTime, lastUpdatedTime)
	return input
}

func DescribeEventsForOrganization(
	ctx context.Context,
	healthClient *health.Client,
	input *health.DescribeEventsForOrganizationInput,
) (*health.DescribeEventsForOrganizationOutput, error) {
	eventsResp, err := healthClient.DescribeEventsForOrganization(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving events: %w", err)
	}
	return eventsResp, nil
}

func GetAffectedAccounts(ctx context.Context, healthClient *health.Client,
	eventArn string) ([]string, error) {

	var accounts []string
	describePaginator := health.NewDescribeAffectedAccountsForOrganizationPaginator(
		healthClient, &health.DescribeAffectedAccountsForOrganizationInput{
			EventArn: aws.String(eventArn),
		})

	for describePaginator.HasMorePages() {
		page, err := describePaginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get affected accounts for event ARN %s: %w", eventArn, err)
		}

		accounts = append(accounts, page.AffectedAccounts...)
	}

	return accounts, nil
}

func GetAffectedEntities(ctx context.Context, healthClient *health.Client,
	eventArn, accountId string) ([]types.AffectedEntity, error) {

	var entities []types.AffectedEntity
	entityPaginator := health.NewDescribeAffectedEntitiesForOrganizationPaginator(
		healthClient, &health.DescribeAffectedEntitiesForOrganizationInput{
			OrganizationEntityFilters: []types.EventAccountFilter{
				{
					AwsAccountId: aws.String(accountId),
					EventArn:     aws.String(eventArn),
				},
			},
		})

	for entityPaginator.HasMorePages() {
		entityPage, err := entityPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entityPage.Entities...)
	}

	return entities, nil
}
