package bot

import (
	"strings"
	"testing"
	"unicode/utf8"

	"grimoire/internal/domain/player"
)

type fakeRepo struct{}

func (fakeRepo) SavePlayer(p *player.Player) error { return nil }

func (fakeRepo) LoadPlayers(names []string) (map[string]*player.Player, error) {
	m := make(map[string]*player.Player)
	for _, n := range names {
		m[n] = player.New(n)
	}
	return m, nil
}

func TestRenderTableFocusLine(t *testing.T) {
	names := []string{"Ada"}
	stats := map[string]*player.Player{"Ada": player.New("Ada")}
	b := NewGrimoireBot(names, stats, fakeRepo{})

	outEmpty := b.RenderTable("")
	if !strings.Contains(outEmpty, "Selecione um jogador.") {
		t.Fatal("expected hint when no focus")
	}
	if strings.Contains(outEmpty, "> Ada") {
		t.Fatal("expected no focused row when focus empty")
	}

	outFocus := b.RenderTable("Ada")
	if !strings.Contains(outFocus, "> Ada") {
		t.Fatalf("expected focused row in output: %q", outFocus)
	}
	if strings.Contains(outFocus, "Selecione") {
		t.Fatal("expected no hint when focus is set")
	}
}

func TestRenderTableColumnAlignment(t *testing.T) {
	names := []string{"Ada", "BobinhoGrande"}
	stats := make(map[string]*player.Player)
	for _, n := range names {
		stats[n] = player.New(n)
	}
	b := NewGrimoireBot(names, stats, fakeRepo{})
	nameW := b.maxNameWidth()
	inner := tableInnerWidth(nameW, 0)
	h := formatHeaderPlain(nameW, 0)
	p := stats["Ada"]
	row := formatDataRowPlain(p, nameW, false)
	rowFocus := formatDataRowPlain(p, nameW, true)
	if utf8.RuneCountInString(h) != utf8.RuneCountInString(row) {
		t.Fatalf("header width %d != row width %d\n%q\n%q", utf8.RuneCountInString(h), utf8.RuneCountInString(row), h, row)
	}
	if utf8.RuneCountInString(row) != utf8.RuneCountInString(rowFocus) {
		t.Fatalf("focus toggles width: %d vs %d", utf8.RuneCountInString(row), utf8.RuneCountInString(rowFocus))
	}
	if utf8.RuneCountInString(h) != inner {
		t.Fatalf("inner %d != header width %d", inner, utf8.RuneCountInString(h))
	}
}

func TestRenderTableColumnAlignmentWithNotaColumn(t *testing.T) {
	stats := map[string]*player.Player{"Ada": player.New("Ada")}
	stats["Ada"].SetCustom("obs")
	b := NewGrimoireBot([]string{"Ada"}, stats, fakeRepo{})
	nameW := b.maxNameWidth()
	cw := b.customColumnRunes()
	if cw != customColMaxRunes {
		t.Fatalf("custom column width: got %d want %d", cw, customColMaxRunes)
	}
	inner := tableInnerWidth(nameW, cw)
	h := formatHeaderPlain(nameW, cw)
	p := stats["Ada"]
	row := formatDataRowPlain(p, nameW, false) + " " + padRunesRight(truncateRunes(strings.TrimSpace(p.Custom()), cw), cw)
	if utf8.RuneCountInString(h) != utf8.RuneCountInString(row) {
		t.Fatalf("header %d != row %d\n%q\n%q", utf8.RuneCountInString(h), utf8.RuneCountInString(row), h, row)
	}
	if utf8.RuneCountInString(h) != inner {
		t.Fatalf("inner %d != width %d", inner, utf8.RuneCountInString(h))
	}
}

func TestRenderTableWithinDiscordMessageLimit(t *testing.T) {
	names := []string{"Gustavo", "Mariana", "Pedro", "Joao", "Janis", "Catti", "Maria", "Eric", "Andre"}
	stats := make(map[string]*player.Player)
	for _, n := range names {
		p := player.New(n)
		p.AddNat20()
		stats[n] = p
	}
	stats["Maria"].SetCustom("ex.: anotação livre por jogador")
	b := NewGrimoireBot(names, stats, fakeRepo{})
	out := b.RenderTable("Maria")
	if len(out) > 2000 {
		t.Fatalf("message len %d exceeds Discord 2000 limit", len(out))
	}
}

func TestUndoStackRestoresPlayer(t *testing.T) {
	names := []string{"Ada"}
	stats := map[string]*player.Player{"Ada": player.New("Ada")}
	b := NewGrimoireBot(names, stats, fakeRepo{})
	msgID := "m1"
	p := stats["Ada"]
	b.recordUndo(msgID, "Ada", p.Snapshot())
	p.AddNat20()
	if p.SucessoCritico() != 1 {
		t.Fatalf("expected increment, got n20=%d", p.SucessoCritico())
	}
	if !b.popUndo(msgID) {
		t.Fatal("expected undo entry")
	}
	if p.SucessoCritico() != 0 {
		t.Fatalf("after undo n20=%d want 0", p.SucessoCritico())
	}
	if b.popUndo(msgID) {
		t.Fatal("expected empty stack")
	}
}

func TestParseModalID(t *testing.T) {
	id, k, ok := parseModalID("modal_data:123456789")
	if !ok || k != modalKindDanoCura || id != "123456789" {
		t.Fatalf("modal_data: got id=%q kind=%v ok=%v", id, k, ok)
	}
	id2, k2, ok2 := parseModalID("modal_custom:999")
	if !ok2 || k2 != modalKindCustom || id2 != "999" {
		t.Fatalf("modal_custom: got id=%q kind=%v ok=%v", id2, k2, ok2)
	}
	id3, k3, ok3 := parseModalID("modal_edit_full:888")
	if !ok3 || k3 != modalKindEditFull || id3 != "888" {
		t.Fatalf("modal_edit_full: got id=%q kind=%v ok=%v", id3, k3, ok3)
	}
	_, _, ok4 := parseModalID("other")
	if ok4 {
		t.Fatal("expected false for unknown prefix")
	}
}

func TestValidPanelMessageID(t *testing.T) {
	if !validPanelMessageID("12345678901234567") {
		t.Fatal("expected valid snowflake length")
	}
	if validPanelMessageID("123") || validPanelMessageID("abc123456789012345") {
		t.Fatal("expected invalid id")
	}
}
