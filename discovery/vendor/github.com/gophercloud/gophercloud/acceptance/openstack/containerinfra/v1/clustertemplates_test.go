// +build acceptance containerinfra

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	//"github.com/gophercloud/gophercloud/acceptance/tools"
	//"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestClusterTemplatesCRUD(t *testing.T) {
	client, err := clients.NewContainerInfraV1Client()
	th.AssertNoErr(t, err)

	clusterTemplate, err := CreateClusterTemplate(t, client)
	th.AssertNoErr(t, err)
	t.Log(clusterTemplate.Name)
	/*
		defer DeleteClusterTemplate(t, client, clusterTemplate.UUID)

		// Test clusters list
		allPages, err := clustertemplates.List(client, nil).AllPages()
		th.AssertNoErr(t, err)

		allClusterTemplates, err := clustertemplates.ExtractClusterTemplates(allPages)
		th.AssertNoErr(t, err)

		var found bool
		for _, v := range allClusterTemplates {
			if v.UUID == clusterTemplate.UUID {
				found = true
			}
		}

		th.AssertEquals(t, found, true)

		// Test cluster update
		updateOpts := []clustertemplates.UpdateOptsBuilder{
			clustertemplates.UpdateOpts{
				Path:  "/master_lb_enabled",
				Value: "false",
				Op:    "replace",
			},
			clustertemplates.UpdateOpts{
				Path:  "/registry_enabled",
				Value: "false",
				Op:    "replace",
			},
			clustertemplates.UpdateOpts{
				Path:  "/labels/test",
				Value: "test",
				Op:    "add",
			},
		}

		updateClusterTemplate, err := clustertemplates.Update(client, clusterTemplate.UUID, updateOpts).Extract()
		th.AssertNoErr(t, err)
		th.AssertEquals(t, false, updateClusterTemplate.MasterLBEnabled)
		th.AssertEquals(t, false, updateClusterTemplate.RegistryEnabled)
		th.AssertEquals(t, "test", updateClusterTemplate.Labels["test"])
		tools.PrintResource(t, updateClusterTemplate)

	*/
}
