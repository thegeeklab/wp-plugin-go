package docs

import (
	"bytes"
	"embed"
	"html/template"
	"reflect"
	"regexp"
	"sort"
	"strings"

	plugin_template "github.com/thegeeklab/wp-plugin-go/v5/template"

	"github.com/urfave/cli/v3"
)

type PluginArg struct {
	Name        string
	EnvVars     []string
	Description string
	Default     string
	Type        string
	Required    bool
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

// ToMarkdown creates a markdown string for the `*App`
// The function errors if either parsing or writing of the string fails.
func ToMarkdown(app *cli.Command) (string, error) {
	var w bytes.Buffer

	tpls, err := template.New("cli").Funcs(plugin_template.LoadFuncMap()).ParseFS(templateFs, "**/*.tmpl")
	if err != nil {
		return "", err
	}

	if err := tpls.ExecuteTemplate(&w, "markdown.md.tmpl", GetTemplateData(app)); err != nil {
		return "", err
	}

	return w.String(), nil
}

func GetTemplateData(app *cli.Command) *CliTemplate {
	return &CliTemplate{
		Name:        app.Name,
		Version:     app.Version,
		Description: prepareMultilineString(app.Description),
		Usage:       prepareMultilineString(app.Usage),
		UsageText:   prepareMultilineString(app.UsageText),
		GlobalArgs:  prepareArgsWithValues(app.VisibleFlags()),
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

func prepareArgsWithValues(flags []cli.Flag) []*PluginArg {
	return parseFlags(flags)
}

func parseFlags(flags []cli.Flag) []*PluginArg {
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
		modArg.Default = flag.GetDefaultText()

		if rf, _ := f.(cli.RequiredFlag); ok {
			modArg.Required = rf.IsRequired()
		}

		modArg.Type = parseType(reflect.TypeOf(f).String())

		args = append(args, modArg)
	}

	sort.SliceStable(args, func(i, j int) bool {
		return args[i].Name < args[j].Name
	})

	return args
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
