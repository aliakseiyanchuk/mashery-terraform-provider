package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var MemberResource *ResourceTemplate

func init() {
	MemberResource = &ResourceTemplate{
		Mapper: mashschema.MemberMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			typedIdentifier := identifier.(masherytypes.MemberIdentifier)
			return client.GetMember(ctx, typedIdentifier)
		},
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, _ mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			typedUpsertable := upsertable.(masherytypes.Member)
			return client.CreateMember(ctx, typedUpsertable)
		},
		DoUpdate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable) (mashschema.Upsertable, error) {
			typedUpsertable := upsertable.(masherytypes.Member)
			return client.UpdateMember(ctx, typedUpsertable)
		},
		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (int64, error) {
			memberId := identifier.(masherytypes.MemberIdentifier)
			return client.CountApplicationsOfMember(ctx, memberId)
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			memberId := identifier.(masherytypes.MemberIdentifier)
			return client.DeleteMember(ctx, memberId)
		},
	}
}
