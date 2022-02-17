---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
Defines Mashery service
---

# Resource: `mashery_package_plan` 

The `mashery_package_plan` represents the [`/packages/{packageId}/plans/{planId}`](https://developer.mashery.com/docs/read/mashery_api/30/resources/packages/plans)
resource of the Mashery V3 API.

## Example Usage

```hcl
resource "mashery_package_plan" "lspwd2_Default" {
  package_id = mashery_package.lspwd2_pack.id
  name = "Default"
}
```

## Argument Reference

* `package_id`: identifier of the `mashery_package`, to which this package belongs
* `name`
* `description`
* `extended_attribute_values`
* `self_service_provisioning`
* `admin_provisioning`
* `notes`
* `max_keys`
* `keys_before_review`
* `qps`
* `unlimited_qps`
* `qps_override`
* `quota`
* `unlimited_quota`
* `quota_override`
* `quota_period`
* `response_filter_override`
* `status`
* `email_template_set` ID of the email template set

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `plan_id`: compound identifier of this plan 
* `created`
* `updated`