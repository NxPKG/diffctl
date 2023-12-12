package enumeration

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/khulnasoft-lab/driftctl/enumeration/diagnostic"
	"github.com/khulnasoft-lab/driftctl/enumeration/resource"
)

type RefreshInput struct {
	// Resources to refresh
	Resources map[string][]*resource.Resource
}

type RefreshOutput struct {
	Resources   map[string][]*resource.Resource
	Diagnostics diagnostic.Diagnostics
}

type GetSchemasOutput struct {
	Schema *terraform.ProviderSchema
}

type Refresher interface {
	Refresh(input *RefreshInput) (*RefreshOutput, error)
	GetSchema() (*GetSchemasOutput, error)
}
