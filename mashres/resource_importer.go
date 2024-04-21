package mashres

import (
	"context"
	"errors"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

type ResourceImporter[ParentIdent any, Ident any, MType any] struct {
	schema.ResourceImporter
	resourceTemplate *ResourceTemplate[ParentIdent, Ident, MType]
}

func newResourceImporter[ParentIdent any, Ident any, MType any](rt *ResourceTemplate[ParentIdent, Ident, MType]) *ResourceImporter[ParentIdent, Ident, MType] {
	rv := &ResourceImporter[ParentIdent, Ident, MType]{
		resourceTemplate: rt,
	}

	return rv
}

func (ri *ResourceImporter[ParentIdent, Ident, MType]) AsSchemaImporter() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: ri.Import,
	}
}

func (ri *ResourceImporter[ParentIdent, Ident, MType]) Import(ctx context.Context, rd *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if identity, err := ri.resourceTemplate.ImportIdentityParser(rd.Id()); err != nil {
		return nil, err
	} else {
		if curState, exists, err := ri.resourceTemplate.DoRead(ctx, m.(v3client.Client), identity); err != nil {
			return nil, err
		} else if !exists {
			return nil, errors.New("referenced resource is missing")
		} else {
			if diag := ri.resourceTemplate.Mapper.RemoteToSchema(&curState, rd); diag.HasError() {
				for _, d := range diag {
					tflog.Error(ctx, "error in persisting data", map[string]interface{}{
						"msg":    d.Summary,
						"detail": d.Detail,
					})
				}
				return nil, errors.New("returned data cannot be persisted")
			}
			_ = ri.resourceTemplate.Mapper.AssignIdentity(identity, rd)

			return []*schema.ResourceData{rd}, nil
		}
	}
}

func regexpImportIdentityParser[Ident any](expr string, mismatched Ident, conv func(items []string) Ident) ImportIdentityFunc[Ident] {
	return func(s string) (Ident, error) {
		rg := regexp.MustCompile(expr)
		if !rg.MatchString(s) {
			return mismatched, errors.New("supplied identifier does not match the expected expression")
		}
		m := rg.FindStringSubmatch(s)
		return conv(m), nil
	}
}
