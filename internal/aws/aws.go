package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
)

func LoadAWSConfig(ctx context.Context, profile string) (aws.Config, error) {
	loadOptions := []func(*config.LoadOptions) error{
		config.WithRegion("us-east-1"),
	}

	if profile != "" {
		loadOptions = append(loadOptions, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

func CreateServices(cfg aws.Config) (*health.Client, *organizations.Client) {
	healthClient := health.NewFromConfig(cfg)
	orgClient := organizations.NewFromConfig(cfg)
	return healthClient, orgClient
}
