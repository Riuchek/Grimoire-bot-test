package bot

import (
	"fmt"
	"grimoire/internal/status"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const (
	modalDataPrefix   = "modal_data:"
	modalCustomPrefix = "modal_custom:"
)

type PlayerRepository interface {
	SavePlayer(p *status.Player) error
	LoadPlayers(names []string) (map[string]*status.Player, error)
}

type GrimoireBot struct {
	Players      []string
	PlayersStats map[string]*status.Player
	Repo         PlayerRepository
	activeByMsg  map[string]string
	Mu           sync.Mutex
}

func NewGrimoireBot(names []string, players map[string]*status.Player, repo PlayerRepository) *GrimoireBot {
	return &GrimoireBot{
		Players:      names,
		PlayersStats: players,
		Repo:         repo,
		activeByMsg:  make(map[string]string),
	}
}

func interactionMessageID(ic *discordgo.Interaction) string {
	if ic.Message != nil {
		return ic.Message.ID
	}
	return ""
}

func parseModalCustomID(id string) (msgID string, statsModal bool, ok bool) {
	if rest, ok := strings.CutPrefix(id, modalDataPrefix); ok {
		return rest, true, true
	}
	if rest, ok := strings.CutPrefix(id, modalCustomPrefix); ok {
		return rest, false, true
	}
	return "", false, false
}

func (b *GrimoireBot) persistSave(p *status.Player) {
	_ = b.Repo.SavePlayer(p)
}

func (b *GrimoireBot) RespondSlashGrimoire(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "grimoire" {
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    b.RenderTable(""),
			Components: b.CreateComponents(),
		},
	})
}

func (b *GrimoireBot) RenderTable(focus string) string {
	var sb strings.Builder
	sb.WriteString("```ansi\n")
	sb.WriteString("\x1b[1;34m==================================================\x1b[0m\n")
	sb.WriteString("           📖 \x1b[1;37mGRIMOIRE: AUTOS DA AVENTURA\x1b[0m\n")
	sb.WriteString("\x1b[1;34m==================================================\x1b[0m\n")
	sb.WriteString("\x1b[1;33mJOGADOR  | N20 | N1 | D.TOTAL | D.MAX | C.TOTAL | C.MAX | Q | M \x1b[0m\n")
	sb.WriteString("--------------------------------------------------\n")

	for _, name := range b.Players {
		p := b.PlayersStats[name]
		row := fmt.Sprintf("%-8s | %-3d | %-2d | %-7d | %-5d | %-7d | %-5d | %-1d | %-1d\n",
			p.Name(), p.Nat20(), p.Nat1(), p.DanoTotal(), p.DanoMax(), p.CuraTotal(), p.CuraMax(), p.Quedas(), p.Mortes())
		if p.Custom() != "" {
			row += fmt.Sprintf(" └─ \x1b[0;32m%s\x1b[0m\n", p.Custom())
		}
		sb.WriteString(row)
	}
	sb.WriteString("--------------------------------------------------\n")
	if focus != "" {
		sb.WriteString(fmt.Sprintf("\x1b[1;32mFoco Atual: %s\x1b[0m\n", focus))
	}
	sb.WriteString("```")
	return sb.String()
}

func (b *GrimoireBot) CreateComponents() []discordgo.MessageComponent {
	var options []discordgo.SelectMenuOption
	for _, name := range b.Players {
		options = append(options, discordgo.SelectMenuOption{Label: name, Value: name})
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{CustomID: "select_player", Placeholder: "👤 Selecionar Jogador", Options: options},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "N20", CustomID: "add_n20", Style: discordgo.SuccessButton},
			discordgo.Button{Label: "N1", CustomID: "add_n1", Style: discordgo.DangerButton},
			discordgo.Button{Label: "Queda", CustomID: "add_q", Style: discordgo.SecondaryButton},
			discordgo.Button{Label: "Morte", CustomID: "add_m", Style: discordgo.SecondaryButton},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "📝 Registrar Dano/Cura", CustomID: "open_modal", Style: discordgo.PrimaryButton},
			discordgo.Button{Label: "⚙️ Custom", CustomID: "open_custom", Style: discordgo.SecondaryButton},
		}},
	}
}

func (b *GrimoireBot) HandleComponents(s *discordgo.Session, i *discordgo.InteractionCreate) {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	msgID := interactionMessageID(i.Interaction)
	if msgID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Use botoes e menus nesta mensagem do painel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	id := i.MessageComponentData().CustomID
	focus := b.activeByMsg[msgID]

	if id == "select_player" {
		b.activeByMsg[msgID] = i.MessageComponentData().Values[0]
		focus = b.activeByMsg[msgID]
	} else if focus == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Selecione um jogador primeiro!", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	p := b.PlayersStats[focus]
	needsSave := false
	switch id {
	case "add_n20":
		p.AddNat20()
		needsSave = true
	case "add_n1":
		p.AddNat1()
		needsSave = true
	case "add_q":
		p.AddQueda()
		needsSave = true
	case "add_m":
		p.AddMorte()
		needsSave = true
	case "open_modal":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: modalDataPrefix + msgID, Title: "Registrar para " + focus,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_dano_total", Label: "Valor de Dano Total", Style: discordgo.TextInputShort, Placeholder: "0", Value: strconv.Itoa(p.DanoTotal())},
					}},
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_dano_max", Label: "Valor de Dano Maximo", Style: discordgo.TextInputShort, Placeholder: "0", Value: strconv.Itoa(p.DanoMax())},
					}},
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_cura_total", Label: "Valor de Cura Total", Style: discordgo.TextInputShort, Placeholder: "0", Value: strconv.Itoa(p.CuraTotal())},
					}},
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_cura_max", Label: "Valor de Cura Maximo", Style: discordgo.TextInputShort, Placeholder: "0", Value: strconv.Itoa(p.CuraMax())},
					}},
				},
			},
		})
		return
	case "open_custom":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: modalCustomPrefix + msgID, Title: "Anota\u00e7\u00e3o Custom: " + focus,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_custom", Label: "Texto (ex: Sorte: 2)", Style: discordgo.TextInputShort, Value: p.Custom()},
					}},
				},
			},
		})
		return
	}

	if needsSave {
		b.persistSave(p)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{Content: b.RenderTable(focus), Components: b.CreateComponents()},
	})
}

func (b *GrimoireBot) HandleModals(s *discordgo.Session, i *discordgo.InteractionCreate) {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	msgID, statsModal, ok := parseModalCustomID(i.ModalSubmitData().CustomID)
	if !ok {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Formulario invalido.", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	focus := b.activeByMsg[msgID]
	if focus == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Selecione um jogador primeiro!", Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	p := b.PlayersStats[focus]
	d := i.ModalSubmitData()

	if statsModal {
		dano_total, _ := strconv.Atoi(d.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		dano_max, _ := strconv.Atoi(d.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		cura_total, _ := strconv.Atoi(d.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		cura_max, _ := strconv.Atoi(d.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		p.UpdateStats(dano_total, dano_max, cura_total, cura_max)
	} else {
		p.SetCustom(d.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
	}

	b.persistSave(p)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{Content: b.RenderTable(focus), Components: b.CreateComponents()},
	})
}
