---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_system_domains"
description: |-
Retrieves the Mashery ID of a named email template set.
---

# Data Source: mashery_email_template_set

Use this data source to query existing email template id by name that is already configured  in your account. The email
template set configures automatic email notifications for Mashery packages.

## Example Usage
```hcl
data mashery_email_template_set "default" {
  search = {
    name: "Default"
  }
}
```

Most Mashery-managed API programmes will define relatively small number of email notifications sets. These sets
are then re-used across multiple API Packages managed by multiple teams. A typical Terraform project would not
query the id by email template set name.

## Argument Reference
- `search`: key-value map specifying [search criteria](https://developer.mashery.com/docs/read/mashery_api/30/resources/emailtemplatesets), where
  the only meaningful parameter is `name`.
- `required`: boolean indicating whether at last one match must exist. Defaults to `true`
  
## Attribute Reference
- `id`: string, id of this email template set;
- `created`: date/time this template set was created;
- `updated`: date/time this template set was last updated;
- `name`: string, name of this email template set;
- `type`: Type of this template set.