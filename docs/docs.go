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

type PluginArg struct {
	Name            string
	EnvVars         []string
	Description     string
	LongDescription string
	Default         string
	Type            string
	Required        bool
}

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
		GlobalArgs:  prepareArgsWithValues(app.VisibleFlags(), LongDescriptionsFor(sourcePath)),
	}
}

// LongDescriptionsFor returns the long descriptions extracted from the Go
// source file at sourcePath, keyed by normalized flag name. It returns an
// empty map when sourcePath is empty or the file cannot be parsed so callers
// can pass an unset path without failing the whole pipeline.
func LongDescriptionsFor(sourcePath string) map[string]string {
	if sourcePath == "" {
		return map[string]string{}
	}

	longs, err := LongDescriptions(sourcePath)
	if err != nil {
		return map[string]string{}
	}

	normalized := make(map[string]string, len(longs))
	for name, desc := range longs {
		normalized[normalizeFlagName(name)] = desc
	}

	return normalized
}

func normalizeFlagName(name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(name, "-", "_"), ".", "_")
}

func prepareMultilineString(s string) string {
	return strings.TrimRight(
		strings.TrimSpace(
			strings.ReplaceAll(s, "\n", " "),
		),
		".\r\n\t",
	)
}

func prepareArgsWithValues(flags []cli.Flag, longDescriptions map[string]string) []*PluginArg {
	return parseFlags(flags, longDescriptions)
}

func parseFlags(flags []cli.Flag, longDescriptions map[string]string) []*PluginArg {
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
		modArg.LongDescription = formatLongDescription(longDescriptions[modArg.Name])

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

// formatLongDescription converts a long description from LongDescriptions
// (paragraphs separated by "\n\n") into a markdown-ready string with each
// paragraph prefixed by "&emsp;". Returns an empty string for empty input.
func formatLongDescription(s string) string {
	if s == "" {
		return ""
	}

	paragraphs := strings.Split(s, "\n\n")

	for i, p := range paragraphs {
		paragraphs[i] = "&emsp;" + p
	}

	return strings.Join(paragraphs, "\n\n")
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
