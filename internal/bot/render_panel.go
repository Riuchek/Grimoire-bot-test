package bot

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"grimoire/internal/domain/player"
)

const (
	ansiReset  = "\x1b[0m"
	ansiBrand  = "\x1b[1;36m"
	ansiHeader = "\x1b[1;33m"
	ansiLabel  = "\x1b[0;37m"
	ansiFocus  = "\x1b[1;32m"
	ansiNote   = "\x1b[0;36m"
)

const (
	maxNameRunes       = 12
	customColMaxRunes  = 24
	colN20             = 3
	colN1              = 3
	colDtot            = 5
	colDmax            = 4
	colCtot            = 5
	colCmax            = 4
	colQ               = 2
	colM               = 2
)

func (b *GrimoireBot) RenderTable(focus string) string {
	nameW := b.maxNameWidth()
	customW := b.customColumnRunes()
	inner := tableInnerWidth(nameW, customW)
	rule := strings.Repeat("─", inner)

	var sb strings.Builder
	sb.WriteString("```ansi\n")
	sb.WriteString(ansiBrand)
	sb.WriteString("GRIMOIRE")
	sb.WriteString(ansiReset)
	sb.WriteString("\n")
	sb.WriteString(ansiDimLine(rule))
	sb.WriteString("\n")

	headerPlain := formatHeaderPlain(nameW, customW)
	sb.WriteString(ansiHeader)
	sb.WriteString(headerPlain)
	sb.WriteString(ansiReset)
	sb.WriteString("\n")

	for _, name := range b.Players {
		p := b.PlayersStats[name]
		leftPlain := formatDataRowPlain(p, nameW, name == focus)
		sb.WriteString(styleDataRow(leftPlain, name == focus))
		if customW > 0 {
			c := strings.TrimSpace(p.Custom())
			cell := padRunesRight(truncateRunes(c, customW), customW)
			sb.WriteString(" ")
			if c == "" {
				sb.WriteString(ansiLabel + cell + ansiReset)
			} else {
				sb.WriteString(ansiNote + cell + ansiReset)
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(ansiDimLine(rule))
	sb.WriteString("\n")
	if focus == "" {
		sb.WriteString(ansiLabel)
		sb.WriteString(fillLine("Selecione um jogador.", inner))
		sb.WriteString(ansiReset)
	}
	sb.WriteString("\n```")
	return sb.String()
}

func (b *GrimoireBot) customColumnRunes() int {
	for _, name := range b.Players {
		if strings.TrimSpace(b.PlayersStats[name].Custom()) != "" {
			return customColMaxRunes
		}
	}
	return 0
}

func tableInnerWidth(nameW int, customColW int) int {
	base := tableBaseWidth(nameW)
	if customColW > 0 {
		return base + 1 + customColW
	}
	return base
}

func tableBaseWidth(nameW int) int {
	marker := 2
	afterName := 1
	cells := colN20 + colN1 + colDtot + colDmax + colCtot + colCmax + colQ + colM
	gapsBetween := 7
	return marker + nameW + afterName + cells + gapsBetween
}

func formatHeaderPlain(nameW int, customColW int) string {
	label := padRunesRight("Jogador", nameW)
	cells := []string{
		alignRightCol("N20", colN20),
		alignRightCol("N1", colN1),
		alignRightCol("D·Σ", colDtot),
		alignRightCol("D·↑", colDmax),
		alignRightCol("C·Σ", colCtot),
		alignRightCol("C·↑", colCmax),
		alignRightCol("Q", colQ),
		alignRightCol("M", colM),
	}
	h := "  " + label + " " + strings.Join(cells, " ")
	if customColW > 0 {
		h += " " + padRunesRight("nota", customColW)
	}
	return h
}

func formatDataRowPlain(p *player.Player, nameW int, isFocus bool) string {
	marker := "  "
	if isFocus {
		marker = "> "
	}
	name := padRunesRight(truncateRunes(p.Name(), nameW), nameW)
	nums := fmt.Sprintf("%3d %3d %5d %4d %5d %4d %2d %2d",
		p.SucessoCritico(), p.FalhaCritica(),
		p.DanoTotal(), p.DanoMax(),
		p.CuraTotal(), p.CuraMax(),
		p.Quedas(), p.Mortes(),
	)
	return marker + name + " " + nums
}

func styleDataRow(plain string, isFocus bool) string {
	if isFocus {
		return ansiFocus + plain + ansiReset
	}
	return ansiLabel + plain + ansiReset
}

func fillLine(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if n > width {
		return truncateRunes(s, width)
	}
	return s + strings.Repeat(" ", width-n)
}

func alignRightCol(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if n > width {
		return truncateRunes(s, width)
	}
	return strings.Repeat(" ", width-n) + s
}

func (b *GrimoireBot) maxNameWidth() int {
	w := 7
	for _, name := range b.Players {
		n := utf8.RuneCountInString(name)
		if n > w {
			w = n
		}
	}
	if w > maxNameRunes {
		return maxNameRunes
	}
	if w < 7 {
		return 7
	}
	return w
}

func padRunesRight(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if n >= width {
		return truncateRunes(s, width)
	}
	return s + strings.Repeat(" ", width-n)
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	r := []rune(s)
	if max <= 2 {
		return string(r[:1]) + "."
	}
	return string(r[:max-1]) + "…"
}

func ansiDimLine(line string) string {
	return "\x1b[0;34m" + line + ansiReset
}
