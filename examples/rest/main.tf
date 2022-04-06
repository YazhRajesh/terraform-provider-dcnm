terraform {
  required_providers {
    dcnm = {
      source = "CiscoDevNet/dcnm"
    }
  }
}

provider "dcnm" {
  # Cisco DCNM/NDFC user name
  username = "{{enter_username}}"
  # Cisco DCNM/NDFC password
  password = "{{enter_password}}"
  # Cisco DCNM/NDFC url
  url      = "{{enter_url_of_dcnm_or_ndfc}}"
  insecure = true
  # Used to select DCNM or NDFC for authentication purposes
  #nd for NDFC and dcnm for DCNM
  platform = "nd"
}

resource "dcnm_rest" "new_network" {
  #To add a new network using REST API for NDFC
  path    = "/appcenter/cisco/ndfc/api/v1/lan-fabric/rest/top-down/fabrics/fab1/networks"
  method  = "POST"
  payload = <<EOF

    {
      "displayName": "Terraform_Network",
      "fabric": "fab1",
      "hierarchicalKey": "fab1",
      "id": 2,
      "interfaceGroups": null,
      "networkExtensionTemplate": "Default_Network_Extension_Universal",
      "networkId": 30003,
      "networkName": "Terraform_Network",
      "networkTemplate": "Default_Network_Universal",
      "networkTemplateConfig": "{\"suppressArp\":\"false\",\"loopbackId\":\"\",\"enableL3OnBorder\":\"false\",\"SVI_NETFLOW_MONITOR\":\"\",\"enableIR\":\"false\",\"isLayer2Only\":\"false\",\"vrfDhcp3\":\"\",\"segmentId\":\"30003\",\"ENABLE_NETFLOW\":\"false\",\"dhcpServerAddr3\":\"\",\"vrfDhcp2\":\"\",\"dhcpServerAddr2\":\"\",\"tag\":\"12345\",\"dhcpServerAddr1\":\"\",\"nveId\":\"1\",\"vrfDhcp\":\"\",\"vlanId\":\"\",\"gatewayIpAddress\":\"172.30.102.1/24\",\"vlanName\":\"\",\"mtu\":\"\",\"intfDescription\":\"\",\"mcastGroup\":\"239.1.1.1\",\"trmEnabled\":\"\",\"VLAN_NETFLOW_MONITOR\":\"\",\"vrfName\":\"Terraform_VRF\"}",
      "serviceNetworkTemplate": null,
      "source": null,
      "tenantName": null,
      "vrf": "Terraform_VRF"
    }

  EOF
}
