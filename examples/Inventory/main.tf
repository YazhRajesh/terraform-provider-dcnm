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

resource "dcnm_inventory" "add_new_switch" {
  fabric_name     = "{{enter_fabric_name}}"
  username        = "{{enter_switch_username}}"
  password        = "{{enter_switch_password}}"
  preserve_config = "false"
  max_hops = 0
  config_timeout  = 30          #Takes a long time to load switches -increase acc
  #Add Switch Details below
  switch_config {
    ip   = "172.31.217.102"
    role = "leaf"
  }
  switch_config {
    ip   = "172.31.217.122"
    role = "leaf"
  }
  switch_config {
    ip   = "172.31.219.123"
    role = "leaf"
  }
  switch_config {
    ip   = "172.31.217.9"
    role = "leaf"
  }
  switch_config {
    ip   = "172.31.186.154"
    role = "spine"
  }
  switch_config {
    ip   = "172.31.186.153"
    role = "spine"
  }

}
