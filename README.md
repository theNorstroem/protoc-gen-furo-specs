> MOVED TO https://github.com/eclipse/eclipsefuro

# protoc-gen-furo-specs


## Use Case
- If you have a bunch of services and messages specified in proto and want to use them in a furo client project,
this generator comes very handy.
  

## Parameters
#### [exclude] 
Optional regex to match target files that should not be built.

Maybe you do not want all request and response types from the services...

`--furo-specs_out=exclude=".*(Response)|(Request).type.spec":. *.proto`

## Installation

``` 
go get github.com/theNorstroem/protoc-gen-furo-specs
```

Add protoc-gen-furo-specs to your tools.go file if you want.

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

