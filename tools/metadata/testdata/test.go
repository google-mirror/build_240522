package main

import (
	"os/exec"
	"testing"
)

func TestMetadataGeneration(t *testing.T) {
	cmd := exec.Command("metadata", "-inputFile", "testdata/inputFiles.txt", "-outputFile", "testdata/generatedOutputFile.txt")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Error running metadata command: %s", err)
	}

	// Compare the generated file contents with the expected output file
	expectedOutput, err := exec.Command("cat", "testdata/expectedOutputFile.txt").Output()
	if err != nil {
		t.Fatalf("Error reading expected output file: %s", err)
	}

	generatedOutput, err := exec.Command("cat", "testdata/generatedOutputFile.txt").Output()
	if err != nil {
		t.Fatalf("Error reading generated output file: %s", err)
	}

	if string(expectedOutput) != string(generatedOutput) {
		t.Errorf("Generated file contents do not match the expected output")
	}
}
