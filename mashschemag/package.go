package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform/helper/hashcode"
	"strings"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var PackageResourceSchemaBuilder = tfmapper.NewSchemaBuilder[tfmapper.Orphan, masherytypes.PackageIdentifier, masherytypes.Package]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.PackageIdentifier]{
		IdentityFunc: func() masherytypes.PackageIdentifier {
			return masherytypes.PackageIdentifier{}
		},
	}).
	RootIdentity(&tfmapper.RootParentIdentity{})

// ------------------------------------------------------------------------------------------------------------------
// Field mappers
// ------------------------------------------------------------------------------------------------------------------

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackageId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 package identifier of this package",
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.DateMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was updated",
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.SerOrPrefixedFieldMapper[masherytypes.Package]{
		StringFieldMapper: tfmapper.StringFieldMapper[masherytypes.Package]{
			Locator: func(in *masherytypes.Package) *string {
				return &in.Name
			},

			FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
				Key: mashschema.MashObjName,
			},
		},

		PrefixKey: mashschema.MashObjNamePrefix,
		CompositeMapperBase: tfmapper.CompositeMapperBase{
			CompositeSchema: map[string]*schema.Schema{
				mashschema.MashObjName: {
					Type:          schema.TypeString,
					Optional:      true,
					Computed:      true,
					Description:   "Package name",
					ConflictsWith: []string{mashschema.MashObjNamePrefix},
				},
				mashschema.MashObjNamePrefix: {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "Prefix for the package name",
					ConflictsWith: []string{mashschema.MashObjName},
				},
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *string {
			return &in.Description
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackDescription,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Managed by Terraform",
				Description: "Package description",
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *string {
			return &in.NotifyDeveloperPeriod
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyDeveloperPeriod,
			Schema: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "hour",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &mashschema.NotifyDeveloperPeriodEnum)
				},
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *bool {
			return &in.NotifyDeveloperNearQuota
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyDeveloperNearQuota,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Notify developer when approaching quota",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *bool {
			return &in.NotifyDeveloperOverQuota
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyDeveloperOverQuota,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Notify developer when quota exceeded",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *bool {
			return &in.NotifyDeveloperOverThrottle
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyDeveloperOverThrottle,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Notify developer when throttle exceeded",
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *string {
			return &in.NotifyAdminPeriod
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyAdminPeriod,
			Schema: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "hour",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &mashschema.NotifyDeveloperPeriodEnum)
				},
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *bool {
			return &in.NotifyAdminNearQuota
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyAdminNearQuota,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Notify admin when approaching quota",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *bool {
			return &in.NotifyAdminOverQuota
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyAdminOverQuota,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Notify admin when quota exceeded",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *bool {
			return &in.NotifyAdminOverThrottle
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyAdminOverThrottle,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Notify admin when throttle exceeded",
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Package]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNotifyAdminEmails,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Email addresses to send admin notifications",
				// TODO: Reference to the string element set is repeated
				// It could be placed in a shared module.
				Set: func(i interface{}) int {
					return hashcode.String(i.(string))
				},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.Package) {
			val := strings.Join(mashschema.ExtractStringArray(state, key, &mashschema.EmptyStringArray), ",")
			remote.NotifyAdminEmails = val
		},
		RemoteToSchemaFunc: func(remote *masherytypes.Package, key string, state *schema.ResourceData) *diag.Diagnostic {
			if len(remote.NotifyAdminEmails) == 0 {
				var emptyArray []string

				_ = state.Set(key, emptyArray)
				return nil
			}

			split := strings.Split(remote.NotifyAdminEmails, ",")
			rv := make([]interface{}, len(split))
			for i, v := range split {
				rv[i] = v
			}

			return tfmapper.SetKeyWithDiag(state, key, rv)
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.IntPointerFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) **int {
			return &in.NearQuotaThreshold
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackNearQuotaThreshold,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Percentage of quota when approaching limit notifications will be sent",
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.EAVFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) **masherytypes.EAV {
			return &in.Eav
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackEAVs,
			Schema: &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Extended attribute values",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) *string {
			return &in.KeyAdapter
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackKeyAdapter,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom adapter for key generation",
			},
		},
	})
}

func init() {
	PackageResourceSchemaBuilder.Add(&tfmapper.IntPointerFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) **int {
			return &in.KeyLength
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackKeyLength,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Length of keys for this package",
			},
		},
	}).Add(&tfmapper.IntPointerFieldMapper[masherytypes.Package]{
		Locator: func(in *masherytypes.Package) **int {
			return &in.SharedSecretLength
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Package]{
			Key: mashschema.MashPackSharedSecretLength,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Length of shared secret for this package",
			},
		},
	})
}
