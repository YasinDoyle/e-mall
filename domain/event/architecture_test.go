package event

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEventHandlersDoNotImportService(t *testing.T) {
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
			violations = append(violations, filepath.Join("domain/event", file))
		}
	}
	if len(violations) > 0 {
		t.Fatalf("domain event handlers must use dao/repository adapters instead of service package:\n%s", strings.Join(violations, "\n"))
	}
}
