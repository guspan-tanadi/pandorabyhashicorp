package generator

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func (s *ServiceGenerator) models(data ServiceGeneratorData) error {
	if len(data.models) == 0 {
		return nil
	}

	for modelName, model := range data.models {
		// model_{name}.go
		// arguably we could enhance this to make `MyThing` be `my_thing` but this is fine for now
		fileName := fmt.Sprintf("model_%s.go", strings.ToLower(modelName))
		gen := modelsTemplater{
			name:  modelName,
			model: model,
		}
		if err := s.writeToPath(data.outputPath, fileName, gen, data); err != nil {
			return fmt.Errorf("templating model for %q: %+v", modelName, err)
		}
	}

	return nil
}

var _ templater = modelsTemplater{}

type modelsTemplater struct {
	name  string
	model resourcemanager.ModelDetails
}

func (c modelsTemplater) template(data ServiceGeneratorData) (*string, error) {
	structCode, err := c.structCode(data)
	if err != nil {
		return nil, fmt.Errorf("generating struct code: %+v", err)
	}

	methods, err := c.methods(data)
	if err != nil {
		return nil, fmt.Errorf("generating functions: %+v", err)
	}

	template := fmt.Sprintf(`package %[1]s

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

%[2]s
%[3]s
`, data.packageName, *structCode, *methods)
	return &template, nil
}

func (c modelsTemplater) structCode(data ServiceGeneratorData) (*string, error) {
	// if this is an Abstract/Type Hint, we output an Interface with a manual unmarshal func that gets called wherever it's used
	if c.model.TypeHintIn != nil && c.model.ParentTypeName == nil {
		out := fmt.Sprintf(`type %[1]s interface {
}`, c.name)
		return &out, nil
	}

	fields := make([]string, 0)
	for fieldName := range c.model.Fields {
		fields = append(fields, fieldName)
	}
	sort.Strings(fields)

	structLines := make([]string, 0)
	for _, fieldName := range fields {
		fieldDetails := c.model.Fields[fieldName]
		fieldTypeName := "FIXME"
		fieldTypeVal, err := golangTypeNameForObjectDefinition(fieldDetails.ObjectDefinition)
		if err != nil {
			return nil, fmt.Errorf("determining type information for %q: %+v", fieldName, err)
		}
		fieldTypeName = *fieldTypeVal

		structLine, err := c.structLineForField(fieldName, fieldTypeName, fieldDetails, data)
		if err != nil {
			return nil, err
		}

		if c.model.TypeHintIn != nil && *c.model.TypeHintIn == fieldName {
			// this isn't user configurable (and is hard-coded) so there's no point outputting this
			continue
		}

		structLines = append(structLines, *structLine)
	}

	out := fmt.Sprintf(`
type %[1]s struct {
%[2]s
}
`, c.name, strings.Join(structLines, "\n"))
	return &out, nil
}

