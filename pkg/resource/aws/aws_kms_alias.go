package aws

import (
	"github.com/khulnasoft-lab/driftctl/enumeration/resource"
	dctlresource "github.com/khulnasoft-lab/driftctl/pkg/resource"
)

const AwsKmsAliasResourceType = "aws_kms_alias"

func initAwsKmsAliasMetaData(resourceSchemaRepository dctlresource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetNormalizeFunc(AwsKmsAliasResourceType, func(res *resource.Resource) {
		val := res.Attrs
		val.SafeDelete([]string{"name"})
		val.SafeDelete([]string{"name_prefix"})
	})
}
