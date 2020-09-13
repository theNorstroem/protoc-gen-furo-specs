package generator

import (
	"github.com/iancoleman/strcase"
	"github.com/theNorstroem/protoc-gen-furo-specs/pkg/protoast"
	"github.com/theNorstroem/spectools/pkg/orderedmap"
	"github.com/theNorstroem/spectools/pkg/specSpec"
	"github.com/theNorstroem/spectools/pkg/specSpec/furo"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"gopkg.in/yaml.v3"
	"path"
	"strings"
)

func Generate(protoAST *protoast.ProtoAST) error {

	for protofilename, descriptor := range protoAST.ProtoMap {

		// in this example we want to generate a patch.go file for every type in the protos
		for MessageIndex, Message := range descriptor.MessageType {
			if shouldGenerateTypeSpec(protoAST, *Message.Name, descriptor, Message) {
				SourceInfo := protoast.GetSourceInfo(descriptor)
				filename, packagename := FileAndPackageNameToGenerate(descriptor, Message)

				description := packagename + " does not have a description"
				if SourceInfo.Messages[MessageIndex].Info.LeadingComments != nil {
					description = cleanDescription(*SourceInfo.Messages[MessageIndex].Info.LeadingComments)
				}

				typeSpec := specSpec.Type{
					Name:        *Message.Name,
					Type:        *Message.Name,
					Description: description,
					XProto: &specSpec.Typeproto{
						Package:    strings.Join(strings.Split(path.Dir(protofilename), "/"), "."), // package is protofilename with . and is not the go package
						Targetfile: path.Base(protofilename),                                       // is base of packagename
						Imports:    descriptor.Dependency,
						Options:    getProtoOptions(descriptor.Options),
					},
					Fields: getFields(SourceInfo.Messages[MessageIndex]),
				}

				var responseFile pluginpb.CodeGeneratorResponse_File
				responseFile.Name = &filename

				content, _ := yaml.Marshal(typeSpec)
				s := string(content)
				responseFile.Content = &s

				protoAST.Response.File = append(protoAST.Response.File, &responseFile)
			}

		}
	}

	return nil
}

func cleanDescription(s string) string {
	res := s[1 : len(s)-1]
	strings.Replace(s, "\n", "\\n", -1)
	return res
}

func getFields(messageInfo protoast.MessageInfo) *orderedmap.OrderedMap {
	omap := orderedmap.New()
	for _, f := range messageInfo.FieldInfos {

		fielddescription := ""
		if f.Info.LeadingComments != nil {
			fielddescription = cleanDescription(*f.Info.LeadingComments)
		}

		field := specSpec.Field{
			Type:        extractTypeFromField(&f),
			Description: fielddescription,
			XProto: &specSpec.Fieldproto{
				Number: *f.Field.Number,
				Oneof:  "",
			},
			XUi: &specSpec.Fielduiextension{},
			Meta: &furo.FieldMeta{
				Options: &furo.Fieldoption{},
				Label:   "label." + messageInfo.Name + "." + *f.Field.Name,
			},
			Constraints: nil,
		}

		// set repeated, must be false on maps!
		// repeated is in f.Field.Label
		isRepeated := false
		if *f.Field.Label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			isRepeated = !strings.HasPrefix(field.Type, "map<")
		}
		field.Meta.Repeated = isRepeated

		omap.Set(f.Name, field)
	}

	return omap
}

func extractTypeFromField(fieldinfo *protoast.FieldInfo) string {
	// If type_name is set, this need not be set.  If both this and type_name
	// are set, this must be one of TYPE_ENUM, TYPE_MESSAGE or TYPE_GROUP.
	// --> Type *FieldDescriptorProto_Type `protobuf:"varint,5,opt,name=type,enum=google.protobuf.FieldDescriptorProto_Type" json:"type,omitempty"`
	// For message and enum types, this is the name of the type.  If the name
	// starts with a '.', it is fully-qualified.  Otherwise, C++-like scoping
	// rules are used to find the type (i.e. first the nested types within this
	// message are searched, then within the parent, on up to the root
	// namespace).
	// --> TypeName *string `protobuf:"bytes,6,opt,name=type_name,json=typeName" json:"type_name,omitempty"`

	// get primitive types first
	// vendor/google.golang.org/protobuf/types/descriptorpb/descriptor.pb.go Line 54
	field := fieldinfo.Field

	if field.Type != nil {
		t := field.Type.String()
		if !(*field.Type == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE ||
			*field.Type == descriptorpb.FieldDescriptorProto_TYPE_ENUM ||
			*field.Type == descriptorpb.FieldDescriptorProto_TYPE_GROUP) {
			return strings.ToLower(t[5:len(t)])
		}
		// if we have message, we look in Typename
		if *field.Type == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			// check for nested type map<string,xxx>
			if fieldinfo.Message.NestedType == nil {
				// must be type
				f := *field.TypeName
				return f[1:len(f)]
			}
			for _, nested := range fieldinfo.Message.NestedType {
				if nested.Options != nil {
					if *nested.Options.MapEntry {
						if strings.Title(fieldinfo.Name)+"Entry" == *nested.Name {
							// this is a map
							maptype := "not_evaluated"
							if !(*nested.Field[1].Type == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE ||
								*nested.Field[1].Type == descriptorpb.FieldDescriptorProto_TYPE_ENUM ||
								*nested.Field[1].Type == descriptorpb.FieldDescriptorProto_TYPE_GROUP) {
								t := nested.Field[1].Type.String()
								maptype = strings.ToLower(t[5:len(t)])
							} else {
								// can be a message or a primitive
								if *nested.Field[1].Type == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
									// message
									m := *nested.Field[1].TypeName
									maptype = m[1:len(m)]
								}
							}
							return "map<string," + maptype + ">"
						}
					}
				}
			}

			f := *field.TypeName
			return f[1:len(f)]
		}
	}

	return "unknown"
}

