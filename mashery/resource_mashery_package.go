package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var PackageResource *ResourceTemplate

func init() {
	PackageResource = &ResourceTemplate{
		Mapper: mashschema.PackageMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.GetPackage(ctx, identifier.(masherytypes.PackageIdentifier))
		},
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreatePackage(ctx, upsertable.(masherytypes.Package))
		},
		DoUpdate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable) (mashschema.Upsertable, error) {
			return client.UpdatePackage(ctx, upsertable.(masherytypes.Package))
		},
		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (int64, error) {
			return client.CountPlans(ctx, identifier.(masherytypes.PackageIdentifier))
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeletePackage(ctx, identifier.(masherytypes.PackageIdentifier))
		},
	}
}
