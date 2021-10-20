// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/errordeveloper/kuegen/compiler"
)

const (
	templateKey = "template"
	defaultsKey = "defaults"
	resourceKey = "resource"
)

type Generator struct {
	InputDirectory string

	template *cue.Instance
}

// TODO: move this package to kue
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

	return nil
}

func (g *Generator) with(key string, obj interface{}) (*Generator, error) {
	result, err := g.template.Fill(obj, key)
	if err != nil {
		return nil, fmt.Errorf("cannot fill %q: %w", key, err)
	}
	if result.Err != nil {
		return nil, fmt.Errorf("error after filling %q: %w", key, err)
	}
	return &Generator{
		InputDirectory: g.InputDirectory,
		template:       result,
	}, nil
}

func (g *Generator) WithDefaults(obj interface{}) (*Generator, error) {
	return g.with(defaultsKey, obj)
}

func (g *Generator) WithResource(obj interface{}) (*Generator, error) {
	return g.with(resourceKey, obj)
}

func (g *Generator) RenderJSON() ([]byte, error) {
	value := g.template.Lookup(templateKey)
	if err := value.Err(); err != nil {
		return nil, fmt.Errorf("unable to lookup %q: %w", templateKey, err)
	}
	return value.MarshalJSON()
}
