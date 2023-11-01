package azurerm

import (
	"github.com/khulnasoft-lab/driftctl/enumeration/resource"
	dctlresource "github.com/khulnasoft-lab/driftctl/pkg/resource"
)

const AzurePostgresqlDatabaseResourceType = "azurerm_postgresql_database"

func initAzurePostgresqlDatabaseMetadata(resourceSchemaRepository dctlresource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetHumanReadableAttributesFunc(AzurePostgresqlDatabaseResourceType, func(res *resource.Resource) map[string]string {
		val := res.Attrs
		attrs := make(map[string]string)
		if name := val.GetString("name"); name != nil && *name != "" {
			attrs["Name"] = *name
		}
		return attrs
	})
}
