package aws_test

import (
	"testing"

	"github.com/khulnasoft-lab/driftctl/test"
	"github.com/khulnasoft-lab/driftctl/test/acceptance"
)

func TestAcc_Aws_ApiGatewayRestApi(t *testing.T) {
	acceptance.Run(t, acceptance.AccTestCase{
		TerraformVersion: "0.15.5",
		Paths:            []string{"./testdata/acc/aws_api_gateway_rest_api"},
		Args:             []string{"scan"},
		Checks: []acceptance.AccCheck{
			{
				Env: map[string]string{
					"AWS_REGION": "us-east-1",
				},
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
