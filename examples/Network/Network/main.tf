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

resource "dcnm_network" "overlay_network" {
  fabric_name     = "{{enter_fabric_name}}"     #Enter Fabric Name
  name            = "{{enter_network_name}}"
  network_id      = 30002
  description     = "${var.overlay_network.network_name}_Terraform"
  vrf_name        = "{{enter_vrf_name}}"
  vlan_id         = "{{enter_vlan_id}}"
  vlan_name       = "{{enter_vlan_name}}"
  ipv4_gateway    = "{{enter_ip4_gateway}}"
  deploy = true
  attachments {
    serial_number = "{{enter_serial_number_of_switch}}"
    vlan_id       = "{{enter_vlan_id}}"
    attach        = true
    switch_ports = ["Ethernet1/11"]
  }
  attachments {
    serial_number = "{{enter_serial_number_of_switch}}"
    vlan_id       = "{{enter_vlan_id}}"
    attach        = true
    switch_ports = ["Ethernet1/11"]
  }
}
