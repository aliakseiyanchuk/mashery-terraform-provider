package tfmapper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"terraform-provider-mashery/funcsupport"
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

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) MapperAndTestData() (*Mapper[ParentIdent, Ident, MType], *schema.ResourceData) {
	return sb.Mapper(), sb.TestResourceData()
}

func MergeSchemas(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	var rv = map[string]*schema.Schema{}

	for _, s := range schemas {
		for k, v := range s {
			ensureUniqueSchemaKey(rv, k)
			rv[k] = v
		}
	}

	return rv
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

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) ParentIdentity(mapper ParentIdentityMapper[ParentIdent]) *SchemaBuilder[ParentIdent, Ident, MType] {
	ensureUniqueSchemaKey(sb.resourceSchema, mapper.GetKey())
	sb.resourceSchema[mapper.GetKey()] = mapper.GetSchema()

	sb.parentIdentityMapper = mapper
	return sb
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) RootIdentity(mapper IdentityMapper[ParentIdent]) *SchemaBuilder[ParentIdent, Ident, MType] {
	sb.parentIdentityMapper = mapper
	return sb
}

func (sb *SchemaBuilder[ParentIdent, Ident, MType]) Add(field FieldMapper[MType]) *SchemaBuilder[ParentIdent, Ident, MType] {
	resourceSchema := sb.resourceSchema
	if comp, ok := field.(CompositeFieldMapper); ok {
		for k, v := range comp.GetCompositeSchema() {
			ensureUniqueSchemaKey(resourceSchema, k)
			resourceSchema[k] = v
		}
	} else {
		ensureUniqueSchemaKey(resourceSchema, field.GetKey())
		resourceSchema[field.GetKey()] = field.GetSchema()
	}

	sb.fields = append(sb.fields, field)
	return sb
}

func ensureUniqueSchemaKey(resourceSchema TFResourceSchema, k string) {
	if resourceSchema[k] != nil {
		panic(fmt.Sprintf("duplicate key %s; change code or use composite", k))
	}
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

func (m *Mapper[ParentIdent, Ident, MType]) StateModified(state *schema.ResourceData) bool {
	for _, fld := range m.fields {
		if m.mapperHasModifications(state, fld) {
			return true
		}
	}

	return false
}

func (m *Mapper[ParentIdent, Ident, MType]) Identity(state *schema.ResourceData) (Ident, error) {
	return m.identityMapper.Identity(state)
}

func (m *Mapper[ParentIdent, Ident, MType]) ParentIdentity(state *schema.ResourceData) (ParentIdent, error) {
	return m.parentIdentityMapper.Identity(state)
}

// TestSetPrentIdentity set the identity of the resource. This method should be used only within the context of the
// unit tests. For actual terraform scripts, these need to be set using references
func (m *Mapper[ParentIdent, Ident, MType]) TestSetPrentIdentity(ident ParentIdent, state *schema.ResourceData) error {
	return m.parentIdentityMapper.Assign(ident, state)
}

func (m *Mapper[ParentIdent, Ident, MType]) AssignIdentity(ident Ident, state *schema.ResourceData) error {
	return m.identityMapper.Assign(ident, state)
}

func (m *Mapper[ParentIdent, Ident, MType]) ResetIdentity(state *schema.ResourceData) {
	state.SetId("")
}

func (m *Mapper[ParentIdent, Ident, MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) diag.Diagnostics {
	rv := diag.Diagnostics{}

	type consumer func(fld FieldMapper[MType], remote *MType, state *schema.ResourceData) *diag.Diagnostic
	var c consumer

	// Handle gracefully nil values by delegating the assignment of nil values to specific methods. This is needed
	// to avoid null pointer exceptions in the individual mappers.

	if remote != nil {
		c = func(fld FieldMapper[MType], remote *MType, state *schema.ResourceData) *diag.Diagnostic {
			return fld.RemoteToSchema(remote, state)
		}
	} else {
		c = func(fld FieldMapper[MType], remote *MType, state *schema.ResourceData) *diag.Diagnostic {
			return fld.NilRemote(state)
		}
	}

	for _, k := range m.fields {
		if dg := c(k, remote, state); dg != nil {
			rv = append(rv, *dg)
		}
	}

	return rv
}

func (m *Mapper[ParentIdent, Ident, MType]) IsStateValid(state *schema.ResourceData) diag.Diagnostics {
	var rv diag.Diagnostics

	for _, k := range m.fields {
		if valid, msg := k.IsStateValid(state); !valid {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("invalid input for field %s", k.GetKey()),
				Detail:   fmt.Sprintf("field validation has returned the following error: %s", msg),
			})
		}
	}

	return rv
}

