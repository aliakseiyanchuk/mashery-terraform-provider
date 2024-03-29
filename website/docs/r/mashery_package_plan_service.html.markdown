---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
    Defines Mashery service
---

# Resource: `mashery_package_plan_service`

The `mashery_package_plan_service` represents the [`/packages/{packageId}/plans/{planId}/services/{serviceId}`](https://developer.mashery.com/docs/read/mashery_api/30/resources/packages/plans/services)
V3 resource. It is used to link a service with the package plan.

## Example Usage

```hcl
resource mashery_package_plan_service "lspwds_Default_service" {
  package_plan_ref  = mashery_package_plan.lspwd_Default.id
  service_ref = mashery_service.lspwd2-first.id
}
```

## Argument Reference
* `package_plan_ref` package plan reference
* `service_ref` service reference
* 
## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id` Reference to the service included in the plan
