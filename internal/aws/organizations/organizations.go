package organizations

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
)

func GetAccountsMapping(
	ctx context.Context,
	client *organizations.Client,
) (map[string]string, error) {
	accountsMapping := make(map[string]string)
	paginator := organizations.NewListAccountsPaginator(
		client,
		&organizations.ListAccountsInput{},
	)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving accounts: %w", err)
		}
		for _, account := range page.Accounts {
			accountsMapping[aws.ToString(account.Id)] = aws.ToString(account.Name)
		}
	}
	return accountsMapping, nil
}
