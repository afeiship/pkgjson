package pkgjson_test

import (
	"encoding/json"
	pkgjson "github.com/afeiship/pkgjson"
	"os"
	"path/filepath"
	"testing"

	"github.com/iancoleman/orderedmap"
)

func TestPackageJSON(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "package.json")

	// Create initial package.json content
	initialData := orderedmap.New()
	initialData.Set("name", "test-package")
	initialData.Set("version", "1.0.0")

	bytes, err := json.MarshalIndent(initialData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal initial data: %v", err)
	}

	if err := os.WriteFile(testFile, bytes, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test NewPackageJSON
	pj, err := pkgjson.NewPackageJSON(testFile)
	if err != nil {
		t.Fatalf("NewPackageJSON failed: %v", err)
	}

	// Test Read
	name, _ := pj.Data.Get("name")
	if name != "test-package" {
		t.Errorf("Expected name to be 'test-package', got %v", name)
	}

	// Test Update
	if err := pj.Update("description", "Test package"); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	desc, exists := pj.Data.Get("description")
	if !exists || desc != "Test package" {
		t.Error("Update did not set description correctly")
	}

	// Test Delete
	if err := pj.Delete("description"); err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, exists = pj.Data.Get("description")
	if exists {
		t.Error("Delete did not remove description")
	}

	// Test Save
	if err := pj.Save(); err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// Verify saved content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	newPJ := orderedmap.New()
	if err := json.Unmarshal(content, newPJ); err != nil {
		t.Fatalf("Failed to unmarshal saved content: %v", err)
	}

	name, _ = newPJ.Get("name")
	if name != "test-package" {
		t.Errorf("Saved content does not match expected: got %v", name)
	}
}