func (c modelsTemplater) methods(data ServiceGeneratorData) (*string, error) {
	code := make([]string, 0)

	dateFunctions, err := c.codeForDateFunctions()
	if err != nil {
		return nil, fmt.Errorf("generating date functions: %+v", err)
	}
	code = append(code, *dateFunctions)

	marshalFunctions, err := c.codeForMarshalFunctions()
	if err != nil {
		return nil, fmt.Errorf("generating marshal functions: %+v", err)
	}
	code = append(code, *marshalFunctions)

	unmarshalFunctions, err := c.codeForUnmarshalFunctions(data)
	if err != nil {
		return nil, fmt.Errorf("generating unmarshal functions: %+v", err)
	}
	code = append(code, *unmarshalFunctions)

	// TODO: validation functions (#58)

	//requiresMarshalFunctionForImplementation := c.model.TypeHintIn != nil && c.model.TypeHintValue != nil
	//requiresUnmarshalFunctionForImplementation := false

	////requiresImplementationUnmarshalFunc := false
	//for fieldName, fieldDetails := range c.model.Fields {
	//	topLevelObject := topLevelObjectDefinition(fieldDetails.ObjectDefinition)
	//	if topLevelObject.Type == resourcemanager.ReferenceApiObjectDefinitionType {
	//		// if this is a Model, is that Model a TypeHint? (The Ref. can also be a Constant, so 'ok' can be false)
	//		if model, ok := data.models[*topLevelObject.ReferenceName]; ok && model.TypeHintIn != nil {
	//			a = true
	//		}
	//	}
	//
	//	// TODO: implement this
	//	//if fieldDetails.IsTypeHint {
	//	//	requiresUnmarshalFunc = true
	//	//}
	//}

	// TODO: implement this (from above)
	//if requiresImplementationUnmarshalFunc {
	//
	//}

	//if requiresUnmarshalFunc {
	//	unmarshalFunc, err := c.unmarshalerFunc()
	//	if err != nil {
	//		return nil, err
	//	}
	//	methods = append(methods, *unmarshalFunc)
	//}
	//
	//if c.model.TypeHintIn != nil && c.model.ParentTypeName != nil {
	//	// this'll only be an implementation/child so this can generate the marshaler func
	//	marshalFunc, err := c.marshalFuncForChildClass()
	//	if err != nil {
	//		return nil, err
	//	}
	//	methods = append(methods, *marshalFunc)
	//}

	output := strings.Join(code, "\n")
	return &output, nil
}

func (c modelsTemplater) structLineForField(fieldName, fieldType string, fieldDetails resourcemanager.FieldDetails, data ServiceGeneratorData) (*string, error) {
	jsonDetails := fieldDetails.JsonName

	isOptional := false
	if fieldDetails.Optional {
		isOptional = true

		// however if the immediate (not top-level) object definition is a Reference to a Parent it's Optional
		// by default since Parent types are output as an interface (which is implied nullable)
		if fieldDetails.ObjectDefinition.Type == resourcemanager.ReferenceApiObjectDefinitionType {
			model := data.models[*fieldDetails.ObjectDefinition.ReferenceName]
			if model.TypeHintIn != nil && model.ParentTypeName == nil {
				isOptional = false
			}
		}
	}

	if isOptional {
		fieldType = fmt.Sprintf("*%s", fieldType)
		jsonDetails += ",omitempty"
	}

	line := fmt.Sprintf("\t%s %s `json:\"%s\"`", fieldName, fieldType, jsonDetails)
	return &line, nil
}

func dateFormatString(input resourcemanager.DateFormat) string {
	switch input {
	case resourcemanager.RFC3339:
		return time.RFC3339

	case resourcemanager.RFC3339Nano:
		return time.RFC3339Nano

	default:
		panic(fmt.Errorf("unsupported date format %q", string(input)))
	}
}

func (c modelsTemplater) codeForDateFunctions() (*string, error) {
	fieldsRequiringDateFunctions := make([]string, 0)
	for fieldName, fieldDetails := range c.model.Fields {
		if fieldDetails.DateFormat != nil {
			fieldsRequiringDateFunctions = append(fieldsRequiringDateFunctions, fieldName)
		}
	}

	sort.Strings(fieldsRequiringDateFunctions)
	lines := make([]string, 0)
	for _, fieldName := range fieldsRequiringDateFunctions {
		fieldDetails := c.model.Fields[fieldName]

		dateFormat := dateFormatString(*fieldDetails.DateFormat)

		lines := []string{
			fmt.Sprintf("\tfunc (o %[1]s) Get%[2]sAsTime() (*time.Time, error) {", c.name, fieldName),
		}

		// Get{Name}AsTime method for getting *time.Time from a string
		if fieldDetails.Optional {
			lines = append(lines, fmt.Sprintf("\t\tif o.%s == nil {", fieldName))
			lines = append(lines, fmt.Sprintf("\t\t\treturn nil, nil"))
			lines = append(lines, fmt.Sprintf("\t\t}"))
			lines = append(lines, fmt.Sprintf("\t\treturn formatting.ParseAsDateFormat(o.%s, %q)", fieldName, dateFormat))
		} else {
			lines = append(lines, fmt.Sprintf("\t\treturn formatting.ParseAsDateFormat(&o.%s, %q)", fieldName, dateFormat))
		}

		lines = append(lines, fmt.Sprintf("\t}\n"))

		// Set{Name}AsTime method - for setting time.Time -> string
		lines = append(lines, fmt.Sprintf("\tfunc (o %[1]s) Set%[2]sAsTime(input time.Time) {", c.name, fieldName))
		lines = append(lines, fmt.Sprintf("\t\tformatted := input.Format(%q)", dateFormat))
		if fieldDetails.Optional {
			lines = append(lines, fmt.Sprintf("\t\to.%s = &formatted", fieldName))
		} else {
			lines = append(lines, fmt.Sprintf("\t\to.%s = formatted", fieldName))
		}
		lines = append(lines, fmt.Sprintf("\t}\n"))
	}

	output := strings.Join(lines, "\n")
	return &output, nil
}

