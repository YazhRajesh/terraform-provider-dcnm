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

resource "dcnm_interface" "PortChannelTerraform" {
  policy        = "int_port_channel_access_host_11_1"
  type          = "port-channel"
  name          = "Po300"
  fabric_name   = "var.fabric"    #Enter fabric name
  switch_name_1 = var.switch_name #Enter Switch_Name

  mode            = "active"
  bpdu_guard_flag = "true"
  mtu             = "jumbo"
  allowed_vlans   = "none"
  access_vlans    = "10"
  pc_interface    = ["eth1/10", "eth1/12"]
}
