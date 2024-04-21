package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ServiceEndpointMethodFilterResource *ResourceTemplate[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter]

func init() {
	ServiceEndpointMethodFilterResource = &ResourceTemplate[masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethodFilterIdentifier, masherytypes.ServiceEndpointMethodFilter]{
		Schema: mashschemag.ServiceEndpointMethodFilterResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ServiceEndpointMethodFilterResourceSchemaBuilder.Mapper(),

		ImportIdentityParser: regexpImportIdentityParser("/services/(.+)/endpoints/(.+)/methods/(.+)/filters/(.+)",
			masherytypes.ServiceEndpointMethodFilterIdentifier{},
			func(items []string) masherytypes.ServiceEndpointMethodFilterIdentifier {
				rv := masherytypes.ServiceEndpointMethodFilterIdentifier{}
				rv.ServiceId = items[1]
				rv.EndpointId = items[2]
				rv.MethodId = items[3]
				rv.FilterId = items[4]

				return rv
			}),

		UpsertableFunc: func() masherytypes.ServiceEndpointMethodFilter {
			return masherytypes.ServiceEndpointMethodFilter{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointMethodFilterIdentifier) (masherytypes.ServiceEndpointMethodFilter, bool, error) {
			return client.GetEndpointMethodFilter(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, methIdent masherytypes.ServiceEndpointMethodIdentifier, m masherytypes.ServiceEndpointMethodFilter) (masherytypes.ServiceEndpointMethodFilter, masherytypes.ServiceEndpointMethodFilterIdentifier, error) {
			if createdFilter, err := client.CreateEndpointMethodFilter(ctx, methIdent, m); err != nil {
				return masherytypes.ServiceEndpointMethodFilter{}, masherytypes.ServiceEndpointMethodFilterIdentifier{}, err
			} else {
				rvIdent := createdFilter.Identifier()
				return createdFilter, rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointMethodFilterIdentifier, m masherytypes.ServiceEndpointMethodFilter) (masherytypes.ServiceEndpointMethodFilter, error) {
			m.Id = identifier.FilterId
			m.ServiceEndpointMethod = identifier.ServiceEndpointMethodIdentifier

			if updatedFilter, err := client.UpdateEndpointMethodFilter(ctx, m); err != nil {
				return masherytypes.ServiceEndpointMethodFilter{}, err
			} else {
				return updatedFilter, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointMethodFilterIdentifier) error {
			return client.DeleteEndpointMethodFilter(ctx, identifier)
		},
	}
}
