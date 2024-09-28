package generator

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const defaultTemplate = `<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title> {{ Title }} </title>
</head>

<body>
{{ Content }}
</body>

</html>`

func ValidateTemplateFile(file string) error {
	info, err := os.Stat(file)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open template file: %s", err.Error()))
	}

	if info.IsDir() {
		return errors.New("Template file is a directory")
	}

	template, err := ReadFileS(file)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read template file: %s", err.Error()))
	}
	template = strings.TrimSpace(template)

	templateTags := []string{
		"{{ Title }}", "{{ Content }}",
	}
	for _, tag := range templateTags {
		if !strings.Contains(template, tag) {
			return errors.New(fmt.Sprintf("Invalid template: %s tag not found.", tag))
		}
	}

	htmlTags := []string{
		"<!DOCTYPE html>",
		"<html", ">",
		"<head>", "</head>",
		"<body>", "</body>",
		"</html>",
	}
	for _, tag := range htmlTags {
		index := strings.Index(template, tag)
		if index == -1 {
			return errors.New(fmt.Sprintf("Invalid template: %s tag not found.", tag))
		}
		template = template[index:]
	}

	return nil
}

// Splice body and title into template string
// body: Body of html page. Inserted at {{ Content }} tag.
// title: Title of html page. Inserted at {{ Title }} tag.
// template: (Optional) Template string. If empty, defaultTemplate is used.
func PopulateTemplate(body, title, template string) string {
	var result string

	if template == "" {
		template = defaultTemplate
	}

	result = strings.Join(
		strings.Split(template, "{{ Title }}"),
		title,
	)
	result = strings.Join(
		strings.Split(result, "{{ Content }}"),
		body,
	)

	return result
}
