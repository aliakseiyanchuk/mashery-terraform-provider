---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_system_domains"
description: |-
Outputs system domains that are associated with this account
---

# Data Source: mashery_system_domains

Use this data source to query the system domains that are assigned to your account. The system domains
are the domains to which the Mashery traffic manager can forward API calls.

## Example Usage

```hcl
data mashery_public_domains "myArea" {
}
```

## Attribute Reference
- `domains` - set of strings listing unique system domains assigned to your Mashery customer account.