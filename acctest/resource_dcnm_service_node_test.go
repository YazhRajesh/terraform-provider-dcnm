package acctest

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerServiceNode *schema.Provider
var listParams = []string{"switches"}

const serviceFabricName = "external"
const attachedFabricName = "terraform"
const attachedSwitchInterfaceName = "Ethernet1/9"
const interfaceName = "node"
const sw1 = "9LMU8W6W8VG"
const sw2 = "9CIWTMB13GP"
const sw3 = "9FTD5XQKR46"
const nodeTypeDefault = "Firewall"

func TestAccDCNMServiceNode_Basic(t *testing.T) {
	var serviceNodeDefaultId string
	var serviceNodeUpdatedId string
	swtiches := []string{sw1}
	resourceName := "dcnm_service_node.test"
	rName := acctest.RandString(5)
	rNameOther := acctest.RandString(5)
	var nodeTypeValues [3]string
	nodeTypeValues[0] = "Firewall"
	nodeTypeValues[1] = "ADC"
	nodeTypeValues[2] = "VNF"
	m := make(map[string]interface{})
	m["name"] = rName
	m["node_type"] = nodeTypeValues[0]
	m["service_fabric"] = serviceFabricName
	m["attached_fabric"] = attachedFabricName
	m["attached_switch_interface_name"] = attachedSwitchInterfaceName
	m["interface_name"] = interfaceName
	m["switches"] = swtiches
	optmap := make(map[string]interface{})
	optmap["name"] = rName
	optmap["node_type"] = nodeTypeValues[0]
	optmap["service_fabric"] = serviceFabricName
	optmap["attached_fabric"] = attachedFabricName
	optmap["attached_switch_interface_name"] = attachedSwitchInterfaceName
	optmap["interface_name"] = interfaceName
	optmap["switches"] = swtiches
	optmap["form_factor"] = "Physical"
	optmap["bpdu_guard_flag"] = "true"
	optmap["porttype_fast_enabled"] = "false"
	optmap["admin_state"] = "false"
	optmap["policy_description"] = "sample_policy_description"
	optmap["mtu"] = "default"
	optmap["speed"] = "10Mb"
	optmap["allowed_vlans"] = "all"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerServiceNode),
		CheckDestroy:      testAccCheckDCNMServiceNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateServiceNode([]string{"name"}, m, listParams),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceNode([]string{"node_type"}, m, listParams),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceNode([]string{"service_fabric"}, m, listParams),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceNode([]string{"attached_fabric"}, m, listParams),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceNode([]string{"attached_switch_interface_name"}, m, listParams),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceNode([]string{"interface_name"}, m, listParams),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceNode([]string{"switches"}, m, listParams),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateServiceNode([]string{}, m, listParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, nodeTypeValues[0], rName, &serviceNodeDefaultId),
					resource.TestCheckResourceAttr(resourceName, "admin_state", "true"),
					resource.TestCheckResourceAttr(resourceName, "allowed_vlans", "none"),
					resource.TestCheckResourceAttr(resourceName, "bpdu_guard_flag", "no"),
					resource.TestCheckResourceAttr(resourceName, "form_factor", "Virtual"),
					resource.TestCheckResourceAttr(resourceName, "link_template_name", "service_link_trunk"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "jumbo"),
					resource.TestCheckResourceAttr(resourceName, "porttype_fast_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "speed", "Auto"),
				),
			},
			{

				Config: CreateServiceNode([]string{}, optmap, listParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, nodeTypeValues[0], rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "form_factor", "Physical"),
					resource.TestCheckResourceAttr(resourceName, "bpdu_guard_flag", "true"),
					resource.TestCheckResourceAttr(resourceName, "porttype_fast_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "admin_state", "false"),
					resource.TestCheckResourceAttr(resourceName, "policy_description", "sample_policy_description"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "default"),
					resource.TestCheckResourceAttr(resourceName, "speed", "10Mb"),
					resource.TestCheckResourceAttr(resourceName, "allowed_vlans", "all"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"policy_description"},
				ImportStateIdFunc:       testServiceNodeImportStateIdFunc(resourceName),
			},
			{
				Config: CreateServiceNodeByReplacingValueOfKey(m, "name", rNameOther, listParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, nodeTypeValues[0], rNameOther, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "name", rNameOther),
					testAccCheckDCNMServiceNodeIdNotEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNode([]string{}, m, listParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, nodeTypeValues[0], rName, &serviceNodeDefaultId),
				),
			},
			{
				Config: CreateServiceNodeByReplacingValueOfKey(m, "node_type", nodeTypeValues[1], listParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, nodeTypeValues[1], rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "node_type", nodeTypeValues[1]),
					testAccCheckDCNMServiceNodeIdNotEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNode([]string{}, m, listParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, nodeTypeValues[0], rName, &serviceNodeDefaultId),
				),
			},
			{
				Config: CreateServiceNodeByReplacingValueOfKey(m, "node_type", nodeTypeValues[2], listParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, nodeTypeValues[2], rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "node_type", nodeTypeValues[2]),
					testAccCheckDCNMServiceNodeIdNotEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
		},
	})
}