func (c modelsTemplater) codeForMarshalFunctions() (*string, error) {
	output := ""

	if c.model.TypeHintValue != nil {
		if c.model.TypeHintIn == nil {
			return nil, fmt.Errorf("model %q must contain a TypeHintIn when a TypeHintValue is present", c.name)
		}

		output = fmt.Sprintf(`
var _ json.Marshaler = %[1]s{}

func (s %[1]s) MarshalJSON() ([]byte, error) {
	type wrapper %[1]s
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling %[1]s: %%+v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling %[1]s: %%+v", err)
	}
	decoded[%[2]q] = %[3]q

	encoded, err = json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling %[1]s: %%+v", err)
	}

	return encoded, nil
}
`, c.name, *c.model.TypeHintIn, *c.model.TypeHintValue)
	}

	return &output, nil
}

func (c modelsTemplater) codeForUnmarshalFunctions(data ServiceGeneratorData) (*string, error) {
	unmarshalFunction, err := c.codeForUnmarshalStructFunction(data)
	if err != nil {
		return nil, fmt.Errorf("generating code for unmarshal struct function: %+v", err)
	}

	parentFunction, err := c.codeForUnmarshalParentFunction(data)
	if err != nil {
		return nil, fmt.Errorf("generating code for unmarshal parent function: %+v", err)
	}

	output := fmt.Sprintf(`
%s
%s
`, *unmarshalFunction, *parentFunction)
	return &output, nil
}

func (c modelsTemplater) codeForUnmarshalParentFunction(data ServiceGeneratorData) (*string, error) {
	// if this is a Discriminated Type (e.g. Parent) then we need to generate a unmarshal{Name}Implementations
	// function which can be used in any usages
	lines := make([]string, 0)
	if c.model.TypeHintIn != nil && c.model.ParentTypeName == nil {
		modelsImplementingThisClass := make([]string, 0)
		for modelName, model := range data.models {
			if model.ParentTypeName == nil || model.TypeHintIn == nil || model.TypeHintValue == nil || modelName == c.name {
				continue
			}

			// sanity-checking
			if *model.ParentTypeName != c.name {
				continue
			}

			if *model.TypeHintIn != *c.model.TypeHintIn {
				return nil, fmt.Errorf("implementation %q uses a different discriminated field (%q) than parent %q (%q)", modelName, *model.TypeHintIn, c.name, *c.model.TypeHintIn)
			}

			modelsImplementingThisClass = append(modelsImplementingThisClass, modelName)
		}

		// sanity-checking
		if len(modelsImplementingThisClass) == 0 {
			return nil, fmt.Errorf("model %q is a discriminated parent type with no implementations", c.name)
		}
		jsonFieldName := c.model.Fields[*c.model.TypeHintIn].JsonName
		lines = append(lines, fmt.Sprintf(`
func unmarshal%[1]sImplementation(input []byte) (%[1]s, error) {
	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling %[1]s into map[string]interface: %%+v", err)
	}

	value, ok := temp[%[2]q].(string)
	if !ok {
		return nil, fmt.Errorf("missing field '%[2]s' needed to discriminate %[1]s type")
	}
`, c.name, jsonFieldName))

		sort.Strings(modelsImplementingThisClass)
		for _, implementationName := range modelsImplementingThisClass {
			model := data.models[implementationName]

			lines = append(lines, fmt.Sprintf(`
	if strings.EqualFold(value, %[1]q) {
		var out %[2]s
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into %[2]s: %%+v", err)
		}
		return out, nil
	}
`, *model.TypeHintValue, implementationName))
		}

		// if it doesn't match - we generate and deserialize into a 'Raw{Name}Impl' type - named intentionally
		// so that we don't conflict with a generated 'Raw{Name}' type which exists in a handful of Swaggers
		jsonIgnoreTag := "`json:\"-\"`"
		lines = append(lines, fmt.Sprintf(`
	type Raw%[1]sImpl struct {
		Type string %[2]s
		Values map[string]interface{} %[2]s
	}
	out := Raw%[1]sImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil
`, c.name, jsonIgnoreTag, *c.model.TypeHintIn))

		lines = append(lines, "}")
	}

	output := strings.Join(lines, "\n")
	return &output, nil
}

