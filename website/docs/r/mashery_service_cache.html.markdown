---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service_cache"
description: |-
  Defines Mashery service cache
---

# Resource: `mashery_service_cache`

The resource represents [API caceh definitions](https://support.mashery.com/docs/read/mashery_api/30/resources/services/cache)
configured using Mashery V3 API resource. 

> Note: this resource is meaningful if caching is changed frequently without affecting
> the service definition. In case caching is changed infrequently, it may be 
> advantageous to declare the `cache_ttl` within the service definition itself.


## Example Usage

```hcl
resource mashery_service_cache "demo-service" {
  service_ref = mashery_service.terraform-service.id
  cache_ttl = "1h"
}
```

## Argument Reference
The service accepts the following arguments:
- `service_ref` service id to which this cache belongs
- `cache_ttl`: cache TTL, e.g. `1h`. Must be a valid Golang duration expression.


## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id`: service identified
* `created`: timestamp cache was first deployed in Mashery
* `updated`: timestamp cache was last updated in Mashery
