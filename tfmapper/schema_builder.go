package tfmapper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type TFResourceSchema = map[string]*schema.Schema

type Orphan int

// LocatorFunc a function that locates a field in the mapping struct
type LocatorFunc[I any, O any] func(in *I) *O

type SchemaBuilder[ParentIdent any, Ident any, MType any] struct {
	resourceSchema       TFResourceSchema
	identityMapper       IdentityMapper[Ident]
	parentIdentityMapper IdentityMapper[ParentIdent]
	fields               []FieldMapper[MType]
}

func NewSchemaBuilder[ParentIdent any, Ident any, MType any]() *SchemaBuilder[ParentIdent, Ident, MType] {
	return &SchemaBuilder[ParentIdent, Ident, MType]{
		resourceSchema: map[string]*schema.Schema{},
		fields:         []FieldMapper[MType]{},
	}
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) Identity(mapper IdentityMapper[Ident]) *SchemaBuilder[ParentIdent, Ident, MType] {
	sb.identityMapper = mapper
	return sb
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) ParentIdentity(key string, fieldSchema *schema.Schema, mapper IdentityMapper[ParentIdent]) *SchemaBuilder[ParentIdent, Ident, MType] {
	sb.resourceSchema[key] = fieldSchema
	sb.parentIdentityMapper = mapper
	return sb
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) RootIdentity(mapper IdentityMapper[ParentIdent]) *SchemaBuilder[ParentIdent, Ident, MType] {
	sb.parentIdentityMapper = mapper
	return sb
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) Add(field FieldMapper[MType]) *SchemaBuilder[ParentIdent, Ident, MType] {
	if comp, ok := field.(CompositeFieldMapper); ok {
		for k, v := range comp.GetCompositeSchema() {
			if sb.resourceSchema[k] != nil {
				panic(fmt.Sprintf("duplicate key %s; change code or use composite", k))
			}
			sb.resourceSchema[k] = v
		}
	} else {
		if sb.resourceSchema[field.GetKey()] != nil {
			panic(fmt.Sprintf("duplicate key %s; change code or use composite", field.GetKey()))
		}
		sb.resourceSchema[field.GetKey()] = field.GetSchema()
	}

	sb.fields = append(sb.fields, field)
	return sb
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) AddComposite(field FieldMapper[MType], extras map[string]*schema.Schema) *SchemaBuilder[ParentIdent, Ident, MType] {
	sb.resourceSchema[field.GetKey()] = field.GetSchema()
	for k, v := range extras {
		sb.resourceSchema[k] = v
	}

	sb.fields = append(sb.fields, field)
	return sb
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) ResourceSchema() TFResourceSchema {
	return sb.resourceSchema
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) TestResourceData() *schema.ResourceData {
	res := schema.Resource{
		Schema: sb.resourceSchema,
	}

	return res.TestResourceData()
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) Mapper() *Mapper[ParentIdent, Ident, MType] {
	return &Mapper[ParentIdent, Ident, MType]{
		fields:               sb.fields,
		identityMapper:       sb.identityMapper,
		parentIdentityMapper: sb.parentIdentityMapper,
	}
}

type Mapper[ParentIdent any, Ident any, MType any] struct {
	identityMapper       IdentityMapper[Ident]
	parentIdentityMapper IdentityMapper[ParentIdent]
	fields               []FieldMapper[MType]
}

func (m *Mapper[ParentIdent, Ident, MType]) Identity(state *schema.ResourceData) (Ident, error) {
	return m.identityMapper.Identity(state)
}

func (m *Mapper[ParentIdent, Ident, MType]) ParentIdentity(state *schema.ResourceData) (ParentIdent, error) {
	return m.parentIdentityMapper.Identity(state)
}

func (m *Mapper[ParentIdent, Ident, MType]) AssignIdentity(ident Ident, state *schema.ResourceData) error {
	return m.identityMapper.Assign(ident, state)
}

func (m *Mapper[ParentIdent, Ident, MType]) ResetIdentity(state *schema.ResourceData) {
	state.SetId("")
}

func (m *Mapper[ParentIdent, Ident, MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) diag.Diagnostics {
	rv := diag.Diagnostics{}

	for _, k := range m.fields {
		if dg := k.RemoteToSchema(remote, state); dg != nil {
			rv = append(rv, *dg)
		}
	}

	return rv
}

func (m *Mapper[ParentIdent, Ident, MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	for _, k := range m.fields {
		fmt.Println("Mapping " + k.GetKey())
		k.SchemaToRemote(state, remote)
	}
}

type FieldMapper[MType any] interface {
	GetKey() string
	GetSchema() *schema.Schema

	RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic
	SchemaToRemote(state *schema.ResourceData, remote *MType)
}

type CompositeFieldMapper interface {
	GetCompositeSchema() map[string]*schema.Schema
}

type FieldMapperBase struct {
	Key    string
	Schema *schema.Schema

	ParentIdentityKey string
}

type PluggableFiledMapperBase[MType any] struct {
	FieldMapperBase

	RemoteToSchemaFunc func(remote *MType, key string, state *schema.ResourceData) *diag.Diagnostic
	SchemaToRemoteFunc func(state *schema.ResourceData, key string, remote *MType)
}

func (pfmb *PluggableFiledMapperBase[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	if pfmb.RemoteToSchemaFunc != nil {
		return pfmb.RemoteToSchemaFunc(remote, pfmb.Key, state)
	} else {
		return nil
	}
}

func (pfmb *PluggableFiledMapperBase[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	if pfmb.SchemaToRemoteFunc != nil {
		pfmb.SchemaToRemoteFunc(state, pfmb.Key, remote)
	}
}

type CompositeMapperBase struct {
	CompositeSchema map[string]*schema.Schema
}

func (cmb *CompositeMapperBase) GetCompositeSchema() map[string]*schema.Schema {
	return cmb.CompositeSchema
}

func (fmb *FieldMapperBase) GetKey() string {
	return fmb.Key
}

func (fmb *FieldMapperBase) GetSchema() *schema.Schema {
	return fmb.Schema
}

type IdentityMapper[Ident any] interface {
	Identity(state *schema.ResourceData) (Ident, error)
	Assign(ident Ident, state *schema.ResourceData) error
}

type JsonIdentityMapper[Ident any] struct {
	Key          string
	IdentityFunc func() Ident
}

func (im *JsonIdentityMapper[Ident]) Identity(state *schema.ResourceData) (Ident, error) {
	rv := im.IdentityFunc()
	if len(im.Key) > 0 {
		v := mashschema.ExtractString(state, im.Key, "")
		err := unwrapJSON(v, &rv)
		return rv, err
	} else {

		err := unwrapJSON(state.Id(), &rv)
		return rv, err
	}

}

func (im *JsonIdentityMapper[Ident]) Assign(ident Ident, state *schema.ResourceData) error {
	val := wrapJSON(ident)
	if len(im.Key) > 0 {
		return state.Set(im.Key, val)
	} else {
		state.SetId(val)
		return nil
	}
}

func unwrapJSON(inp string, receiver interface{}) error {
	if sEnc, err := base64.StdEncoding.DecodeString(inp); err != nil {
		return err
	} else {
		return json.Unmarshal(sEnc, receiver)
	}
}

func wrapJSON(origin interface{}) string {
	v, _ := json.Marshal(origin)
	return base64.StdEncoding.EncodeToString(v)
}
