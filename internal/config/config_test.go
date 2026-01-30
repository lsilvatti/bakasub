package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("Default() returned nil")
	}

	if cfg.InterfaceLang != "en" {
		t.Errorf("expected InterfaceLang 'en', got %q", cfg.InterfaceLang)
	}

	if cfg.TargetLang != "PT-BR" {
		t.Errorf("expected TargetLang 'PT-BR', got %q", cfg.TargetLang)
	}

	if cfg.AIProvider != "openrouter" {
		t.Errorf("expected AIProvider 'openrouter', got %q", cfg.AIProvider)
	}

	if cfg.Temperature != 0.3 {
		t.Errorf("expected Temperature 0.3, got %f", cfg.Temperature)
	}

	if cfg.RemoveHITags != true {
		t.Error("expected RemoveHITags to be true")
	}

	if cfg.TouchlessMode != false {
		t.Error("expected TouchlessMode to be false")
	}

	if cfg.ActiveProfile != "anime" {
		t.Errorf("expected ActiveProfile 'anime', got %q", cfg.ActiveProfile)
	}

	if cfg.AutoCheckUpdates != true {
		t.Error("expected AutoCheckUpdates to be true")
	}

	if cfg.LogLevel != "info" {
		t.Errorf("expected LogLevel 'info', got %q", cfg.LogLevel)
	}
}

func TestDefaultTouchlessRules(t *testing.T) {
	cfg := Default()
	if cfg.TouchlessRules.MultipleSubtitles != "largest" {
		t.Errorf("expected MultipleSubtitles 'largest', got %q", cfg.TouchlessRules.MultipleSubtitles)
	}

	if cfg.TouchlessRules.DefaultProfile != "anime" {
		t.Errorf("expected DefaultProfile 'anime', got %q", cfg.TouchlessRules.DefaultProfile)
	}

	if cfg.TouchlessRules.MuxingStrategy != "replace" {
		t.Errorf("expected MuxingStrategy 'replace', got %q", cfg.TouchlessRules.MuxingStrategy)
	}
}

func TestGetFactoryProfiles(t *testing.T) {
	profiles := GetFactoryProfiles()
	expectedProfiles := []string{"anime", "movie", "series", "documentary", "youtube"}

	for _, name := range expectedProfiles {
		profile, ok := profiles[name]
		if !ok {
			t.Errorf("expected factory profile %q to exist", name)
			continue
		}

		if !profile.IsFactory {
			t.Errorf("profile %q should be marked as factory", name)
		}

		if !profile.IsLocked {
			t.Errorf("profile %q should be marked as locked", name)
		}

		if profile.SystemPrompt == "" {
			t.Errorf("profile %q should have a system prompt", name)
		}

		if profile.Name == "" {
			t.Errorf("profile %q should have a display name", name)
		}

		if profile.MediaType == "" {
			t.Errorf("profile %q should have a media type", name)
		}
	}

	if len(profiles) != 5 {
		t.Errorf("expected 5 factory profiles, got %d", len(profiles))
	}
}

func TestFactoryProfilesTemperature(t *testing.T) {
	profiles := GetFactoryProfiles()

	docTemp := profiles["documentary"].Temperature
	if docTemp != 0.2 {
		t.Errorf("expected documentary temperature 0.2, got %f", docTemp)
	}

	animeTemp := profiles["anime"].Temperature
	if animeTemp != 0.7 {
		t.Errorf("expected anime temperature 0.7, got %f", animeTemp)
	}
}

func TestExists(t *testing.T) {
	originalPath := configPath
	configPath = "nonexistent_config_test.json"
	defer func() { configPath = originalPath }()

	if Exists() {
		t.Error("Exists() should return false for non-existent file")
	}

	tmpDir := t.TempDir()
	tmpConfig := filepath.Join(tmpDir, "config.json")
	configPath = tmpConfig
	if err := os.WriteFile(tmpConfig, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	if !Exists() {
		t.Error("Exists() should return true for existing file")
	}
}

func TestCloneProfile(t *testing.T) {
	cfg := Default()
	err := cfg.CloneProfile("anime", "My Anime Profile")

	if err != nil {
		t.Fatalf("CloneProfile failed: %v", err)
	}

	cloned, ok := cfg.PromptProfiles["user_My Anime Profile"]
	if !ok {
		t.Fatal("cloned profile not found")
	}

	if cloned.IsFactory {
		t.Error("cloned profile should not be marked as factory")
	}

	if cloned.IsLocked {
		t.Error("cloned profile should not be locked")
	}

	if cloned.Name != "My Anime Profile" {
		t.Errorf("expected name 'My Anime Profile', got %q", cloned.Name)
	}

	original := cfg.PromptProfiles["anime"]
	if cloned.SystemPrompt != original.SystemPrompt {
		t.Error("cloned profile should have the same system prompt")
	}
}

func TestCloneProfileNonExistent(t *testing.T) {
	cfg := Default()
	err := cfg.CloneProfile("nonexistent", "Test")
	if err == nil {
		t.Error("CloneProfile should fail for non-existent source profile")
	}
}

func TestDeleteProfile(t *testing.T) {
	cfg := Default()
	cfg.CloneProfile("movie", "Test Profile")

	err := cfg.DeleteProfile("user_Test Profile")
	if err != nil {
		t.Fatalf("DeleteProfile failed: %v", err)
	}

	if _, ok := cfg.PromptProfiles["user_Test Profile"]; ok {
		t.Error("profile should have been deleted")
	}
}

func TestDeleteFactoryProfile(t *testing.T) {
	cfg := Default()
	err := cfg.DeleteProfile("anime")
	if err == nil {
		t.Error("DeleteProfile should fail for factory profiles")
	}
}

func TestDeleteNonExistentProfile(t *testing.T) {
	cfg := Default()
	err := cfg.DeleteProfile("nonexistent")
	if err == nil {
		t.Error("DeleteProfile should fail for non-existent profile")
	}
}

func TestUpdateProfile(t *testing.T) {
	cfg := Default()
	cfg.CloneProfile("series", "Editable")

	updated := cfg.PromptProfiles["user_Editable"]
	updated.Temperature = 0.9
	updated.SystemPrompt = "Custom prompt"
	err := cfg.UpdateProfile("user_Editable", updated)

	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}

	result := cfg.PromptProfiles["user_Editable"]
	if result.Temperature != 0.9 {
		t.Errorf("expected temperature 0.9, got %f", result.Temperature)
	}

	if result.SystemPrompt != "Custom prompt" {
		t.Errorf("expected custom prompt, got %q", result.SystemPrompt)
	}
}

