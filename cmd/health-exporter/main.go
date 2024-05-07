package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hayao-k/aws-health-exporter/internal/aws"
	"github.com/hayao-k/aws-health-exporter/internal/aws/health"
	"github.com/hayao-k/aws-health-exporter/internal/aws/organizations"
	"github.com/hayao-k/aws-health-exporter/internal/csv"
	"github.com/hayao-k/aws-health-exporter/internal/ui"
	"github.com/urfave/cli/v2"
)

var version = "v0.0.0"

func main() {
	app := &cli.App{
		Name:    "health-exporter",
		Usage:   "Describe AWS Health events for your organization",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "event-filter",
				Aliases: []string{"filter", "f"},
				Usage:   "Filter events by multiple criteria, e.g., --event-filter=\"service=RDS,status=open\"",
			},
			&cli.StringFlag{
				Name:    "status-code",
				Aliases: []string{"c"},
				Usage:   "Filter entity by status code. Possible values are IMPAIRED, UNIMPAIRED, UNKNOWN, PENDING and RESOLVED",
			},
			&cli.StringFlag{
				Name:    "account-id",
				Aliases: []string{"i"},
				Usage:   "Specify a single account ID to process",
			},
			&cli.BoolFlag{
				Name:    "echo",
				Aliases: []string{"e"},
				Usage:   "Echo CSV content to standard output",
			},
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "AWS profile name to use",
			},
			&cli.StringFlag{
				Name:    "output-file",
				Aliases: []string{"file-name", "o"},
				Usage:   "Specify the output CSV file name",
			},
		},
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			eventFilter := c.String("event-filter")
			filters := parseEventFilter(eventFilter)

			service := filters["service"]
			eventStatus := filters["status"]
			eventCategory := filters["category"]
			region := filters["region"]
			startTime := filters["startTime"]
			endTime := filters["endTime"]
			lastUpdatedTime := filters["lastUpdatedTime"]

			statusCode := c.String("status-code")
			echoToStdout := c.Bool("echo")
			profile := c.String("profile")
			specifiedAccountId := c.String("account-id")
			specifiedFileName := c.String("output-file")

			cfg, err := aws.LoadAWSConfig(ctx, profile)
			if err != nil {
				return err
			}

			healthClient, orgClient := aws.CreateServices(cfg)

			input := health.DescribeEventsForOrganizationInput(
				service,
				eventStatus,
				eventCategory,
				region,
				specifiedAccountId,
				startTime,
				endTime,
				lastUpdatedTime,
			)
			eventsResp, err := health.DescribeEventsForOrganization(ctx, healthClient, input)
			if err != nil {
				return err
			}

			eventChoices, eventsMap := ui.PrepareEventChoicesAndMap(eventsResp)
			selectedEvent, err := ui.PromptEventSelection(eventChoices, eventsMap)
			if err != nil {
				return err
			}

			eventArn := awssdk.ToString(selectedEvent.Arn)
			accountsMapping, err := organizations.GetAccountsMapping(ctx, orgClient)
			if err != nil {
				return err
			}

			var eventFileName string
			if specifiedFileName != "" {
				eventFileName = specifiedFileName
			} else {
				eventFileName = csv.GenerateEventFileName(
					selectedEvent,
					specifiedAccountId,
					statusCode,
				)
			}

			err = csv.WriteEventDetailsToCsv(
				ctx,
				healthClient,
				eventArn,
				eventFileName,
				specifiedAccountId,
				statusCode,
				accountsMapping,
				selectedEvent,
				echoToStdout,
			)
			if err != nil {
				return err
			}

			fmt.Printf("Event details have been written to %s.\n", eventFileName)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseEventFilter(filter string) map[string]string {
	result := make(map[string]string)
	// Regular expression pattern to extract key-value pairs
	// Example: key=value, key={value1,value2}
	pattern := regexp.MustCompile(`([^,=]+)=({[^}]*}|[^,]*)`)
	matches := pattern.FindAllStringSubmatch(filter, -1)

	for _, match := range matches {
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])
			result[key] = value
		}
	}
	return result
}
