// Copyright (c) 2018, Drone.IO Inc
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package template

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

// Render parses and executes a template, returning the results in string
// format. Trailing or leading spaces or new-lines are not getting truncated. It
// is able to read templates from remote paths, local files or directly from the
// string.
func Render(ctx context.Context, client http.Client, tmpl string, payload interface{}) (string, error) {
	var outString bytes.Buffer

	tpl := new(template.Template).Funcs(LoadFuncMap())

	templateURL, err := url.Parse(tmpl)
	if err == nil {
		switch templateURL.Scheme {
		case "http", "https":
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, templateURL.String(), nil)
			if err != nil {
				return "", fmt.Errorf("failed to create request: %w", err)
			}

			res, err := client.Do(req)
			if err != nil {
				return "", fmt.Errorf("failed to fetch: %w", err)
			}

			defer res.Body.Close()

			out, err := io.ReadAll(res.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read: %w", err)
			}

			tmpl = string(out)
		case "file":
			out, err := os.ReadFile(templateURL.Path)
			if err != nil {
				return "", fmt.Errorf("failed to read: %w", err)
			}

			tmpl = string(out)
		}
	}

	tpl, err = tpl.Parse(tmpl)
	if err != nil {
		return "", err
	}

	err = tpl.Execute(&outString, payload)

	return outString.String(), err
}

// RenderTrim parses and executes a template, returning the results in string
// format. The result is trimmed to remove left and right padding and newlines
// that may be added unintentially in the template markup.
func RenderTrim(ctx context.Context, client http.Client, tmpl string, playload interface{}) (string, error) {
	out, err := Render(ctx, client, tmpl, playload)

	return strings.Trim(out, " \n"), err
}

// AddPrefix prepends the prefix string to the input string if the input
// string does not already have the prefix. It trims any leading or trailing
// space from the input before checking for the prefix.
func AddPrefix(prefix, input string) string {
	if strings.TrimSpace(input) == "" {
		return input
	}

	if strings.HasPrefix(input, prefix) {
		return input
	}

	return prefix + input
}
