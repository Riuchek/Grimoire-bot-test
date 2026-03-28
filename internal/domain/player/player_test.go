package player

import "testing"

func TestPlayerUpdateStats(t *testing.T) {
	p := New("x")
	p.UpdateStats(10, 9, 8, 7)
	if p.DanoTotal() != 10 || p.DanoMax() != 9 || p.CuraTotal() != 8 || p.CuraMax() != 7 {
		t.Fatalf("UpdateStats: got dano=%d/%d cura=%d/%d", p.DanoTotal(), p.DanoMax(), p.CuraTotal(), p.CuraMax())
	}
}

func TestPlayerLoadStats(t *testing.T) {
	p := New("y")
	p.LoadStats(1, 2, 3, 4, 5, 6, 7, 8, "note")
	if p.SucessoCritico() != 1 || p.FalhaCritica() != 2 || p.DanoTotal() != 3 || p.Mortes() != 8 || p.Custom() != "note" {
		t.Fatal("LoadStats did not restore fields")
	}
}

func TestPlayerIncrements(t *testing.T) {
	p := New("z")
	p.AddNat20()
	p.AddNat20()
	p.AddNat1()
	p.AddQueda()
	p.AddMorte()
	if p.SucessoCritico() != 2 || p.FalhaCritica() != 1 || p.Quedas() != 1 || p.Mortes() != 1 {
		t.Fatalf("increments: n20=%d n1=%d q=%d m=%d", p.SucessoCritico(), p.FalhaCritica(), p.Quedas(), p.Mortes())
	}
}

func TestPlayerClearAll(t *testing.T) {
	p := New("Ada")
	p.AddNat20()
	p.UpdateStats(5, 4, 3, 2)
	p.SetCustom("x")
	p.ClearAll()
	if p.SucessoCritico() != 0 || p.DanoTotal() != 0 || p.Custom() != "" {
		t.Fatalf("ClearAll: n20=%d dano=%d custom=%q", p.SucessoCritico(), p.DanoTotal(), p.Custom())
	}
}

func TestPlayerSnapshotRoundTrip(t *testing.T) {
	p := New("z")
	p.AddNat20()
	p.UpdateStats(3, 2, 1, 0)
	p.SetCustom("x")
	before := p.Snapshot()
	p.AddNat1()
	p.AddQueda()
	p.RestoreSnapshot(before)
	if p.SucessoCritico() != 1 || p.FalhaCritica() != 0 || p.DanoTotal() != 3 || p.Custom() != "x" || p.Quedas() != 0 {
		t.Fatalf("after restore: n20=%d n1=%d d=%d custom=%q q=%d", p.SucessoCritico(), p.FalhaCritica(), p.DanoTotal(), p.Custom(), p.Quedas())
	}
}
