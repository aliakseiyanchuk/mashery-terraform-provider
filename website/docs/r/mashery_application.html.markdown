---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
Defines Mashery service
---

# Resource: `mashery_application`
 
The resource represents [`/applications`](https://developer.mashery.com/docs/read/mashery_api/30/resources/applications)
V3 API resource which declares an application. An application is owned by a speicfic member and contains
API keys.

## Example Usage

```hcl
resource "mashery_application" "myApp" {
  owner = mashery_member.myUser.id
  name_prefix = "myApp"
}

```

## Argument Reference

* `name` application name, or
* `name_prefix` name prefix to generate the application name
* `owner` the user that will own
* `description`
* `type`
* `commercial`
* `ads`
* `ads_system`
* `usage_model`
* `tags`
* `notes`
* `how_did_you_hear`
* `preferred_protocol`
* `external_id`
* `uri`
* `oauth_redirect_uri`
* `eav`: the EAV

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id`
* `created`
* `updated`
* `owner_username`