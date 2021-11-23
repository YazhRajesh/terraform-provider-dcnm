package acctest

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerServicePolicy *schema.Provider

func TestAccDCNMServiceNode_Basic(t *testing.T) {
	var serviceNodeDefaultId string
	var serviceNodeUpdatedId string
	resourceName := "dcnm_service_node.first"
	rName := acctest.RandString(5)
	nodeTypeInd:=acctest.RandStringFromCharSet(1,"012")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerServicePolicy),
		// CheckDestroy:      testAccCheckDCNMPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// terraform will try to create service policy without required argument policy_name
				Config:      CreateAccServicePolicyWithoutPolicy_name(rName, ip), // configuration to check creation of service policy without policy_name
				ExpectError: regexp.MustCompile(`Missing required argument`),     // test step expect error which should be match with defined regex
			},
			{
				// terraform will try to create service policy without required argument service_fabric
				Config:      CreateAccServicePolicyWithoutService_fabric(rName, ip), // configuration to check creation of service policy without service_fabric
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// terraform will try to create service policy without required argument attached_fabric
				Config:      CreateAccServicePolicyWithoutAttached_fabric(rName, ip), // configuration to check creation of service policy without attached_fabric
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// terraform will try to create service policy without required argument dest_network
				Config:      CreateAccServicePolicyWithoutDest_network(rName, ip), // configuration to check creation of service policy without dest_network
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// terraform will try to create service policy without required argument dest_vrf_name
				Config:      CreateAccServicePolicyWithoutDest_vrf_name(rName, ip), // configuration to check creation of service policy without dest_vrf_name
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// terraform will try to create service policy without required argument next_hop_ip
				Config:      CreateAccServicePolicyWithoutNext_hop_ip(rName), // configuration to check creation of service policy without next_hop_ip
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// terraform will try to create service policy without required argument peering_name
				Config:      CreateAccServicePolicyWithoutPeering_name(rName, ip), // configuration to check creation of service policy without peering_name
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// terraform will try to create service policy without required argument service_node_name
				Config:      CreateAccServicePolicyWithoutService_node_name(rName, ip), // configuration to check creation of service policy without service_node_name
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// terraform will try to create service policy without required argument source_network
				Config:      CreateAccServicePolicyWithoutSource_network(rName, ip), // configuration to check creation of service policy without source_network
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				// terraform will try to create service policy without required argument source_vrf_name
				Config:      CreateAccServicePolicyWithoutSource_vrf_name(rName, ip), // configuration to check creation of service policy without source_vrf_name
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip), // configuration to create ServicePolicy with required fields only
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_default), // this function will check whether any resource is exist or not in state file with given resource name
					// now will compare value of all attributes with default value for given resource
					resource.TestCheckResourceAttr(resourceName, "policy_template_name", "service_pbr"),
					resource.TestCheckResourceAttr(resourceName, "reverse_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "ip"),
					resource.TestCheckResourceAttr(resourceName, "src_port", "any"),
					resource.TestCheckResourceAttr(resourceName, "dest_port", "any"),
					resource.TestCheckResourceAttr(resourceName, "route_map_action", "permit"),
					resource.TestCheckResourceAttr(resourceName, "next_hop_action", "none"),
					resource.TestCheckResourceAttr(resourceName, "fwd_direction", "true"),
					resource.TestCheckResourceAttr(resourceName, "deploy", "false"),
					resource.TestCheckResourceAttr(resourceName, "deploy_timeout", "300"),
				),
			},
			{
				// this step will import state of particular resource and will test state file with configuration file
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// in this step all optional attribute are given for the same resource and then compared
				Config: CreateAccServicePolicyConfigWithOptionalValues(rName, ip), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					// now will compare value of all optional attributes with updated values from configuration
					resource.TestCheckResourceAttr(resourceName, "policy_template_name", "service_template"),
					resource.TestCheckResourceAttr(resourceName, "reverse_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "src_port", "3200"),
					resource.TestCheckResourceAttr(resourceName, "dest_port", "3300"),
					resource.TestCheckResourceAttr(resourceName, "route_map_action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "next_hop_action", "drop"),
					resource.TestCheckResourceAttr(resourceName, "fwd_direction", "false"),
					resource.TestCheckResourceAttr(resourceName, "deploy", "false"),
					resource.TestCheckResourceAttr(resourceName, "deploy_timeout", "200"),
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccServicePolicyWithInavalidIP(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccServicePolicyConfigUpdatedName(rName, ip, longerName), // passing invalid name for service policy
				ExpectError: regexp.MustCompile(fmt.Sprintf("property policy_name of sp-%s failed validation for value '%s'", longerName, longerName)),
			},
			{
				Config: CreateAccServicePolicyWithPolicyIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "policy_name", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyWithServicefabricIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "service_fabric", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyWithAttachedfabricIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "attached_fabric", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyWithDestnetworkIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "dest_network", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyWithDest_vrf_nameIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "dest_vrf_name", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyWithPeeringnameIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "peering_name", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyWithServicenodenameIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "service_node_name", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyWithSource_networkIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "source_network", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyConfig(rName, ip),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyWithSource_vrf_nameIPConfig(rName, rOtherName, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),
					resource.TestCheckResourceAttr(resourceName, "source_vrf_name", rOtherName),
					testAccCheckDCNMServicePolicyIdNotEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // checking whether id or dn of both resource are different because policyname changed and terraform need to create another resource
				),
			},
		},
	})
}

