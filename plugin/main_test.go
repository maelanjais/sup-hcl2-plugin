package main

import (
	"os"
	"strings"
	"testing"
)

func TestHCL2Parser_ParseFile(t *testing.T) {
	// Créer un fichier HCL2 temporaire pour le test
	hclContent := `
version = "0.5"

network "test-net" {
  hosts = ["127.0.0.1"]
  env {
    DEBUG = "true"
  }
}

command "echo-test" {
  desc = "A simple echo command"
  run  = "echo hello"
}
`
	tmpFile, err := os.CreateTemp("", "test-*.hcl")
	if err != nil {
		t.Fatalf("Impossible de créer le fichier temp: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(hclContent)
	if err != nil {
		t.Fatalf("Impossible d'écrire dans le fichier temp: %v", err)
	}
	tmpFile.Close()

	// Tester le parseur
	parser := &HCL2Parser{}
	yamlBytes, err := parser.ParseFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Erreur inattendue de ParseFile: %v", err)
	}

	yamlStr := string(yamlBytes)

	// Vérifications du YAML généré
	if !strings.Contains(yamlStr, "version: \"0.5\"") {
		t.Errorf("Le YAML ne contient pas la bonne version:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "hosts:") || !strings.Contains(yamlStr, "- 127.0.0.1") {
		t.Errorf("Le YAML ne contient pas les hosts attendus:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "run: echo hello") {
		t.Errorf("Le YAML ne contient pas la commande run attendue:\n%s", yamlStr)
	}
}

func TestExtractEnvFromBody_Empty(t *testing.T) {
	env := extractEnvFromBody(nil)
	if env != nil {
		t.Errorf("Attendu nil pour un body vide, reçu %v", env)
	}
}
