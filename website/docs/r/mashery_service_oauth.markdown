---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
  Defines Mashery service
---

# Resource: `mashery_service_oauth`

### OAuth Policy
The OAuth policy object represents configuration entered in Mashery [API Definition Security
Settings UI](http://docs.mashery.com/design/GUID-2D1DEABD-0630-41BA-807C-FD139B80482B.html).

## Example Usage
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

## Argument Reference

The resource requires `service_ref` to identify the service to which this OAuth
policy is applied.

The optional parameters configure the desired security settings:
- `access_token_ttl_enabled` whether access token TTL is enabled
- `access_token_ttl` the token TTl, e.g. `1h`. Must be a valid golang duration expression
- `access_token_type` access token type, either `bearer` or `mac`
- `allow_multiple_token` whether an application can request multiple tokens
- `authorization_code_ttl` authorization code TTL, e.g. `5m`. Must be valid Golang duration expression
`forwarded_headers` a set of forwarded OAuth headers to back-end, options are:
  - `access-token`
  - `client-id`
  - `scope`
  - `user-context`
- `mashery_token_api_enabled` whther Mashery's own client-credentials endpoint (`/token`) should be added to this service
- `refresh_token_enabled` whether access tokens should receive a refresh token
- `refresh_token_ttl` a refresh token TTL, e.g. `345h`. Must be a valid Golang duration expression"
- `enable_refresh_token_ttl` whether refresh tokens should receive TTL
- `token_based_rate_limits_enabled` rate limits are specific to token, not to all applications
- `force_oauth_redirect_url` verify redirect URL in the token request
- `force_ssl_redirect_url_enabled` the redirect URL must contain https
- `grant_types` grant types, option include:
  - `authorization-code`
  - `implicit`
  - `password`
  - `client-credentials`
- `mac_algorithm` a MAC algorithm for MAC tokens, either  `hmac-sha-1` or `hmac-sha-256`
- `qps_limit_ceiling` QPS limit ceiling
- `rate_limit_ceiling` rate limit ceiling
- `secure_tokens_enabled` whether secure tokens are enabled

# Attribute Reference
In addition to all arguments above, the following attributes are exposed:

* `id`: compound id of this error set.
* `created`: date this OAuth security profile was created;
* `updated`: date this OAuth security profile was last updated.