func TestAccDCNMPolicy_Update(t *testing.T) {
	var servicePolicy_default models.ServicePolicy
	var servicePolicy_updated models.ServicePolicy
	resourceName := "dcnm_service_policy.first"
	rName := acctest.RandString(5)
	//rOtherName := acctest.RandString(5)
	ip, _ := acctest.RandIpAddress("10.20.0.0/16")
	ip = fmt.Sprintf("%s/16", ip)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerServicePolicy),
		// CheckDestroy:      testAccCheckDCNMPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccServicePolicyConfig(rName, ip), // configuration to create ServicePolicy with required fields only
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_default), // this function will check whether any resource is exist or not in state file with given resource name
				),
			},
			{
				// this step will import state of particular resource and will test state file with configuration file
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "policy_template_name", "updated policy_template_name for terraform test"), // updating only policy_template_name parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                                               // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "policy_template_name", "updated policy_template_name for terraform test"), // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated),                      // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "reverse_enabled", "true"), // updating only reverse_enabled parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "reverse_enabled", "true"),                            // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "protocol", "tcp"), // updating only protocol parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),                                    // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "src_port", "2800"), // updating only src_port parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "src_port", "2800"),                                   // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "dest_port", "3800"), // updating only dest_port parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "dest_port", "3800"),                                  // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "route_map_action", "permit"), // updating only route_map_action parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "route_map_action", "permit"),                         // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "next_hop_action", "drop"), // updating only next_hop_action parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "next_hop_action", "drop"),                            // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "next_hop_action", "drop-on-fail"), // updating only next_hop_action parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "next_hop_action", "drop-on-fail"),                    // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "fwd_direction", "false"), // updating only fwd_direction parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "fwd_direction", "false"),                             // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "deploy", "true"), // updating only deploy parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "deploy", "true"),                                     // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccServicePolicyUpdatedAttr(rName, ip, "deploy_timeout", "10"), // updating only deploy_timeout parameter
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServicePolicyExists(resourceName, &servicePolicy_updated),                          // checking whether resource is exist or not in state file
					resource.TestCheckResourceAttr(resourceName, "deploy_timeout", "10"),                               // checking value updated value of description parameter
					testAccCheckDCNMServicePolicyIdEqual(resourceName, &servicePolicy_default, &servicePolicy_updated), // this function will check whether id or dn of both resource are same or not to make sure updation is performed on the same resource
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

/*
func TestAccApplicationProfile_NegativeCases(t *testing.T) {
	resourceName := "aci_application_profile.test"
	rName := acctest.RandString(5)
	longPolicy_template_name := acctest.RandString(129)                                     // creating random string of 129 characters
	longNameAlias := acctest.RandString(64)                                           // creating random string of 64 characters
	randomPrio := acctest.RandString(6)                                               // creating random string of 6 characters
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz") // creating random string of 5 characters (to give as random parameter)
	randomValue := acctest.RandString(5)                                              // creating random string of 5 characters (to give as random value of random parameter)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationProfileConfig(rName), // creating application profile with required arguements only
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccApplicationProfileWithInValidTenantDn(rName),                       // checking application profile creation with invalid tenant_dn value
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvAp (.)+`), // test step expect error which should be match with defined regex
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "description", longDescAnnotation), // checking application profile creation with invalid description value
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "annotation", longDescAnnotation), // checking application profile creation with invalid annotation value
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "name_alias", longNameAlias), // checking application profile creation with invalid name_alias value
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, "prio", randomPrio), // checking application profile creation with invalid prio value
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationProfileUpdatedAttr(rName, randomParameter, randomValue), // checking application profile creation with randomly created parameter and value
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccApplicationProfileConfig(rName), // creating application profile with required arguements only
			},
		},
	})
}
*/

