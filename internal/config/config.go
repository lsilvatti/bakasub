package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// TouchlessRules defines automation behavior when conflicts are detected
type TouchlessRules struct {
	MultipleSubtitles string `json:"multiple_subtitles" mapstructure:"multiple_subtitles"` // "largest", "smallest", "skip"
	DefaultProfile    string `json:"default_profile" mapstructure:"default_profile"`       // Profile name to use
	MuxingStrategy    string `json:"muxing_strategy" mapstructure:"muxing_strategy"`       // "replace", "create_new"
}

// PromptProfile represents a translation prompt configuration
type PromptProfile struct {
	Name         string  `json:"name" mapstructure:"name"`
	SystemPrompt string  `json:"system_prompt" mapstructure:"system_prompt"`
	Temperature  float64 `json:"temperature" mapstructure:"temperature"`
	IsFactory    bool    `json:"is_factory" mapstructure:"is_factory"`
	IsLocked     bool    `json:"is_locked" mapstructure:"is_locked"`
	MediaType    string  `json:"media_type" mapstructure:"media_type"` // anime, movie, series, documentary, youtube
}

// Config represents the application configuration
type Config struct {
	// General Settings
	InterfaceLang string `json:"interface_lang" mapstructure:"interface_lang"` // en, pt-br, es
	TargetLang    string `json:"target_lang" mapstructure:"target_lang"`       // PT-BR, EN-US, etc.
	BinPath       string `json:"bin_path" mapstructure:"bin_path"`             // Path to binaries directory

	// AI Provider Settings
	AIProvider    string  `json:"ai_provider" mapstructure:"ai_provider"`       // openrouter, gemini, openai, local
	APIKey        string  `json:"api_key" mapstructure:"api_key"`               // API key or empty for local
	LocalEndpoint string  `json:"local_endpoint" mapstructure:"local_endpoint"` // For local LLM
	Model         string  `json:"model" mapstructure:"model"`                   // Selected model ID
	Temperature   float64 `json:"temperature" mapstructure:"temperature"`       // AI temperature (0.0-1.0)

	// Processing Settings
	RemoveHITags      bool    `json:"remove_hi_tags" mapstructure:"remove_hi_tags"`
	GlobalTemperature float64 `json:"global_temperature" mapstructure:"global_temperature"`

	// Automation
	TouchlessMode  bool           `json:"touchless_mode" mapstructure:"touchless_mode"`
	TouchlessRules TouchlessRules `json:"touchless_rules" mapstructure:"touchless_rules"`

	// Prompt Profiles
	PromptProfiles map[string]PromptProfile `json:"prompt_profiles" mapstructure:"prompt_profiles"`
	ActiveProfile  string                   `json:"active_profile" mapstructure:"active_profile"`

	// Advanced
	AutoCheckUpdates bool   `json:"auto_check_updates" mapstructure:"auto_check_updates"`
	LogLevel         string `json:"log_level" mapstructure:"log_level"` // info, debug
	SaveRawJSON      bool   `json:"save_raw_json" mapstructure:"save_raw_json"`
}

var (
	configPath = "config.json"
	instance   *Config
)

// GetFactoryProfiles returns the built-in, immutable prompt profiles
func GetFactoryProfiles() map[string]PromptProfile {
	return map[string]PromptProfile{
		"anime": {
			Name: "Anime (Factory Default)",
			SystemPrompt: `You are a professional translator specializing in Japanese animation. Your goal is to translate the JSON payload while preserving honorifics (e.g., -san, -kun) and attack names. Use the glossary terms: {{glossary}}. Avoid localizing jokes.

Context: This is an anime series with typical shounen/slice-of-life elements.
Target Audience: International anime fans who appreciate cultural nuances.

Rules:
1. Keep Japanese honorifics (-san, -kun, -chan, -sama, -senpai)
2. Preserve attack names and special techniques in original language
3. Translate dialogue naturally but maintain Japanese cultural context
4. Keep sound effects (e.g., "Ara ara~", "Yamete~") when culturally relevant
5. Maintain the emotional tone and energy of the original

Output Format: Return ONLY valid JSON array with the same structure as input.`,
			Temperature: 0.7,
			IsFactory:   true,
			IsLocked:    true,
			MediaType:   "anime",
		},
		"movie": {
			Name: "Movie (Factory Default)",
			SystemPrompt: `You are a professional subtitle translator for feature films. Translate the JSON payload with cinematic quality and proper pacing.

Context: This is a feature film requiring natural, flowing dialogue.
Target Audience: General moviegoers expecting professional subtitles.

Rules:
1. Translate dialogue naturally for the target language
2. Adapt idioms and cultural references for broad understanding
3. Maintain the emotional impact and dramatic timing
4. Keep subtitle length appropriate for reading speed
5. Preserve character voice and personality

Output Format: Return ONLY valid JSON array with the same structure as input.`,
			Temperature: 0.5,
			IsFactory:   true,
			IsLocked:    true,
			MediaType:   "movie",
		},
		"series": {
			Name: "Series (Factory Default)",
			SystemPrompt: `You are a professional translator for TV series. Translate the JSON payload maintaining consistency across episodes.

Context: This is a TV series requiring consistent character voices and terminology.
Target Audience: Series viewers expecting coherent, binge-worthy translations.

Rules:
1. Maintain character voice consistency throughout episodes
2. Keep recurring terminology consistent (check glossary: {{glossary}})
3. Translate naturally for episodic viewing
4. Preserve plot-relevant terms and references
5. Adapt humor appropriately for the target culture

Output Format: Return ONLY valid JSON array with the same structure as input.`,
			Temperature: 0.4,
			IsFactory:   true,
			IsLocked:    true,
			MediaType:   "series",
		},
		"documentary": {
			Name: "Documentary (Factory Default)",
			SystemPrompt: `You are a professional translator for documentary content. Translate the JSON payload with accuracy and clarity.

Context: This is educational/documentary content requiring precise translation.
Target Audience: Viewers seeking informative, accurate content.

Rules:
1. Prioritize accuracy and clarity over stylistic flourishes
2. Translate technical terms correctly
3. Maintain formal, educational tone
4. Keep proper nouns and names consistent
5. Preserve factual information exactly

Output Format: Return ONLY valid JSON array with the same structure as input.`,
			Temperature: 0.2,
			IsFactory:   true,
			IsLocked:    true,
			MediaType:   "documentary",
		},
		"youtube": {
			Name: "YouTube (Factory Default)",
			SystemPrompt: `You are a professional translator for online video content. Translate the JSON payload with energy and engagement.

Context: This is YouTube/online video content with casual, dynamic speech.
Target Audience: Online viewers expecting casual, engaging subtitles.

Rules:
1. Keep casual, conversational tone
2. Translate slang and internet expressions appropriately
3. Maintain creator's personality and energy
4. Preserve humor and memes when culturally transferable
5. Keep subtitle style snappy and readable

Output Format: Return ONLY valid JSON array with the same structure as input.`,
			Temperature: 0.6,
			IsFactory:   true,
			IsLocked:    true,
			MediaType:   "youtube",
		},
	}
}

