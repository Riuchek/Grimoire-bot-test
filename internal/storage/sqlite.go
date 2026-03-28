package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"grimoire/internal/domain/player"

	_ "modernc.org/sqlite"
)

var _ player.Repository = (*SQLiteRepo)(nil)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(path string) (*SQLiteRepo, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("sqlite ping: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS players (
		name TEXT PRIMARY KEY,
		nat20 INTEGER, nat1 INTEGER,
		dano_total INTEGER, dano_max INTEGER,
		cura_total INTEGER, cura_max INTEGER,
		quedas INTEGER, mortes INTEGER,
		custom TEXT
	);`
	if _, err := db.Exec(query); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("sqlite migrate: %w", err)
	}

	return &SQLiteRepo{db: db}, nil
}

func (r *SQLiteRepo) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepo) SavePlayer(p *player.Player) error {
	if err := player.ValidateForPersistence(p); err != nil {
		return err
	}
	query := `
	INSERT INTO players (name, nat20, nat1, dano_total, dano_max, cura_total, cura_max, quedas, mortes, custom)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(name) DO UPDATE SET
		nat20=excluded.nat20, nat1=excluded.nat1,
		dano_total=excluded.dano_total, dano_max=excluded.dano_max,
		cura_total=excluded.cura_total, cura_max=excluded.cura_max,
		quedas=excluded.quedas, mortes=excluded.mortes,
		custom=excluded.custom;`

	_, err := r.db.Exec(query, p.Name(), p.SucessoCritico(), p.FalhaCritica(), p.DanoTotal(), p.DanoMax(), p.CuraTotal(), p.CuraMax(), p.Quedas(), p.Mortes(), p.Custom())
	return err
}

func (r *SQLiteRepo) LoadPlayers(names []string) (map[string]*player.Player, error) {
	res := make(map[string]*player.Player)
	for _, name := range names {
		row := r.db.QueryRow(`SELECT nat20, nat1, dano_total, dano_max, cura_total, cura_max, quedas, mortes, custom FROM players WHERE name = ?`, name)

		p := player.New(name)
		var n20, n1, dt, dm, ct, cm, q, m int
		var c string

		err := row.Scan(&n20, &n1, &dt, &dm, &ct, &cm, &q, &m, &c)
		switch {
		case err == nil:
			p.LoadStats(n20, n1, dt, dm, ct, cm, q, m, c)
		case errors.Is(err, sql.ErrNoRows):
		default:
			return nil, fmt.Errorf("load player %q: %w", name, err)
		}
		res[name] = p
	}
	return res, nil
}