// get all known options
func getProtoOptions(options *descriptorpb.FileOptions) map[string]string {
	opts := map[string]string{}

	if options.JavaPackage != nil {
		opts[strcase.ToSnake("JavaPackage ")] = *options.JavaPackage //= {*string | 0xc000011c80} "pro.furo.bigdecimal"
	}

	if options.JavaOuterClassname != nil {
		opts[strcase.ToSnake("JavaOuterClassname")] = *options.JavaOuterClassname // = {*string | 0xc000011c90} "BigdecimalProto"
	}

	if options.JavaMultipleFiles != nil {
		if *options.JavaMultipleFiles {
			opts[strcase.ToSnake("JavaMultipleFiles")] = "true" // = {*bool | 0xc00001cc5f} true
		} else {
			opts[strcase.ToSnake("JavaMultipleFiles")] = "false" // = {*bool | 0xc00001cc5f} true
		}
	}

	if options.JavaStringCheckUtf8 != nil {
		if *options.JavaStringCheckUtf8 {
			opts[strcase.ToSnake("JavaStringCheckUtf8")] = "true" // = {*bool} nil
		} else {
			opts[strcase.ToSnake("JavaStringCheckUtf8")] = "false" // = {*bool | 0xc00001cc5f} true
		}
	}
	if options.GoPackage != nil {
		opts[strcase.ToSnake("GoPackage")] = *options.GoPackage // = {*string | 0xc000011ca0} "github.com/theNorstroem/FuroBaseSpecs/dist/pb/furo/bigdecimal;bigdecimalpb"
	}
	if options.CcGenericServices != nil {
		if *options.CcGenericServices {
			opts[strcase.ToSnake("CcGenericServices")] = "true" // = {*bool} nil
		} else {
			opts[strcase.ToSnake("CcGenericServices")] = "false" // = {*bool | 0xc00001cc5f} true
		}
	}
	if options.JavaGenericServices != nil {
		if *options.JavaGenericServices {
			opts[strcase.ToSnake("JavaGenericServices")] = "true" // = {*bool} nil
		} else {
			opts[strcase.ToSnake("JavaGenericServices")] = "false" // = {*bool | 0xc00001cc5f} true
		}
	}
	if options.PyGenericServices != nil {
		if *options.PyGenericServices {
			opts[strcase.ToSnake("PyGenericServices")] = "true" // = {*bool} nil
		} else {
			opts[strcase.ToSnake("PyGenericServices")] = "false" // = {*bool | 0xc00001cc5f} true
		}
	}
	if options.PhpGenericServices != nil {
		if *options.PhpGenericServices {
			opts[strcase.ToSnake("PhpGenericServices")] = "true" // = {*bool} nil
		} else {
			opts[strcase.ToSnake("PhpGenericServices")] = "false" // = {*bool | 0xc00001cc5f} true
		}
	}
	if options.Deprecated != nil {
		if *options.Deprecated {
			opts[strcase.ToSnake("Deprecated")] = "true" // = {*bool} nil
		} else {
			opts[strcase.ToSnake("Deprecated")] = "false" // = {*bool | 0xc00001cc5f} true
		}
	}
	if options.CcEnableArenas != nil {
		if *options.CcEnableArenas {
			opts[strcase.ToSnake("CcEnableArenas")] = "true" // = {*bool | 0xc00001cc60} true
		} else {
			opts[strcase.ToSnake("CcEnableArenas")] = "false" // = {*bool | 0xc00001cc5f} true
		}
	}

	if options.ObjcClassPrefix != nil {
		opts[strcase.ToSnake("ObjcClassPrefix")] = *options.ObjcClassPrefix // = {*string | 0xc000011cb0} "FPB"
	}
	if options.CsharpNamespace != nil {
		opts[strcase.ToSnake("CsharpNamespace")] = *options.CsharpNamespace // = {*string | 0xc000011cc0} "Furo.Bigdecimal"
	}
	if options.SwiftPrefix != nil {
		opts[strcase.ToSnake("SwiftPrefix")] = *options.SwiftPrefix // = {*string} nil
	}
	if options.PhpClassPrefix != nil {
		opts[strcase.ToSnake("PhpClassPrefix")] = *options.PhpClassPrefix // = {*string} nil
	}
	if options.PhpNamespace != nil {
		opts[strcase.ToSnake("PhpNamespace")] = *options.PhpNamespace // = {*string} nil
	}
	if options.PhpMetadataNamespace != nil {
		opts[strcase.ToSnake("PhpMetadataNamespace")] = *options.PhpMetadataNamespace // = {*string} nil
	}
	if options.RubyPackage != nil {
		opts[strcase.ToSnake("RubyPackage")] = *options.RubyPackage // = {*string} nil
	}
	return opts
}
