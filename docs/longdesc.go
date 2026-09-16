package docs

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// LongDescriptions parses a Go source file and returns a map from each
// urfave/cli flag's Name field to its long docs description, derived from
// the leading doc comment directly above the flag's composite literal.
//
// Convention: a // comment block directly above an &cli.XFlag{...} literal
// is treated as the long docs description. Lines within a paragraph are
// joined with spaces; paragraphs (separated by an empty // line) are joined
// with "\n\n". A blank line in the source between the comment and the
// literal breaks the association.
//
// This helper is intentionally format-agnostic: callers are responsible
// for merging the result into their docs format of choice (YAML, JSON,
// markdown, etc.).
func LongDescriptions(sourcePath string) (map[string]string, error) {
	out := make(map[string]string)

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, sourcePath, nil, parser.ParseComments)
	if err != nil {
		return out, err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		if !isUrfaveFlagType(cl.Type) {
			return true
		}

		cg := leadingComment(fset, file, cl)
		if cg == nil {
			return true
		}

		long := commentText(cg)
		if long == "" {
			return true
		}

		name := flagName(cl)
		if name == "" {
			return true
		}

		out[name] = long

		return true
	})

	return out, nil
}

func isUrfaveFlagType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok || ident.Name != "cli" {
		return false
	}

	switch sel.Sel.Name {
	case "BoolFlag", "StringFlag", "IntFlag", "StringSliceFlag":
		return true
	}

	return false
}

func flagName(cl *ast.CompositeLit) string {
	for _, e := range cl.Elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Name" {
			continue
		}

		bl, ok := kv.Value.(*ast.BasicLit)
		if !ok || bl.Kind != token.STRING {
			continue
		}

		return strings.Trim(bl.Value, `"`)
	}

	return ""
}

// leadingComment returns the comment group immediately above node, or nil if
// there is a blank line between them. *ast.CompositeLit has no Doc field, so
// a position-based lookup against file.Comments is required.
func leadingComment(fset *token.FileSet, file *ast.File, node ast.Node) *ast.CommentGroup {
	target := node.Pos()

	var leading *ast.CommentGroup

	for _, cg := range file.Comments {
		if cg.End() >= target {
			continue
		}

		if leading == nil || cg.End() > leading.End() {
			leading = cg
		}
	}

	if leading == nil {
		return nil
	}

	endLine := fset.Position(leading.End()).Line
	startLine := fset.Position(node.Pos()).Line

	if startLine > endLine+1 {
		return nil
	}

	return leading
}

// commentText turns an *ast.CommentGroup into a normalized description.
// Line comments are joined with spaces within a paragraph; paragraphs
// (separated by an empty // line in the source) are joined with "\n\n".
func commentText(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}

	var paragraphs []string

	var current []string

	flush := func() {
		if len(current) > 0 {
			paragraphs = append(paragraphs, strings.Join(current, " "))
			current = nil
		}
	}

	for _, c := range cg.List {
		body := strings.TrimPrefix(c.Text, "//")
		body = strings.TrimPrefix(body, " ")

		body = strings.TrimRight(body, " \t")
		if body == "" {
			flush()

			continue
		}

		current = append(current, body)
	}

	flush()

	return strings.Join(paragraphs, "\n\n")
}
