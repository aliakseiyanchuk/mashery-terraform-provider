---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_system_domains"
description: |-
Outputs public domains that are associated with this account
---

# Data Source: mashery_public_domains

Use this data source to query the public domains that are assigned to your account. The public domain
represent the addresses that Mashery will be responding to.

## Example Usage

```hcl
data mashery_public_domains "myArea" {
}
```

## Attribute Reference
- `domains` - set of strings listing unique public domains assigned to your Mashery customer account.