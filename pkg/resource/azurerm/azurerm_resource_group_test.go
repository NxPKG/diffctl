package azurerm_test

import (
	"testing"

	"github.com/khulnasoft-lab/driftctl/test"
	"github.com/khulnasoft-lab/driftctl/test/acceptance"
)

func TestAcc_Azure_ResourceGroup(t *testing.T) {
	acceptance.Run(t, acceptance.AccTestCase{
		TerraformVersion: "0.15.5",
		Paths:            []string{"./testdata/acc/azurerm_resource_group"},
		Args: []string{
			"scan",
			"--to", "azure+tf",
			"--filter", "contains(Id, 'acc-test-res-group-')",
		},
		Checks: []acceptance.AccCheck{
			{
				Check: func(result *test.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertInfrastructureIsInSync()
					result.AssertManagedCount(2)
				},
			},
		},
	})
}
