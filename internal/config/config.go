package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"grimoire/internal/domain/player"

	"github.com/joho/godotenv"
)

type Config struct {
	Token  string
	DBPath string
	Names  []string
}

func Load() (Config, error) {
	loadDotEnv()
	names := playerNames()
	for _, n := range names {
		if err := player.ValidateName(n); err != nil {
			return Config{}, fmt.Errorf("GRIMOIRE_PLAYERS name %q: %w", n, err)
		}
	}
	return Config{
		Token:  strings.TrimSpace(os.Getenv("DISCORD_TOKEN")),
		DBPath: dbPath(),
		Names:  names,
	}, nil
}

func dbPath() string {
	if p := strings.TrimSpace(os.Getenv("GRIMOIRE_DB_PATH")); p != "" {
		return p
	}
	return "./grimoire.db"
}

func defaultPlayerNames() []string {
	return []string{
		"Gustavo", "Mariana", "Pedro", "Joao", "Janis", "Catti", "Maria", "Eric", "Andre",
	}
}

func playerNames() []string {
	raw := os.Getenv("GRIMOIRE_PLAYERS")
	if strings.TrimSpace(raw) == "" {
		return defaultPlayerNames()
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return defaultPlayerNames()
	}
	return out
}

func loadDotEnv() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	dir := wd
	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			_ = godotenv.Load(candidate)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}
