package main

import (
	"context"
	"fmt"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hayao-k/health-exporter/internal/aws"
	"github.com/hayao-k/health-exporter/internal/aws/health"
	"github.com/hayao-k/health-exporter/internal/aws/organizations"
	"github.com/hayao-k/health-exporter/internal/csv"
	"github.com/hayao-k/health-exporter/internal/ui"
	"github.com/urfave/cli/v2"
	"os"
)

var version = "v0.0.0"

func main() {
	app := &cli.App{
		Name:    "health-exporter",
		Usage:   "Describe AWS Health events for your organization",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "service",
				Aliases: []string{"s"},
				Usage:   "Filter events by service name, e.g., RDS",
			},
			&cli.StringFlag{
				Name:    "status",
				Aliases: []string{"t"},
				Usage:   "Filter events by status. Possible values are open, closed and upcoming",
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
				Name:    "account-id",
				Aliases: []string{"i"},
				Usage:   "Specify a single account ID to process",
			},
		},
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			service := c.String("service")
			status := c.String("status")
			echoToStdout := c.Bool("echo")
			profile := c.String("profile")
			specifiedAccountId := c.String("account-id")

			cfg, err := aws.LoadAWSConfig(ctx, profile)
			if err != nil {
				return err
			}

			healthClient, orgClient := aws.CreateServices(cfg)

			input := health.DescribeEventsForOrganizationInput(service, status, specifiedAccountId)
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

			eventFileName := csv.GenerateEventFileName(selectedEvent, specifiedAccountId)
			err = csv.WriteEventDetailsToCsv(
				ctx,
				healthClient,
				eventArn,
				accountsMapping,
				selectedEvent,
				eventFileName,
				echoToStdout,
				specifiedAccountId,
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
