terraform {
  required_providers {
    dcnm = {
      source = "CiscoDevNet/dcnm"
    }
  }
}

provider "dcnm" {
  # Cisco DCNM/NDFC user name
  username = var.user.username  #Enter Username
  # Cisco DCNM/NDFC password
  password = var.user.password  #Enter Password
  # Cisco DCNM/NDFC url
  url      = var.user.url       #NDFC URL
  insecure = true
  # Used to select DCNM or NDFC for authentication purposes
  #nd for NDFC and dcnm for DCNM
  platform = "nd"
}

resource "dcnm_interface" "EthernetTerraform" {
  policy        = "int_trunk_host_11_1"
  type          = "ethernet"
  name          = "Ethernet1/14"
  fabric_name   = var.fabric      # Enter fabric_name
  switch_name_1 = "93216TC-FX2-L1-S3"

  ethernet_speed  = "Auto"
  bpdu_guard_flag = "no"
  allowed_vlans   = "none"
  mtu             = "jumbo"
  port_fast_flag  = true
}
