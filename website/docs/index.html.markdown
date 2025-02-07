---
layout: "dcnm"
page_title: "Provider: DCNM"
sidebar_current: "docs-dcnm-index"
description: |-
  The Cisco DCNM provider is used to interact with the resources provided by Cisco DCNM/NDFC.
  The provider needs to be configured with the proper credentials before it can be used.
---
  

Overview
--------------------------------------------------
Terraform provider DCNM is a Terraform plugin which will be used to manage the DCNM/NDFC Constructs on the Cisco DCNM/NDFC platform with leveraging advantages of Terraform. DCNM Terraform provider lets users represent the infrastructure as a code and provides a way to enforce state on the infrastructure managed by the Terraform provider. Customers can use this provider to integrate the Terraform configuration with the DevOps pipeline to manage the DCNM/NDFC fabric policies in a more flexible, consistent and reliable way.

Cisco DCNM Provider
------------
The Cisco DCNM provider is used to interact with the resources provided by Cisco DCNM and Cisco NDFC. The provider needs to be configured with the proper credentials before it can be used.

Authentication
-------------- 

Authentication with user-id and password.  
 example:  

 ```hcl

terraform {
  required_providers {
    dcnm = {
      source = "CiscoDevNet/dcnm"
    }
  }
}

provider "dcnm" {
  # cisco-dcnm/ndfc user name
  username = "admin"
  # cisco-dcnm/ndfc password
  password = "password"
  # cisco-dcnm/ndfc url
  url      = "https://my-cisco-dcnm.com"
  insecure = true
  platform = "dcnm"
}

 ```

Example Usage
------------
```hcl

terraform {
  required_providers {
    dcnm = {
      source = "CiscoDevNet/dcnm"
    }
  }
}

#configure provider with your cisco dcnm/ndfc credentials.
provider "dcnm" {
  # cisco-dcnm/ndfc user name
  username = "admin"
  # cisco-dcnm/ndfc password
  password = "password"
  # cisco-dcnm/ndfc url
  url      = "https://my-cisco-dcnm.com"
  insecure = true
  platform = "dcnm"
}

resource "dcnm_vrf" "test-vrf" {
  fabric_name = "fab1"
  name = "MyVRF"
  description = "This vrf is created by terraform"
}

```

Argument Reference
------------------
Following arguments are supported with Cisco DCNM terraform provider.

 * `username` - (Required) This is the Cisco DCNM/NDFC username, which is required to authenticate with CISCO DCNM/NDFC.
 * `password` - (Required) Password of the user mentioned in username argument. It is required when you want to use token-based authentication.
 * `url` - (Required) URL for CISCO DCNM/NDFC.
 * `insecure` - (Optional) This determines whether to use insecure HTTP connection or not. Default value is `true`.
 * `platform` - (Optional) NDFC/DCNM Platform information(Nexus-Dashboard/DCNM). Allowed values are "nd" and "dcnm". Default value is "dcnm".