---
layout: "dcnm"
page_title: "DCNM: dcnm_policy"
sidebar_current: "docs-dcnm-resource-policy"
description: |-
  Manages DCNM policy modules
---

# dcnm_policy #
Manages DCNM policy modules

## Example Usage ##

```hcl
resource "dcnm_policy" "second" {
    serial_number   =   "9BH270169LJ" 
    template_name   =   "aaa_radius_deadtime"
    template_props  =   {
                            "DTIME" : "3"
                            "AAA_GROUP" : "management"
                        }
    priority        =   500  
    description     =   "This is demo policy."

}
```

## Common Argument Reference ##

* `serial_number` - (Required) Serial number of switch under which policy will be created.
* `template_name` - (Required)  A unique name identifying the template. Please note that a template name can be used by multiple policies and hence a template name does not identify a policy uniquely.
* `template_props` - (Required) Properties of the templates related to template name.
* `priority` - (Optional) Priority of the policy.Default value is 500.
* `description`- (Optional) Description of the policy. The description may include the details regarding the policy.Default value is "".



## Attribute Reference

*  `policy_id` - (Optional) A unique ID identifying a policy.
    NOTE: User can specify only empty string value.

## Importing ##

An existing policy can be [imported][docs-import] into this resource via its policy id using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import dcnm_policy.example <policyId>
```
