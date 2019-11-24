package prodinspect

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

type Inspector struct {
	base WithStacker
	fset Filer
}

func New(base WithStacker, fset Filer) *Inspector {
	return &Inspector{
		base: base,
		fset: fset,
	}
}

func (i *Inspector) Preorder(types []ast.Node, f func(n ast.Node)) {
	c := containsFile(types)

	if !c {
		types = append(types, (*ast.File)(nil))
	}

	i.base.WithStack(types, func(n ast.Node, push bool, _ []ast.Node) bool {
		if !push {
			return false
		}

		if f, ok := n.(*ast.File); ok {
			if ignored(f, i.fset) {
				return false
			}

			if !c {
				return true
			}
		}

		f(n)

		return true
	})
}

func (i *Inspector) Nodes(types []ast.Node, f func(n ast.Node, push bool) (prune bool)) {
	c := containsFile(types)

	if !c {
		types = append(types, (*ast.File)(nil))
	}

	i.base.WithStack(types, func(n ast.Node, push bool, _ []ast.Node) bool {
		if f, ok := n.(*ast.File); ok {
			if ignored(f, i.fset) {
				return false
			}

			if !c {
				return true
			}
		}

		return f(n, push)
	})
}

func (i *Inspector) WithStack(types []ast.Node, f func(n ast.Node, push bool, stack []ast.Node) (prune bool)) {
	c := containsFile(types)

	if !c {
		types = append(types, (*ast.File)(nil))
	}

	i.base.WithStack(types, func(n ast.Node, push bool, stack []ast.Node) bool {
		if f, ok := n.(*ast.File); ok {
			if ignored(f, i.fset) {
				return false
			}

			if !c {
				return true
			}
		}

		return f(n, push, stack)
	})
}

func containsFile(types []ast.Node) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if _, ok := t.(*ast.File); ok {
			return true
		}
	}
	return false
}

func ignored(n *ast.File, fset Filer) bool {
	f := fset.File(n.Pos())
	name := f.Name()

	return strings.HasSuffix(name, "_test.go") || generated(n)
}

// https://github.com/golang/go/issues/13560#issuecomment-288457920
var pattern = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

func generated(f *ast.File) bool {
	for _, c := range f.Comments {
		for _, l := range c.List {
			if pattern.MatchString(l.Text) {
				return true
			}
		}
	}
	return false
}

type WithStacker interface {
	WithStack(types []ast.Node, f func(n ast.Node, push bool, stack []ast.Node) (prune bool))
}

type Filer interface {
	File(p token.Pos) (f *token.File)
}
