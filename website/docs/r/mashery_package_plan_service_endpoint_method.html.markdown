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
resource "mashery_endpoint_method" "my_method" {
  endpoint_ref = mashery_endpoint.lspwd2-enp-a.id
  name = "lspwd2 method"
  sample_json = file("./my_method_sample.json")
}
```

## Argument Reference

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id`
* `created`
* `updated`