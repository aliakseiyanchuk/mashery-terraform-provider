package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var PackagePlanServiceEndpointMethodResource *ResourceTemplate[masherytypes.PackagePlanServiceEndpointIdentifier, masherytypes.PackagePlanServiceEndpointMethodIdentifier, mashschemag.PackagePlanServiceEndpointMethodParam]

func init() {
	PackagePlanServiceEndpointMethodResource = &ResourceTemplate[masherytypes.PackagePlanServiceEndpointIdentifier, masherytypes.PackagePlanServiceEndpointMethodIdentifier, mashschemag.PackagePlanServiceEndpointMethodParam]{
		Schema: mashschemag.PackagePlanServiceEndpointMethodResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.PackagePlanServiceEndpointMethodResourceSchemaBuilder.Mapper(),

		ImportIdentityParser: regexpImportIdentityParser("/packages/(.+)/plans/(.+)/services/(.+)/endpoints/(.+)/methods/(.+)",
			masherytypes.PackagePlanServiceEndpointMethodIdentifier{},
			func(items []string) masherytypes.PackagePlanServiceEndpointMethodIdentifier {
				rv := masherytypes.PackagePlanServiceEndpointMethodIdentifier{}
				rv.PackageId = items[1]
				rv.PlanId = items[2]
				rv.ServiceId = items[3]
				rv.EndpointId = items[4]
				rv.MethodId = items[5]

				return rv
			}),

		UpsertableFunc: func() mashschemag.PackagePlanServiceEndpointMethodParam {
			return mashschemag.PackagePlanServiceEndpointMethodParam{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointMethodIdentifier) (mashschemag.PackagePlanServiceEndpointMethodParam, bool, error) {
			rv := mashschemag.PackagePlanServiceEndpointMethodParam{
				ServiceEndpointMethod: identifier.ServiceEndpointMethodIdentifier,
			}

			if _, methodExists, err := client.GetPackagePlanMethod(ctx, identifier); err != nil || !methodExists {
				return rv, methodExists, err
			} else {

				// Try reading a method to detect drift.
				if filterState, filterExists, filterErr := client.GetPackagePlanMethodFilter(ctx, identifier); filterErr != nil || !filterExists {
					return rv, filterExists, filterErr
				} else if filterExists {
					filterIdent := filterState.Identifier()

					desiredFilterIdent := masherytypes.ServiceEndpointMethodFilterIdentifier{}
					desiredFilterIdent.ServiceId = filterIdent.ServiceId
					desiredFilterIdent.EndpointId = filterIdent.EndpointId
					desiredFilterIdent.MethodId = filterIdent.MethodId
					desiredFilterIdent.FilterId = filterIdent.FilterId

					rv.ServiceEndpointMethodFilterDesired = desiredFilterIdent
					return rv, filterExists, nil
				} else {
					return rv, false, nil
				}
			}
		},

		DoCreate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointIdentifier, m mashschemag.PackagePlanServiceEndpointMethodParam) (mashschemag.PackagePlanServiceEndpointMethodParam, masherytypes.PackagePlanServiceEndpointMethodIdentifier, error) {

			ident := masherytypes.PackagePlanServiceEndpointMethodIdentifier{
				ServiceEndpointMethodIdentifier: m.ServiceEndpointMethod,
				PackagePlanIdentifier:           identifier.PackagePlanIdentifier,
			}

			if _, err := client.CreatePackagePlanMethod(ctx, ident); err != nil {
				return m, ident, err
			} else {
				if len(m.ServiceEndpointMethodFilterDesired.FilterId) > 0 {
					filterIdent := masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{
						ServiceEndpointMethodFilterIdentifier: m.ServiceEndpointMethodFilterDesired,
						PackagePlanIdentifier:                 identifier.PackagePlanIdentifier,
					}

					if _, err := client.CreatePackagePlanMethodFilter(ctx, filterIdent); err != nil {
						return m, ident, err
					}
				}

				return m, ident, nil
			}
		},

		// Update is not required: it will be delete-only

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointMethodIdentifier, param mashschemag.PackagePlanServiceEndpointMethodParam) (mashschemag.PackagePlanServiceEndpointMethodParam, error) {
			rv := mashschemag.PackagePlanServiceEndpointMethodParam{}

			if param.ServiceEndpointMethodFilterChanged {

				if len(param.ServiceEndpointMethodFilterPrevious.FilterId) > 0 {
					if err := client.DeleteEndpointMethodFilter(ctx, param.ServiceEndpointMethodFilterPrevious); err != nil {
						return rv, err
					}
				}
				if len(param.ServiceEndpointMethodFilterDesired.FilterId) > 0 {
					ident := masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{
						ServiceEndpointMethodFilterIdentifier: param.ServiceEndpointMethodFilterDesired,
						PackagePlanIdentifier:                 identifier.PackagePlanIdentifier,
					}

					if _, err := client.CreatePackagePlanMethodFilter(ctx, ident); err != nil {
						return rv, err
					}
				}
			}

			return rv, nil
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointMethodIdentifier) error {
			return client.DeletePackagePlanMethod(ctx, identifier)
		},

		// Offending methods is not required: it's a sub-part of the endpoint, and it doesn't need to be specifically
		// tracked.
	}
}
