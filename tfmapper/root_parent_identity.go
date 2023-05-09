package tfmapper

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type RootParentIdentity struct {
	IdentityMapper[Orphan]
}

func (r RootParentIdentity) Identity(_ *schema.ResourceData) (Orphan, error) {
	return 0, nil
}

func (r RootParentIdentity) Assign(_ Orphan, _ *schema.ResourceData) error {
	// Nothing to do; ignored
	return nil
}
