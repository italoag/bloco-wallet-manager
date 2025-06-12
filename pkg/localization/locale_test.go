package localization

import (
	"blocowallet/pkg/config"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"os"
	"path/filepath"
	"testing"
)

func TestInitLocalization(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "locale_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config
	cfg := &config.Config{
		AppDir:    tempDir,
		Language:  "en",
		LocaleDir: filepath.Join(tempDir, "locale"),
	}

	// Initialize localization
	err = InitLocalization(cfg)
	if err != nil {
		t.Fatalf("InitLocalization failed: %v", err)
	}

	// Check if the locale directory was created
	if _, err := os.Stat(cfg.LocaleDir); os.IsNotExist(err) {
		t.Errorf("Locale directory was not created")
	}

	// Check if the default language files were created
	languages := []string{"en", "pt", "es"}
	for _, lang := range languages {
		filename := fmt.Sprintf("language.%s.toml", lang)
		filePath := filepath.Join(cfg.LocaleDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Language file %s was not created", filename)
		}
	}
}

func TestT(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "locale_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config
	cfg := &config.Config{
		AppDir:    tempDir,
		Language:  "en",
		LocaleDir: filepath.Join(tempDir, "locale"),
	}

	// Initialize localization
	err = InitLocalization(cfg)
	if err != nil {
		t.Fatalf("InitLocalization failed: %v", err)
	}

	// Test translation without template data
	translation := T("welcome_message", nil)
	if translation == "" || translation == "welcome_message" {
		t.Errorf("Translation failed for welcome_message")
	}

	// Test translation with template data
	data := map[string]interface{}{
		"View": "Wallets",
	}
	translation = T("status_bar_instructions", data)
	if translation == "" || translation == "status_bar_instructions" {
		t.Errorf("Translation with template data failed for status_bar_instructions")
	}
	if translation == "View: {{.View}} | Press 'esc' to return | Press 'q' to quit" {
		t.Errorf("Template was not processed correctly")
	}
}

func TestTP(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "locale_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config
	cfg := &config.Config{
		AppDir:    tempDir,
		Language:  "en",
		LocaleDir: filepath.Join(tempDir, "locale"),
	}

	// Initialize localization
	err = InitLocalization(cfg)
	if err != nil {
		t.Fatalf("InitLocalization failed: %v", err)
	}

	// Create a test plural message
	testMsg := `
[cats]
one = "{{.Name}} has {{.Count}} cat."
other = "{{.Name}} has {{.Count}} cats."
`
	testFile := filepath.Join(cfg.LocaleDir, "test_plural.toml")
	err = os.WriteFile(testFile, []byte(testMsg), 0644)
	if err != nil {
		t.Fatalf("Failed to create test plural message: %v", err)
	}

	// Load the test message file directly into the bundle
	_, err = bundle.LoadMessageFile(testFile)
	if err != nil {
		t.Fatalf("Failed to load test message file: %v", err)
	}

	// Test plural translation with count=1
	data := map[string]interface{}{
		"Name": "John",
	}

	// Debug: Print the message file content
	fileContent, _ := os.ReadFile(testFile)
	t.Logf("Test file content: %s", fileContent)

	// Debug: Try to get the message directly from the localizer
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "cats",
			One: "{{.Name}} has {{.Count}} cat.",
			Other: "{{.Name}} has {{.Count}} cats.",
		},
		TemplateData: map[string]interface{}{
			"Name":  "John",
			"Count": 1,
		},
		PluralCount: 1,
	})
	if err != nil {
		t.Logf("Error localizing directly: %v", err)
	} else {
		t.Logf("Direct localizer localization: %s", msg)
	}

	// Use our TP function
	translation := TP("cats", 1, data)
	t.Logf("TP function result: %s", translation)

	if translation == "" || translation == "cats" {
		t.Errorf("Plural translation failed for cats with count=1")
	}
	if translation != "John has 1 cat." {
		t.Errorf("Incorrect plural translation for count=1: %s", translation)
	}

	// Test plural translation with count=2
	translation = TP("cats", 2, data)
	t.Logf("TP function result (count=2): %s", translation)

	if translation == "" || translation == "cats" {
		t.Errorf("Plural translation failed for cats with count=2")
	}
	if translation != "John has 2 cats." {
		t.Errorf("Incorrect plural translation for count=2: %s", translation)
	}
}

func TestChangeLanguage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "locale_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config
	cfg := &config.Config{
		AppDir:    tempDir,
		Language:  "en",
		LocaleDir: filepath.Join(tempDir, "locale"),
	}

	// Initialize localization
	err = InitLocalization(cfg)
	if err != nil {
		t.Fatalf("InitLocalization failed: %v", err)
	}

	// Get English translation
	enTranslation := T("welcome_message", nil)

	// Change language to Portuguese
	ChangeLanguage("pt")

	// Get Portuguese translation
	ptTranslation := T("welcome_message", nil)

	// Translations should be different
	if enTranslation == ptTranslation {
		t.Errorf("Language change did not affect translations")
	}

	// Change language to Spanish
	ChangeLanguage("es")

	// Get Spanish translation
	esTranslation := T("welcome_message", nil)

	// Translations should be different
	if enTranslation == esTranslation || ptTranslation == esTranslation {
		t.Errorf("Language change did not affect translations")
	}
}
