# Grimoire

Bot de Discord para acompanhar estatísticas da mesa (críticos, dano, cura, quedas, mortes e anotações por jogador), com persistência em SQLite.

## Documentação

- [Índice da documentação](docs/README.md)
- [Ideia e propósito do projeto](docs/ideia.md)
- [Arquitetura (pacotes, fluxo, dados)](docs/arquitetura.md)

## Requisitos

- Go (versão indicada em `go.mod`)
- Token de aplicação Discord com permissões adequadas ao bot

## Variáveis de ambiente

| Variável | Descrição |
|----------|------------|
| `DISCORD_TOKEN` | Token do bot (obrigatório) |
| `GRIMOIRE_DB_PATH` | Caminho do arquivo SQLite (padrão: `./grimoire.db`) |
| `GRIMOIRE_PLAYERS` | Nomes separados por vírgula (opcional; há lista padrão no código) |

O carregamento de `.env` sobe pelos diretórios pais a partir do diretório de trabalho atual.

## Executar

```bash
go run ./cmd/grimoire
```
