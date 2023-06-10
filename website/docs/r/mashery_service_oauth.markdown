---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
Defines Mashery service
---

# Resource: mashery_service_oauth

### OAuth Policy
The OAuth policy object represents configuration entered in Mashery [API Definition Security
Settings UI](http://docs.mashery.com/design/GUID-2D1DEABD-0630-41BA-807C-FD139B80482B.html).

```hcl
resource mashery_service "demo-service" {
  service_ref = "ey...=="

  access_token_ttl_enabled        = true
  access_token_ttl                = "1h"
  access_token_type               = "bearer"
  allow_multiple_token            = true
  authorization_code_ttl          = 300
  forwarded_headers               = toset(["access-token", "client-id", "scope", "user-context"])
  mashery_token_api_enabled       = false
  refresh_token_enabled           = true
  refresh_token_ttl               = "345h"
  enable_refresh_token_ttl        = true
  token_based_rate_limits_enabled = true
  force_oauth_redirect_url        = true
  force_ssl_redirect_url_enabled  = false
  grant_types                     = toset(["authorization-code", "implicit", "password", "client-credentials"])
  mac_algorithm                   = ""
  qps_limit_ceiling               = -1
  rate_limit_ceiling              = -1
  secure_tokens_enabled           = false
}
```