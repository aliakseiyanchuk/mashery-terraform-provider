package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	mashschema "terraform-provider-mashery/mashschema"
)

var PackageKeyResource *ResourceTemplate

func init() {
	PackageKeyResource = &ResourceTemplate{
		Mapper: mashschema.PackageKeyMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.GetPackageKey(ctx, identifier.(masherytypes.PackageKeyIdentifier))
		},
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreatePackageKey(ctx, identifier.(masherytypes.ApplicationIdentifier), upsertable.(masherytypes.PackageKey))
		},
		DoUpdate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable) (mashschema.Upsertable, error) {
			return client.UpdatePackageKey(ctx, upsertable.(masherytypes.PackageKey))
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeletePackageKey(ctx, identifier.(masherytypes.PackageKeyIdentifier))
		},
	}
}