func TestUpdateFactoryProfile(t *testing.T) {
	cfg := Default()
	updated := cfg.PromptProfiles["anime"]
	updated.Temperature = 1.0
	err := cfg.UpdateProfile("anime", updated)
	if err == nil {
		t.Error("UpdateProfile should fail for factory profiles")
	}
}

func TestUpdateNonExistentProfile(t *testing.T) {
	cfg := Default()
	err := cfg.UpdateProfile("nonexistent", PromptProfile{})
	if err == nil {
		t.Error("UpdateProfile should fail for non-existent profile")
	}
}

func TestConfigSave(t *testing.T) {
	tmpDir := t.TempDir()
	tmpConfig := filepath.Join(tmpDir, "config.json")
	originalPath := configPath
	configPath = tmpConfig
	defer func() { configPath = originalPath }()

	cfg := Default()
	cfg.TargetLang = "ES"
	cfg.AIProvider = "gemini"
	cfg.Model = "gemini-1.5-pro"
	err := cfg.Save()

	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(tmpConfig); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	content, err := os.ReadFile(tmpConfig)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if len(content) == 0 {
		t.Error("config file should not be empty")
	}
}

func TestTouchlessRulesStruct(t *testing.T) {
	rules := TouchlessRules{
		MultipleSubtitles: "smallest",
		DefaultProfile:    "movie",
		MuxingStrategy:    "create_new",
	}

	if rules.MultipleSubtitles != "smallest" {
		t.Errorf("unexpected MultipleSubtitles: %q", rules.MultipleSubtitles)
	}

	if rules.DefaultProfile != "movie" {
		t.Errorf("unexpected DefaultProfile: %q", rules.DefaultProfile)
	}

	if rules.MuxingStrategy != "create_new" {
		t.Errorf("unexpected MuxingStrategy: %q", rules.MuxingStrategy)
	}
}

func TestPromptProfileStruct(t *testing.T) {
	profile := PromptProfile{
		Name:         "Test Profile",
		SystemPrompt: "You are a translator",
		Temperature:  0.5,
		IsFactory:    false,
		IsLocked:     false,
		MediaType:    "anime",
	}

	if profile.Name != "Test Profile" {
		t.Errorf("unexpected Name: %q", profile.Name)
	}

	if profile.Temperature != 0.5 {
		t.Errorf("unexpected Temperature: %f", profile.Temperature)
	}

	if profile.IsFactory {
		t.Error("IsFactory should be false")
	}

	if profile.IsLocked {
		t.Error("IsLocked should be false")
	}
}

func TestConfigStruct(t *testing.T) {
	cfg := &Config{
		InterfaceLang:     "pt-br",
		TargetLang:        "EN-US",
		BinPath:           "/usr/local/bin",
		AIProvider:        "openai",
		APIKey:            "sk-test-key",
		LocalEndpoint:     "http://localhost:8080",
		Model:             "gpt-4o",
		Temperature:       0.5,
		GlobalTemperature: 0.5,
		RemoveHITags:      false,
		TouchlessMode:     true,
		AutoCheckUpdates:  false,
		LogLevel:          "debug",
		SaveRawJSON:       true,
	}

	if cfg.InterfaceLang != "pt-br" {
		t.Errorf("unexpected InterfaceLang: %q", cfg.InterfaceLang)
	}

	if cfg.TouchlessMode != true {
		t.Error("TouchlessMode should be true")
	}

	if cfg.RemoveHITags != false {
		t.Error("RemoveHITags should be false")
	}

	if cfg.SaveRawJSON != true {
		t.Error("SaveRawJSON should be true")
	}
}
