# Terraform Provider for TIBCO Cloud Mashery

This is an experimental Terraform provider for TIBCO Cloud Mashery. It is written 
largely for the proof-of-concept/educational purposes to investigate the possibility of 
managing Mashery deployment from the Terraform-based CI/CD.

Using this provider requires being Mashery customer and understanding how to authenticate
the Terraform process.

The terraform configuration largely repeats the [V3 Mashery API](https://support.mashery.com/docs/read/mashery_api/30)
structure. The provider can be used to create objects and detect configuration drift
applied via Mashery administrative portal. This includes the following resources:
- [service](./website/docs/r/mashery_service.html.markdown), including
  - [OAuth configuration](./website/docs/r/mashery_service_oauth.markdown)
  - [Error set](./website/docs/r/mashery_service_error_set.html.markdown)
  - [Endpoint](./website/docs/r/mashery_service_endpoint.html.markdown)
    - [Endpoint method](./website/docs/r/mashery_service_endpoint_method.html.markdown)
      - [Endpoint method filter](./website/docs/r/mashery_service_endpoint_method_filter.html.markdown)
- [package](./website/docs/r/mashery_package.html.markdown), including
  - [package plan](./website/docs/r/mashery_package_plan.html.markdown)
    - [package plan service](./website/docs/r/mashery_package_plan_service.html.markdown)
      - [package plan service endpoint](./website/docs/r/package_plan_service_endpoint_method.html.markdown)
        - [package plan service endpoint method and filter](./website/docs/r/mashery_package_plan_service_endpoint_method.html.markdown)

The provider allows querying for [roles](./website/docs/d/mashery_role.html.markdown), 
[organizations](./website/docs/d/mashery_organization.html.markdown), and 
[email templates](./website/docs/d/mashery_email_template_set.html.markdown). In practice relatively
few instances of these would be created within the API programme. Typical deployment
would only require referring to existing roles, organizations, and email templates.

The provider does not include creating members, applications, and package keys as these
need to be provisioned via different channels e.g. via the developer portal or via
Mashery's Web administrative console.

Further information can be found [here](./website/docs/index.html.markdown).

## Building and Running
This provider requires a local installation as [this Terraform documentation describes](https://developer.hashicorp.com/terraform/cli/config/config-file#provider-installation).
The provider can be built either using supplied Makefile
```shell
make install
```
> Note: you may need to edit `OS_ARCH` variable to match the processor architecture of
> your system.

On Windows machines, the supplied `makefile.bat` provides an alternative.
```shell
makefile.bat
```