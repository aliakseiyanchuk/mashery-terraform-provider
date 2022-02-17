---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
Defines Mashery service
---

# Resource: `mashery_package_key`

The package key resrouce represents [`/applications/{applicationId}/packageKeys`](https://developer.mashery.com/docs/read/mashery_api/30/resources/packagekeys)
V3 resource.

Within Mashery, the (package) keys are requested by the developers and do not need to be managed with 
Terraform. This resource allows the deployer creating a package key for its own use, e.g. to perform
probing in production.

## Example Usage

```hcl
resource "mashery_package_key" "my-key" {
  plan_id = mashery_package_plan.lspwd_Default.id
  application_id = mashery_application.myApp.id

  // Cannot update secret.
  secret = "enterYourSecret"
  quota = mashery_package_plan.lspwd_Default.quota
  qps = mashery_package_plan.lspwd_Default.qps
}
```

## Argument Reference

* `plan_id` ID of the package plan
* `application_id` owner application id
* `secret` secret value of the package key
* `quota` quota assigned to this key
* `quota_exempt` whether to exempt any quota from this ke
* `qos` qps assigned to this package key
* `qos_exempt` whether ot exempt this package key from the QPS controls
* `stauts` package key status: `waiting`, `active`, `disabled`
* `limits` limits for quota and qps for this package key

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id`
* `created`
* `updated`
* `api_key` the API key