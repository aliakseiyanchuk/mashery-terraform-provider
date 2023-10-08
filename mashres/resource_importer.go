package mashres

import (
	"context"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceImporter[ParentIdent any, Ident any, MType any] struct {
	schema.ResourceImporter
	resourceTemplate ResourceTemplate[ParentIdent, Ident, MType]
}

func newResourceImporter[ParentIdent any, Ident any, MType any](rt *ResourceTemplate[ParentIdent, Ident, MType]) *ResourceImporter[ParentIdent, Ident, MType] {
	rv := &ResourceImporter[ParentIdent, Ident, MType]{}
	rv.StateContext = rv.Import

	return rv
}

func (ri *ResourceImporter[ParentIdent, Ident, MType]) Import(ctx context.Context, rd *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if identity, err := ri.resourceTemplate.Mapper.Identity(rd); err != nil {
		return nil, err
	} else {
		if curState, err := ri.resourceTemplate.DoRead(ctx, m.(v3client.Client), identity); err != nil {
			return nil, err
		} else {
			if diag := ri.resourceTemplate.Mapper.RemoteToSchema(curState, rd); diag.HasError() {
				return nil, errors.New("returned data cannot be persisted")
			}

			return []*schema.ResourceData{rd}, nil
		}
	}
}