func testAccCheckDCNMServicePolicyExists(name string, servicePolicy *models.ServicePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("ServicePolicy %s not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ServicePolicy dn was set")
		}

		dcnmClient := (*providerServicePolicy).Meta().(*client.Client)
		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticserviceapi/fabrics/testService/servicenodes/SN-1/Firewall/policies/Test_fabric_1/%s", rs.Primary.ID))
		log.Printf("[DEBUG] before err %s", cont)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] after err %s", cont)
		servicePolicyTest := &models.ServicePolicy{}
		servicePolicyTest.PolicyName = stripQuotes(cont.S("policyName").String())
		servicePolicyTest.FabricName = stripQuotes(cont.S("fabricName").String())
		servicePolicyTest.AttachedFabricName = stripQuotes(cont.S("attachedFabricName").String())
		servicePolicyTest.ServiceNodeName = stripQuotes(cont.S("serviceNodeName").String())
		servicePolicyTest.NextHopIp = stripQuotes(cont.S("nextHopIp").String())
		servicePolicyTest.PeeringName = stripQuotes(cont.S("peeringName").String())
		servicePolicyTest.DestinationNetwork = stripQuotes(cont.S("destinationNetwork").String())
		servicePolicyTest.DestinationVrfName = stripQuotes(cont.S("destinationVrfName").String())
		servicePolicyTest.SourceNetwork = stripQuotes(cont.S("sourceNetwork").String())
		servicePolicyTest.SourceVrfName = stripQuotes(cont.S("sourceVrfName").String())
		*servicePolicy = *servicePolicyTest
		return nil
	}
}

func CreateAccServicePolicyConfig(rName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rName, rName, ip, rName, rName, rName)
}

func CreateAccServicePolicyWithInavalidIP(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing with invalid IP")
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s0"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rName, rName, ip, rName, rName, rName)
}

func CreateAccServicePolicyConfigUpdatedName(rName, ip, longerName string) string {
	fmt.Println("=== STEP  Basic: testing servicePolicy creation with invalid name")
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, longerName, rName, rName, ip, rName, rName, rName)
}

func CreateAccServicePolicyWithPolicyIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rOtherName, rName, rName, ip, rName, rName, rName)
}

func CreateAccServicePolicyWithServicefabricIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "%s"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rOtherName, rName, rName, ip, rName, rName, rName)
}

func CreateAccServicePolicyWithAttachedfabricIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "%s"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rOtherName, rName, rName, ip, rName, rName, rName)
}

func CreateAccServicePolicyWithDestnetworkIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rOtherName, rName, ip, rName, rName, rName)
}

func CreateAccServicePolicyWithDest_vrf_nameIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rName, rOtherName, ip, rName, rName, rName)
}

func CreateAccServicePolicyWithPeeringnameIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rName, rName, ip, rOtherName, rName, rName)
}

func CreateAccServicePolicyWithServicenodenameIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "%s"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rName, rName, ip, rName, rOtherName, rName, rName)
}

func CreateAccServicePolicyWithSource_networkIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rName, rName, ip, rName, rOtherName, rName)
}

func CreateAccServicePolicyWithSource_vrf_nameIPConfig(rName, rOtherName, ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s" 
  }
  `, rName, rName, ip, rName, rName, rOtherName)
}

func CreateAccServicePolicyConfigWithOptionalValues(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing servicePolicy creation with optional parameters")
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s"
		
		policy_template_name 	 = "service_template"
		reverse_enabled 	 	 = true
		protocol			 	 = "tcp"
		src_port			 	 = "3200"
		dest_port			 	 = "3300"
		route_map_action	 	 = "deny"
		next_hop_action	 	 	 = "drop"
		fwd_direction		 	 = false
		deploy				 	 = false
		deploy_timeout		 	 = 200
  }
  `, rName, rName, ip, rName, rName, rName)
}

func testAccCheckDCNMPolicyDestroy(s *terraform.State) error {
	dcnmClient := (*providerServicePolicy).Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_policy" {
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/control/policies/%s", "test-demo-1"))
			if err == nil {
				return fmt.Errorf("Policy still exists!!")
			}
		}
	}

	return nil
}

func testAccCheckDCNMServicePolicyIdEqual(name string, sp1, sp2 *models.ServicePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("ServicePolicy %s not found", name)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No ServicePolicy dn was set")
		}
		rs2, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("ServicePolicy %s not found", name)
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No ServicePolicy dn was set")
		}

		if rs1.Primary.ID != rs2.Primary.ID {
			return fmt.Errorf("ServicePolicy ids are not equal")
		}
		return nil
	}
}

func testAccCheckDCNMServicePolicyIdNotEqual(name string, sp1, sp2 *models.ServicePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("ServicePolicy %s not found", name)
		}
		if rs1.Primary.ID == "" {
			return fmt.Errorf("No ServicePolicy dn was set")
		}
		rs2, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("ServicePolicy %s not found", name)
		}
		if rs2.Primary.ID == "" {
			return fmt.Errorf("No ServicePolicy dn was set")
		}

		if rs1.Primary.ID == rs2.Primary.ID {
			return fmt.Errorf("ServicePolicy id are equal")
		}
		return nil
	}
}

