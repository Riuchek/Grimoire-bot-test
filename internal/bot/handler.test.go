package bot

import (
	"strings"
	"testing"

	"grimoire/internal/status"
)

type fakeRepo struct{}

func (fakeRepo) SavePlayer(p *status.Player) error { return nil }

func (fakeRepo) LoadPlayers(names []string) (map[string]*status.Player, error) {
	m := make(map[string]*status.Player)
	for _, n := range names {
		m[n] = status.NewPlayer(n)
	}
	return m, nil
}

func TestRenderTableFocusLine(t *testing.T) {
	names := []string{"Ada"}
	stats := map[string]*status.Player{"Ada": status.NewPlayer("Ada")}
	b := NewGrimoireBot(names, stats, fakeRepo{})

	outEmpty := b.RenderTable("")
	if strings.Contains(outEmpty, "Foco Atual") {
		t.Fatal("expected no focus line when focus empty")
	}

	outFocus := b.RenderTable("Ada")
	if !strings.Contains(outFocus, "Foco Atual: Ada") {
		t.Fatalf("expected focus line in output: %q", outFocus)
	}
}

func TestParseModalCustomID(t *testing.T) {
	id, stats, ok := parseModalCustomID("modal_data:123456789")
	if !ok || !stats || id != "123456789" {
		t.Fatalf("modal_data: got id=%q stats=%v ok=%v", id, stats, ok)
	}
	id2, stats2, ok2 := parseModalCustomID("modal_custom:999")
	if !ok2 || stats2 || id2 != "999" {
		t.Fatalf("modal_custom: got id=%q stats=%v ok=%v", id2, stats2, ok2)
	}
	_, _, ok3 := parseModalCustomID("other")
	if ok3 {
		t.Fatal("expected false for unknown prefix")
	}
}
