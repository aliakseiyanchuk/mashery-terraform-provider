package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
	"terraform-provider-mashery/tfmapper"
)

var MemberResource *ResourceTemplate[tfmapper.Orphan, masherytypes.MemberIdentifier, masherytypes.Member]

func init() {
	MemberResource = &ResourceTemplate[tfmapper.Orphan, masherytypes.MemberIdentifier, masherytypes.Member]{
		Schema: mashschemag.MemberResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.MemberResourceSchemaBuilder.Mapper(),

		ImportIdentityParser: regexpImportIdentityParser("/members/(.+)", masherytypes.MemberIdentifier{}, func(items []string) masherytypes.MemberIdentifier {
			return masherytypes.MemberIdentifier{MemberId: items[1]}
		}),

		UpsertableFunc: func() masherytypes.Member {
			return masherytypes.Member{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.MemberIdentifier) (masherytypes.Member, bool, error) {
			return client.GetMember(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, orphan tfmapper.Orphan, m masherytypes.Member) (masherytypes.Member, masherytypes.MemberIdentifier, error) {
			if createdMember, err := client.CreateMember(ctx, m); err != nil {
				return masherytypes.Member{}, masherytypes.MemberIdentifier{}, err
			} else {
				rvIdent := createdMember.Identifier()
				return createdMember, rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.MemberIdentifier, m masherytypes.Member) (masherytypes.Member, error) {
			m.Id = identifier.MemberId

			return client.UpdateMember(ctx, m)
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.MemberIdentifier) error {
			return client.DeleteMember(ctx, identifier)
		},

		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier masherytypes.MemberIdentifier) (int64, error) {
			return client.CountApplicationsOfMember(ctx, identifier)
		},
	}
}
