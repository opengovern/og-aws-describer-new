package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetRoleArnFromName(accountId string, roleName string) string {
	if roleName == "" {
		return ""
	}
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, roleName)
}

func GetPolicyArnFromName(accountId string, policyName string) string {
	if policyName == "" {
		return ""
	}
	return fmt.Sprintf("arn:aws:iam::%s:policy/%s", accountId, policyName)
}

// GetConfig loads the AWS credentionals and returns the configuration to be used by the AWS services client.
// If the awsAccessKey is specified, the config will be created for the combination of awsAccessKey, awsSecretKey, awsSessionToken.
// Else it will use the default AWS SDK logic to load the configuration. See https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
// If assumeRoleArn is provided, it will use the evaluated configuration to then assume the specified role.
func GetConfig(ctx context.Context, awsAccessKey, awsSecretKey, awsSessionToken, assumeRoleArn string, externalId *string) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error

	if awsAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, awsSessionToken)))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	if externalId != nil && *externalId == "" {
		externalId = nil
	}

	if assumeRoleArn != "" {
		cfg, err = config.LoadDefaultConfig(
			context.Background(),
			config.WithCredentialsProvider(
				stscreds.NewAssumeRoleProvider(
					sts.NewFromConfig(cfg),
					assumeRoleArn,
					func(o *stscreds.AssumeRoleOptions) {
						o.ExternalID = externalId
					},
				),
			),
		)
		if err != nil {
			return aws.Config{}, fmt.Errorf("failed to assume role: %w", err)
		}
	}

	return cfg, nil
}

type AccountConfig struct {
	AccountID            string   `json:"accountId"`
	Regions              []string `json:"regions"`
	SecretKey            string   `json:"secretKey"`
	AccessKey            string   `json:"accessKey"`
	SessionToken         string   `json:"sessionToken"`
	AssumeRoleName       string   `json:"assumeRoleName"`
	ExternalID           *string  `json:"externalId,omitempty"`
	AssumeAdminRoleName  string   `json:"assumeAdminRoleName,omitempty"`
	AssumeRolePolicyName string   `json:"assumeRolePolicyName,omitempty"`
}

func AccountConfigFromMap(m map[string]any) (AccountConfig, error) {
	mj, err := json.Marshal(m)
	if err != nil {
		return AccountConfig{}, err
	}

	var c AccountConfig
	err = json.Unmarshal(mj, &c)
	if err != nil {
		return AccountConfig{}, err
	}

	return c, nil
}
