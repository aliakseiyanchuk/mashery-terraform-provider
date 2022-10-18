package mashschema

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type SchemaBuilder struct {
	schema *TFResourceSchema
}

func (sb *SchemaBuilder) AddComputedString(key, desc string) *SchemaBuilder {
	(*sb.schema)[key] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: desc,
	}
	return sb
}

func (sb *SchemaBuilder) AddComputedBoolean(key, desc string) *SchemaBuilder {
	(*sb.schema)[key] = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: desc,
	}
	return sb
}

func (sb *SchemaBuilder) AddComputedOptionalString(key, desc string) *SchemaBuilder {
	(*sb.schema)[key] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Optional:    true,
		Description: desc,
	}
	return sb
}

func (sb *SchemaBuilder) AddComputedOptionalInt(key, desc string) *SchemaBuilder {
	(*sb.schema)[key] = &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Optional:    true,
		Description: desc,
	}
	return sb
}

func (sb *SchemaBuilder) AddOptionalString(key, desc string) *SchemaBuilder {
	(*sb.schema)[key] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: desc,
	}
	return sb
}

func (sb *SchemaBuilder) AddRequiredString(key, desc string) *SchemaBuilder {
	(*sb.schema)[key] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: desc,
	}
	return sb
}

func (sb *SchemaBuilder) AddRequiredInt(key, desc string) *SchemaBuilder {
	(*sb.schema)[key] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: desc,
	}
	return sb
}

func (sb *SchemaBuilder) AddOptionalBoolean(key, desc string) *SchemaBuilder {
	(*sb.schema)[key] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: desc,
	}
	return sb
}
