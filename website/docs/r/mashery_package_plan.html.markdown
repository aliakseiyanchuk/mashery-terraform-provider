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

* `package_ref`: identifier of the `mashery_package`, to which this package belongs
* `name` plan name
* `description` plans description
* `extended_attribute_values` extended attribute values
* `self_service_provisioning` if developer can get key from the dev portal
* `admin_provisioning` if admin can generate key for the developer
* `notes` notes
* `max_keys` maximum number of keys a single user can hold on this package
* `keys_before_review` number of keys before requiring validation
* `qps` qps of this key
* `unlimited_qps` set true for unlimited qps
* `qps_override` set `true` to allow key to override the settings 
* `quota` number of calls the key can ame
* `unlimited_quota` set `true` to allow any number of calls
* `quota_override` allow the key to override the quota value set in the package
* `quota_period` a period of time over which the quota is calculated; possible values are:
  * `minute`
  * `hour`
  * `day`
  * `month`
* `response_filter_override` whether response filter can be overridden
* `status` status of this plan, either `active` or `inactive`. Defaults to `active`
* `developer_email_template_set` ID of the email template set to use to send notifications to developers
* `admin_email_template_set` ID of the email template set to use to send notifications to admins
* `portal_access_roles` portal access roles

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `plan_id`: compound identifier of this plan 
* `created` date when this package plan was created
* `updated` date when this package plan was last updated