func (c modelsTemplater) codeForUnmarshalStructFunction(data ServiceGeneratorData) (*string, error) {
	lines := make([]string, 0)
	// fields either require unmarshaling or can be explicitly assigned, determine which
	fieldsRequiringAssignment := make([]string, 0)
	fieldsRequiringUnmarshalling := make([]string, 0)
	for fieldName, fieldDetails := range c.model.Fields {
		topLevelObject := topLevelObjectDefinition(fieldDetails.ObjectDefinition)
		if topLevelObject.Type == resourcemanager.ReferenceApiObjectDefinitionType {
			model, ok := data.models[*topLevelObject.ReferenceName]
			if ok && model.TypeHintIn != nil {
				fieldsRequiringUnmarshalling = append(fieldsRequiringUnmarshalling, fieldName)
				continue
			}
		}

		fieldsRequiringAssignment = append(fieldsRequiringAssignment, fieldName)
	}

	// we only need a custom unmarshal function when there's something needing unmarshaling
	// else the default unmarshaler will be fine
	if len(fieldsRequiringUnmarshalling) > 0 {
		lines = append(lines, fmt.Sprintf(`
var _ json.Unmarshaler = &%[1]s{}

func (s *%[1]s) UnmarshalJSON(bytes []byte) error {`, c.name))

		// first for each regular field, decode & assign that
		if len(fieldsRequiringAssignment) > 0 {
			lines = append(lines, fmt.Sprintf(`type alias %[1]s
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into %[1]s: %%+v", err)
	}
`, c.name))

			sort.Strings(fieldsRequiringAssignment)
			for _, fieldName := range fieldsRequiringAssignment {
				lines = append(lines, fmt.Sprintf("s.%[1]s = decoded.%[1]s", fieldName))
			}
		}

		lines = append(lines, fmt.Sprintf(`
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling %[1]s into map[string]json.RawMessage: %%+v", err)
	}
`, c.name))

		sort.Strings(fieldsRequiringUnmarshalling)
		for _, fieldName := range fieldsRequiringUnmarshalling {
			fieldDetails := c.model.Fields[fieldName]
			topLevelObjectDef := topLevelObjectDefinition(fieldDetails.ObjectDefinition)

			supportedDiscriminatorWrappers := map[resourcemanager.ApiObjectDefinitionType]struct{}{
				resourcemanager.DictionaryApiObjectDefinitionType: {},
				resourcemanager.ListApiObjectDefinitionType:       {},
				resourcemanager.ReferenceApiObjectDefinitionType:  {},
			}
			if _, supported := supportedDiscriminatorWrappers[fieldDetails.ObjectDefinition.Type]; !supported {
				return nil, fmt.Errorf("discriminators can only be unwrapped for Dictionaries, Lists and References but got %q for field %q in model %q", fieldDetails.ObjectDefinition.Type, fieldName, c.name)
			}

			if fieldDetails.ObjectDefinition.Type == resourcemanager.DictionaryApiObjectDefinitionType {
				if fieldDetails.ObjectDefinition.NestedItem == nil {
					return nil, fmt.Errorf("dictionaries of discriminators require a NestedItem but didn't get one for field %q in model %q", fieldName, c.name)
				}
				if fieldDetails.ObjectDefinition.NestedItem.Type != resourcemanager.ReferenceApiObjectDefinitionType {
					return nil, fmt.Errorf("dictionaries of discriminators only support a single level deep but got a non-Reference type for field %q in model %q", fieldName, c.name)
				}

				// if the Dictionary is optional, we need to assign the pointer value
				assignmentPrefix := ""
				if fieldDetails.Optional {
					assignmentPrefix = "&"
				}

				// TODO: handle the Dictionary Element being Optional if necessary
				lines = append(lines, fmt.Sprintf(`
	if v, ok := temp[%[5]q]; ok {
		var dictionaryTemp map[string]json.RawMessage
		if err := json.Unmarshal(v, &dictionaryTemp); err != nil {
			return fmt.Errorf("unmarshaling %[1]s into dictionary map[string]json.RawMessage: %%+v", err)
		}

		output := make(map[string]%[2]s)
		for key, val := range dictionaryTemp {
			impl, err := unmarshal%[2]sImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling key %%q field '%[1]s' for '%[3]s': %%+v", key, err)
			}
			output[key] = impl
		}
		s.%[1]s = %[4]soutput
	}`, fieldName, *topLevelObjectDef.ReferenceName, c.name, assignmentPrefix, fieldDetails.JsonName))
			}

			if fieldDetails.ObjectDefinition.Type == resourcemanager.ListApiObjectDefinitionType {
				if fieldDetails.ObjectDefinition.NestedItem == nil {
					return nil, fmt.Errorf("lists of discriminators require a NestedItem but didn't get one for field %q in model %q", fieldName, c.name)
				}
				if fieldDetails.ObjectDefinition.NestedItem.Type != resourcemanager.ReferenceApiObjectDefinitionType {
					return nil, fmt.Errorf("lists of discriminators only support a single level deep but got a non-Reference type for field %q in model %q", fieldName, c.name)
				}

				// if the List is optional, we need to assign the pointer value
				assignmentPrefix := ""
				if fieldDetails.Optional {
					assignmentPrefix = "&"
				}

				// TODO: handle the List Element being Optional if necessary

				lines = append(lines, fmt.Sprintf(`
	if v, ok := temp[%[5]q]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling %[1]s into list []json.RawMessage: %%+v", err)
		}

		output := make([]%[2]s, 0)
		for i, val := range listTemp {
			impl, err := unmarshal%[2]sImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %%d field '%[1]s' for '%[3]s': %%+v", i, err)
			}
			output = append(output, impl)
		}
		s.%[1]s = %[4]soutput
	}`, fieldName, *topLevelObjectDef.ReferenceName, c.name, assignmentPrefix, fieldDetails.JsonName))
			}

			if fieldDetails.ObjectDefinition.Type == resourcemanager.ReferenceApiObjectDefinitionType {
				lines = append(lines, fmt.Sprintf(`
	if v, ok := temp[%[4]q]; ok {
		impl, err := unmarshal%[2]sImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field '%[1]s' for '%[3]s': %%+v", err)
		}
		s.%[1]s = impl
	}`, fieldName, *topLevelObjectDef.ReferenceName, c.name, fieldDetails.JsonName))
			}
		}

		lines = append(lines, "return nil")
		lines = append(lines, "}")
	}

	output := strings.Join(lines, "\n")
	return &output, nil
}
