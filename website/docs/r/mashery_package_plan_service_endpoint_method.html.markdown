---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_endpoint_method"
description: |-
  Defines Mashery endpoint methos
---

# Resource: `mashery_package_plan_service_endpoint_method`

The `mashery_package_plan_service_endpoint_method` resource represents [`/packages/{packageId}/plans/{planId}/services/{serviceId}/endpoints/{endpointId}/methods/{methodId}`](https://developer.mashery.com/docs/read/mashery_api/30/resources/packages/plans/services/endpoints/methods)
V3 API resource that creates endpoint method

## Example Usage

```hcl
resource "mashery_package_plan_service_endpoint_method" "meth_abc" {
  package_plan_service_endpoint_ref = mashery_package_plan_service_endpoint.endp.id
  service_endpoint_method_ref = mashery_service_endpoint_method.meth_abc.id
  service_endpoint_method_filter_ref = mashery_service_endpoint_method_filter.abc_filter.id
}
```

## Argument Reference

- `package_plan_service_endpoint_ref` package plan service endpoint identifier
- `method_ref` method identifier; must be from the same service endpoint
- `service_filter_ref` optional method filter to apply to this method.

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id` Terraform object id
* `created` date when this method was created
* `updated` date when this method was last updated