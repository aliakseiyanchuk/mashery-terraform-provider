---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_system_domains"
description: |-
Retrieves the Mashery Id of a given a role (aka portal access group) name.
---

# Data Source: mashery_role

This data source queries Mashery area for the details of the particular [organization](https://docs.mashery.com/manage/GUID-EAD30F7B-689D-4BC5-9B25-28CD6BD400A7.html). 

Organizations can be associated with services which controls the access to the data via the Mashery portal.

## Example usage

```hcl
data "mashery_organization" "tf_org" {
  search = {
    "name": "Terraform"
  }
}

resource "mashery_service" "srv" {
  name_prefix="tf-demo"
  organization = data.mashery_organization.tf_org.id
}

```

## Argument Reference
The data source accepts the following arguments:
- `search`: map of search parameters on the [V3 Organization fields](https://support.mashery.com/docs/read/mashery_api/30/resources/organizations)
- `required` whether matching role should be found. Defaults to `true`.

## Attribute Reference
The data source provides the following attributes:
- `id`: Mashery Id of this organization;
- `created`: date when this organization was created;
- `updated`: date when this organization was last updated;
- `name`: organization name.

