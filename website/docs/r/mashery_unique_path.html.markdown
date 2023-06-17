---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
  Creates a uniquely prefixed path
---

# Resource: `mashery_unique_path`

This resource is used to generate a uniquely-prefixed path, allowing creating non-intersecting
copies of the API definition services.

Within Mashery, each endpoint serving traffic should be bound to a given public domain name, 
request path alias, and process set of http verbs. These settings should be unique within the area.

During active development, it is not unusual that a product team would like to test various API
implementation version. Where only a single copy of an API service definition exists, this creates
a contention between various requests.

The unique path resource offers a solution to this problem. It allows writing Terraform specification
in a way that allows creating non-overlapping (even temporary) of the APIs by adding 
the timestamp-based suffixes to the request path alias prefix.

## Example Usage

```hcl
resource "mashery_unique_path" "ctx" {
  prefix = "test/featureApi"
}

locals {
  # A randomly generated prefix will be copied in to local value
  ctx = mashery_unique_path.ctx.path
}

resource "mashery_endpoint" "demo-endpoint" {
  # ... other endpoint settings ...
  # Use the local value to derive the prefix of this endpoint 
  request_path_alias = "${local.ctx}/a"
  #  ... other endpoint settings ...
}

# Output the value that an API developer could use to test the API calls.
output "mounting_context" {
  value = local.ctx
}
```

The value set to `locals.ctx` could be `/test/featureApi_sXEbVb`, where `sXEbVb` is a timestamp-based
suffix.

A more sophisticated usage could include variables to control the request path alias prefix:

```hcl
# There variables define:
# - A "production" context to deploy;
# - Whether custom prefix path should be used instead of production path.
# Both settings are easily set within the variables.tf file.

variable "ctx" {
  type = string
  default = "/production"
}

variable "usePrefixPath" {
  type = bool
  default = false
}

resource "mashery_unique_path" "ctx" {
  prefix = "testApi"
}

locals {
  ctx = var.usePrefixPath ? mashery_unique_path.ctx.path : var.ctx
}

resource "mashery_service_endpoint" "demo-endpoint" {
  # ... other endpoint settings ...
  # Use the va 
  request_path_alias = "${local.ctx}/a"
  #  ... other endpoint settings ...
}
```

## Argument Reference
The resources require a single argument, `prefix`. Leading slash is not required as it will be
automatically added to the output `path` attribute.

## Attribute Reference
The resource outputs a single attribute, `path`, which is formatted as `/{prefix}_{suffix}`
