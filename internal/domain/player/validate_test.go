package player

import (
	"strings"
	"testing"
)

func TestSanitizeCustomStripsControls(t *testing.T) {
	got := SanitizeCustom("  hello\x00\x01world  ")
	if got != "helloworld" {
		t.Fatalf("got %q", got)
	}
}

func TestSanitizeCustomTruncatesLongText(t *testing.T) {
	long := strings.Repeat("a", MaxCustomRunes+50)
	got := SanitizeCustom(long)
	if len([]rune(got)) > MaxCustomRunes {
		t.Fatalf("len %d > max %d", len([]rune(got)), MaxCustomRunes)
	}
}

func TestValidateName(t *testing.T) {
	if err := ValidateName("Ada"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateName(""); err == nil {
		t.Fatal("expected error for empty")
	}
	if err := ValidateName(strings.Repeat("x", MaxNameRunes+1)); err == nil {
		t.Fatal("expected error for long name")
	}
}

func TestParseModalIntFields(t *testing.T) {
	got, err := ParseModalIntFields("1 2 3 4", 4)
	if err != nil || len(got) != 4 || got[0] != 1 || got[3] != 4 {
		t.Fatalf("got %v err=%v", got, err)
	}
	if _, err := ParseModalIntFields("1 2", 4); err == nil {
		t.Fatal("expected error for wrong field count")
	}
}

func TestParseModalStatInt(t *testing.T) {
	v, err := ParseModalStatInt(" 42 ")
	if err != nil || v != 42 {
		t.Fatalf("got %d err=%v", v, err)
	}
	v2, err := ParseModalStatInt("")
	if err != nil || v2 != 0 {
		t.Fatalf("empty: got %d err=%v", v2, err)
	}
	if _, err := ParseModalStatInt("12x"); err == nil {
		t.Fatal("expected error for non-numeric")
	}
	if _, err := ParseModalStatInt("-1"); err == nil {
		t.Fatal("expected error for negative")
	}
}

func TestValidateForPersistence(t *testing.T) {
	p := New("Bob")
	p.LoadStats(1, 1, 1, 1, 1, 1, 1, 1, "ok")
	if err := ValidateForPersistence(p); err != nil {
		t.Fatal(err)
	}
	bad := New("")
	if err := ValidateForPersistence(bad); err == nil {
		t.Fatal("expected error for empty name")
	}
}
