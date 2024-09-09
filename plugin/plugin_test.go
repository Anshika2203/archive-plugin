package plugin

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestZipPlugin_Exec_Zip(t *testing.T) {
	sourceDir := t.TempDir()
	targetFile := filepath.Join(t.TempDir(), "output.zip")

	file, err := os.Create(filepath.Join(sourceDir, "test.txt"))
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	file.WriteString("hello, world")
	file.Close()

	plugin := ZipPlugin{
		Source: sourceDir,
		Target: targetFile,
	}

	err = plugin.Exec(context.Background())
	if err != nil {
		t.Fatalf("failed to zip: %v", err)
	}

	_, err = zip.OpenReader(targetFile)
	if err != nil {
		t.Fatalf("output is not a valid zip file: %v", err)
	}
}

func TestZipPlugin_Exec_Unzip(t *testing.T) {
	sourceZip := filepath.Join(t.TempDir(), "test.zip")
	targetDir := t.TempDir()

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	fileWriter, err := writer.Create("test.txt")
	if err != nil {
		t.Fatalf("failed to create file in zip: %v", err)
	}
	fileWriter.Write([]byte("hello, world"))
	writer.Close()

	err = os.WriteFile(sourceZip, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("failed to write zip file: %v", err)
	}

	plugin := ZipPlugin{
		Source: sourceZip,
		Target: targetDir,
	}

	// Execute the plugin
	err = plugin.Exec(context.Background())
	if err != nil {
		t.Fatalf("failed to unzip: %v", err)
	}

	extractedFile := filepath.Join(targetDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(content) != "hello, world" {
		t.Errorf("extracted content mismatch: got %v, want %v", string(content), "hello, world")
	}
}

func TestZipPlugin_Zip(t *testing.T) {
	sourceDir := t.TempDir()
	targetFile := filepath.Join(t.TempDir(), "output.zip")

	file, err := os.Create(filepath.Join(sourceDir, "test.txt"))
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	file.WriteString("dummy content")
	file.Close()

	plugin := ZipPlugin{
		Source: sourceDir,
		Target: targetFile,
	}

	err = plugin.Zip()
	if err != nil {
		t.Fatalf("failed to zip: %v", err)
	}

	_, err = zip.OpenReader(targetFile)
	if err != nil {
		t.Fatalf("output is not a valid zip file: %v", err)
	}
}

func TestZipPlugin_Unzip(t *testing.T) {
	sourceZip := filepath.Join(t.TempDir(), "test.zip")
	targetDir := t.TempDir()

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	fileWriter, err := writer.Create("test.txt")
	if err != nil {
		t.Fatalf("failed to create file in zip: %v", err)
	}
	fileWriter.Write([]byte("test content"))
	writer.Close()

	// Write zip data to source file
	err = os.WriteFile(sourceZip, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("failed to write zip file: %v", err)
	}

	plugin := ZipPlugin{
		Source: sourceZip,
		Target: targetDir,
	}

	err = plugin.Unzip()
	if err != nil {
		t.Fatalf("failed to unzip: %v", err)
	}

	extractedFile := filepath.Join(targetDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("extracted content mismatch: got %v, want %v", string(content), "test content")
	}
}

func TestZipPlugin_Exec_InvalidSource(t *testing.T) {
	plugin := ZipPlugin{
		Source: "invalid/path",
		Target: t.TempDir(),
	}

	err := plugin.Exec(context.Background())
	if err == nil {
		t.Fatalf("expected error for invalid source, got nil")
	}

	if !strings.Contains(err.Error(), "error accessing source") {
		t.Errorf("unexpected error message: %v", err)
	}
}
