---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
Defines Mashery package
---

# Resource: `mashery_package`

The `mashery_package` represent a [`/packages/{packageId}](https://developer.mashery.com/docs/read/mashery_api/30/resources/packages)
V3 API resource. The package defines package plans and services.

## Example Usage

```hcl
resource "mashery_package" "lspwd_pack" {
  name_prefix = "x002924"
}
```

## Argument Reference

* `name`
* `name_prefix`
* `description`
* `tags`
* `notify_developer_period`
* `notify_developer_near_quota`
* `notify_developer_over_quota`
* `notify_developer_over_throttle`
* `notify_admin_period`
* `notify_admin_near_quota`
* `notify_admin_over_throttle`
* `notify_admin_emails`
* `near_quota_threshold`
* `extended_attribute_values`
* `key_adapter`
* `key_length`
* `shared_secret_length`

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `package_id`
* `created`
* `updated`