package aws

import (
	"github.com/khulnasoft-lab/driftctl/enumeration/resource"
	"github.com/khulnasoft-lab/driftctl/pkg/helpers"
	dctlresource "github.com/khulnasoft-lab/driftctl/pkg/resource"
)

const AwsSnsTopicPolicyResourceType = "aws_sns_topic_policy"

func initSnsTopicPolicyMetaData(resourceSchemaRepository dctlresource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetNormalizeFunc(AwsSnsTopicPolicyResourceType, func(res *resource.Resource) {
		val := res.Attrs
		val.SafeDelete([]string{"owner"})
		jsonString, err := helpers.NormalizeJsonString((*val)["policy"])
		if err != nil {
			return
		}
		_ = val.SafeSet([]string{"policy"}, jsonString)
	})
	resourceSchemaRepository.UpdateSchema(AwsSnsTopicPolicyResourceType, map[string]func(attributeSchema *resource.AttributeSchema){
		"policy": func(attributeSchema *resource.AttributeSchema) {
			attributeSchema.JsonString = true
		},
	})
}
