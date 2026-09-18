package docs

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// LongDescription is a structured representation of the leading doc
// comment above a flag literal. Source layout is preserved so consumers
// can adapt it to any output format (markdown, YAML literal blocks,
// AsciiDoc, etc.) without re-parsing a flattened string.
//
// Paragraphs split on empty // lines — the standard godoc/JSDoc
// convention. Each Paragraph keeps its original source lines (the "// "
// comment prefix is stripped, but additional leading whitespace such as
// the two-space indent used for markdown list-item continuations is
// preserved verbatim).
type LongDescription struct {
	// Paragraphs indexes paragraphs (outer) × source lines (inner).
	Paragraphs [][]string
}

// IsZero reports whether the receiver has no paragraphs.
func (d *LongDescription) IsZero() bool {
	if d == nil {
		return true
	}

	return len(d.Paragraphs) == 0
}

// Flat returns the description as a single string — the canonical
// lossless serialisation. Lines inside a paragraph are joined with "\n"
// and paragraphs are joined with "\n\n".
func (d *LongDescription) Flat() string {
	if d.IsZero() {
		return ""
	}

	parts := make([]string, len(d.Paragraphs))
	for i, p := range d.Paragraphs {
		parts[i] = strings.Join(p, "\n")
	}

	return strings.Join(parts, "\n\n")
}

// String implements fmt.Stringer and is equivalent to Flat.
func (d *LongDescription) String() string {
	return d.Flat()
}

// FlagTypeMatcher reports whether an AST type expression refers to a
// flag composite literal that the caller wants documented. See
// DefaultFlagTypeMatcher for the built-in matcher and LongDescriptionsWith
// for usage.
type FlagTypeMatcher func(ast.Expr) bool

// LongDescriptionsWith parses sourcePath and extracts a LongDescription
// for every flag composite literal whose type matches at least one of
// the supplied matchers (OR semantics; evaluation short-circuits at the
// first match). When matchers is empty, only DefaultFlagTypeMatcher is
// applied.
//
// A // comment block directly above an &<flag>Flag{...} literal is
// treated as the long docs description; a blank line in the source
// between the comment and the literal breaks the association.
//
// A typical call site that wants both urfave core flags and wp-plugin-go
// custom flag types:
//
//	docs.LongDescriptionsWith(
//	    "plugin/plugin.go",
//	    docs.DefaultFlagTypeMatcher,
//	    docs.SelectorMatcher("plugin_cli", "StringMapFlag", "DeepStringMapFlag"),
//	)
func LongDescriptionsWith(sourcePath string, matchers ...FlagTypeMatcher) (map[string]*LongDescription, error) {
	out := make(map[string]*LongDescription)

	if len(matchers) == 0 {
		matchers = []FlagTypeMatcher{DefaultFlagTypeMatcher}
	}

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

		if !matchesAny(cl.Type, matchers) {
			return true
		}

		cg := leadingComment(fset, file, cl)
		if cg == nil {
			return true
		}

		desc := parseComment(cg)
		if desc.IsZero() {
			return true
		}

		name := flagName(cl)
		if name == "" {
			return true
		}

		out[name] = desc

		return true
	})

	return out, nil
}

// LongDescriptionsForWith is the matcher-aware template-data convenience
// wrapper around LongDescriptionsWith. It normalises flag names so
// "upload.metadata" matches the env-derived "upload_metadata", and
// returns an empty map (instead of an error) when sourcePath is empty or
// unparseable, so a missing source never breaks a docs pipeline.
func LongDescriptionsForWith(sourcePath string, matchers ...FlagTypeMatcher) map[string]*LongDescription {
	if sourcePath == "" {
		return map[string]*LongDescription{}
	}

	descs, err := LongDescriptionsWith(sourcePath, matchers...)
	if err != nil {
		return map[string]*LongDescription{}
	}

	normalized := make(map[string]*LongDescription, len(descs))
	for name, d := range descs {
		normalized[normalizeFlagName(name)] = d
	}

	return normalized
}

// LongDescriptions is the v6.4.0-compatible flat-string adapter around
// LongDescriptionsWith. It returns descriptions serialised via Flat() —
// the same canonical string form documented on LongDescription.Flat.
//
// Retained to avoid breaking the v6.x API. New callers should prefer
// LongDescriptionsWith, which exposes the structured form and lets
// callers choose their own formatter.
func LongDescriptions(sourcePath string) (map[string]string, error) {
	descs, err := LongDescriptionsWith(sourcePath)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(descs))
	for name, d := range descs {
		out[name] = d.Flat()
	}

	return out, nil
}

// LongDescriptionsFor is the v6.4.0-compatible flat-string adapter
// around LongDescriptionsForWith. It mirrors the normalisation behaviour
// of LongDescriptionsForWith (so "upload.metadata" matches
// "upload_metadata") and returns descriptions serialised via Flat().
//
// Retained to avoid breaking the v6.x API. New callers should prefer
// LongDescriptionsForWith.
func LongDescriptionsFor(sourcePath string) map[string]string {
	descs := LongDescriptionsForWith(sourcePath)

	out := make(map[string]string, len(descs))
	for name, d := range descs {
		out[name] = d.Flat()
	}

	return out
}

// DefaultFlagTypeMatcher matches the urfave/cli/v3 core flag composite
// literals: BoolFlag, StringFlag, IntFlag and StringSliceFlag. Custom
// flag types are intentionally not matched here — see LongDescriptionsWith
// for the extension mechanism.
func DefaultFlagTypeMatcher(expr ast.Expr) bool {
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

// SelectorMatcher returns a FlagTypeMatcher that matches selector
// expressions whose package ident equals pkg and whose selector name is
// in names. Use it to build a matcher for custom flag types:
//
//	docs.SelectorMatcher("plugin_cli", "StringMapFlag", "DeepStringMapFlag")
func SelectorMatcher(pkg string, names ...string) FlagTypeMatcher {
	set := make(map[string]struct{}, len(names))
	for _, n := range names {
		set[n] = struct{}{}
	}

	return func(expr ast.Expr) bool {
		sel, ok := expr.(*ast.SelectorExpr)
		if !ok {
			return false
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != pkg {
			return false
		}

		_, ok = set[sel.Sel.Name]

		return ok
	}
}

func matchesAny(expr ast.Expr, matchers []FlagTypeMatcher) bool {
	for _, m := range matchers {
		if m == nil {
			continue
		}

		if m(expr) {
			return true
		}
	}

	return false
}

func normalizeFlagName(name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(name, "-", "_"), ".", "_")
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

// leadingComment returns the comment group immediately above node, or
// nil if a blank line separates them. *ast.CompositeLit has no Doc field
// on the AST, so a position-based lookup against file.Comments is
// required.
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

// parseComment converts an *ast.CommentGroup into a structured
// LongDescription. Paragraphs are split on empty // lines; within each
// paragraph, source line breaks and the two-space indent used for
// markdown list continuations are preserved verbatim. Only the
// conventional "// " comment prefix is stripped.
func parseComment(cg *ast.CommentGroup) *LongDescription {
	if cg == nil {
		return &LongDescription{}
	}

	desc := &LongDescription{}
	var current []string

	flush := func() {
		if len(current) > 0 {
			desc.Paragraphs = append(desc.Paragraphs, current)
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

	return desc
}
