package storage

import (
	"testing"

	"grimoire/internal/domain/player"
)

func TestSQLiteSaveLoadRoundTrip(t *testing.T) {
	repo, err := NewSQLiteRepo(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	names := []string{"Ada", "Bob"}
	loaded, err := repo.LoadPlayers(names)
	if err != nil {
		t.Fatal(err)
	}
	p := loaded["Ada"]
	p.AddNat20()
	p.AddNat20()
	p.UpdateStats(10, 5, 3, 2)
	p.SetCustom("note")
	if err := repo.SavePlayer(p); err != nil {
		t.Fatal(err)
	}

	again, err := repo.LoadPlayers(names)
	if err != nil {
		t.Fatal(err)
	}
	q := again["Ada"]
	if q.SucessoCritico() != 2 || q.DanoTotal() != 10 || q.Custom() != "note" {
		t.Fatalf("unexpected state: n20=%d dano=%d custom=%q", q.SucessoCritico(), q.DanoTotal(), q.Custom())
	}
}

func TestSQLiteLoadMissingRowIsFreshPlayer(t *testing.T) {
	repo, err := NewSQLiteRepo(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	loaded, err := repo.LoadPlayers([]string{"Nobody"})
	if err != nil {
		t.Fatal(err)
	}
	p := loaded["Nobody"]
	if p == nil {
		t.Fatal("expected player entry")
	}
	if p.SucessoCritico() != 0 || p.Name() != "Nobody" {
		t.Fatalf("expected fresh player, got n20=%d name=%q", p.SucessoCritico(), p.Name())
	}
}

func TestSQLiteSavePlayerRejectsInvalid(t *testing.T) {
	repo, err := NewSQLiteRepo(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	if err := repo.SavePlayer(nil); err == nil {
		t.Fatal("expected error for nil player")
	}
	if err := repo.SavePlayer(player.New("")); err == nil {
		t.Fatal("expected error for empty name")
	}
}
