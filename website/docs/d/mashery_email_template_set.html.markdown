---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_system_domains"
description: |-
Retrieves the Mashery ID of a named email template set.
---

# Data Source: mashery_email_template_set

Use this data source to query existing email template id by name that is already configured  in your account. The email
template set configures automatic email notifications for Mashery package plans.

Most Mashery-managed API programmes will define relatively small number of email notifications sets. These sets
are then re-used across multiple API packages managed by multiple teams. A typical Terraform project would
query the email template set by its well-known name.

## Example Usage
```hcl
# Query Mashery IDs of existing email template sets using a well-known name
data "mashery_email_template_set" "admin_email_template" {
  search = {
    name : "Terraform X1"
  }
}

data "mashery_email_template_set" "user_email_template" {
  search = {
    name : "Terraform X2"
  }
}


resource "mashery_package_plan" "default" {
  package_ref = mashery_package.oauth.id
  name = "Default"
  admin_provisioning = true

  # Specify email templates to be used for the desired plan
  email_template_set = data.mashery_email_template_set.user_email_template.id
  admin_email_template_set = data.mashery_email_template_set.admin_email_template.id
}
```

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