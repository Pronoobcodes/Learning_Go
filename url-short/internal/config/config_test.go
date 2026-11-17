package config

import "testing"

func TestLoadUsesDefaultDatabaseURLWhenEnvUnset(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("PORT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}
	if cfg.DatabaseURL != defaultDatabaseURL {
		t.Fatalf("expected default database URL %q, got %q", defaultDatabaseURL, cfg.DatabaseURL)
	}
	if cfg.Port != defaultPort {
		t.Fatalf("expected default port %q, got %q", defaultPort, cfg.Port)
	}
}

func TestLoadPrefersEnvironmentOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://custom:secret@localhost:5432/app")
	t.Setenv("PORT", "9000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}
	if cfg.DatabaseURL != "postgres://custom:secret@localhost:5432/app" {
		t.Fatalf("expected overridden database URL to be used, got %q", cfg.DatabaseURL)
	}
	if cfg.Port != "9000" {
		t.Fatalf("expected overridden port to be used, got %q", cfg.Port)
	}
}
