---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_system_domains"
description: |-
Retrieves the Mashery Id of a given a role (aka portal access group) name.
---

# Data Source: mashery_role

This data source queries Mashery area for the details of the particular role. Mashery refers also to
these as "Portal Access Groups".

The primary purpose of the role is to supply access controls on Mashery developer portal. In practice,
the provider assumes that the roles will be set up once per program (or very infrequently)
e.g. by an API programme architect. After these have been established, the package plan configuration
would manage the configuration via these portal groups.

This data source allows retrieving Mashery Id of a role by role's name (as seen in Mashery console).

## Example usage

```hcl
data "mashery_role" "my_role" {
  search = {
    name: "Acme Department"
  }
}

# Specify IODocs access roles for the portal
resource "mashery_service" "srv" {
   name_prefix="tf-demo"
   iodocs_accessed_by = toset([data.mashery_role.my_role.id])
}

# Specify portal access role for a package plan (the parent package is omitted in this example)
resource "mashery_package_plan" "default" {
   package_ref = mashery_package.oauth.id
   name = "Default"
   portal_access_roles = toset([data.mashery_role.my_role.id])
}

```

## Argument Reference
The data source accepts the following arguments:
- `search`: map of search parameters on the [V3 Role fields](https://developer.mashery.com/docs/read/mashery_api/30/resources/roles)
- `required` whether matching role should be found. Defaults to `true`.

## Attribute Reference
The data source provides the following attributes:
- `id`: Mashery Id of this role;
- `created`: date when this role was created;
- `updated`: date when this role was last updated;
- `name`: role name;
- `description`: role's description;
- `predefined_role`: whether this role is pre-defined;
- `org_role`: whether this role is an org-role;
- `assignable_role`: whether this role is assignable.

