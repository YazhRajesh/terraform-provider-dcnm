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


resource "dcnm_interface" "vPCTerraform" {
  policy        = "int_vpc_trunk_host_11_1"
  type          = "vpc"
  name          = "vPC_001"
  fabric_name   = "{{enter_fabric_name}}"     #Enter Fabric Name
  switch_name_1 = "((enter_switch_name))"   #Enter Switch Name

  switch_name_2           = "93216TC-FX2-L2-S3"
  vpc_peer1_id            = "505"
  vpc_peer2_id            = "506"
  mode                    = "active"
  bpdu_guard_flag         = "true"
  mtu                     = "jumbo"
  vpc_peer1_allowed_vlans = "none"
  vpc_peer2_allowed_vlans = "none"
  vpc_peer1_access_vlans  = "10"
  vpc_peer2_access_vlans  = "20"
  vpc_peer1_interface     = ["e1/5", "eth1/7"]
  vpc_peer2_interface     = ["e1/5", "eth1/7"]

  deploy = false
}
