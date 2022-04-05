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

resource "dcnm_vrf" "vrf_lite" {
  fabric_name             = "{{enter_fabric_name}}"
  name                    = "{{enter_vrflite_name}}"
  description             = "VRF-Lite created by Terraform"
  intf_description        = "${var.vrf_lite.name}_Terraform_VRF-Lite"
  deploy                  = true
  attachments {
    serial_number = "{{enter_switch_serial_number}}"
    attach        = true
    vrf_lite {
      peer_vrf_name = "{{enter_vrflite_name}}"
      interface_name = var.vrf_lite.attachment_interface
      dot1q_id = "{{enter_vrflite_dot1q_number}}"
      ip_mask = "{{enter_ipv4_adress}}"
      neighbor_ip = "{{enter_vrflite_name_ipv4_address}}"
    }
  }
}
