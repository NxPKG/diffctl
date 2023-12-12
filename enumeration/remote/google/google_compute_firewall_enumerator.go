package google

import (
	remoteerror "github.com/khulnasoft-lab/driftctl/enumeration/remote/error"
	"github.com/khulnasoft-lab/driftctl/enumeration/remote/google/repository"
	"strings"

	"github.com/khulnasoft-lab/driftctl/enumeration/resource"
	"github.com/khulnasoft-lab/driftctl/enumeration/resource/google"
	"github.com/sirupsen/logrus"
)

type GoogleComputeFirewallEnumerator struct {
	repository repository.AssetRepository
	factory    resource.ResourceFactory
}

func NewGoogleComputeFirewallEnumerator(repo repository.AssetRepository, factory resource.ResourceFactory) *GoogleComputeFirewallEnumerator {
	return &GoogleComputeFirewallEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *GoogleComputeFirewallEnumerator) SupportedType() resource.ResourceType {
	return google.GoogleComputeFirewallResourceType
}

func (e *GoogleComputeFirewallEnumerator) Enumerate() ([]*resource.Resource, error) {
	resources, err := e.repository.SearchAllFirewalls()

	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(resources))

	for _, res := range resources {
		splittedName := strings.Split(res.GetName(), "/")
		if len(splittedName) != 8 {
			logrus.WithField("name", res.GetName()).Error("Unable to decode project from firewall name")
			continue
		}
		project := splittedName[4]
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				trimResourceName(res.GetName()),
				map[string]interface{}{
					"name":    res.DisplayName,
					"project": project,
				},
			),
		)
	}

	return results, err
}
