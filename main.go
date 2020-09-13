package main

import (
	"github.com/theNorstroem/protoc-gen-furo-specs/pkg/generator"
	"github.com/theNorstroem/protoc-gen-furo-specs/pkg/protoast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"os"
)

func main() {
	// os.Stdin will contain data which will unmarshal into the following object:
	// https://godoc.org/github.com/golang/protobuf/protoc-gen-go/plugin#CodeGeneratorRequest
	req := &pluginpb.CodeGeneratorRequest{}

	//data, err := ioutil.ReadAll(os.Stdin)
	// debug mode
	data, err := ioutil.ReadFile("protocdata")
	//ioutil.WriteFile("protocdata",data, 666)

	if err != nil {
		panic(err)
	}

	proto.Unmarshal(data, req)

	if err != nil {
		panic(err)
	}

	Ast := protoast.NewProtoAST(req)

	err = generator.Generate(Ast)
	if err != nil {
		panic(err)
	}

	marshalled, err := proto.Marshal(Ast.Response)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(marshalled)
}
