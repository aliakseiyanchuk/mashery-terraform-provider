package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var ApplicationResource *ResourceTemplate

func init() {
	ApplicationResource = &ResourceTemplate{
		Mapper: mashschema.ApplicationMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.GetApplication(ctx, identifier.(masherytypes.ApplicationIdentifier))
		},
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, parentObject mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreateApplication(ctx, upsertable.(masherytypes.Application), parentObject.(masherytypes.MemberIdentifier))
		},
		DoUpdate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable) (mashschema.Upsertable, error) {
			return client.UpdateApplication(ctx, upsertable.(masherytypes.Application))
		},
		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (int64, error) {
			return client.CountApplicationPackageKeys(ctx, identifier.(masherytypes.ApplicationIdentifier))
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeleteApplication(ctx, identifier.(masherytypes.ApplicationIdentifier))
		},
	}
}
