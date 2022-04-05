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


resource "dcnm_vrf" "vrf" {
  fabric_name             = "{{enter_fabric_name}}"
  name                    = "{{enter_vrf_name}}"
  vlan_id                 = 2002
  segment_id              = "50016"
  vlan_name               = "{{enter_vlan_name}}"
  description             = "vrf creation"
  intf_description        = "vrf"
  tag                     = "1250"
  max_bgp_path            = 2
  max_ibgp_path           = 4
  trm_enable              = false
  rp_external_flag        = false
  rp_address              = "1.1.1.2"
  loopback_id             = 15
  mutlicast_address       = "10.0.0.2"
  mutlicast_group         = "224.0.0.1/4"
  ipv6_link_local_flag    = "true"
  trm_bgw_msite_flag      = false
  advertise_host_route    = false
  advertise_default_route = "true"
  static_default_route    = false
  deploy                  = true
  attachments {
    serial_number = "{{enter_switch_serial_number}}"
    vlan_id       = 2300
    attach        = true
    loopback_id   = 70
    loopback_ipv4 = "1.1.1.3"
  }
}
