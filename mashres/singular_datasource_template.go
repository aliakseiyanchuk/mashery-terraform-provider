package mashres

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sort"
	"strings"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

// QueryFunc Queries root elements of the data source that do not require parent object context
type QueryFunc[Ident any, MType any] func(context.Context, v3client.Client, map[string]string) (Ident, *MType, error)

// ParentQueryFunc query function that needs a parent identifier
type ParentQueryFunc[ParentIdent any, Ident any, MType any] func(context.Context, v3client.Client, ParentIdent, map[string]string) (Ident, *MType, error)

type DatasourceTemplate[ParentIdent any, Ident any, MType any] interface {
	TestState() *schema.ResourceData
	DataSourceSchema() *schema.Resource
	Query(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics
	DatasourceMapper() *tfmapper.Mapper[ParentIdent, Ident, MType]
}

type SingularDatasourceTemplate[ParentIdent any, Ident any, MType any] struct {
	Schema map[string]*schema.Schema
	Mapper *tfmapper.Mapper[ParentIdent, Ident, MType]

	DoQuery       QueryFunc[Ident, MType]
	DoParentQuery ParentQueryFunc[ParentIdent, Ident, MType]

	ResourcePrefix string
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) CacheKeyOf(ident interface{}, query map[string]string) string {
	sb := strings.Builder{}

	if ident != nil {
		b, _ := json.Marshal(ident)
		sb.Write(b)
		sb.WriteString("---")
	}

	keys := make([]string, len(query))

	idx := 0
	for k := range query {
		keys[idx] = k
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(query[k])
		sb.WriteString(";;")
	}

	h := sha256.New()
	h.Write([]byte(sb.String()))
	bs := h.Sum(nil)

	return fmt.Sprintf("%s_%x", sdt.ResourcePrefix, bs)
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) DatasourceMapper() *tfmapper.Mapper[ParentIdent, Ident, MType] {
	return sdt.Mapper
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) TestState() *schema.ResourceData {
	res := schema.Resource{
		Schema: sdt.Schema,
	}

	return res.TestResourceData()
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) DataSourceSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: sdt.Query,
		Schema:      sdt.Schema,
	}
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) isMatchRequired(d *schema.ResourceData) bool {
	rv := mashschema.ExtractBool(d, mashschema.MashDataSourceRequired, true)
	return rv
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) DoCachedParentQuery(ctx context.Context, v3Cl v3client.Client, parentIdent ParentIdent, searchParams map[string]string) (Ident, *MType, error) {
	cacheKey := sdt.CacheKeyOf(parentIdent, searchParams)

	cachedIdent, cachedData := GetFromCache[Ident, MType](ctx, cacheKey)
	if cachedIdent != nil && cachedData != nil {
		tflog.Info(ctx, fmt.Sprintf("%s of %s with search params %s is served from cache", sdt.ResourcePrefix, parentIdent, searchParams))
		return *cachedIdent, cachedData, nil
	}

	rvIdent, rvData, err := sdt.DoParentQuery(ctx, v3Cl, parentIdent, searchParams)
	if err == nil && rvData != nil {
		tflog.Info(ctx, fmt.Sprintf("%s of %s with search params %s is now stored in cache", sdt.ResourcePrefix, parentIdent, searchParams))
		StoreInCacheDefault(ctx, cacheKey, rvIdent, rvData)
	}

	return rvIdent, rvData, err
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) DoCacheQuery(ctx context.Context, v3Cl v3client.Client, searchParams map[string]string) (Ident, *MType, error) {
	cacheKey := sdt.CacheKeyOf(nil, searchParams)

	cachedIdent, cachedData := GetFromCache[Ident, MType](ctx, cacheKey)
	if cachedIdent != nil && cachedData != nil {
		tflog.Info(ctx, fmt.Sprintf("%s with search params %s is served from cache", sdt.ResourcePrefix, searchParams))

		return *cachedIdent, cachedData, nil
	}

	rvIdent, rvData, err := sdt.DoQuery(ctx, v3Cl, searchParams)
	if err == nil && rvData != nil {
		tflog.Info(ctx, fmt.Sprintf("%s with search params %s is now stored in cache", sdt.ResourcePrefix, searchParams))
		StoreInCacheDefault(ctx, cacheKey, rvIdent, rvData)
	}

	return rvIdent, rvData, err
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) ExecQuery(ctx context.Context, v3Cl v3client.Client, data *schema.ResourceData) (Ident, *MType, error) {
	query := mashschema.ExtractStringMap(data, mashschema.MashDataSourceSearch)
	if len(query) == 0 {
		rvIdent := new(Ident)
		rvType := new(MType)
		return *rvIdent, rvType, errors.New("search criteria must be specified")
	}

	if sdt.DoParentQuery != nil {
		if parentIdent, identErr := sdt.Mapper.ParentIdentity(data); identErr != nil {
			rvIdent := new(Ident)
			rvType := new(MType)
			return *rvIdent, rvType, errors.New("identity must be read at this point")
		} else {
			return sdt.DoCachedParentQuery(ctx, v3Cl, parentIdent, query)
		}
	}

	return sdt.DoCacheQuery(ctx, v3Cl, query)
}

