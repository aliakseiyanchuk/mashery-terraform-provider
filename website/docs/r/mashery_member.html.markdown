---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_member"
description: |-
Defines Mashery service
---

# Resource: `mashery_service`

The member resource represents `/members` V3 resource. 

Under normal operation, a member would register via a developer portal and, thus, would not require management
with Terraform. A member would only be required where a deployer wishes to issue keys to himself e.g. to
perform probing in production environment

## Example Usage

```hcl
resource "mashery_member" "myUser" {
  username_prefix = "lspwd2_terraform"
  email = "aliaksei.yanchuk@github.com"
  display_name = "Deployer"
}
```

## Argument Reference

* `username` a unique user name. Must be supplied, otherwise `username_prefix` must be speified
* `username_prefix` prefix for generating a unique user name
* `email` email address. Email cannot be changed via V3 API after it was created
* `area_status` status of the user in the area. Possible values are: `active`, `waiting`, or `disabled`
* `display_name` display name to show in the Mashery portal
* `uri`
* `blog`
* `im`
* `imsvc`
* `phone`
* `company`
* `address_1`
* `address_2`
* `locality`
* `region`
* `postal_code`
* `country_code`
* `first_name`
* `last_name`
* `external_id`

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id` 
* `created`
* `updated`