package mashschemag

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var ServiceErrorSetResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.ErrorSetIdentifier]{
		IdentityFunc: func() masherytypes.ErrorSetIdentifier {
			return masherytypes.ErrorSetIdentifier{}
		},
	})

//go:embed service_error_test_set_default.json
var defaultErrorSetContent []byte
var defaultErrorMessages map[string]masherytypes.MasheryErrorMessage

func isDefaultMessage(message *masherytypes.MasheryErrorMessage) bool {
	if v, ok := defaultErrorMessages[message.Id]; ok {
		return v.Code == message.Code &&
			v.Status == message.Status &&
			v.DetailHeader == message.DetailHeader &&
			v.ResponseBody == message.ResponseBody
	}

	return false
}

func init() {
	var msgArr []masherytypes.MasheryErrorMessage
	if e := json.Unmarshal(defaultErrorSetContent, &msgArr); e != nil {
		panic(fmt.Sprintf("cannot parse default messages: %s", e.Error()))
	}

	defaultErrorMessages = make(map[string]masherytypes.MasheryErrorMessage, len(msgArr))
	for _, m := range msgArr {
		defaultErrorMessages[m.Id] = m
	}
}

func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.ServiceIdentifier]{
		Key: mashschema.MashSvcId,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Service Id, to which this OAuth security profile belongs",
		},
		IdentityFunc: func() masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.ServiceIdentifier) bool {
			return len(inp.ServiceId) > 0
		},
	}

	ServiceErrorSetResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

func init() {
	ServiceErrorSetResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.ErrorSet]{
		Locator: func(in *masherytypes.ErrorSet) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ErrorSet]{
			Key: mashschema.MashSvcErrorSetId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 service identifier of this error set",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.ErrorSet]{
		Locator: func(in *masherytypes.ErrorSet) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ErrorSet]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.ErrorSet]{
		Locator: func(in *masherytypes.ErrorSet) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ErrorSet]{
			Key: mashschema.MashPackUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was updated",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.ErrorSet]{
		Locator: func(in *masherytypes.ErrorSet) *string {
			return &in.Name
		},

		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ErrorSet]{
			Key: mashschema.MashObjName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this error set",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.ErrorSet]{
		Locator: func(in *masherytypes.ErrorSet) *string {
			return &in.Type
		},

		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ErrorSet]{
			Key: mashschema.MashSvcErrorSetType,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of this error set",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.ErrorSet]{
		Locator: func(in *masherytypes.ErrorSet) *bool {
			return &in.JSONP
		},

		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ErrorSet]{
			Key: mashschema.MashSvcErrorSetJsonp,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable the error return format as JSONP. All JSONP HTTP response codes will automatically be 200",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.ErrorSet]{
		Locator: func(in *masherytypes.ErrorSet) *string {
			return &in.JSONPType
		},

		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ErrorSet]{
			Key: mashschema.MashSvcErrorSetJsonpType,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The error format to use if the JSONP callback parameter is not available and the jsonp option is enabled",
			},
		},
	})
}

func ErrorMessageToTerraform(inp masherytypes.MasheryErrorMessage) map[string]interface{} {
	return map[string]interface{}{
		mashschema.MashSvcErrorSetMessageErrorId:      inp.Id,
		mashschema.MashSvcErrorSetMessageStatus:       inp.Status,
		mashschema.MashSvcErrorSetMessageDetailHeader: inp.DetailHeader,
		mashschema.MashSvcErrorSetMessageResponseBody: inp.ResponseBody,
	}
}

