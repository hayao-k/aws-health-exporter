package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/health/types"
	awshealth "github.com/hayao-k/health-exporter/internal/aws/health"
)

func GenerateEventFileName(selectedEvent types.OrganizationEvent, specifiedAccountId string) string {
	eventTypeCode := strings.ReplaceAll(
		strings.ReplaceAll(aws.ToString(selectedEvent.EventTypeCode), " ", "_"), "/", "_")
	startTimeStr := "N-A"
	if selectedEvent.StartTime != nil {
		startTimeStr = selectedEvent.StartTime.Format("2006-01-02_15-04-05")
	}
	eventRegion := aws.ToString(selectedEvent.Region)
	if specifiedAccountId != "" {
		return fmt.Sprintf("%s_%s_%s_%s.csv", eventTypeCode, startTimeStr, eventRegion, specifiedAccountId)
	}
	return fmt.Sprintf("%s_%s_%s.csv", eventTypeCode, startTimeStr, eventRegion)
}

func createCsvFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSV file: %w", err)
	}
	return file, nil
}

func writeHeader(writer *csv.Writer, echoToStdout bool) error {
	header := []string{"Account ID", "Account Name", "Region", "Identifier",
		"Status", "Last Updated"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header to CSV: %w", err)
	}
	if echoToStdout {
		fmt.Println(strings.Join(header, ",")) // ヘッダーを標準出力に出力
	}
	return nil
}

func writeEntriesToCsv(writer *csv.Writer, account, accountName, region string,
	entities []types.AffectedEntity, echoToStdout bool) error {

	for _, entity := range entities {
		statusCode := "N/A"
		if entity.StatusCode != "" {
			statusCode = string(entity.StatusCode)
		}

		lastUpdatedTime := "N/A"
		if entity.LastUpdatedTime != nil {
			lastUpdatedTime = entity.LastUpdatedTime.Format("2006-01-02 15:04:05")
		}

		record := []string{account, accountName, region,
			aws.ToString(entity.EntityValue), statusCode, lastUpdatedTime}
		if err := writer.Write(record); err != nil {
			return err
		}

		if echoToStdout {
			fmt.Println(strings.Join(record, ",")) // レコードを標準出力に出力
		}
	}

	return nil
}

func WriteEventDetailsToCsv(ctx context.Context, healthClient *health.Client,
	eventArn string, accountsMapping map[string]string,
	selectedEvent types.OrganizationEvent, csvFileName string,
	echoToStdout bool, specifiedAccountId string) error {

	file, err := createCsvFile(csvFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writeHeader(writer, echoToStdout); err != nil {
		return err
	}

	var affectedAccountIds []string
	if specifiedAccountId != "" {
		affectedAccountIds = []string{specifiedAccountId} // Process only the specified accounts
	} else {
		affectedAccountIds, err = awshealth.GetAffectedAccounts(ctx, healthClient, eventArn)
		if err != nil {
			return err
		}
	}

	for _, accountId := range affectedAccountIds {
		accountName, exists := accountsMapping[accountId]
		if !exists {
			accountName = "Unknown"
		}

		entities, err := awshealth.GetAffectedEntities(ctx, healthClient, eventArn, accountId)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error retrieving entities for account %s: %v\n", accountId, err)
			continue
		}

		if err := writeEntriesToCsv(writer, accountId, accountName,
			aws.ToString(selectedEvent.Region), entities, echoToStdout); err != nil {

			fmt.Fprintf(os.Stderr,
				"Failed to write records for accountId %s: %v\n", accountId, err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("Failed to complete writing to CSV: %w", err)
	}

	return nil
}
