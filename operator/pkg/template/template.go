package template

import (
	"cuelang.org/go/cue"
	"github.com/errordeveloper/kue/pkg/compiler"
)

const (
	templateFilename = "template.cue"

	templateKey = "template"
	resourceKey = "resource"
)

type Generator struct {
	InputDirectory string

	template *cue.Instance
}

func (g *Generator) CompileAndValidate() error {
	// TODO: produce meanigful compilation errors
	// TODO: produce meanigful validation errors
	// TODO: validate the types match expections, e.g. objets not string

	c := compiler.NewCompiler(g.InputDirectory)

	template, err := c.BuildAll()
	if err != nil {
		return err
	}

	g.template = template

	if err := g.template.Value().Err(); err != nil {
		return err
	}

	return nil
}

func (g *Generator) RenderJSON(obj interface{}) ([]byte, error) {
	result, err := g.template.Fill(obj, resourceKey)
	if err != nil {
		return nil, err
	}

	return result.Lookup(templateKey).MarshalJSON()
}
