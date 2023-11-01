package aws

import (
	"github.com/khulnasoft-lab/driftctl/enumeration/remote/aws/repository"
	remoteerror "github.com/khulnasoft-lab/driftctl/enumeration/remote/error"
	"github.com/khulnasoft-lab/driftctl/enumeration/resource"
	"github.com/khulnasoft-lab/driftctl/enumeration/resource/aws"
)

type ApiGatewayV2ApiEnumerator struct {
	repository repository.ApiGatewayV2Repository
	factory    resource.ResourceFactory
}

func NewApiGatewayV2ApiEnumerator(repo repository.ApiGatewayV2Repository, factory resource.ResourceFactory) *ApiGatewayV2ApiEnumerator {
	return &ApiGatewayV2ApiEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *ApiGatewayV2ApiEnumerator) SupportedType() resource.ResourceType {
	return aws.AwsApiGatewayV2ApiResourceType
}

func (e *ApiGatewayV2ApiEnumerator) Enumerate() ([]*resource.Resource, error) {
	apis, err := e.repository.ListAllApis()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(apis))

	for _, api := range apis {
		a := api
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				*a.ApiId,
				map[string]interface{}{},
			),
		)
	}
	return results, err
}
