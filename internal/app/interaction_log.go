package app

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func logDiscordInteraction(i *discordgo.InteractionCreate) {
	uid := ""
	switch {
	case i.Member != nil && i.Member.User != nil:
		uid = i.Member.User.ID
	case i.User != nil:
		uid = i.User.ID
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		slog.Info("discord interaction",
			"kind", "application_command",
			"command", i.ApplicationCommandData().Name,
			"user_id", uid,
			"interaction_id", i.ID,
		)
	case discordgo.InteractionMessageComponent:
		slog.Info("discord interaction",
			"kind", "message_component",
			"custom_id", i.MessageComponentData().CustomID,
			"user_id", uid,
			"interaction_id", i.ID,
		)
	case discordgo.InteractionModalSubmit:
		slog.Info("discord interaction",
			"kind", "modal_submit",
			"custom_id", i.ModalSubmitData().CustomID,
			"user_id", uid,
			"interaction_id", i.ID,
		)
	default:
		slog.Info("discord interaction",
			"kind", "other",
			"type", i.Type,
			"user_id", uid,
			"interaction_id", i.ID,
		)
	}
}
