---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_system_domains"
description: |-
Retrieves the Mashery Id of a given a role (aka portal access group) name.
---

# Data Source: mashery_role

This data source queries Mashery's for the details of the particular role. Mashery refers also to
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
    name: var.role_name
  }
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
- `read_permission`: a role permission object that can be passed to service I/O docs parameter (`iodocs_accessed_by`) and to 
   package plan visibility (`plan_visible_to` parameter).
