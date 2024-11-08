package aws

import (
	"context"

	opengovernance "github.com/opengovern/og-describer-aws/SDK/generated"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEcrRegistry(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ecr_registry",
		Description: "AWS ECR Registry",
		Get: &plugin.GetConfig{
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"RepositoryNotFoundException", "RepositoryPolicyNotFoundException", "LifecyclePolicyNotFoundException"}),
			},
			Hydrate: opengovernance.GetECRRegistry,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListECRRegistry,
		},

		Columns: awsKaytuRegionalColumns([]*plugin.Column{
			{
				Name:        "registry_id",
				Description: "The Registry ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RegistryId")},
			{
				Name:        "rules",
				Description: "The Registry Rules.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.ReplicationRules"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description.RegistryId")},
		}),
	}
}
