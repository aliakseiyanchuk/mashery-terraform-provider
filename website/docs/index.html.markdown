---
layout: "mashery"
page_title: "Provider: Mashery"
description: |-
The TIBCO Cloud Mashery provider is used to interact with the Mashery services and packages. The provider needs to be configured with the proper credentials before it can be used.
---

# TIBCO Cloud Mashery Terraform Provider

The TIBCO Cloud Mashery provider is used to interact with the Mashery services and packages. The provider needs to be
configured with the proper credentials before it can be used.

## Example Usage

Terraform 0.13 and later:
```hcl
terraform {
  required_providers {
    mashery = {
      version = "0.0.1"
      source = "yanchuk.nl/aliakseiyanchuk/mashery"
    }
  }
}

provider "mashery" {
}
```

The provider needs to authenticate itself to Mashery V3 API. Method vary depending on the deployment topology. Consult
various authentication options explained below.

Where no explicit configuration is supplied, the provider will attempt to infer the authentication from:
1. Environment variable `TF_MASHERY_V3_ACCESS_TOKEN`. If not set, 
2. Environment variable `TF_MASHERY_V3_TOKEN_FILE` indicating file containing saved access token. Otherwise, file 
   `.mashery-logon` in the user's home directory will be consulted for the access token. If this file would not be
   found, or if the token is already expired, then
3. Environment variable `TF_MASHERY_V3_CREDS` pointing to the credentials file. Otherwise, `.mashery-v3-credentials` 
   file is present in user's home directory wll be checked.
   
Where no default authentication would be possible, an explicit configuration is required.

## Supported parameters
The provider accepts the following options:
- `log_file`: a full path to the *directory* where to store log files from individual provider runs. This is an optional
   parameter that you may wish to set to troubleshoot a problem;
- `token`: an active Mashery V3 access token, e.g. retrieved from HashiCorp vault. It can also be passed as `TF_MASHERY_V3_ACCESS_TOKEN`
  environment variable;
- `token_file`: a path to the token file (e.g. obtained using Mashery V3 client library). It can also be passed
  as `TF_MASHERY_V3_TOKEN_FILE` environment variable;
- `deploy_duration`: a typical deployment duration. It can also be set by `TF_MASHERY_DEPLOY_DURATION` environment
  variable. This setting is used to determine whether the remaining Mashery V3 token time sufficient to perform
  necessary steps. By default, the provider assumes that the deployment would take about 3 minutes and will expect
  that at least 3 minutes are left before the token will expire.
- `qps`:  number of calls per second configured for these credentials. It can also be passed as `TF_MASHERY_V3_QPS` 
  environment variable;
- `network_latency`: average network latency between your location and Mashery V3 API gateway. It can also be passed
  as `TF_MASHERY_V3_NETWORK_LATENCY` environment variable
- `creds_file`: a path to the credentials file that will be used to obtain. The path can also be passed as `TF_MASHERY_V3_CREDS`
  environment variable.