func (sdt *SingularDatasourceTemplate[ParentIdent, Ident, MType]) Query(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Cl := m.(v3client.Client)

	if ident, obj, err := sdt.ExecQuery(ctx, v3Cl, d); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("query has returned an error: %s", err.Error()),
		}}
	} else if obj == nil && sdt.isMatchRequired(d) {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "no matching object was found, however the configuration requires a match",
		}}
	} else {
		if obj != nil {
			_ = sdt.Mapper.AssignIdentity(ident, d)
			tfmapper.SetKeyWithDiag(d, mashschema.MashDataSourceExists, true)
		} else {
			d.SetId("nil")
			tfmapper.SetKeyWithDiag(d, mashschema.MashDataSourceExists, false)
		}
		return sdt.Mapper.RemoteToSchema(obj, d)
	}
}

var singularDataSourceSchema = map[string]*schema.Schema{
	mashschema.MashDataSourceSearch: {
		Type:        schema.TypeMap,
		Required:    true,
		Description: "Search conditions for this resource, typically name = value",
		Elem:        mashschema.StringElem(),
	},
	mashschema.MashDataSourceRequired: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "If true (default), then exactly a single object must be found. If an element doesn't exist, the error is generated",
	}, mashschema.MashDataSourceExists: {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "True if match has been found; and false otherwise",
	},
}

func CreateSingularDataSource[ParentIdent any, Ident any, MType any](
	resourcePrefix string,
	builder *tfmapper.SchemaBuilder[ParentIdent, Ident, MType],
	queryFunc QueryFunc[Ident, MType]) *SingularDatasourceTemplate[ParentIdent, Ident, MType] {

	mapperSchema := tfmapper.MergeSchemas(mashschema.CloneAsComputed(builder.ResourceSchema()), singularDataSourceSchema)

	rv := SingularDatasourceTemplate[ParentIdent, Ident, MType]{
		Schema:         mapperSchema,
		Mapper:         builder.Mapper(),
		DoQuery:        queryFunc,
		ResourcePrefix: resourcePrefix,
	}

	return &rv
}

func CreateSingularParentScopedDataSource[ParentIdent any, Ident any, MType any](
	resourcePrefix string,
	builder *tfmapper.SchemaBuilder[ParentIdent, Ident, MType],
	parentElementSchemaKey string,
	queryFunc ParentQueryFunc[ParentIdent, Ident, MType]) *SingularDatasourceTemplate[ParentIdent, Ident, MType] {

	mapperSchema := tfmapper.MergeSchemas(mashschema.CloneAsComputedExcept(builder.ResourceSchema(), parentElementSchemaKey), singularDataSourceSchema)

	rv := SingularDatasourceTemplate[ParentIdent, Ident, MType]{
		Schema:         mapperSchema,
		Mapper:         builder.Mapper(),
		DoParentQuery:  queryFunc,
		ResourcePrefix: resourcePrefix,
	}

	return &rv
}
