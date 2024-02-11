---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: -|
  Defines an filter within Mashery service endpoint method
---

# Resource: `mashery_service_endpoint_method_filter`

An endpoint method object representing [`/services/{serviceId}/endpoints/{endpointId}/method/{methdoId}/filter/{filterId}`](https://support.mashery.com/docs/read/mashery_api/30/resources/services/endpoints/methods/responsefilters)
Mashery V3 API resource.

## Example Usage

```hcl
resource "mashery_service_endpoint_method_filter" "abc_filter" {
  service_endpoint_method_ref = mashery_service_endpoint_method.meth_abc.id
  name = "abc filter"
  json_fields = "/a"
}
```

## Argument Reference

* `service_endpoint_method_ref` endpoint identifier, to which this method belongs
* `name` method name
* `nodes` filter notes
* `xml_fields` xml fields to include in this filter
* `json_fields` json fields to include in this filter


## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `service_endpoint_method_filter_id` Mashery method id
* `created` date method was created
* `updated` date method was updated