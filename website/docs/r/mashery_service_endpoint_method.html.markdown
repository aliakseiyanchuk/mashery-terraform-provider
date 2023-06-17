---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_endpoint_method"
description: -|
  Defines an endpoint method within Mashery service endpoint
---

# Resource: `mashery_service_endpoint_method`

An endpoint method object representing [`/services/{serviceId}/endpoints/{endpointId}/method/{methdoId}`](https://support.mashery.com/docs/read/mashery_api/30/resources/services/endpoints/methods)
Mashery V3 API resource.

## Example Usage

```hcl
resource "mashery_service_endpoint_method" "meth_abc" {
  service_endpoint_ref = mashery_service_endpoint.endp.id
  name = "do something good"
  sample_json = file("${path.module}/meth_abc.json")
}
```

## Argument Reference

* `endpoint_ref` endpoint identifier, to which this method belongs
* `name` method name
* `sample_json` sample json payload; necessary to define a method filter
* `sample_xml` sample xml payload; necessary to define a method filter



## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `method_id` Mashery method id
* `created` date method was created
* `updated` date method was updated