package google_test

import (
	"testing"
	"time"

	"github.com/khulnasoft-lab/driftctl/test"
	"github.com/khulnasoft-lab/driftctl/test/acceptance"
)

func TestAcc_Google_StorageBucket(t *testing.T) {
	t.Skip("flake")

	acceptance.Run(t, acceptance.AccTestCase{
		TerraformVersion: "0.15.5",
		Paths:            []string{"./testdata/acc/google_storage_bucket"},
		Args: []string{
			"scan",
			"--to", "gcp+tf",
		},
		Checks: []acceptance.AccCheck{
			{
				// New resources are not visible immediately through GCP API after an apply operation.
				ShouldRetry: acceptance.LinearBackoff(10 * time.Minute),
				Check: func(result *test.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertInfrastructureIsInSync()
					result.AssertManagedCount(3)
				},
			},
		},
	})
}
