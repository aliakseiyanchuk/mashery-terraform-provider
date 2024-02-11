package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var MemberDataSource = CreateSingularDataSource(mashschemag.MemberResourceSchemaBuilder, queryMember)

func queryMember(ctx context.Context, client v3client.Client, m map[string]string) (masherytypes.MemberIdentifier, *masherytypes.Member, error) {
	if sets, err := client.ListMembersFiltered(ctx, m); err != nil {
		return masherytypes.MemberIdentifier{}, nil, err
	} else if len(sets) == 1 {
		return sets[0].Identifier(), &sets[0], nil
	} else {
		return masherytypes.MemberIdentifier{}, nil, nil
	}
}
