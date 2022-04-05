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

resource "dcnm_policy" "dcnm_policy" {
  serial_number = "{{enter_switch_serial_number}}"
  template_name = "aaa_radius_deadtime"   #Enter Template Name
  template_props = {
    "DTIME" : "3"
    "AAA_GROUP" : "management"
  }
  priority    = 500
  description = "This is demo policy."

}
