package prodinspect

import (
	"go/ast"
	"go/token"
	"testing"

	"golang.org/x/tools/go/ast/inspector"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	base := inspector.New(nil)
	var fset MockFiler

	i := New(base, &fset)

	assert.Equal(base, i.base)
	assert.Equal(&fset, i.fset)
}

func TestFilter_Preorder(t *testing.T) {
	t.Run("empty types", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)

		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var result []ast.Node
		i.Preorder(nil, func(n ast.Node) {
			result = append(result, n)
		})

		assert.Equal([]ast.Node{
			files[0],
			files[0].Name,
			files[0].Decls[0],
			files[0].Decls[0].(*ast.FuncDecl).Name,
			files[0].Decls[0].(*ast.FuncDecl).Type,
			files[0].Decls[0].(*ast.FuncDecl).Type.Params,
			files[0].Decls[0].(*ast.FuncDecl).Body,
			files[0].Decls[0].(*ast.FuncDecl).Body.List[0],
			files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X,
			files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Fun,
			files[0].Decls[1],
			files[0].Decls[1].(*ast.FuncDecl).Name,
			files[0].Decls[1].(*ast.FuncDecl).Type,
			files[0].Decls[1].(*ast.FuncDecl).Type.Params,
			files[0].Decls[1].(*ast.FuncDecl).Body,
		}, result)
	})

	t.Run("with file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var result []ast.Node
		i.Preorder([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node) {
			result = append(result, n)
		})

		assert.Equal([]ast.Node{
			files[0],
			files[0].Decls[0],
			files[0].Decls[1],
		}, result)
	})

	t.Run("with test file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo_test.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		i.Preorder([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node) {
			assert.Fail("shouldn't be called")
		})
	})

	t.Run("with generated file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Comments: []*ast.CommentGroup{
					{
						List: []*ast.Comment{
							{
								Text: "// Code generated by a generator; DO NOT EDIT.",
							},
						},
					},
				},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		i.Preorder([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node) {
			assert.Fail("shouldn't be called")
		})
	})

	t.Run("without file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var result []ast.Node
		i.Preorder([]ast.Node{
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node) {
			result = append(result, n)
		})

		assert.Equal([]ast.Node{
			files[0].Decls[0],
			files[0].Decls[1],
		}, result)
	})
}

func TestFilter_Nodes(t *testing.T) {
	type event struct {
		push bool
		node ast.Node
	}

	t.Run("empty types", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.Nodes(nil, func(n ast.Node, push bool) bool {
			events = append(events, event{push: push, node: n})
			return true
		})

		assert.Equal([]event{
			{push: true, node: files[0]},
			{push: true, node: files[0].Name},
			{push: false, node: files[0].Name},
			{push: true, node: files[0].Decls[0]},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Name},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Name},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Type},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Type.Params},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Type.Params},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Type},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Body},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0]},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Fun},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Fun},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0]},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Body},
			{push: false, node: files[0].Decls[0]},
			{push: true, node: files[0].Decls[1]},
			{push: true, node: files[0].Decls[1].(*ast.FuncDecl).Name},
			{push: false, node: files[0].Decls[1].(*ast.FuncDecl).Name},
			{push: true, node: files[0].Decls[1].(*ast.FuncDecl).Type},
			{push: true, node: files[0].Decls[1].(*ast.FuncDecl).Type.Params},
			{push: false, node: files[0].Decls[1].(*ast.FuncDecl).Type.Params},
			{push: false, node: files[0].Decls[1].(*ast.FuncDecl).Type},
			{push: true, node: files[0].Decls[1].(*ast.FuncDecl).Body},
			{push: false, node: files[0].Decls[1].(*ast.FuncDecl).Body},
			{push: false, node: files[0].Decls[1]},
			{push: false, node: files[0]},
		}, events)
	})

	t.Run("with file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.Nodes([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node, push bool) bool {
			events = append(events, event{push: push, node: n})
			return true
		})

		assert.Equal([]event{
			{push: true, node: files[0]},
			{push: true, node: files[0].Decls[0]},
			{push: false, node: files[0].Decls[0]},
			{push: true, node: files[0].Decls[1]},
			{push: false, node: files[0].Decls[1]},
			{push: false, node: files[0]},
		}, events)
	})

	t.Run("with test file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo_test.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.Nodes([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node, push bool) bool {
			events = append(events, event{push: push, node: n})
			return true
		})

		assert.Nil(events)
	})

	t.Run("with generated file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Comments: []*ast.CommentGroup{
					{
						List: []*ast.Comment{
							{
								Text: "// Code generated by a generator; DO NOT EDIT.",
							},
						},
					},
				},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.Nodes([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node, push bool) bool {
			events = append(events, event{push: push, node: n})
			return true
		})

		assert.Nil(events)
	})

	t.Run("without file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.Nodes([]ast.Node{
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node, push bool) bool {
			events = append(events, event{push: push, node: n})
			return true
		})

		assert.Equal([]event{
			{push: true, node: files[0].Decls[0]},
			{push: false, node: files[0].Decls[0]},
			{push: true, node: files[0].Decls[1]},
			{push: false, node: files[0].Decls[1]},
		}, events)
	})
}

