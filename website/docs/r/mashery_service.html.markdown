---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
  Defines Mashery service
---

# Resource: `mashery_service`

The resource represents [API definitions](https://developer.mashery.com/docs/read/mashery_api/30/resources/services)
configured using Mashery V3 API resource. The resource allows configuring:
- basic attributes;
- service cache
- service roles.

Service OAuth security profile is configured as a separate resource, `mashery_service_oauth`. Error sets
are configured via resource `mashery_service_error_set`.


## Example Usage

```hcl
resource mashery_service "demo-service" {
  name_prefix = "TFDemo"
}
```

## Argument Reference
The service accepts the following arguments:
- `name` or `name_prefix`: defines a (preferably unique) name of this service within the area
- `description`: service description that would be shown in Mashery Control Center
- `version`: service version as it would appear in Mashery Control Center  
- `qps_limit_overall`: maximum number of calls this service can handle. Defaults to -1 for unlimited
- `rfc3986_encode`: whether Mashery should use RFC 3986 specification for URL syntax
- `cache_ttl`: whether to set caching TTL for included endpoints. Defaults to 0
- `iodocs_accessed_by`: a set of roles permissions that are granted access to the I/O docs.
- `organization`: an id of organization this service belongs.



### Enabling caching
If your area supports caching, the caching is enabled for the service by `cache_ttl` key:
```hcl
resource mashery_service "demo-service" {
  name_prefix = "TFDemo"
  cache_ttl = 30
}
```

### Granting access to IODocs

To grant access to IODocs for this service, `iodocs_accessed_by` needs to be supplied a set of 
role (or portal access groups) whose members can use the IODocs of this service. The easiest way
to grant such permissions is to use `mashery_role` data source as the example below
illustrates:

```hcl
data "mashery_role" "my_role" {
  search = {
    name: "Internal Developer"
  }
}

resource mashery_service "demo-service" {
  name_prefix = "TFDemo"
  iodocs_accessed_by = toset([data.mashery_role.my_role.read_permission])
}
```

The `iodocs_accessed_by` is a set of objects comprising the following required keys:
- `id`: the ID of the role,
- `action`: string that should be set to `read`.


## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id`: service identified
* `created`: timestamp service was first deployed in Mashery
* `updated`: timestamp service was last updated in Mashery
* `editor_handle`: identity of the user that performed latest modification
* `revision_number`: counter, how many times a service has been deployed
* `robots_policy`: robots policy of this service
* `crossdomain_policy`: a crossdomain policy for APIs accessed via Flash