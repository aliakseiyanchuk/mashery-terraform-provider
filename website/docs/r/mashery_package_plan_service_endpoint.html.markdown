---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
  Defines Mashery service
---

# Resource: `mashery_plan_service_endpoint`

The `masher_plan_service_endpoint` represents [`/packages/{packageId}/plans/{planId}/services/{serviceId}/endpoints/{endpointId}`](https://developer.mashery.com/docs/read/mashery_api/30/resources/packages/plans/services/endpoints)
V3 API resource.

## Example Usage

```hcl
resource "mashery_package_plan_endpoint" "fff" {
  plan_service_id = mashery_package_plan_service.lspwd_Default_service.id
  endpoint_id = mashery_endpoint.lspwd2-enp-a.id
}
```

## Argument Reference

* `plan_service_ref` Package plan service id
* `endpoint_ref` endpoint id to be included. Must be from the same service as a referred by `plan_service_ref` attribute

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id` compound identifier of the included service 