func TestFilter_WithStack(t *testing.T) {
	type event struct {
		push  bool
		node  ast.Node
		stack []ast.Node
	}

	t.Run("empty types", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.WithStack(nil, func(n ast.Node, push bool, stack []ast.Node) bool {
			events = append(events, event{push: push, node: n, stack: append(stack[:0:0], stack...)})
			return true
		})

		assert.Equal([]event{
			{push: true, node: files[0], stack: []ast.Node{
				files[0],
			}},
			{push: true, node: files[0].Name, stack: []ast.Node{
				files[0],
				files[0].Name,
			}},
			{push: false, node: files[0].Name, stack: []ast.Node{
				files[0],
				files[0].Name,
			}},
			{push: true, node: files[0].Decls[0], stack: []ast.Node{
				files[0],
				files[0].Decls[0],
			}},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Name, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Name,
			}},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Name, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Name,
			}},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Type, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Type,
			}},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Type.Params, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Type,
				files[0].Decls[0].(*ast.FuncDecl).Type.Params,
			}},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Type.Params, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Type,
				files[0].Decls[0].(*ast.FuncDecl).Type.Params,
			}},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Type, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Type,
			}},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Body, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Body,
			}},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0], stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Body,
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0],
			}},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Body,
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0],
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X,
			}},
			{push: true, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Fun, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Body,
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0],
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X,
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Fun,
			}},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Fun, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Body,
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0],
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X,
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Fun,
			}},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Body,
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0],
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ExprStmt).X,
			}},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Body.List[0], stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Body,
				files[0].Decls[0].(*ast.FuncDecl).Body.List[0],
			}},
			{push: false, node: files[0].Decls[0].(*ast.FuncDecl).Body, stack: []ast.Node{
				files[0],
				files[0].Decls[0],
				files[0].Decls[0].(*ast.FuncDecl).Body,
			}},
			{push: false, node: files[0].Decls[0], stack: []ast.Node{
				files[0],
				files[0].Decls[0],
			}},
			{push: true, node: files[0].Decls[1], stack: []ast.Node{
				files[0],
				files[0].Decls[1],
			}},
			{push: true, node: files[0].Decls[1].(*ast.FuncDecl).Name, stack: []ast.Node{
				files[0],
				files[0].Decls[1],
				files[0].Decls[1].(*ast.FuncDecl).Name,
			}},
			{push: false, node: files[0].Decls[1].(*ast.FuncDecl).Name, stack: []ast.Node{
				files[0],
				files[0].Decls[1],
				files[0].Decls[1].(*ast.FuncDecl).Name,
			}},
			{push: true, node: files[0].Decls[1].(*ast.FuncDecl).Type, stack: []ast.Node{
				files[0],
				files[0].Decls[1],
				files[0].Decls[1].(*ast.FuncDecl).Type,
			}},
			{push: true, node: files[0].Decls[1].(*ast.FuncDecl).Type.Params, stack: []ast.Node{
				files[0],
				files[0].Decls[1],
				files[0].Decls[1].(*ast.FuncDecl).Type,
				files[0].Decls[1].(*ast.FuncDecl).Type.Params,
			}},
			{push: false, node: files[0].Decls[1].(*ast.FuncDecl).Type.Params, stack: []ast.Node{
				files[0],
				files[0].Decls[1],
				files[0].Decls[1].(*ast.FuncDecl).Type,
				files[0].Decls[1].(*ast.FuncDecl).Type.Params,
			}},
			{push: false, node: files[0].Decls[1].(*ast.FuncDecl).Type, stack: []ast.Node{
				files[0],
				files[0].Decls[1],
				files[0].Decls[1].(*ast.FuncDecl).Type,
			}},
			{push: true, node: files[0].Decls[1].(*ast.FuncDecl).Body, stack: []ast.Node{
				files[0],
				files[0].Decls[1],
				files[0].Decls[1].(*ast.FuncDecl).Body,
			}},
			{push: false, node: files[0].Decls[1].(*ast.FuncDecl).Body, stack: []ast.Node{
				files[0],
				files[0].Decls[1],
				files[0].Decls[1].(*ast.FuncDecl).Body,
			}},
			{push: false, node: files[0].Decls[1], stack: []ast.Node{
				files[0],
				files[0].Decls[1],
			}},
			{push: false, node: files[0], stack: []ast.Node{
				files[0],
			}},
		}, events)
	})

	t.Run("with file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.WithStack([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node, push bool, stack []ast.Node) bool {
			events = append(events, event{push: push, node: n, stack: append(stack[:0:0], stack...)})
			return true
		})

		assert.Equal([]event{
			{push: true, node: files[0], stack: []ast.Node{
				files[0],
			}},
			{push: true, node: files[0].Decls[0], stack: []ast.Node{
				files[0],
				files[0].Decls[0],
			}},
			{push: false, node: files[0].Decls[0], stack: []ast.Node{
				files[0],
				files[0].Decls[0],
			}},
			{push: true, node: files[0].Decls[1], stack: []ast.Node{
				files[0],
				files[0].Decls[1],
			}},
			{push: false, node: files[0].Decls[1], stack: []ast.Node{
				files[0],
				files[0].Decls[1],
			}},
			{push: false, node: files[0], stack: []ast.Node{
				files[0],
			}},
		}, events)
	})

	t.Run("with test file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo_test.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.WithStack([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node, push bool, stack []ast.Node) bool {
			events = append(events, event{push: push, node: n, stack: stack})
			return true
		})

		assert.Nil(events)
	})

	t.Run("with generated file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Comments: []*ast.CommentGroup{
					{
						List: []*ast.Comment{
							{
								Text: "// Code generated by a generator; DO NOT EDIT.",
							},
						},
					},
				},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.WithStack([]ast.Node{
			(*ast.File)(nil),
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node, push bool, stack []ast.Node) bool {
			events = append(events, event{push: push, node: n, stack: stack})
			return true
		})

		assert.Nil(events)
	})

	t.Run("without file", func(t *testing.T) {
		assert := assert.New(t)

		fs := token.NewFileSet()
		fs.AddFile("foo.go", 1, 1)

		files := []*ast.File{
			{
				Name: &ast.Ident{},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Name: &ast.Ident{},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{},
					},
				},
			},
		}

		base := inspector.New(files)
		fset := MockFiler{
			file: fs.File(1),
		}

		i := Inspector{
			base: base,
			fset: &fset,
		}

		var events []event
		i.WithStack([]ast.Node{
			(*ast.FuncDecl)(nil),
		}, func(n ast.Node, push bool, stack []ast.Node) bool {
			events = append(events, event{push: push, node: n, stack: append(stack[:0:0], stack...)})
			return true
		})

		assert.Equal([]event{
			{push: true, node: files[0].Decls[0], stack: []ast.Node{
				files[0],
				files[0].Decls[0],
			}},
			{push: false, node: files[0].Decls[0], stack: []ast.Node{
				files[0],
				files[0].Decls[0],
			}},
			{push: true, node: files[0].Decls[1], stack: []ast.Node{
				files[0],
				files[0].Decls[1],
			}},
			{push: false, node: files[0].Decls[1], stack: []ast.Node{
				files[0],
				files[0].Decls[1],
			}},
		}, events)
	})
}

type MockFiler struct {
	file *token.File
}

func (m *MockFiler) File(p token.Pos) (f *token.File) {
	return m.file
}