// Mapper of the error messages
func init() {
	ServiceErrorSetResourceSchemaBuilder.Add(&tfmapper.PluggableFiledMapperBase[masherytypes.ErrorSet]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ErrorSet]{
			Key: mashschema.MashSvcErrorSetMessage,
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				if rw, ok := in.GetOk(key); ok {
					stateSet := mashschema.UnwrapStructArrayFromTerraformSet(rw)

					for _, mp := range stateSet {
						rv := schemaToErrorSetMessage(mp)
						if isDefaultMessage(&rv) {
							return false, fmt.Sprintf("error message for %s matches Mashery-default", rv.Id)
						} else if mashschema.FindInArray(rv.Id, &mashschema.ValidErrorSetMessageId) < 0 {
							return false, fmt.Sprintf("unsupported error code: %s", rv.Id)
						}

					}
				}

				return true, ""
			},

			Schema: &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Set:      hashMasheryErrorMessage,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mashschema.MashSvcErrorSetMessageErrorId: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Error Id of this message",
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								return mashschema.ValidateStringValueInSet(i, path, &mashschema.ValidErrorSetMessageId)
							},
						},
						mashschema.MashSvcErrorSetMessageStatus: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "String passed in the status field",
						},
						mashschema.MashSvcErrorSetMessageDetailHeader: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Detailed header contents",
						},
						mashschema.MashSvcErrorSetMessageResponseBody: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Response body value",
						},
					},
				},
			},
		},
		NilRemoteToSchemaFunc: func(key string, state *schema.ResourceData) *diag.Diagnostic {
			if err := state.Set(key, []map[string]interface{}{}); err != nil {
				return &diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("unable to set machery error messages: %s", err.Error()),
				}
			} else {
				return nil
			}
		},
		RemoteToSchemaFunc: func(remote *masherytypes.ErrorSet, key string, state *schema.ResourceData) *diag.Diagnostic {
			var val []interface{}

			if remote.ErrorMessages != nil {
				for _, msg := range *remote.ErrorMessages {
					if !isDefaultMessage(&msg) {
						val = append(val, ErrorMessageToTerraform(msg))
					}
				}
			}

			if err := state.Set(key, val); err != nil {
				return &diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("failed to set key %s due to error: %s", key, err.Error()),
				}
			} else {
				return nil
			}
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.ErrorSet) {
			output := make(map[string]masherytypes.MasheryErrorMessage, len(defaultErrorMessages))
			// Do we need defaults?
			//for k, v := range defaultErrorMessages {
			//	output[k] = v
			//}

			if d, ok := state.GetOk(key); ok {
				stateSet := mashschema.UnwrapStructArrayFromTerraformSet(d)

				for _, mp := range stateSet {
					rv := schemaToErrorSetMessage(mp)
					output[rv.Id] = rv
				}
			}

			flattenedArray := make([]masherytypes.MasheryErrorMessage, len(output))
			idx := 0
			for _, v := range output {
				flattenedArray[idx] = v
				idx++
			}

			remote.ErrorMessages = &flattenedArray
		},
	})
}

func schemaToErrorSetMessage(mp map[string]interface{}) masherytypes.MasheryErrorMessage {
	rv := masherytypes.MasheryErrorMessage{}

	mashschema.ExtractKeyFromMap(mp, mashschema.MashSvcErrorSetMessageErrorId, &rv.Id)
	if mashschema.ErrorParsePattern.MatchString(rv.Id) {
		match := mashschema.ErrorParsePattern.FindStringSubmatch(rv.Id)
		i, _ := strconv.ParseInt(match[1], 10, 32)
		rv.Code = int(i)
	}

	mashschema.ExtractKeyFromMap(mp, mashschema.MashSvcErrorSetMessageStatus, &rv.Status)
	mashschema.ExtractKeyFromMap(mp, mashschema.MashSvcErrorSetMessageDetailHeader, &rv.DetailHeader)
	mashschema.ExtractKeyFromMap(mp, mashschema.MashSvcErrorSetMessageResponseBody, &rv.ResponseBody)

	return rv
}

func hashMasheryErrorMessage(i interface{}) int {
	val := i.(map[string]interface{})
	if errId, ok := val[mashschema.MashSvcErrorSetMessageErrorId]; ok {
		return mashschema.StringHashcode(errId)
	} else {
		return mashschema.StringHashcode("")
	}
}
