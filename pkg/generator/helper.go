package generator

import (
	"fmt"

	"github.com/theNorstroem/protoc-gen-furo-specs/pkg/protoast"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
	"path"
	"regexp"
	"strings"
)

// decide if something should be generated
// the complete descriptor is given
func shouldGenerateTypeSpec(protoAST *protoast.ProtoAST, name string, descriptor *descriptorpb.FileDescriptorProto, message *descriptorpb.DescriptorProto) bool {
	// we generate types only
	// todo: generate enums
	// todo: generate services ???
	if len(message.Field) == 0 {
		return false
	}

	if descriptor.Options == nil || descriptor.Options.GoPackage == nil {
		// give a warning because go package option is missing
		os.Stderr.WriteString(fmt.Sprintf("Go package is missing: %s", descriptor.Name))
		return false
	}

	// check for excludes
	rgx, found := protoAST.GetParameter("exclude")
	if found {
		// filter out excluded files based on their target name
		// because a message can have multiple messages
		filename, _ := FileAndPackageNameToGenerate(descriptor, message)
		//rgx = ".*(Entity)|(Collection).go"
		match, _ := regexp.MatchString(rgx, filename)
		if match {
			return false
		}
	}

	return protoAST.FileProtoMap[*descriptor.Name] != nil && descriptor.Options != nil && descriptor.Options.GoPackage != nil

}

func FileAndPackageNameToGenerate(descriptor *descriptorpb.FileDescriptorProto, message *descriptorpb.DescriptorProto) (filename string, packagename string) {
	name := *message.Name
	p := strings.Split(*descriptor.Options.GoPackage, ";")
	pkg := p[len(p)-1]

	filename = path.Join(path.Dir(*descriptor.Name), name+".type.spec")
	return filename, pkg
}
