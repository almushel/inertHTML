package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestValidateTemplateFile(t *testing.T) {
	template := strings.Clone(defaultTemplate)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err = ValidateTemplateFile(""); err == nil {
		t.Fatalf("Validation succeeded with empty file path")
	}

	removeList := []string{
		"",
		"<html>", "</html>",
		"<body>", "</body>",
		"{{ Title }}", "{{ Content }}",
	}

	for _, removeStr := range removeList {
		t.Run(fmt.Sprintf("Template missing %s", removeStr),
			func(t *testing.T) {
				testFile, err := os.CreateTemp(wd, "template_file")
				if err != nil {
					t.Fatal(err)
				}

				_, err = testFile.WriteString(strings.Replace(template, removeStr, "", -1))
				if err != nil {
					t.Fatalf("Failed to write to testfile %s", testFile.Name())
				}

				err = ValidateTemplateFile(testFile.Name())
				if removeStr == "" {
					if err != nil {
						t.Fatalf("Validation failed with default template")
					}
				} else if err == nil {
					t.Fatalf("Validation passed with missing %s", removeStr)
				}

				os.Remove(testFile.Name())
			})
	}
}
