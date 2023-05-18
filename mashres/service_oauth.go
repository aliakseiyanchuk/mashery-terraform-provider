package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ServiceOAuthResource *ResourceTemplate[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth]

func init() {
	ServiceOAuthResource = &ResourceTemplate[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth]{
		Schema: mashschemag.ServiceOAuthResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ServiceOAuthResourceSchemaBuilder.Mapper(),

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) (*masherytypes.MasheryOAuth, error) {
			return client.GetServiceOAuthSecurityProfile(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, ident masherytypes.ServiceIdentifier, m masherytypes.MasheryOAuth) (*masherytypes.MasheryOAuth, *masherytypes.ServiceIdentifier, error) {
			m.ParentService.ServiceId = ident.ServiceId

			if createOAuth, err := client.CreateServiceOAuthSecurityProfile(ctx, m); err != nil {
				return nil, nil, err
			} else {
				rvIdent := masherytypes.ServiceIdentifier{
					ServiceId: ident.ServiceId,
				}
				return createOAuth, &rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier, m masherytypes.MasheryOAuth) (*masherytypes.MasheryOAuth, error) {
			m.ParentService.ServiceId = identifier.ServiceId

			if updatedOAuth, err := client.UpdateServiceOAuthSecurityProfile(ctx, m); err != nil {
				return nil, err
			} else {
				return updatedOAuth, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) error {
			return client.DeleteServiceOAuthSecurityProfile(ctx, identifier)
		},

		// Offending count is not required for Service OAuth, as it can be deleted at any moment.
	}
}
