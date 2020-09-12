# protoc-gen-furo-specs
protoc-gen-furo-specs is a protoc plugin to create a patch method for your pb.go messages. 
The method applys a delta for every field to an existing target of the same
type according to the update_mask. 

It generates the method if your message type have a field **update_mask**.

## Use Case
- Client sends a PATCH with update_mask (*google.protobuf.types.known.FieldMask*)
- The Server loads the full data and applies the changes with original.PatchWithUpdateMask(deltapb)
  
 
## Parameters
#### [exclude] 
Optional regex to match target files that should not be built.

If you generate files which have a update_mask but do not need this method you can add a exclude regex which must 
match to not be generated. 

## Installation

``` 
go get github.com/theNorstroem/protoc-gen-furo-specs
```

Add protoc-gen-furo-specs to your tools.go file.

```go
//+build tools

package tools

import (
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "github.com/theNorstroem/protoc-gen-furo-specs"
)

```

## Using the plugin
Like every other protoc generator... Nothing special here.
```
protoc -I. \
-I/usr/local/include \
-I$GOPATH/src  \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  \
--go-patch_out=output:. **/*.proto
```