## Authentication
This provider is using Mashery V3 API internally. Mashery V3 requires the caller to [authenticate](https://developer.mashery.com/docs/read/mashery_api/30/Authentication)
and pass an access token with each V3 API call. The access token is obtained based on five data elements:
- Mashery V3 API key and secret, 
- Mashery username and password, and
- Mashery Area Id where changes are being made.

> For security reasons, the Mashery provider is designed in such a way that it
> actively avoids storing any of these credentials anywhere in the Terraform configuration files.
> These credentials are **long-lived**. Should they leak outside of secure perimeter, it may be possible to retrieve
> API key and secret for all API consumers.
>
> Prefer passing these parameters indirectly wherever possible. **NEVER** store these in your 
> Terraform project.

Multiple options to supply authentication exist, depending on specifics of the deployment and the
physical location of where the `terraform` command is executed:

1. Highly advised for production use: Vault-issued, quickly revoked access token;
2. Obtain a token and pass it to `terraform` as an environment variable;
3. Obtain an access token and save it in the file on the same machine;
4. Save the Mashery V3 credentials on the machine where the `terraform` command
   is running.
   > Note: this option should be reserved for proof-of-concept exercises and/or developer workstations.
   > Using this option in shared environments is highly discouraged for obvious security
   > considerations.

Where the authentication method implies saving the authentication token on disk, it is implied that a script/program
needs to be written/deployed/executed prior to running the `terraform` command. An example of such program is
[mashery-connect](https://github.com/aliakseiyanchuk/mashery-v3-go-client/wiki) achieving the required. This program
is provided to help you started quickly. 

## Authentication with Mashery V2/V3 Secret Engine
The recommended way for supplying authentication is to obtain the access token from a HashiCorp vault, first Mashery V2/V3 secret engine needs to be 
deployed and configured in the server/cluster that you will be obtaining the access token from.
Once configured, then the Mashery access token could be read with the standard vault provider:

```hcl
terraform {
  required_providers {
    mashery = {
      version = "0.0.1"
      source = "yanchuk.nl/aliakseiyanchuk/mashery"
    }
    vault = {
      source = "hashicorp/vault"
      version = "~> 2.17.0"
    }

  }
}

provider "vault" {
  address = var.vault_url
  // Other configuration as required
}

data "vault_generic_secret" "v3Token" {
  path = "mash-auth/auth/testSite/v3"
}

provider "mashery" {
  token = "${data.vault_generic_secret.v3Token.data["access_token"]}"
  qps = "${data.vault_generic_secret.v3Token.data["qps"]}"
}
```

Mashery V3 access token have a validity of 1 hour. Obtaining the tokens from V2/V3 secret engine
allows configuring a shorter duration of the access token, after which it will be automatically revoked.

> Note: with this approach, the maximum possible duration of the `terraform plan` or `terraform apply`
> command is limited to 1 hour. This will be sufficient for most applications. For situations where
> the `plan` or `apply  may take longer than 60 minutes, different technique should be used. These
> are explained in this guide below.

## Passing access token using environment variable
This method assumes that an V3 access token will be retrieved (any way you want) and made available 
in the`TF_MASHERY_V3_ACCESS_TOKEN` environment variable. If you are using mash-connect application mentioned
earlier, then the following commands would achieve this:

```shell
$ mash-connect init
$ eval "$(./mash-connect.exe export TF_MASHERY_V3_ACCESS_TOKEN)"
```

On Windows it could be achieved by executing the following sequence of commands:
```shell
> mash-connect init
> mash-connect export -win TF_MASHERY_V3_ACCESS_TOKEN > mash_token.bat
> call mash_token.bat
```
Why so complicated? Keep in mind that tokens typed on the console could be logged. It's a good security habit
to never type confidential token in the command line.

Once the token is initialized, the Mashery provider is workable with the example configuration.

> Note: using this method it is not possible to determine whether the access token is valid as well as how long will
> the token be valid. Remember to refresh V3 access token frequently as they expire after 1 hour.
 
## Pre-saving access token in a file (as substitute for Vault)
If you are running terraform on a long-lived host or on your personal laptop, there is a fallback for this provider
to read the V3 access token from the (refreshable) file. Storing a JSON token is a variant of the fixed access token 
where a Terraform developer also wants to assert the active time remaining for this token.

This approach would also be beneficial in situations where `terraform` runs plan and deploy commands frequently.

### Saved token file JSON schema

A token file needs to be saved in a JSON file containing timestamp, access token, and expiry time (in seconds) of 
this given access token. The following gives a schematic overview:
```json
{
  "obtained":"2021-02-06T22:59:23.3492214+01:00",
  "access_token":"accessToken1234567890123",
  "qps": 2,
  "expires_in":3600
}
```
This file os is in *json* as it will be an output of a program and ***should not*** be editable manually.

### Helper tools
The [mashery-connect](https://github.com/aliakseiyanchuk/mashery-v3-go-client/wiki) sample illustrates the process of
saving this data into this file.

The  `mash-connect init` command is executed to save the file in the default location (such that it would be read
by the provider without explicit configuration). 

### Provider configuration

By default, the provider will be looking for a file called `.mashery-logon` in user's home directory. If required, any 
other path can be specified by providing `token_file` or `TF_MASHERY_V3_TOKEN_FILE` environment variable.

```hcl
provider "mashery" {
  token_file="/home/terraform-deployer/.v3/.production_area"
}
```

Where  [mashery-connect](https://github.com/aliakseiyanchuk/mashery-v3-go-client/wiki) would be used, then the
following command obtains this file: 
```text
mash-connect init --token-file /home/terraform-deployer/.v3/.production_area
```

### Avoiding token expiry during plan/apply

Where the access token is sourced from file, a check is made to ensure that there is sufficient access time left to
complete the operation run before the token will expire. The default Terraform operation duration is set to be 3 minutes.
Setting `deploy_duration` key or `TF_MASHERY_DEPLOY_DURATION` environment variable allows supplying desired guaranteed
time remaining in `<X>m<Y>s` for X minutes and Y seconds, e.g. `4m27s` will correspond to 4 minute and 27 seconds. 
> Seconds could be omitted, e.g. instead of `4m0s` value `4m` needs to be supplied. 

### Keeping access token alive

On a long-lived host, it is possible to start a daemon process that would proactively obtain the access token, 
thus ensuring that the access token will always be refreshed. The example command
[`mashery-connect keep-live`](https://github.com/aliakseiyanchuk/mashery-v3-go-client/wiki) gives an example of
such application.

Referring to the example above, such process could be started using the following command.

```text
$ nohup mash-connect keep-alive --token-file /home/terraform-deployer/.v3/.production_area 2>&1 &
```

> Note: be sure to understand the risk and implications of running this command. It could be acceptable for a 
> Terraform developer workstation who will start the program at the start of the working day. It is certainly
> less acceptable for production environments.
> 
> If you start the process, be sure you include means to stop, e.g. include a `crontab` that will automatically
> kill this process if it is still running after the working day. 

## Authenticating via credentials file. Encrypting credentials files

The credentials file, as it's name implies, specifies credentials in plain-text that can be used to obtain Mashery
access token. By default, the provider will be looking for a file called `.mashery-v3-credentials` located in the 
current user's home directory. This path can be changed by supplying `creds_file` or `TF_MASHERY_V3_CREDS` environment
variable.

> Given the sensitivity of this data, the provider requires that this file be **encrypted** prior
> this provider will accept it for operation.

For this guide, let's assume that you want to store authentication in the file called
`/home/terraform-deployer/.v3/.my_v3_creds`. The procedure below illustrates steps that could be 
followed in a Linux terminal or Cygwin bash shell. Similar procedure is possible also for Windows
command line or PowerShell.

1. Create this YAML file `home/terraform-deployer/.v3/.my_v3_creds` and specify the following keys
   as relevant to the area you want to administer with the provider:
```yaml
areaId: <area id>
apiKey: <v3 api key>
secret: <v3 api secret>
username: <user name>
password: <password>
```
2. Create a 32-character secret file. This can be achieved e.g. by creating a file accessible
   only to the privileged user. Methods to enforce this varies depending on the operating system
   capabilities. The password should be ideally random-generated.
```text
$   cat /dev/urandom | tr -cd '[a-zA-Z0-9]' | fold -w 32 | head -n 1 > /home/secrets/.my_v3_creds.pass
```
3. Encrypt YAML file using `mash-connect` program:
```text
$ cat /home/secrets/.my_v3_creds.pass | mash-connect encrypt --credentials /home/terraform-deployer/.v3/.my_v3_creds
```
4. Specify the path to this file in your HCL configuration
```hcl
provider "mashery" {
  creds_file="/home/terraform-deployer/.v3/.my_v3_creds"
}
```
5. Ensure you load the value of the secret as `TF_MASHERY_CREDS_PASS` environment variable prior running
the terraform command.
```text
$ export TF_MASHERY_CREDS_PASS=$(cat /home/secrets/.my_v3_creds.pass)
```
6. Move password to the secure place. 

## Configuring multiple Mashery Areas in the same Terraform project

Due to limitations of Mashery V3 API, single provider instance can make changes in a single Mashery area.
Should your project require making changes across multiple Masherey areas, multiple
providers need to be created using the provider alias as [this terraform documentation](https://www.terraform.io/docs/language/providers/configuration.html)
explains. 