func CreateAccServicePolicyWithoutPolicy_name(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving policy_name")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
  		service_fabric           = "%s"
  		attached_fabric    		 = "%s"
  		dest_network             = "%s"
  		dest_vrf_name            = "%s"
  		next_hop_ip              = "%s"
  		peering_name             = "%s"
  		service_node_name        = "%s"
  		source_network           = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, ip, rName, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutService_fabric(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving service_fabric")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		attached_fabric    		 = "%s"
  		dest_network             = "%s"
  		dest_vrf_name            = "%s"
  		next_hop_ip              = "%s"
  		peering_name             = "%s"
  		service_node_name        = "%s"
  		source_network           = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, ip, rName, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutAttached_fabric(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving attached_fabric")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		service_fabric           = "%s"
  		dest_network             = "%s"
  		dest_vrf_name            = "%s"
  		next_hop_ip              = "%s"
  		peering_name             = "%s"
  		service_node_name        = "%s"
  		source_network           = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, ip, rName, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutDest_network(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving dest_network")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		service_fabric           = "%s"
  		attached_fabric    		 = "%s"
  		dest_vrf_name            = "%s"
  		next_hop_ip              = "%s"
  		peering_name             = "%s"
  		service_node_name        = "%s"
  		source_network           = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, ip, rName, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutDest_vrf_name(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving dest_vrf_name")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		service_fabric           = "%s"
  		attached_fabric    		 = "%s"
  		dest_network             = "%s"
  		next_hop_ip              = "%s"
  		peering_name             = "%s"
  		service_node_name        = "%s"
  		source_network           = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, ip, rName, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutNext_hop_ip(rName string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving next_hop_ip")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		service_fabric           = "%s"
  		attached_fabric    		 = "%s"
  		dest_network             = "%s"
  		dest_vrf_name            = "%s"
  		peering_name             = "%s"
  		service_node_name        = "%s"
  		source_network           = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, rName, rName, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutPeering_name(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving peering_name")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		service_fabric           = "%s"
  		attached_fabric    		 = "%s"
  		dest_network             = "%s"
  		dest_vrf_name            = "%s"
  		next_hop_ip              = "%s"
  		service_node_name        = "%s"
  		source_network           = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, rName, ip, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutService_node_name(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving service_node_name")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		service_fabric           = "%s"
  		attached_fabric    		 = "%s"
  		dest_network             = "%s"
  		dest_vrf_name            = "%s"
  		next_hop_ip              = "%s"
  		peering_name             = "%s"
  		source_network           = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, rName, ip, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutSource_network(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving source_network")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		service_fabric           = "%s"
  		attached_fabric    		 = "%s"
  		dest_network             = "%s"
  		dest_vrf_name            = "%s"
  		next_hop_ip              = "%s"
  		peering_name             = "%s"
  		service_node_name        = "%s"
  		source_vrf_name          = "%s"
	}
	`, rName, rName, rName, rName, rName, ip, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyWithoutSource_vrf_name(rName, ip string) string {
	fmt.Println("=== STEP  Basic: testing ServicePolicy creation without giving source_vrf_name")
	resource := fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "%s"
  		service_fabric           = "%s"
  		attached_fabric    		 = "%s"
  		dest_network             = "%s"
  		dest_vrf_name            = "%s"
  		next_hop_ip              = "%s"
  		peering_name             = "%s"
  		service_node_name        = "%s"
  		source_network           = "%s"
	}
	`, rName, rName, rName, rName, rName, ip, rName, rName, rName)
	return resource
}

func CreateAccServicePolicyUpdatedAttr(rName, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing attribute: %s=%s \n", attribute, value)
	return fmt.Sprintf(`
	resource "dcnm_service_policy" "first" {
		policy_name              = "SP-2"
		service_fabric           = "testService"
		attached_fabric    		 = "Test_fabric_1"
		dest_network             = "%s"
		dest_vrf_name            = "%s"
		next_hop_ip              = "%s"
		peering_name             = "%s"
		service_node_name        = "SN-1"
		source_network           = "%s"
		source_vrf_name          = "%s"
		
		%s 	 = "%s"
		reverse_enabled 	 	 = true
		protocol			 	 = "tcp"
		src_port			 	 = "3200"
		dest_port			 	 = "3300"
		route_map_action	 	 = "deny"
		next_hop_action	 	 	 = "drop"
		fwd_direction		 	 = false
		deploy				 	 = false
		deploy_timeout		 	 = 200
  }
  `, rName, rName, ip, rName, rName, rName, attribute, value)
}