func TestAccDCNMServiceNode_Update(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "dcnm_service_node.test"
	var serviceNodeDefaultId string
	var serviceNodeUpdatedId string
	switches := []string{sw1}
	m := make(map[string]interface{})
	m["name"] = rName
	m["node_type"] = "Firewall"
	m["service_fabric"] = serviceFabricName
	m["attached_fabric"] = attachedFabricName
	m["attached_switch_interface_name"] = attachedSwitchInterfaceName
	m["interface_name"] = interfaceName
	m["switches"] = switches
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerServiceNode),
		CheckDestroy:      testAccCheckDCNMServiceNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateServiceNode([]string{}, m, listParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeDefaultId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "100Mb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "100Mb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "1Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "1Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "2.5Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "2.5Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "5Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "5Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "10Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "10Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "25Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "25Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "40Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "40Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "50Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "50Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "100Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "100Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "200Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "200Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "speed", "400Gb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "speed", "400Gb"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "allowed_vlans", "1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "allowed_vlans", "1"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "allowed_vlans", "200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "allowed_vlans", "200"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "allowed_vlans", "500"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "allowed_vlans", "500"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "allowed_vlans", "2000"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "allowed_vlans", "2000"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "allowed_vlans", "3000"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "allowed_vlans", "3000"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
			{
				Config: CreateServiceNodeByAddingParamAndValue(rName, "bpdu_guard_flag", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDCNMServiceNodeExists(resourceName, serviceFabricName, attachedFabricName, "Firewall", rName, &serviceNodeUpdatedId),
					resource.TestCheckResourceAttr(resourceName, "bpdu_guard_flag", "false"),
					testAccCheckDCNMServiceNodeIdEqual(&serviceNodeDefaultId, &serviceNodeUpdatedId),
				),
			},
		},
	})
}

