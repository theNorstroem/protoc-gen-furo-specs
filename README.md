# protoc-gen-furo-specs


## Use Case


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
go build . && protoc --plugin protoc-gen-furo-specs -I../furoBaseSpecs/dist/proto/Messages/ -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/googleapis/googleapis --furo-specs_out=:./out ../furoBaseSpecs/dist/proto/Messages/**/*.proto

```

