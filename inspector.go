package inspect

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

type Filter struct {
	base Inspector
	fset Filer
}

func New(base Inspector, fset Filer) *Filter {
	return &Filter{
		base: base,
		fset: fset,
	}
}

func (i *Filter) Preorder(types []ast.Node, f func(n ast.Node)) {
	c := containsFile(types)

	if !c {
		types = append(types, (*ast.File)(nil))
	}

	i.base.Preorder(types, func(n ast.Node) {
		if f, ok := n.(*ast.File); ok && (!c || ignored(f, i.fset)) {
			return
		}

		f(n)
	})
}

func (i *Filter) Nodes(types []ast.Node, f func(n ast.Node, push bool) (prune bool)) {
	c := containsFile(types)

	if !c {
		types = append(types, (*ast.File)(nil))
	}

	i.base.Nodes(types, func(n ast.Node, push bool) bool {
		if push {
			if f, ok := n.(*ast.File); ok && (!c || ignored(f, i.fset)) {
				return false
			}
		}

		return f(n, push)
	})
}

func (i *Filter) WithStack(types []ast.Node, f func(n ast.Node, push bool, stack []ast.Node) (prune bool)) {
	c := containsFile(types)

	if !c {
		types = append(types, (*ast.File)(nil))
	}

	i.base.WithStack(types, func(n ast.Node, push bool, stack []ast.Node) bool {
		if push {
			if f, ok := n.(*ast.File); ok && (!c || ignored(f, i.fset)) {
				return false
			}
		}

		return f(n, push, stack)
	})
}

func containsFile(types []ast.Node) bool {
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

type Inspector interface {
	Preorder(types []ast.Node, f func(n ast.Node))
	Nodes(types []ast.Node, f func(n ast.Node, push bool) (prune bool))
	WithStack(types []ast.Node, f func(n ast.Node, push bool, stack []ast.Node) (prune bool))
}

type Filer interface {
	File(p token.Pos) (f *token.File)
}
