package google

import (
	"github.com/khulnasoft-lab/driftctl/enumeration/resource"
	dctlresource "github.com/khulnasoft-lab/driftctl/pkg/resource"
)

const GoogleBigqueryDatasetResourceType = "google_bigquery_dataset"

func initGoogleBigqueryDatasetMetadata(resourceSchemaRepository dctlresource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetHumanReadableAttributesFunc(GoogleBigqueryDatasetResourceType, func(res *resource.Resource) map[string]string {
		return map[string]string{
			"name": *res.Attrs.GetString("friendly_name"),
		}
	})
}