func TestAccDCNMServiceNode_NegativeCases(t *testing.T) {
	rName := acctest.RandString(5)
	randomVal := acctest.RandString(5)
	longName := acctest.RandString(64)
	randomParam := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuv")
	switches := []string{sw1}
	m := make(map[string]interface{})
	m["name"] = rName
	m["node_type"] = "Firewall"
	m["service_fabric"] = serviceFabricName
	m["attached_fabric"] = attachedFabricName
	m["attached_switch_interface_name"] = attachedSwitchInterfaceName
	m["interface_name"] = interfaceName
	m["switches"] = switches
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerServiceNode),
		CheckDestroy:      testAccCheckDCNMServiceNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateServiceNodeByReplacingValueOfKey(m, "node_type", randomVal, listParams),
				ExpectError: regexp.MustCompile(`expected node_type to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateServiceNodeByAddingParamAndValue(rName, "form_factor", randomVal),
				ExpectError: regexp.MustCompile(`expected form_factor to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateServiceNodeByAddingParamAndValue(rName, "bpdu_guard_flag", randomVal),
				ExpectError: regexp.MustCompile(`expected bpdu_guard_flag to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateServiceNodeByAddingParamAndValue(rName, "porttype_fast_enabled", randomVal),
				ExpectError: regexp.MustCompile(`expected porttype_fast_enabled to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateServiceNodeByAddingParamAndValue(rName, "admin_state", randomVal),
				ExpectError: regexp.MustCompile(`expected admin_state to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateServiceNodeByAddingParamAndValue(rName, randomParam, randomVal),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateServiceNodeByReplacingValueOfKey(m, "name", longName, listParams),
				ExpectError: regexp.MustCompile(`(.)*(value too long for type VARCHAR)(.)*`),
			},
			{
				Config:      CreateServiceNodeByReplacingValueOfKey(m, "service_fabric", randomVal, listParams),
				ExpectError: regexp.MustCompile(`(.)*(Cannot find the specified external fabric)(.)*`),
			},
			{
				Config:      CreateServiceNodeByAddingParamAndValue(rName, "speed", randomVal),
				ExpectError: regexp.MustCompile(`(.)*(Validation failed for following fields:)(.)*`),
			},
			{
				Config:      CreateServiceNodeByAddingParamAndValue(rName, "mtu", randomVal),
				ExpectError: regexp.MustCompile(`(.)*(Validation failed for following fields:)(.)*`),
			},
			{
				Config:      CreateServiceNodeByAddingParamAndValue(rName, "allowed_vlans", randomVal),
				ExpectError: regexp.MustCompile(`(.)*(Invalid values found in the 'Trunk Allowed Vlans' field)(.)*`),
			},
			{
				Config:      CreateServiceNodeByReplacingValueOfKey(m, "switches", []string{sw1, sw2, sw3}, listParams),
				ExpectError: regexp.MustCompile(`(.)*(Upto 2 switches only allowed)`),
			},
			{
				Config: CreateServiceNode([]string{}, m, listParams),
			},
		},
	})
}

func testAccCheckDCNMServiceNodeExists(name, serviceFabricName, attachedFabricName, nodeType, rName string, nodeId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Service Node %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service Node id was set")
		}

		*nodeId = rs.Primary.ID
		expectedServiceNodeId := fmt.Sprintf("%s/%s/%s/%s", serviceFabricName, attachedFabricName, nodeType, rName)
		if expectedServiceNodeId != rs.Primary.ID {
			return fmt.Errorf("Service node with id %s doesn't exist", expectedServiceNodeId)
		}
		return nil
	}
}

func testAccCheckDCNMServiceNodeDestroy(s *terraform.State) error {
	dcnmClient := (*providerServiceNode).Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_service_node" {
			primaryIdArr := strings.Split(rs.Primary.ID, "/")
			serviceNodeId := primaryIdArr[3]
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/appcenter/cisco/ndfc/api/v1/elastic-service/fabrics/%s/service-nodes/%s", serviceFabricName, serviceNodeId))
			if err == nil {
				return fmt.Errorf("Service Node still exists")
			}
		}
	}

	return nil
}

func CreateServiceNode(excludeParams []string, m map[string]interface{}, listParams []string) string {
	if len(excludeParams) > 0 {
		fmt.Printf("=== STEP  testing Service Node creation without %v\n", excludeParams)
	} else {
		fmt.Println("=== STEP  creating Service Node with all required parameters")
	}
	res := "resource \"dcnm_service_node\" \"test\"{ \n"
	for k, v := range m {
		exclude := false
		list := false
		for _, excludeParam := range excludeParams {
			if k == excludeParam {
				exclude = true
			}
		}
		for _, listParam := range listParams {
			if k == listParam {
				list = true
			}
		}
		if !exclude && list {
			valList := v.([]string)
			valListArr := convertToQuotedStringArray(valList)
			res = res + string(k) + "=" + valListArr + "\n"
		} else if !exclude {
			res = res + string(k) + "=" + QuotedString(v.(string)) + "\n"
		}
	}
	res = res + "}"
	return res
}

func CreateServiceNodeByReplacingValueOfKey(m map[string]interface{}, key string, value interface{}, listParams []string) string {
	fmt.Printf("=== STEP  testing by updating value of parameter %s with %v\n", key, value)
	res := "resource \"dcnm_service_node\" \"test\"{ \n"
	for k, v := range m {
		list := false
		for _, listParam := range listParams {
			if k == listParam {
				list = true
			}
		}
		if k == key {
			v = value
		}
		if list {
			valList := v.([]string)
			valListArr := convertToQuotedStringArray(valList)
			res = res + string(k) + "=" + valListArr + "\n"
		} else {
			res = res + string(k) + "=" + QuotedString(v.(string)) + "\n"
		}
	}
	res = res + "}"
	return res
}

func CreateServiceNodeByAddingParamAndValue(rName, key, value string) string {

	fmt.Printf("=== STEP  testing Service Node creation with %s : %s\n", key, value)

	res := fmt.Sprintf(`
	resource "dcnm_service_node" "test" {
		name                           = "%s"
		node_type                      = "%s"
		service_fabric                 = "%s"
		attached_fabric                = "%s"
		attached_switch_interface_name = "%s"
		interface_name                 = "%s"
		switches                       = %s
		%s						   = "%s"
	  }
	`, rName, nodeTypeDefault, serviceFabricName, attachedFabricName, attachedSwitchInterfaceName, interfaceName, convertToQuotedStringArray([]string{sw1}), key, value)
	return res
}

func testAccCheckDCNMServiceNodeIdEqual(id1, id2 *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *id1 != *id2 {
			return fmt.Errorf("Ids of service nodes are different")
		}
		return nil
	}
}

func testAccCheckDCNMServiceNodeIdNotEqual(id1, id2 *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *id1 == *id2 {
			return fmt.Errorf("Ids of service nodes are equal")
		}
		return nil
	}
}

func testServiceNodeImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found %s", resourceName)
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["service_fabric"], rs.Primary.Attributes["name"]), nil
	}
}