func (m *Mapper[ParentIdent, Ident, MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	for _, k := range m.fields {
		// If this mapper handles a single field, then casting schema to remove may be skipped
		// for fully computed fields. The schema will be nil in case the mapper is using multiple
		// fields to map to the output, which will never require this.

		if k.GetSchema() != nil {
			// Fully computed fields do not need to be explicitly mapped
			if k.GetSchema().Computed && !k.GetSchema().Optional {
				continue
			}
		}

		// If the mapper is interested in receiving the modification of the field (e.g. to avoid making unnecessary)
		// calls, it will receive the modification.

		k.ConsumeModification(remote, m.mapperHasModifications(state, k))
		k.SchemaToRemote(state, remote)
	}
}

func (m *Mapper[ParentIdent, Ident, MType]) mapperHasModifications(state *schema.ResourceData, fld FieldMapper[MType]) bool {
	if comp, ok := fld.(CompositeFieldMapper); ok {
		keys := make([]string, len(comp.GetCompositeSchema()))
		idx := 0

		for k := range comp.GetCompositeSchema() {
			keys[idx] = k
			idx++
		}
		return state.HasChanges(keys...)
	} else {
		return state.HasChange(fld.GetKey())
	}
}

func (m *Mapper[ParentIdent, Ident, MType]) TestAssign(key string, state *schema.ResourceData, i interface{}) error {
	for _, fm := range m.fields {
		if fm.GetKey() == key {
			if fms, ok := fm.(FieldMapperSetter); ok {
				return fms.ValueToSchema(i, state)
			}
		}
	}

	return errors.New("could not find mapper for this field")
}

// FieldMapper performs mapping of a single fields between Mashery API and Terraform state
type FieldMapper[MType any] interface {
	GetKey() string
	GetSchema() *schema.Schema

	ConsumeModification(out *MType, mod bool)
	IsStateValid(state *schema.ResourceData) (bool, string)

	// NilRemote set the value in case the remote is nil
	NilRemote(state *schema.ResourceData) *diag.Diagnostic
	RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic
	SchemaToRemote(state *schema.ResourceData, remote *MType)
}

type FieldMapperSetter interface {
	ValueToSchema(i interface{}, state *schema.ResourceData) error
}

type CompositeFieldMapper interface {
	GetCompositeSchema() map[string]*schema.Schema
}

type FieldMapperBase[MType any] struct {
	Key    string
	Schema *schema.Schema

	ParentIdentityKey    string
	ModificationConsumer funcsupport.BiConsumer[*MType, bool]

	// ValidateFunc raw data validation function. If defined, verifies if the value
	// supplied by this mapper is re
	ValidateFunc funcsupport.BiFunctionDual[*schema.ResourceData, string, bool, string]
}

func (fmb *FieldMapperBase[MType]) ConsumeModification(out *MType, how bool) {
	if fmb.ModificationConsumer != nil {
		fmb.ModificationConsumer(out, how)
	}
}
func (fmb *FieldMapperBase[MType]) IsStateValid(state *schema.ResourceData) (bool, string) {
	if fmb.ValidateFunc != nil {
		return fmb.ValidateFunc(state, fmb.Key)
	}

	return true, ""
}

type PluggableFiledMapperBase[MType any] struct {
	FieldMapperBase[MType]

	NilRemoteToSchemaFunc func(key string, state *schema.ResourceData) *diag.Diagnostic
	RemoteToSchemaFunc    func(remote *MType, key string, state *schema.ResourceData) *diag.Diagnostic
	SchemaToRemoteFunc    func(state *schema.ResourceData, key string, remote *MType)
}

func (pfmb *PluggableFiledMapperBase[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	if pfmb.RemoteToSchemaFunc != nil {
		return pfmb.RemoteToSchemaFunc(remote, pfmb.Key, state)
	} else {
		return nil
	}
}

func (pfmb *PluggableFiledMapperBase[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	if pfmb.NilRemoteToSchemaFunc != nil {
		return pfmb.NilRemoteToSchemaFunc(pfmb.Key, state)
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

func (fmb *FieldMapperBase[MType]) GetKey() string {
	return fmb.Key
}

func (fmb *FieldMapperBase[MType]) GetSchema() *schema.Schema {
	return fmb.Schema
}

type IdentityMapper[Ident any] interface {
	Identity(state *schema.ResourceData) (Ident, error)
	Assign(ident Ident, state *schema.ResourceData) error
}

type ParentIdentityMapper[Ident any] interface {
	IdentityMapper[Ident]

	GetKey() string
	GetSchema() *schema.Schema
	ValidateIdent(interface{}, cty.Path) diag.Diagnostics
}

type JsonIdentityMapper[Ident any] struct {
	Key               string
	Schema            schema.Schema
	IdentityFunc      func() Ident
	ValidateIdentFunc func(inp Ident) bool
}

func (im *JsonIdentityMapper[Ident]) PrepareParentMapper() *JsonIdentityMapper[Ident] {
	im.Schema.ValidateDiagFunc = im.ValidateIdent

	return im
}

func (im *JsonIdentityMapper[Ident]) Validate(ident Ident) bool {
	return im.ValidateIdentFunc(ident)
}

func (im *JsonIdentityMapper[Ident]) GetKey() string {
	return im.Key
}

func (im *JsonIdentityMapper[Ident]) GetSchema() *schema.Schema {
	return &im.Schema
}

func (im *JsonIdentityMapper[Ident]) Identity(state *schema.ResourceData) (Ident, error) {
	rv := im.IdentityFunc()
	if len(im.Key) > 0 {
		v := mashschema.ExtractString(state, im.Key, "")
		err := UnwrapJSON(v, &rv)
		return rv, err
	} else {

		err := UnwrapJSON(state.Id(), &rv)
		return rv, err
	}
}

func (im *JsonIdentityMapper[Ident]) ValidateIdent(i interface{}, _ cty.Path) diag.Diagnostics {
	rv := diag.Diagnostics{}

	if str, ok := i.(string); ok {
		ident := im.IdentityFunc()
		if err := UnwrapJSON(str, &ident); err != nil {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("supplied value is not a valid wrapped json"),
				Detail:   fmt.Sprintf("could not parse supplied value as type %s", reflect.TypeOf(ident).Name()),
			})
		} else if !im.ValidateIdentFunc(ident) {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("supplied identity is not valid"),
			})
		}
	} else {
		rv = append(rv, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("supplied value for field %s is not string", im.Key),
		})
	}

	return rv
}

func (im *JsonIdentityMapper[Ident]) Assign(ident Ident, state *schema.ResourceData) error {
	val := WrapJSON(ident)
	if len(im.Key) > 0 {
		return state.Set(im.Key, val)
	} else {
		state.SetId(val)
		return nil
	}
}

func UnwrapJSON(inp string, receiver interface{}) error {
	if sEnc, err := base64.StdEncoding.DecodeString(inp); err != nil {
		return err
	} else {
		return json.Unmarshal(sEnc, receiver)
	}
}

func WrapJSON(origin interface{}) string {
	v, _ := json.Marshal(origin)
	return base64.StdEncoding.EncodeToString(v)
}
