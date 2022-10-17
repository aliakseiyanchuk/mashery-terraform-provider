package mashschema

import "github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"

//--------------------------
// Strings

func CoerceStringArrayToInterfaceType(in []string) []interface{} {
	s := make([]interface{}, len(in))
	for i, v := range in {
		s[i] = v
	}

	return s
}

func CoerceInterfaceArrayToStringArray(in []interface{}) []string {
	s := make([]string, len(in))
	for i, v := range in {
		s[i] = v.(string)
	}

	return s
}

//--------------------------
// Roles

func CoerceRolesArrayToInterfaceType(in []masherytypes.Role) []interface{} {
	s := make([]interface{}, len(in))
	for i, v := range in {
		s[i] = v
	}

	return s
}

func CoerceInterfaceArrayToRolesArray(in []interface{}) []masherytypes.Role {
	s := make([]masherytypes.Role, len(in))
	for i, v := range in {
		s[i] = v.(masherytypes.Role)
	}

	return s
}

//--------------------------
// Mashery email template set

func CoerceEmailTemplateSetArrayToInterfaceType(in []masherytypes.EmailTemplateSet) []interface{} {
	s := make([]interface{}, len(in))
	for i, v := range in {
		s[i] = v
	}

	return s
}
