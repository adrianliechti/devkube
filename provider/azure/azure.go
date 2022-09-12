package azure

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/adrianliechti/devkube/pkg/to"
	"github.com/adrianliechti/devkube/provider"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Provider struct {
	location string

	credential *azidentity.DefaultAzureCredential

	resourcegroups  *armresources.ResourceGroupsClient
	managedclusters *armcontainerservice.ManagedClustersClient
}

func NewFromEnvironment() (provider.Provider, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	if tenantID == "" {
		return nil, errors.New("AZURE_TENANT_ID is not set")
	}

	if subscriptionID == "" {
		return nil, errors.New("AZURE_SUBSCRIPTION_ID is not set")
	}

	//group := os.Getenv("AZURE_DEFAULTS_GROUP")

	location := os.Getenv("AZURE_DEFAULTS_LOCATION")

	if location == "" {
		location = "westeurope"
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, err
	}

	rgClient, err := armresources.NewResourceGroupsClient(subscriptionID, credential, nil)

	if err != nil {
		return nil, err
	}

	mcclient, err := armcontainerservice.NewManagedClustersClient(subscriptionID, credential, nil)

	if err != nil {
		return nil, err
	}

	return &Provider{
		location:   location,
		credential: credential,

		resourcegroups:  rgClient,
		managedclusters: mcclient,
	}, nil
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	var list []string

	pager := p.managedclusters.NewListPager(nil)

	for pager.More() {
		result, err := pager.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		if result.ManagedClusterListResult.Value != nil {
			for _, cluster := range result.ManagedClusterListResult.Value {
				list = append(list, *cluster.Name)
			}
		}
	}

	return list, nil
}

func (p *Provider) Create(ctx context.Context, name string, kubeconfig string) error {
	resourcegroup := groupName(name)

	exists, err := p.resourcegroups.CheckExistence(ctx, resourcegroup, nil)

	if err != nil {
		return err
	}

	if exists.Success {
		return errors.New("resource group already exists")
	}

	if _, err := p.resourcegroups.CreateOrUpdate(ctx, resourcegroup, armresources.ResourceGroup{Location: to.StringPtr(p.location)}, nil); err != nil {
		return err
	}

	identityType := armcontainerservice.ResourceIdentityTypeSystemAssigned

	sku := armcontainerservice.ManagedClusterSKUNameBasic
	skuTier := armcontainerservice.ManagedClusterSKUTierFree

	agentpoolMode := armcontainerservice.AgentPoolModeSystem
	agentpoolSize := 1

	parameters := armcontainerservice.ManagedCluster{
		Location: to.StringPtr(p.location),

		Identity: &armcontainerservice.ManagedClusterIdentity{
			Type: &identityType,
		},

		Properties: &armcontainerservice.ManagedClusterProperties{
			AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
				{
					Mode: &agentpoolMode,
					Name: to.StringPtr("agentpool"),

					Count: to.Int32Ptr(int32(agentpoolSize)),

					VMSize: to.StringPtr("Standard_B8ms"),
				},
			},

			DNSPrefix: to.StringPtr(name),

			ServicePrincipalProfile: &armcontainerservice.ManagedClusterServicePrincipalProfile{
				ClientID: to.StringPtr(resourcegroup),
			},
		},

		SKU: &armcontainerservice.ManagedClusterSKU{
			Name: &sku,
			Tier: &skuTier,
		},
	}

	poller, err := p.managedclusters.BeginCreateOrUpdate(ctx, resourcegroup, name, parameters, nil)

	if err != nil {
		return err
	}

	result, err := poller.PollUntilDone(ctx, nil)

	if err != nil {
		return err
	}

	_ = result
	return p.Export(ctx, name, kubeconfig)
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	resourcegroup := groupName(name)

	poller, err := p.managedclusters.BeginDelete(ctx, resourcegroup, name, nil)

	if err != nil {
		return err
	}

	result, err := poller.PollUntilDone(ctx, nil)

	if err != nil {
		return err
	}

	_ = result
	return nil
}

func (p *Provider) Export(ctx context.Context, name, kubeconfig string) error {
	if kubeconfig == "" {
		home, err := os.UserHomeDir()

		if err != nil {
			return err
		}

		dir := path.Join(home, ".kube")

		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}

		kubeconfig = path.Join(home, ".kube", "config")
	}

	resourcegroup := groupName(name)

	result, err := p.managedclusters.ListClusterAdminCredentials(ctx, resourcegroup, name, nil)

	if err != nil {
		return err
	}

	if len(result.Kubeconfigs) == 0 {
		return errors.New("unable to get kubeconfig")
	}

	data := result.Kubeconfigs[0].Value

	return os.WriteFile(kubeconfig, data, 0600)
}

func groupName(name string) string {
	name = strings.ToLower(name)

	if strings.EqualFold(name, "devkube") || strings.HasPrefix(name, "devkube-") {
		return name
	}

	return "devkube-" + name
}