// Default returns a Config with sensible defaults
func Default() *Config {
	return &Config{
		InterfaceLang:     "en",
		TargetLang:        "PT-BR",
		BinPath:           "./bin",
		AIProvider:        "openrouter",
		APIKey:            "",
		LocalEndpoint:     "http://localhost:11434",
		Model:             "google/gemini-flash-1.5",
		Temperature:       0.3,
		GlobalTemperature: 0.3,
		RemoveHITags:      true,
		TouchlessMode:     false,
		TouchlessRules: TouchlessRules{
			MultipleSubtitles: "largest",
			DefaultProfile:    "anime",
			MuxingStrategy:    "replace",
		},
		PromptProfiles:   GetFactoryProfiles(),
		ActiveProfile:    "anime",
		AutoCheckUpdates: true,
		LogLevel:         "info",
		SaveRawJSON:      false,
	}
}

// Exists checks if config file exists
func Exists() bool {
	_, err := os.Stat(configPath)
	return err == nil
}

// Load reads the configuration from config.json
func Load() (*Config, error) {
	if instance != nil {
		return instance, nil
	}

	// Set config file details
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/bakasub")

	// Try to read existing config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, return default
			instance = Default()
			return instance, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal into struct
	cfg := Default() // Start with defaults
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Merge factory profiles (they might have been updated)
	factoryProfiles := GetFactoryProfiles()
	if cfg.PromptProfiles == nil {
		cfg.PromptProfiles = make(map[string]PromptProfile)
	}
	for key, profile := range factoryProfiles {
		cfg.PromptProfiles[key] = profile
	}

	instance = cfg
	return instance, nil
}

// Save writes the configuration to config.json
func (c *Config) Save() error {
	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if configDir != "." && configDir != "" {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	// Set all values in viper
	viper.Set("interface_lang", c.InterfaceLang)
	viper.Set("target_lang", c.TargetLang)
	viper.Set("bin_path", c.BinPath)
	viper.Set("ai_provider", c.AIProvider)
	viper.Set("api_key", c.APIKey)
	viper.Set("local_endpoint", c.LocalEndpoint)
	viper.Set("model", c.Model)
	viper.Set("temperature", c.Temperature)
	viper.Set("remove_hi_tags", c.RemoveHITags)
	viper.Set("global_temperature", c.GlobalTemperature)
	viper.Set("touchless_mode", c.TouchlessMode)
	viper.Set("touchless_rules", c.TouchlessRules)
	viper.Set("prompt_profiles", c.PromptProfiles)
	viper.Set("active_profile", c.ActiveProfile)
	viper.Set("auto_check_updates", c.AutoCheckUpdates)
	viper.Set("log_level", c.LogLevel)
	viper.Set("save_raw_json", c.SaveRawJSON)

	// Write to file
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// CloneProfile creates a user copy of a profile
func (c *Config) CloneProfile(sourceKey, newName string) error {
	source, ok := c.PromptProfiles[sourceKey]
	if !ok {
		return fmt.Errorf("source profile not found: %s", sourceKey)
	}

	// Create new profile as user profile
	newProfile := source
	newProfile.Name = newName
	newProfile.IsFactory = false
	newProfile.IsLocked = false

	// Generate unique key
	newKey := fmt.Sprintf("user_%s", newName)
	c.PromptProfiles[newKey] = newProfile

	return nil
}

// DeleteProfile removes a user profile (cannot delete factory profiles)
func (c *Config) DeleteProfile(key string) error {
	profile, ok := c.PromptProfiles[key]
	if !ok {
		return fmt.Errorf("profile not found: %s", key)
	}

	if profile.IsFactory {
		return fmt.Errorf("cannot delete factory profile")
	}

	delete(c.PromptProfiles, key)
	return nil
}

// UpdateProfile modifies a user profile (cannot modify factory profiles)
func (c *Config) UpdateProfile(key string, updated PromptProfile) error {
	profile, ok := c.PromptProfiles[key]
	if !ok {
		return fmt.Errorf("profile not found: %s", key)
	}

	if profile.IsFactory {
		return fmt.Errorf("cannot modify factory profile")
	}

	c.PromptProfiles[key] = updated
	return nil
}
