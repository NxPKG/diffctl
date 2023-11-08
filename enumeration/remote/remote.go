package remote

import (
	"github.com/khulnasoft-lab/driftctl/enumeration"
	"github.com/khulnasoft-lab/driftctl/enumeration/alerter"
	"github.com/khulnasoft-lab/driftctl/enumeration/remote/aws"
	"github.com/khulnasoft-lab/driftctl/enumeration/remote/azurerm"
	"github.com/khulnasoft-lab/driftctl/enumeration/remote/common"
	"github.com/khulnasoft-lab/driftctl/enumeration/remote/github"
	"github.com/khulnasoft-lab/driftctl/enumeration/remote/google"
	"github.com/khulnasoft-lab/driftctl/enumeration/resource"
	"github.com/khulnasoft-lab/driftctl/enumeration/terraform"
	"github.com/pkg/errors"
)

var supportedRemotes = []string{
	common.RemoteAWSTerraform,
	common.RemoteGithubTerraform,
	common.RemoteGoogleTerraform,
	common.RemoteAzureTerraform,
}

func IsSupported(remote string) bool {
	for _, r := range supportedRemotes {
		if r == remote {
			return true
		}
	}
	return false
}

func Activate(remote, version string, alerter alerter.AlerterInterface, providerLibrary *terraform.ProviderLibrary, remoteLibrary *common.RemoteLibrary, progress enumeration.ProgressCounter, factory resource.ResourceFactory, configDir string) error {
	switch remote {
	case common.RemoteAWSTerraform:
		return aws.Init(version, alerter, providerLibrary, remoteLibrary, progress, factory, configDir)
	case common.RemoteGithubTerraform:
		return github.Init(version, alerter, providerLibrary, remoteLibrary, progress, factory, configDir)
	case common.RemoteGoogleTerraform:
		return google.Init(version, alerter, providerLibrary, remoteLibrary, progress, factory, configDir)
	case common.RemoteAzureTerraform:
		return azurerm.Init(version, alerter, providerLibrary, remoteLibrary, progress, factory, configDir)

	default:
		return errors.Errorf("unsupported remote '%s'", remote)
	}
}

func GetSupportedRemotes() []string {
	return supportedRemotes
}
