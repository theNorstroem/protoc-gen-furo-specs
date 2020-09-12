package template

import (
	"github.com/Masterminds/sprig"
	"github.com/iancoleman/strcase"
	"html/template"
	"log"
	"strings"
)

// template func
func noescape(str string) template.HTML {
	return template.HTML(str)
}

// template func
func camel(str string) string {
	return strcase.ToCamel(str)
}

func comment(c *string) string {
	if c != nil {
		comment := *c
		return "//" + strings.Replace(comment[0:len(comment)-1], "\n", "\n//", 100)
	}
	return ""
}

func LoadTemplate() *template.Template {
	// prepare template
	fn := sprig.FuncMap()
	fn["noescape"] = noescape
	fn["camel"] = camel
	fn["comment"] = comment

	tmpl, templateError := template.New("proto").Funcs(fn).Parse(Template)
	if templateError != nil {
		log.Fatal(templateError)
	}
	return tmpl
}
