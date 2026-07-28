package application

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplicationDoesNotImportService(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var violations []string
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(content), `"github.com/YasinDoyle/e-mall/service"`) {
			violations = append(violations, filepath.Join("application", file))
		}
	}
	if len(violations) > 0 {
		t.Fatalf("application/usecase layer must not import service package:\n%s", strings.Join(violations, "\n"))
	}
}
