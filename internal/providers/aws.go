package providers

import (
	"context"

	appConfig "github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func CreateAwsConfig(ctx context.Context, awsConfig *appConfig.AWSConfig) (aws.Config, error) {
	var cfg aws.Config
	var err error
	if awsConfig.S3Endpoint != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(awsConfig.Region),
			config.WithBaseEndpoint(awsConfig.S3Endpoint),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				awsConfig.AccessKeyID,
				awsConfig.SecretAccessKey,
				"",
			)),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(awsConfig.Region))
	}

	return cfg, err
}
