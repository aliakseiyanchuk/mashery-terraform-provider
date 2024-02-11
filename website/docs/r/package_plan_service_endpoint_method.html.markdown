---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
  Defines Mashery service
---

# Resource: `mashery_plan_service_endpoint_method`

The `masher_plan_service_endpoint` represents [`/packages/{packageId}/plans/{planId}/services/{serviceId}/endpoints/{endpointId}/methods/{methodId}`](https://developer.mashery.com/docs/read/mashery_api/30/resources/packages/plans/services/endpoints/methods)
V3 API resource. The resource also allows defining response filter corresponding to 
[`/packages/{packageId}/plans/{planId}/services/{serviceId}/endpoints/{endpointId}/methods/{methodId}/responseFilter`](https://developer.mashery.com/docs/read/mashery_api/30/resources/packages/plans/services/endpoints/methods/responsefilter)
Mashery V3 API resource.

## Example Usage

```hcl
resource "mashery_package_plan_endpoint" "fff" {
  plan_service_id = mashery_package_plan_service.lspwd_Default_service.id
  endpoint_id = mashery_endpoint.lspwd2-enp-a.id
}
```

## Argument Reference

* `package_plan_service_endpoint_ref`
* `service_endpoint_method_ref`
* `service_endpoint_method_filter_ref`

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id` compound identifier of the included service 