# prodinspect

[inspect](https://godoc.org/golang.org/x/tools/go/analysis/passes/inspect) but for production code only.

## How to use

First, install the package via `go get`.

```console
$ go get github.com/ichiban/prodinspect
```

Then, require `prodinspect.Analyzer` and use `*prodinspect.Inspector` in your analyzer.

```go
package cyclomatic

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/analysis"

	"github.com/ichiban/prodinspect"
)

var Analyzer = &analysis.Analyzer{
	Name:      "cyclomatic",
	Doc:       `check cyclomatic complexity of functions.`,
	Requires:  []*analysis.Analyzer{prodinspect.Analyzer},
	Run:       run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[prodinspect.Analyzer].(*prodinspect.Inspector)

	// ...

	return nil, nil
}

```

## Definition of production code

Go files except:
- test files (files with `_test.go` suffix)
- generated files (files with `// Code generated * DO NOT EDIT.` comment)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
