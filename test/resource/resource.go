package resource_test

import (
	"github.com/hashicorp/terraform/providers"
	"github.com/khulnasoft-lab/driftctl/pkg/resource"
	"github.com/khulnasoft-lab/driftctl/pkg/resource/schemas"
	testschemas "github.com/khulnasoft-lab/driftctl/test/schemas"
)

func InitFakeSchemaRepository(provider, version string) resource.SchemaRepositoryInterface {
	repo := schemas.NewSchemaRepository()
	schema := make(map[string]providers.Schema)
	if provider != "" {
		s, err := testschemas.ReadTestSchema(provider, version)
		if err != nil {
			// TODO HANDLER ERROR PROPERLY
			panic(err)
		}
		schema = s
	}
	_ = repo.Init(provider, version, schema)
	return repo
}
