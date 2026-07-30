package service

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestServiceMethodsDoNotCallOtherServiceGetters(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	getterCall := regexp.MustCompile(`\bGet[A-Za-z0-9]+Srv\(\)`)
	getterDefinition := regexp.MustCompile(`func\s+Get[A-Za-z0-9]+Srv\(\)`)
	var violations []string
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for lineNo, line := range strings.Split(string(content), "\n") {
			if getterCall.MatchString(line) && !getterDefinition.MatchString(line) {
				violations = append(violations, filepath.Join("service", file)+":"+itoa(lineNo+1)+": "+strings.TrimSpace(line))
			}
		}
	}
	if len(violations) > 0 {
		t.Fatalf("service methods must not call other service getters; use application/usecase orchestration or domain events instead:\n%s", strings.Join(violations, "\n"))
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
