---
layout: "mashery"
page_title: "Provider: Mashery"
description: |-
The TIBCO Cloud Mashery provider is used to interact with the Mashery services and packages. The provider needs to be configured with the proper credentials before it can be used.
---

# TIBCO Cloud Mashery Terraform Provider

The TIBCO Cloud Mashery provider is used to create and manage  Mashery services and packages. 

## Example Usage

Terraform 0.13 and later:
```hcl
terraform {
  required_providers {
    mashery = {
      version = "0.4"
      source = "yanchuk.nl/aliakseiyanchuk/mashery"
    }
  }
}

provider "mashery" {
}
```

## Provider Authentication

The provider needs to authenticate itself to Mashery V3 API, which is achieved by obtaining an access token
for Mashery V3 API. This can be achieved in two ways:

1. Obtaining a V3 access token and passing this to the provider, or
2. Configuring the provider to interact with Mashery via [Vault secrets engine](https://github.com/aliakseiyanchuk/hcvault-mashery-api-auth)
    This is the preferred way of interacting with Mashery, and the only method that can support a long-running build.

### Authentication provider with an explicit V3 Access Token
To authenticate a provider with an explicit V3 access token, follow the steps described in the 
[authentication section](https://support.mashery.com/docs/read/mashery_api/30/Authentication) of the Mashery V3 API Guide.
Then the obtained token should be set to `TF_MASHERY_V3_ACCESS_TOKEN` environment variable, e.g.
```shell
export TF_MASHERY_V3_ACCESS_TOKEN=<obtained_token>
```
No further configuration is required.
> The simplicity of this method has a downside that the token time cannot be controlled. All deployment operations
> **must** complete within the validity period of this token. The moment the access token expires, the calls to
> Mashery V3 API will be rejected, and the deployment will fail.
> 
> The solution to this problem requires deploying HashiCorp Vault with [Mashery secrets engine](https://github.com/aliakseiyanchuk/hcvault-mashery-api-auth).

### Authenticating provider via Vault Mashery secrets engine

The minimal required configuration to configure provider to work via Mashery secrets engine requires identifying
the role of the credentials. A secret engine role captures the logon credentials for a particular Mashery V3 area.

```terraform
terraform {
  required_providers {
    mashery = {
      version = "0.4"
      source = "github.com/aliakseiyanchuk/mashery"
    }
  }
}

provider "mashery" {
  role = "<desired role name>"
}
```
The plugin also requires environment variables `VAULT_ADDR` (used by `vault` CLI tool) and `TF_MASHERY_VAULT_TOKEN`
specifying the Vault access token.

Given highly flexible nature of Vault, describing Vault setup and authentication is beyond the scope of this guide.
To illustrate a simplistic case (not to be used for production or shared environments), the deployer can authenticate
to a Vault using user pass method. Then the following sequence of steps will prepare the terminal session:

```shell
$ vault vault login -method userpass username=tfmash-dev
```
If successful, this command will print the access token:

```text
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                    Value
---                    -----
token                  hvs.CAESI<lots of letters>Gk
token_accessor         isHyWHsM7lE6c8H2DBW3OaQM
token_duration         12h
token_renewable        true
token_policies         ["default" "mash-auth"]
identity_policies      []
policies               ["default" "mash-auth"]
token_meta_username    tfmash-dev
```

The, the value of the token should be exported to `TF_MASHERY_VAULT_TOKEN` environment variable using e.g.
```shell
export TF_MASHERY_VAULT_TOKEN=hvs.CAESI<lots of letters>Gk
```
> Running this engine via Vault allows the deployer support sessions that exceed the V3 token time-to-live of 1 hour.
> As long as the Vault token is alive, Vault secret engine will manage the V3 access token automatically.
> 
> The remaining time of Vault access token can be inspected with `vault token lookup` command. The TTL can be 
> extended (up to the max limits) with `vault token renew`
>
> Remember to revoke your Vault access token using `vault token revoke` when your interaction has been completed.

## Supported parameters
The provider accepts the following options:
- `log_file`: a path  *prefix* where to store log files from individual provider runs. This is an optional
   parameter that you may wish to set to troubleshoot a problem of the provider translating terraform configuration
   into the V3 API calls.
- `v3_token`: an active Mashery V3 access token, e.g. retrieved from HashiCorp vault. It can also be passed as `TF_MASHERY_V3_ACCESS_TOKEN`
  environment variable;
- `qps`: number of calls per second to make towards Mashery API. For practical reasons, it should be set value slightly less 
  than Mashery-specified highest QPS. For example, where Mashery settings allow for 2 calls per second, the provider is
  recommended to throttle calls as 1 per second. It can also be passed as `TF_MASHERY_QPS` environment variable;
- `network_latency`: average network latency between your location and Mashery V3 API gateway. It can also be passed
  as `TF_MASHERY_V3_NETWORK_LATENCY` environment variable
- 'vault_addr': address of the Vault server. Can be passed using `VAULT_ADDR` environment variable
- `vault_mount`: specifies the mounting path of Mashery secrets engine within the Vault server. A default value is
  `mash-auth`. The value can also be passed using `TF_MASHERY_VAULT_MOUNT` environment variable.
- `role`: specifies the role (a wrapper over Mashery area and access credentials) within the Vault secret engine 
  the provider should use to interact with Mashery. Can also be passed using `TF_MASHERY_VAULT_ROLE` environment variable.
- `vault_token`: specifies the Vault token to use. Given sensitive and transient nature of this value, it is highly
  recommended to have this value passed with `TF_MASHERY_VAULT_TOKEN` environment variable.


## Configuring multiple Mashery Areas in the same Terraform project

Due to limitations of Mashery V3 API, single provider instance can make changes in a single Mashery area.
Should your project require making changes across multiple Masherey areas, multiple
providers need to be created using the provider alias as [this terraform documentation](https://www.terraform.io/docs/language/providers/configuration.html)
explains. 


