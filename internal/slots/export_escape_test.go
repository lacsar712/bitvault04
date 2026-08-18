package slots

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestExportImageFileEscapeVariants is an independent regression test for the
// path-traversal fix in ExportImageFile. It is intentionally separate from
// runtime_test.go so the original escape case stays pinned while these cover
// nested traversal, a lone "..", and a non-escaping path that must still work.
func TestExportImageFileEscapeVariants(t *testing.T) {
	root := t.TempDir()

	escapes := []string{
		filepath.Join("..", "boot"),
		filepath.Join(".."),
		filepath.Join("a", "..", "..", "boot"),
		filepath.Join("x", "..", "..", "..", "secret"),
	}
	for _, rel := range escapes {
		rel := rel
		t.Run("escape/"+rel, func(t *testing.T) {
			if _, err := ExportImageFile(root, rel); err == nil {
				t.Fatalf("expected escape to be rejected for %q", rel)
			}
		})
	}

	// Legitimate relative paths must still resolve inside the root.
	for _, rel := range []string{"image.bin", filepath.Join("sub", "image.bin"), "..boot"} {
		got, err := ExportImageFile(root, rel)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", rel, err)
		}
		// The resolved path must stay within root (no leading "..").
		relToRoot, err := filepath.Rel(root, got)
		if err != nil {
			t.Fatalf("Rel error for %q: %v", rel, err)
		}
		if relToRoot == ".." || strings.HasPrefix(relToRoot, ".."+string(filepath.Separator)) {
			t.Fatalf("path %q escaped root (rel=%q)", rel, relToRoot)
		}
	}
}
