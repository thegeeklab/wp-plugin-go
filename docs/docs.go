package docs

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

type CliTemplate struct {
	Name         string
	Description  string
	Usage        string
	UsageText    string
	SectionNum   int
	Commands     []string
	GlobalArgs   []string
	SynopsisArgs []string
}

//go:embed templates
var templateFs embed.FS

// ToMarkdown creates a markdown string for the `*App`
// The function errors if either parsing or writing of the string fails.
func ToMarkdown(app *cli.App) (string, error) {
	var w bytes.Buffer

	tpls, err := template.New("").ParseFS(templateFs, "templates/markdown.md.tmpl")
	if err != nil {
		return "", err
	}

	if err := tpls.Execute(&w, GetTemplateData(app, 0)); err != nil {
		return "", err
	}

	return w.String(), nil
}

func GetTemplateData(app *cli.App, sectionNum int) *CliTemplate {
	return &CliTemplate{
		Name:         app.Name,
		Description:  prepareMultilineString(app.Description),
		Usage:        prepareMultilineString(app.Usage),
		UsageText:    app.UsageText,
		SectionNum:   sectionNum,
		Commands:     prepareCommands(app.Commands, 0),
		GlobalArgs:   prepareArgsWithValues(app.VisibleFlags()),
		SynopsisArgs: prepareArgsSynopsis(app.VisibleFlags()),
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

func prepareCommands(commands []*cli.Command, level int) []string {
	coms := make([]string, 0)

	for _, command := range commands {
		if command.Hidden {
			continue
		}

		usageText := prepareUsageText(command)

		usage := prepareUsage(command, usageText)

		prepared := fmt.Sprintf("%s %s\n\n%s%s",
			strings.Repeat("#", level+2), //nolint:gomnd
			strings.Join(command.Names(), ", "),
			usage,
			usageText,
		)

		flags := prepareArgsWithValues(command.VisibleFlags())
		if len(flags) > 0 {
			prepared += fmt.Sprintf("\n%s", strings.Join(flags, "\n"))
		}

		coms = append(coms, prepared)

		// recursively iterate subcommands
		if len(command.Subcommands) > 0 {
			coms = append(
				coms,
				prepareCommands(command.Subcommands, level+1)...,
			)
		}
	}

	return coms
}

func prepareArgsWithValues(flags []cli.Flag) []string {
	return prepareFlags(flags, ", ", "**", "**", `""`, true)
}

func prepareArgsSynopsis(flags []cli.Flag) []string {
	return prepareFlags(flags, "|", "[", "]", "[value]", false)
}

func prepareFlags(
	flags []cli.Flag,
	sep, opener, closer, value string,
	addDetails bool,
) []string {
	args := []string{}

	for _, f := range flags {
		flag, ok := f.(cli.DocGenerationFlag)
		if !ok {
			continue
		}

		modifiedArg := opener

		for _, s := range flag.Names() {
			trimmed := strings.TrimSpace(s)

			if len(modifiedArg) > len(opener) {
				modifiedArg += sep
			}

			if len(trimmed) > 1 {
				modifiedArg += fmt.Sprintf("--%s", trimmed)
			} else {
				modifiedArg += fmt.Sprintf("-%s", trimmed)
			}
		}

		modifiedArg += closer

		if flag.TakesValue() {
			modifiedArg += fmt.Sprintf("=%s", value)
		}

		if addDetails {
			modifiedArg += flagDetails(flag)
		}

		args = append(args, modifiedArg+"\n")
	}

	sort.Strings(args)

	return args
}

// flagDetails returns a string containing the flags metadata.
func flagDetails(flag cli.DocGenerationFlag) string {
	description := flag.GetUsage()

	if flag.TakesValue() {
		defaultText := flag.GetDefaultText()
		if defaultText == "" {
			defaultText = flag.GetValue()
		}

		if defaultText != "" {
			description += " (default: " + defaultText + ")"
		}
	}

	return ": " + description
}

func prepareUsageText(command *cli.Command) string {
	if command.UsageText == "" {
		return ""
	}

	// Remove leading and trailing newlines
	preparedUsageText := strings.Trim(command.UsageText, "\n")

	var usageText string

	if strings.Contains(preparedUsageText, "\n") {
		// Format multi-line string as a code block using the 4 space schema to allow for embedded markdown such
		// that it will not break the continuous code block.
		for _, ln := range strings.Split(preparedUsageText, "\n") {
			usageText += fmt.Sprintf("    %s\n", ln)
		}
	} else {
		// Style a single line as a note
		usageText = fmt.Sprintf(">%s\n", preparedUsageText)
	}

	return usageText
}

func prepareUsage(command *cli.Command, usageText string) string {
	if command.Usage == "" {
		return ""
	}

	usage := command.Usage + "\n"
	// Add a newline to the Usage IFF there is a UsageText
	if usageText != "" {
		usage += "\n"
	}

	return usage
}
