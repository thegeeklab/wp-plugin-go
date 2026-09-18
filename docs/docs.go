package docs

import (
	"bytes"
	"embed"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"

	plugin_template "github.com/thegeeklab/wp-plugin-go/v6/template"

	"github.com/urfave/cli/v3"
)

// PluginArg is one flag in the rendered CLI doc, populated by
// parseFlags. LongDescription is the markdown-formatted string produced
// by LongDescriptionMarkdown — it is *not* the structured *LongDescription
// type, despite sharing a name.
type PluginArg struct {
	Name            string
	EnvVars         []string
	Description     string
	LongDescription string
	Default         string
	Type            string
	Required        bool
}

// CliTemplate is the template-data root handed to the markdown template
// by GetTemplateDataWithSource.
type CliTemplate struct {
	Name        string
	Version     string
	Description string
	Usage       string
	UsageText   string
	GlobalArgs  []*PluginArg
}

//go:embed templates
var templateFs embed.FS

// ToMarkdown creates a markdown string for the `*App` without long
// descriptions. It is a convenience wrapper around ToMarkdownWithSource that
// passes an empty sourcePath.
func ToMarkdown(app *cli.Command) (string, error) {
	return ToMarkdownWithSource(app, "")
}

// ToMarkdownWithSource creates a markdown string for the `*App`.
// If sourcePath points to a readable Go source file, long descriptions are
// extracted from leading doc comments above each flag's composite literal and
// merged into the rendered output. An empty sourcePath disables long
// description lookup.
// The function errors if either parsing or writing of the string fails.
func ToMarkdownWithSource(app *cli.Command, sourcePath string) (string, error) {
	var w bytes.Buffer

	tpls, err := template.New("cli").Funcs(plugin_template.LoadFuncMap()).ParseFS(templateFs, "**/*.tmpl")
	if err != nil {
		return "", err
	}

	if err := tpls.ExecuteTemplate(&w, "markdown.md.tmpl", GetTemplateDataWithSource(app, sourcePath)); err != nil {
		return "", err
	}

	return w.String(), nil
}

// GetTemplateData returns the template data for the `*App` without long
// descriptions. It is a convenience wrapper around GetTemplateDataWithSource
// that passes an empty sourcePath.
func GetTemplateData(app *cli.Command) *CliTemplate {
	return GetTemplateDataWithSource(app, "")
}

// GetTemplateDataWithSource returns the template data for the `*App`.
// If sourcePath points to a readable Go source file, long descriptions are
// extracted from leading doc comments above each flag's composite literal and
// merged into the returned data. An empty sourcePath disables long
// description lookup.
func GetTemplateDataWithSource(app *cli.Command, sourcePath string) *CliTemplate {
	return &CliTemplate{
		Name:        app.Name,
		Version:     app.Version,
		Description: prepareMultilineString(app.Description),
		Usage:       prepareMultilineString(app.Usage),
		UsageText:   prepareMultilineString(app.UsageText),
		GlobalArgs:  prepareArgsWithValues(app.VisibleFlags(), LongDescriptionsForWith(sourcePath)),
	}
}

func prepareMultilineString(s string) string {
	return strings.TrimRight(
		strings.TrimSpace(
			strings.ReplaceAll(s, "\n", " "),
		),
		".\r\n\t",
	)
}

func prepareArgsWithValues(flags []cli.Flag, longDescriptions map[string]*LongDescription) []*PluginArg {
	return parseFlags(flags, longDescriptions)
}

func parseFlags(flags []cli.Flag, longDescriptions map[string]*LongDescription) []*PluginArg {
	args := make([]*PluginArg, 0)
	namePrefix := "plugin_"

	for _, f := range flags {
		flag, ok := f.(cli.DocGenerationFlag)
		if !ok {
			continue
		}

		modArg := &PluginArg{}

		name := strings.ToLower(strings.TrimSpace(flag.GetEnvVars()[0]))
		if !strings.HasPrefix(name, namePrefix) {
			continue
		}

		modArg.Name = strings.TrimPrefix(name, namePrefix)
		modArg.Description = flag.GetUsage()
		modArg.LongDescription = LongDescriptionMarkdown(longDescriptions[modArg.Name])

		if rf, _ := f.(cli.RequiredFlag); ok {
			modArg.Required = rf.IsRequired()
		}

		if !modArg.Required && flag.IsDefaultVisible() {
			if s := flag.GetDefaultText(); s != "" {
				modArg.Default = s
			} else if flag.TypeName() == "bool" {
				modArg.Default = flag.GetValue()
			} else if flag.TakesValue() && flag.GetValue() != "" {
				modArg.Default = flag.GetValue()
			}
		}

		modArg.Type = parseType(reflect.TypeOf(f).String())

		args = append(args, modArg)
	}

	sort.SliceStable(args, func(i, j int) bool {
		return args[i].Name < args[j].Name
	})

	return args
}

// LongDescriptionMarkdown renders a structured LongDescription as a
// markdown-ready string. Line breaks inside a paragraph collapse to a
// single space so markdown renderers flow the text as one sentence; each
// paragraph is prefixed with "&emsp;" to visually align with the short
// Description line that precedes it.
func LongDescriptionMarkdown(d *LongDescription) string {
	if d.IsZero() {
		return ""
	}

	parts := make([]string, len(d.Paragraphs))
	for i, p := range d.Paragraphs {
		parts[i] = "&emsp;" + strings.Join(p, " ")
	}

	return strings.Join(parts, "\n\n")
}

// LongDescriptionYAMLBlock renders a structured LongDescription as the
// body of a YAML literal block scalar. Every source line is prefixed
// with indent; paragraphs are separated by a blank line so the visual
// structure of the original comment is preserved end-to-end; the
// separator blank lines carry no trailing whitespace.
//
// The function is the format-agnostic building block — it knows about
// YAML's `|` literal-block rules but not about any specific data
// schema. Consumers compose the surrounding "key: |" line themselves
// (the key name and nesting depth vary per consumer and determine the
// required indent).
func LongDescriptionYAMLBlock(d *LongDescription, indent string) string {
	if d.IsZero() {
		return ""
	}

	var b strings.Builder

	for i, p := range d.Paragraphs {
		if i > 0 {
			b.WriteString("\n\n")
		}

		for j, line := range p {
			if j > 0 {
				b.WriteByte('\n')
			}

			b.WriteString(indent)
			b.WriteString(line)
		}
	}

	return b.String()
}

func parseType(raw string) string {
	// Check for slice types
	if strings.Contains(raw, "SliceBase") {
		return "list"
	}

	// Check for map types
	if strings.Contains(raw, "MapFlag") || strings.Contains(raw, "StringMapFlag") {
		return "dict"
	}

	// Extract the type from the FlagBase generic parameters
	re := regexp.MustCompile(`\*cli\.FlagBase\[([^,]+),`)
	match := re.FindStringSubmatch(raw)

	if len(match) > 1 {
		baseType := match[1]

		// Handle array/slice types
		if strings.HasPrefix(baseType, "[]") {
			return "list"
		}

		// Handle map types
		if strings.HasPrefix(baseType, "map[") {
			return "dict"
		}

		// Handle basic types
		switch baseType {
		case "int", "int64", "uint", "uint64":
			return "integer"
		case "float64":
			return "float"
		default:
			return strings.ToLower(baseType)
		}
	}

	return ""
}
