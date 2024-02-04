package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ApplicationResource *ResourceTemplate[masherytypes.MemberIdentifier, mashschemag.ApplicationOfMemberIdentifier, masherytypes.Application]

func init() {
	ApplicationResource = &ResourceTemplate[masherytypes.MemberIdentifier, mashschemag.ApplicationOfMemberIdentifier, masherytypes.Application]{
		Schema: mashschemag.ApplicationResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ApplicationResourceSchemaBuilder.Mapper(),

		UpsertableFunc: func() masherytypes.Application { return masherytypes.Application{} },

		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationOfMemberIdentifier) (masherytypes.Application, bool, error) {
			return client.GetApplication(ctx, identifier.ApplicationIdentifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, parent masherytypes.MemberIdentifier, application masherytypes.Application) (masherytypes.Application, mashschemag.ApplicationOfMemberIdentifier, error) {
			if createApp, err := client.CreateApplication(ctx, parent, application); err != nil {
				return masherytypes.Application{}, mashschemag.ApplicationOfMemberIdentifier{}, err
			} else {
				rvIdent := mashschemag.ApplicationOfMemberIdentifier{
					MemberIdentifier:      parent,
					ApplicationIdentifier: createApp.Identifier(),
				}
				return createApp, rvIdent, err
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationOfMemberIdentifier, application masherytypes.Application) (masherytypes.Application, error) {
			application.Id = identifier.ApplicationId
			return client.UpdateApplication(ctx, application)
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationOfMemberIdentifier) error {
			return client.DeleteApplication(ctx, identifier.ApplicationIdentifier)
		},

		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationOfMemberIdentifier) (int64, error) {
			return client.CountApplicationPackageKeys(ctx, identifier.ApplicationIdentifier)
		},
	}
}
