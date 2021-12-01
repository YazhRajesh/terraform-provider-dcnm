package acctest

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDCNMServiceNodeDataSource_Basic(t *testing.T) {
	datasourcename := "data.dcnm_service_node.test"
	resourceName := "dcnm_service_node.test"
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerServiceNode),
		CheckDestroy:      testAccCheckDCNMServiceNodeDSDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateServiceNodeDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceNodeDSWithoutServiceFabric(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateServiceNodeDSResource(rName),
			},
			{
				Config: CreateServiceNodeDS(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourcename, "node_type", resourceName, "node_type"),
					resource.TestCheckResourceAttrPair(datasourcename, "attached_fabric", resourceName, "attached_fabric"),
					resource.TestCheckResourceAttrPair(datasourcename, "attached_switch_interface_name", resourceName, "attached_switch_interface_name"),
					resource.TestCheckResourceAttrPair(datasourcename, "interface_name", resourceName, "interface_name"),
					resource.TestCheckResourceAttrPair(datasourcename, "switches.#", resourceName, "switches.#"),
					resource.TestCheckTypeSetElemAttrPair(datasourcename, "switches.0", resourceName, "switches.0"),
					resource.TestCheckResourceAttrPair(datasourcename, "admin_state", resourceName, "admin_state"),
					resource.TestCheckResourceAttrPair(datasourcename, "allowed_vlans", resourceName, "allowed_vlans"),
					resource.TestCheckResourceAttrPair(datasourcename, "bpdu_guard_flag", resourceName, "bpdu_guard_flag"),
					resource.TestCheckResourceAttrPair(datasourcename, "form_factor", resourceName, "form_factor"),
					resource.TestCheckResourceAttrPair(datasourcename, "link_template_name", resourceName, "link_template_name"),
					resource.TestCheckResourceAttrPair(datasourcename, "mtu", resourceName, "mtu"),
					resource.TestCheckResourceAttrPair(datasourcename, "porttype_fast_enabled", resourceName, "porttype_fast_enabled"),
					resource.TestCheckResourceAttrPair(datasourcename, "speed", resourceName, "speed"),
				),
			},
		},
	})
}

func CreateServiceNodeDSWithoutServiceFabric(rName string) string {

	fmt.Println("=== STEP  testing Service Node data source without name")

	res := fmt.Sprintf(`
	resource "dcnm_service_node" "test" {
		name                           = "%s"
		node_type                      = "%s"
		service_fabric                 = "%s"
		attached_fabric                = "%s"
		attached_switch_interface_name = "%s"
		interface_name                 = "%s"
		switches                       = %s
	}

	data "dcnm_service_node" "test" {
		name = dcnm_service_node.test.name
	}
	`, rName, nodeTypeDefault, serviceFabricName, attachedFabricName, attachedSwitchInterfaceName, interfaceName, convertToQuotedStringArray([]string{sw1}))
	return res
}

func CreateServiceNodeDSWithoutName(rName string) string {

	fmt.Println("=== STEP  testing Service Node data source without service_fabric")

	res := fmt.Sprintf(`
	resource "dcnm_service_node" "test" {
		name                           = "%s"
		node_type                      = "%s"
		service_fabric                 = "%s"
		attached_fabric                = "%s"
		attached_switch_interface_name = "%s"
		interface_name                 = "%s"
		switches                       = %s
	}

	data "dcnm_service_node" "test" {
		service_fabric = dcnm_service_node.test.service_fabric
	}
	`, rName, nodeTypeDefault, serviceFabricName, attachedFabricName, attachedSwitchInterfaceName, interfaceName, convertToQuotedStringArray([]string{sw1}))
	return res
}

func CreateServiceNodeDSResource(rName string) string {

	fmt.Println("=== STEP  creating resource for data source")

	res := fmt.Sprintf(`
	resource "dcnm_service_node" "test" {
		name                           = "%s"
		node_type                      = "%s"
		service_fabric                 = "%s"
		attached_fabric                = "%s"
		attached_switch_interface_name = "%s"
		interface_name                 = "%s"
		switches                       = %s
	}
	`, rName, nodeTypeDefault, serviceFabricName, attachedFabricName, attachedSwitchInterfaceName, interfaceName, convertToQuotedStringArray([]string{sw1}))
	return res
}

func CreateServiceNodeDS(rName string) string {

	fmt.Println("=== STEP  testing data source for Service Node")

	res := fmt.Sprintf(`
	resource "dcnm_service_node" "test" {
		name                           = "%s"
		node_type                      = "%s"
		service_fabric                 = "%s"
		attached_fabric                = "%s"
		attached_switch_interface_name = "%s"
		interface_name                 = "%s"
		switches                       = %s
	}

	data "dcnm_service_node" "test" {
		name                           = "%s"
		service_fabric                 = "%s"
	}
	`, rName, nodeTypeDefault, serviceFabricName, attachedFabricName, attachedSwitchInterfaceName, interfaceName, convertToQuotedStringArray([]string{sw1}), rName, serviceFabricName)
	return res
}

func testAccCheckDCNMServiceNodeDSDestroy(s *terraform.State) error {
	dcnmClient := (*providerServiceNode).Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_service_node" {
			primaryIdArr := strings.Split(rs.Primary.ID, "/")
			serviceNodeId := primaryIdArr[2]
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/appcenter/cisco/ndfc/api/v1/elastic-service/fabrics/%s/service-nodes/%s", serviceFabricName, serviceNodeId))
			if err == nil {
				return fmt.Errorf("Service Node still exists")
			}
		}
	}

	return nil
}
