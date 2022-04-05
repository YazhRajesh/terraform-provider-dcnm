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

resource "dcnm_interface" "LoopBackTerraform" {
  fabric_name = "fabric3"    #Enter fabric name here
  name        = "loopback5"
  type        = "loopback"
  policy      = "int_loopback_11_1"

  switch_name_1             = "93216TC-FX2-L1-S3"
  ipv4                      = "1.1.1.7"
  loopback_tag              = "1234"
  vrf                       = "BLUE-TF"
  loopback_ls_routing       = "ospf"
  loopback_routing_tag      = "1234"
  loopback_router_id        = "10"
  loopback_replication_mode = "Multicast"
  description               = "creation from terraform"
  ipv6                      = "2001::0"
}
