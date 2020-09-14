package generator

import (
	"github.com/theNorstroem/protoc-gen-furo-specs/pkg/protoast"
	"github.com/theNorstroem/spectools/pkg/orderedmap"
	"github.com/theNorstroem/spectools/pkg/specSpec"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"strings"
)

func getServices(info protoast.ServiceInfo) *orderedmap.OrderedMap {
	omap := orderedmap.New()
	for _, methodInfo := range info.Methods {

		rpcMethodDescription := ""
		if methodInfo.Info.LeadingComments != nil {
			rpcMethodDescription = cleanDescription(*methodInfo.Info.LeadingComments)
		}
		deeplinkDescription := ""

		if methodInfo.HttpRule.Info.LeadingComments != nil {
			deeplinkDescription = cleanDescription(*methodInfo.HttpRule.Info.LeadingComments)
		}
		// extract api options to href method and rel

		// *methodInfo.HttpRule.ApiOptions.Pattern is oneof
		// details: vendor/google.golang.org/genproto/googleapis/api/annotations/http.pb.go:400
		href, verb, rel := extractApiOptionPattern(methodInfo.HttpRule)

		method := specSpec.Rpc{
			Description: rpcMethodDescription,
			Data: &specSpec.Servicereqres{
				Request:  *methodInfo.Method.InputType,
				Response: *methodInfo.Method.OutputType,
			},
			Deeplink: &specSpec.Servicedeeplink{
				Description: deeplinkDescription,
				Href:        href,
				Method:      verb,
				Rel:         rel,
			},
			Query:   nil,
			RpcName: *methodInfo.Method.Name,
		}

		omap.Set(methodInfo.Name, method)
	}

	return omap
}

// get the href, method, rel
func extractApiOptionPattern(info *protoast.ApiOptionInfo) (href string, method string, rel string) {

	pattern := info.ApiOptions.Pattern
	href = "/no/option/given"
	method = "GET"
	rel = "self"

	if info.Info != nil {
		// try first line of comment for the rel
		//   Delete: DELETE /samples/{xxx} google.protobuf.Empty, google.protobuf.Empty #Use this to delete existing samples.
		// becomes delete
		c := strings.Split(*info.Info.LeadingComments, ":")
		if len(c) > 0 && len(strings.TrimSpace(c[0])) > 3 {
			rel = strings.ToLower(strings.TrimSpace(c[0]))
		}
	}

	get, isGet := pattern.(*options.HttpRule_Get)
	if isGet {
		href = get.Get
		method = "GET"
		return href, method, rel
	}

	post, isPost := pattern.(*options.HttpRule_Post)
	if isPost {
		href = post.Post
		method = "POST"
		rel = checkForFallbackRel(rel, method)
		return href, method, rel
	}

	patch, isPatch := pattern.(*options.HttpRule_Patch)
	if isPatch {
		href = patch.Patch
		method = "PATCH"
		rel = checkForFallbackRel(rel, method)
		return href, method, rel
	}

	put, isPut := pattern.(*options.HttpRule_Put)
	if isPut {
		href = put.Put
		method = "PUT"
		rel = checkForFallbackRel(rel, method)
		return href, method, rel
	}

	delete, isDelete := pattern.(*options.HttpRule_Delete)
	if isDelete {
		href = delete.Delete
		method = "DELETE"
		rel = checkForFallbackRel(rel, method)
		return href, method, rel
	}

	// custom is for the verb and not for the custom method...
	custom, isCustom := pattern.(*options.HttpRule_Custom)
	if isCustom {
		href = custom.Custom.Path
		method = custom.Custom.Kind
		rel = checkForFallbackRel(rel, method)
		return href, method, rel
	}

	return href, method, rel
}

func checkForFallbackRel(rel string, method string) string {
	// fallback
	if rel == "self" {
		rel = strings.ToLower(method)
	}
	return rel
}
