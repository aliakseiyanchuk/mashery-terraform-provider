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

* `name` package name, mutually exclusive with `name_prefix`
* `name_prefix` name prefix of the package, mutually exclusive with `name`.
* `description` description
* `tags` tags
* `notify_developer_period`: developer notification period
  * `minute` 
  * `hour`
  * `day`
  * `week`
  * `month`
* `notify_developer_near_quota` set `true` to send notification to developer when key usage nears quota
* `notify_developer_over_quota` set `true` to send notification to developer when application exceeds quote. 
* `notify_developer_over_throttle` set `true` to send notification to developer when application exceeds throttle.
* `notify_admin_period` admin notification period
    * `minute`
    * `hour`
    * `day`
    * `week`
    * `month`
* `notify_admin_near_quota` set `true` to send notification to admin when key usage nears quota
* `notify_admin_over_throttle` set `true` to send notification to admin when key usage exceeds quota
* `notify_admin_emails` administrator emails
* `near_quota_threshold` near quota threshold percentage, e.g. 70
* `extended_attribute_values` extended attributes values
* `key_adapter` key adapter
* `key_length` key length
* `shared_secret_length` shared secret length
* `organization` the organization that owns this package

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `package_id` Mashery Id of the package
* `created` date when this package was created
* `updated` date when this pacakge was last updated