package mashres

import (
	"context"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var PackagePlanServiceEndpointMethodResource *ResourceTemplate[masherytypes.PackagePlanServiceEndpointIdentifier, masherytypes.PackagePlanServiceEndpointMethodIdentifier, mashschemag.PackagePlanServiceEndpointMethodParam]

func init() {
	PackagePlanServiceEndpointMethodResource = &ResourceTemplate[masherytypes.PackagePlanServiceEndpointIdentifier, masherytypes.PackagePlanServiceEndpointMethodIdentifier, mashschemag.PackagePlanServiceEndpointMethodParam]{
		Schema: mashschemag.PackagePlanServiceEndpointMethodResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.PackagePlanServiceEndpointMethodResourceSchemaBuilder.Mapper(),

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointMethodIdentifier) (*mashschemag.PackagePlanServiceEndpointMethodParam, error) {
			if methState, err := client.GetPackagePlanMethod(ctx, identifier); err != nil {
				return nil, err
			} else if methState == nil {
				// This method no longer exists.
				return nil, nil
			} else {
				rv := mashschemag.PackagePlanServiceEndpointMethodParam{
					ServiceEndpointMethod: identifier.ServiceEndpointMethodIdentifier,
				}

				// Try reading a method to detect drift.
				if filterState, err := client.GetPackagePlanMethodFilter(ctx, identifier); err != nil {
					return nil, err
				} else if filterState != nil {
					filterIdent := filterState.Identifier()
					rv.ServiceEndpointMethodFilterDesired = masherytypes.ServiceEndpointMethodFilterIdentifier{
						ServiceEndpointMethodIdentifier: masherytypes.ServiceEndpointMethodIdentifier{
							ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
								ServiceIdentifier: masherytypes.ServiceIdentifier{
									ServiceId: filterIdent.ServiceId,
								},
								EndpointId: filterIdent.EndpointId,
							},
							MethodId: filterIdent.MethodId,
						},
						FilterId: filterIdent.FilterId,
					}
				}

				return &rv, nil
			}
		},

		DoCreate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointIdentifier, m mashschemag.PackagePlanServiceEndpointMethodParam) (*mashschemag.PackagePlanServiceEndpointMethodParam, *masherytypes.PackagePlanServiceEndpointMethodIdentifier, error) {

			ident := masherytypes.PackagePlanServiceEndpointMethodIdentifier{
				ServiceEndpointMethodIdentifier: m.ServiceEndpointMethod,
				PackagePlanIdentifier:           identifier.PackagePlanIdentifier,
			}

			if _, err := client.CreatePackagePlanMethod(ctx, ident); err != nil {
				return nil, nil, err
			} else {
				if len(m.ServiceEndpointMethodFilterDesired.FilterId) > 0 {
					filterIdent := masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{
						ServiceEndpointMethodFilterIdentifier: m.ServiceEndpointMethodFilterDesired,
						PackagePlanServiceIdentifier: masherytypes.PackagePlanServiceIdentifier{
							PackagePlanIdentifier: identifier.PackagePlanIdentifier,
							ServiceIdentifier:     identifier.ServiceIdentifier,
						},
					}

					if obj, err := client.CreatePackagePlanMethodFilter(ctx, filterIdent); err != nil {
						return &m, &ident, err
					} else if obj == nil {
						return &m, &ident, errors.New("the api did not return a filter state; the filter was not created")
					}
				}

				return &m, &ident, nil
			}
		},

		// Update is not required: it will be delete-only

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointMethodIdentifier, param mashschemag.PackagePlanServiceEndpointMethodParam) (*mashschemag.PackagePlanServiceEndpointMethodParam, error) {
			rv := mashschemag.PackagePlanServiceEndpointMethodParam{}

			if param.ServiceEndpointMethodFilterChanged {

				if len(param.ServiceEndpointMethodFilterPrevious.FilterId) > 0 {
					if err := client.DeleteEndpointMethodFilter(ctx, param.ServiceEndpointMethodFilterPrevious); err != nil {
						return nil, err
					}
				}
				if len(param.ServiceEndpointMethodFilterDesired.FilterId) > 0 {
					ident := masherytypes.PackagePlanServiceEndpointMethodFilterIdentifier{
						ServiceEndpointMethodFilterIdentifier: param.ServiceEndpointMethodFilterDesired,
						PackagePlanServiceIdentifier: masherytypes.PackagePlanServiceIdentifier{
							PackagePlanIdentifier: identifier.PackagePlanIdentifier,
							ServiceIdentifier:     identifier.ServiceIdentifier,
						},
					}

					if obj, err := client.CreatePackagePlanMethodFilter(ctx, ident); err != nil {
						return nil, err
					} else if obj == nil {
						return nil, errors.New("the api did not return a filter state; the filter was not created")
					}
				}
			}

			return &rv, nil
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointMethodIdentifier) error {
			return client.DeletePackagePlanMethod(ctx, identifier)
		},

		// Offending methods is not required: it's a sub-part of the endpoint, and it doesn't need to be specifically
		// tracked.
	}
}
