package health

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/health/types"
)

func setEventFilter(input *health.DescribeEventsForOrganizationInput,
	service, eventStatus, region, accountId string) {
	if service != "" || eventStatus != "" || region != "" || accountId != "" {
		input.Filter = &types.OrganizationEventFilter{}
		if service != "" {
			input.Filter.Services = []string{service}
		}
		if eventStatus != "" {
			input.Filter.EventStatusCodes = []types.EventStatusCode{
				types.EventStatusCode(eventStatus),
			}
		}
		if region != "" {
			input.Filter.Regions = []string{region}
		}
		if accountId != "" {
			input.Filter.AwsAccountIds = []string{accountId}
		}
	}
}

func DescribeEventsForOrganizationInput(
	service, eventStatus, region, accountId string) *health.DescribeEventsForOrganizationInput {
	input := &health.DescribeEventsForOrganizationInput{
		MaxResults: aws.Int32(100),
	}
	setEventFilter(input, service, eventStatus, region, accountId)
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
