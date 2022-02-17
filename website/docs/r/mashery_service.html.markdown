---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
Defines Mashery service
---

# Resource: mashery_service

The resource represents [`/services/{serviceId}`](https://developer.mashery.com/docs/read/mashery_api/30/resources/services)
Mashery V3 API resource. The resource allows configuring:
- basic attributes;
- service cache (`/services/{serviceId}/cache`)
- security profile (`/services/{serviceId}/securityProfile/oauth`)
- service roles (`/services/{serviceId}/roles`)

The error sets (`/services/{serviceId}/errorSets`) and (`/services/{serviceId}/endpoints/{endpointId}`)
are configured as a separate resources.

The service (called API Definition in Mashery control center UI) is the starting point for exposing
an API. The service is effectively a collection of endpoints. 

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
- `oauth`: OAuth profile settings for this service. If absent, OAuth support is not enabled.
- `iodocs_accessed_by`: a set of roles permissions that are granted access to the I/O docs.

### OAuth Policy
The OAuth policy object represents configuration entered in Mashery [API Definition Security
Settings UI](http://docs.mashery.com/design/GUID-2D1DEABD-0630-41BA-807C-FD139B80482B.html).

```hcl
resource mashery_service "demo-service" {
  name_prefix = "TFDemo"

  oauth {
    access_token_ttl_enabled = true
    access_token_ttl = "1h"
    access_token_type = "bearer"
    allow_multiple_token = true
    authorization_code_ttl = 300
    forwarded_headers = toset(["access-token", "client-id", "scope", "user-context"])
    mashery_token_api_enabled = false
    refresh_token_enabled = true
    refresh_token_ttl = "345h"
    enable_refresh_token_ttl = true
    token_based_rate_limits_enabled = true
    force_oauth_redirect_url = true
    force_ssl_redirect_url_enabled = false
    grant_types = toset(["authorization-code", "implicit", "password", "client-credentials"])
    mac_algorithm = ""
    qps_limit_ceiling = -1
    rate_limit_ceiling = -1
    secure_tokens_enabled= false
  }
}
```

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