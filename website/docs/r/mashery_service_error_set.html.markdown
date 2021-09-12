---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
Defines a custom error set for a given API service definition
---

# Resource: mashery_error_set

Mashery resource for defining a custom [error set](http://docs.mashery.com/design/GUID-0CD66EB0-A62E-4834-8C8A-FB1F72B3D4CB.html) 
associated with a service. 

Once an error set is created, a default set of message will be pre-populated. The provider allows overriding default
messages as desired. The drift is detected (and corrected) only for messages that are specified within this provider.

## Example Usage

```hcl
resource mashery_error_set "svc_a_tfset" {
  service_id = mashery_service.srv.id
  name = "Terraform Set"

  error_message {
    id = "ERR_400_BAD_REQUEST",
    status = "Bad request",
    detail_header = "x-error-detail",
  }
  error_message {
    // Repeat more error messages as required; or use toset() terraform function
    // ...
  }
}

## Use the error set later on in the endpoint
resource mashery_endpoint "custom_error" {
  service_id = mashery_service.srv.id
  # other configuration options
  
  # Specify the error set to use on this endpoint
  error_set = mashery_error_set.svc_a_tfset.id
}
```

## Argument Reference
The resource requires the following arguments to be provided:
- `service_id`: id of the service containing this error set;
- `name`: name of this set.
  
The resource accepts optional arguments:
- `jsonp`: whether this error set concerts jsonp responses. Defaults to `false`;
- `jsonp_type`: type of the jsonp responses;
- `type`: mime-type of the error set
- `message`: a set of error messages. Each error message should have the following fields
  (while reseble [error mesage setttings](http://docs.mashery.com/design/GUID-0CD66EB0-A62E-4834-8C8A-FB1F72B3D4CB.html) 
  as defiend by Mashery):
    - `id`: the message id, must be one fo the acceptable constants;
    - `status`: the HTTP Status Code Message for the error;
    - `detail_header`: information to be placed in the "X-Error-Detail-Header" Response Header for the error. ;
    - `response_body`: information to be placed in the Response Body for the error. .

### Acceptable Error Message IDs
The error message ids must be withing the following set:
- `ERR_400_BAD_REQUEST`
- `ERR_403_NOT_AUTHORIZED`
- `ERR_403_DEVELOPER_INACTIVE`
- `ERR_403_DEVELOPER_OVER_QPS`
- `ERR_403_DEVELOPER_OVER_RATE`
- `ERR_403_DEVELOPER_UNKNOWN_REFERER`
- `ERR_403_SERVICE_OVER_QPS`
- `ERR_403_SERVICE_REQUIRES_SSL`
- `ERR_414_REQUEST_URI_TOO_LONG`
- `ERR_502_BAD_GATEWAY`
- `ERR_503_SERVICE_UNAVAILABLE`
- `ERR_504_GATEWAY_TIMEOUT`
- `ERR_400_UNSUPPORTED_PARAMETER`
- `ERR_400_UNSUPPORTED_SIGNATURE_METHOD`
- `ERR_400_MISSING_REQUIRED_CONSUMER_KEY`
- `ERR_400_MISSING_REQUIRED_REQUEST_TOKEN`
- `ERR_400_MISSING_REQUIRED_ACCESS_TOKEN`
- `ERR_400_DUPLICATED_OAUTH_PROTOCOL_PARAMETER`
- `ERR_401_TIMESTAMP_IS_INVALID`
- `ERR_401_INVALID_SIGNATURE` 
- `ERR_401_INVALID_OR_EXPIRED_TOKEN`
- `ERR_401_INVALID_CONSUMER_KEY`
-  `ERR_401_INVALID_NONCE`

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `id`: compound id of this error set. The id has `<service_id>::<error-set-uuid>` format;
* `created`: date this error set was first created;
* `updated`: date this error set was last updated.