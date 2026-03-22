package storage

import (
	"testing"
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
	if q.Nat20() != 2 || q.DanoTotal() != 10 || q.Custom() != "note" {
		t.Fatalf("unexpected state: n20=%d dano=%d custom=%q", q.Nat20(), q.DanoTotal(), q.Custom())
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
	if p.Nat20() != 0 || p.Name() != "Nobody" {
		t.Fatalf("expected fresh player, got n20=%d name=%q", p.Nat20(), p.Name())
	}
}
