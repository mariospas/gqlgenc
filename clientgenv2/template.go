package clientgenv2

import (

	// nolint:golint, nolintlint
	"bytes"
	_ "embed" // used to load template file
	"fmt"
	"go/types"
	"strings"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	gqlgencConfig "github.com/Yamashou/gqlgenc/config"
)

//go:embed template.gotpl
var template string

func RenderTemplate(cfg *config.Config, query *Query, mutation *Mutation, fragments []*Fragment, operations []*Operation, operationResponses []*OperationResponse, structSources []*StructSource, generateCfg *gqlgencConfig.GenerateConfig, client config.PackageConfig) error {
	if err := templates.Render(templates.Options{
		PackageName: client.Package,
		Filename:    client.Filename,
		Template:    template,
		Data: map[string]interface{}{
			"Query":               query,
			"Mutation":            mutation,
			"Fragment":            fragments,
			"Operation":           operations,
			"OperationResponse":   operationResponses,
			"GenerateClient":      generateCfg.ShouldGenerateClient(),
			"StructSources":       structSources,
			"ClientInterfaceName": generateCfg.GetClientInterfaceName(),
		},
		Packages:   cfg.Packages,
		PackageDoc: "// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.\n",
		Funcs: map[string]interface{}{
			"genGetters": genGetters,
		},
	}); err != nil {
		return fmt.Errorf("%s generating failed: %w", client.Filename, err)
	}

	return nil
}

// This method returns a string of getters for a struct.
// The idea is to be able to chain calls safely without having to check for nil.
// To make this work we need to return a pointer to the struct if the field is a struct.
func genGetters(name string, p types.Type) string {
	var it *types.Struct
	it, ok := p.(*types.Struct)
	if !ok {
		return ""
	}
	var buf bytes.Buffer

	for i := 0; i < it.NumFields(); i++ {
		field := it.Field(i)

		returns := returnTypeName(field.Type(), false)

		buf.WriteString("func (t *" + name + ") Get" + field.Name() + "() " + returns + "{\n")
		buf.WriteString("if t == nil {\n t = &" + name + "{}\n}\n")

		pointerOrNot := ""
		if _, ok := field.Type().(*types.Named); ok {
			pointerOrNot = "&"
		}

		buf.WriteString("return " + pointerOrNot + "t." + field.Name() + "\n}\n")
	}

	return buf.String()
}

func returnTypeName(t types.Type, nested bool) string {
	switch it := t.(type) {
	case *types.Basic:
		return it.String()
	case *types.Pointer:
		return "*" + returnTypeName(it.Elem(), true)
	case *types.Slice:
		return "[]" + returnTypeName(it.Elem(), true)
	case *types.Named:
		isImported := it.Obj().Parent() != nil

		s := strings.Split(it.String(), ".")
		name := s[len(s)-1]
		if isImported {
			name = s[len(s)-2] + "." + name
		}

		if nested {
			return name
		}

		return "*" + name
	default:
		return fmt.Sprintf("%T", it)
	}
}
