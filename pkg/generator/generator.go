package generator

import (
	"github.com/iancoleman/strcase"
	"github.com/theNorstroem/protoc-gen-furo-specs/pkg/protoast"
	"github.com/theNorstroem/spectools/pkg/orderedmap"
	"github.com/theNorstroem/spectools/pkg/specSpec"
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

				description := ""
				if SourceInfo.Messages[MessageIndex].Info.LeadingComments != nil {
					description = *SourceInfo.Messages[MessageIndex].Info.LeadingComments
				}

				typeSpec := specSpec.Type{
					Name:        packagename,
					Type:        *Message.Name,
					Description: description,
					XProto: &specSpec.Typeproto{
						Package:    strings.Join(strings.Split(path.Dir(protofilename), "/"), "."), // package is protofilename with . is not the go package
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

func getFields(fieldinfo protoast.MessageInfo) *orderedmap.OrderedMap {
	omap := orderedmap.New()
	for _, f := range fieldinfo.FieldInfos {

		fielddescription := ""
		if f.Info.LeadingComments != nil {
			fielddescription = *f.Info.LeadingComments
		}

		// repeated is in f.Field.Label

		field := specSpec.Field{
			Type:        *f.Field.TypeName,
			Description: fielddescription,
			XProto: &specSpec.Fieldproto{
				Number: *f.Field.Number,
				Oneof:  "",
			},
			XUi:         nil,
			Meta:        nil,
			Constraints: nil,
		}
		omap.Set(f.Name, field)
	}

	return omap
}

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
