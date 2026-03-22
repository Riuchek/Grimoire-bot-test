package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"grimoire/internal/bot"
	"grimoire/internal/config"
	"grimoire/internal/storage"

	"github.com/bwmarrin/discordgo"
)

func Run(ctx context.Context, cfg config.Config) error {
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return fmt.Errorf("discord session: %w", err)
	}

	sqliteRepo, err := storage.NewSQLiteRepo(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	defer func() {
		if err := sqliteRepo.Close(); err != nil {
			slog.Error("close database", "err", err)
		}
	}()

	repo := &bot.LoggingPlayerRepository{Inner: sqliteRepo}

	loaded, err := repo.LoadPlayers(cfg.Names)
	if err != nil {
		return fmt.Errorf("load players: %w", err)
	}

	g := bot.NewGrimoireBot(cfg.Names, loaded, repo)

	var registerCmd sync.Once
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		logDiscordInteraction(i)
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			g.RespondSlashGrimoire(s, i)
		case discordgo.InteractionMessageComponent:
			g.HandleComponents(s, i)
		case discordgo.InteractionModalSubmit:
			g.HandleModals(s, i)
		}
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		registerCmd.Do(func() {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
				Name:        "grimoire",
				Description: "Abre o painel Grimoire",
			})
			if err != nil {
				slog.Error("register slash command", "err", err)
				return
			}
			slog.Info("grimoire ready")
		})
	})

	if err := dg.Open(); err != nil {
		return fmt.Errorf("discord open: %w", err)
	}
	defer dg.Close()

	select {
	case <-ctx.Done():
		return nil
	}
}
