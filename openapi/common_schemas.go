package openapi

var StringSchema = &Schema{
	Type: StringType,
}

var IntegerSchema = &Schema{
	Type: IntegerType,
}

var UUIDSchema = &Schema{
	Type:   StringType,
	Format: "uuid",
}

func RegExpSchema(pattern string) *Schema {
	return &Schema{
		Type:    StringType,
		Pattern: pattern,
	}
}
