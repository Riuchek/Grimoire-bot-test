package player

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	MaxNameRunes    = 64
	MaxCustomRunes  = 512
	maxStatValue    = 999_999_999
	maxCounterValue = 999_999
)

var (
	ErrInvalidName   = errors.New("invalid player name")
	ErrInvalidStats  = errors.New("invalid stat values")
	ErrInvalidCustom = errors.New("invalid custom text")
)

func SanitizeCustom(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToValidUTF8(s, "")
	if s == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(utf8.RuneCountInString(s))
	for _, r := range s {
		if r == 0 {
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if utf8.RuneCountInString(out) <= MaxCustomRunes {
		return out
	}
	r := []rune(out)
	if MaxCustomRunes <= 0 {
		return ""
	}
	if MaxCustomRunes == 1 {
		return string(r[0])
	}
	return string(r[:MaxCustomRunes-1]) + "…"
}

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidName
	}
	if strings.Contains(name, "\x00") {
		return ErrInvalidName
	}
	if utf8.RuneCountInString(name) > MaxNameRunes {
		return ErrInvalidName
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return ErrInvalidName
		}
	}
	if !utf8.ValidString(name) {
		return ErrInvalidName
	}
	return nil
}

func ParseModalIntFields(s string, want int) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(s))
	if len(fields) != want {
		return nil, ErrInvalidStats
	}
	out := make([]int, want)
	for i := range fields {
		v, err := ParseModalStatInt(fields[i])
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func ParseModalStatInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, ErrInvalidStats
	}
	if v < 0 {
		return 0, ErrInvalidStats
	}
	if v > int64(maxStatValue) {
		return 0, ErrInvalidStats
	}
	return int(v), nil
}

func clampStat(v int) int {
	if v < 0 {
		return 0
	}
	if v > maxStatValue {
		return maxStatValue
	}
	return v
}

func clampCounter(v int) int {
	if v < 0 {
		return 0
	}
	if v > maxCounterValue {
		return maxCounterValue
	}
	return v
}

func ValidateForPersistence(p *Player) error {
	if p == nil {
		return ErrInvalidName
	}
	if err := ValidateName(p.Name()); err != nil {
		return err
	}
	if p.SucessoCritico() < 0 || p.SucessoCritico() > maxCounterValue {
		return ErrInvalidStats
	}
	if p.FalhaCritica() < 0 || p.FalhaCritica() > maxCounterValue {
		return ErrInvalidStats
	}
	if p.DanoTotal() < 0 || p.DanoTotal() > maxStatValue {
		return ErrInvalidStats
	}
	if p.DanoMax() < 0 || p.DanoMax() > maxStatValue {
		return ErrInvalidStats
	}
	if p.CuraTotal() < 0 || p.CuraTotal() > maxStatValue {
		return ErrInvalidStats
	}
	if p.CuraMax() < 0 || p.CuraMax() > maxStatValue {
		return ErrInvalidStats
	}
	if p.Quedas() < 0 || p.Quedas() > maxCounterValue {
		return ErrInvalidStats
	}
	if p.Mortes() < 0 || p.Mortes() > maxCounterValue {
		return ErrInvalidStats
	}
	c := p.Custom()
	if len(c) > 0 && utf8.RuneCountInString(c) > MaxCustomRunes {
		return ErrInvalidCustom
	}
	if !utf8.ValidString(c) {
		return ErrInvalidCustom
	}
	for _, r := range c {
		if unicode.IsControl(r) {
			return ErrInvalidCustom
		}
	}
	return nil
}
