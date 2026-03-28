# Grimoire

Bot de Discord para acompanhar estatísticas de mesa de RPG: críticos, dano, cura, quedas, mortes e anotações por jogador, com persistência em SQLite e painel atualizado na própria mensagem.

## Funcionalidades

- **Comando `/grimoire`** — Abre o painel: tabela estilizada (ANSI) e componentes de interação.
- **Seleção de jogador** — Menu para escolher quem recebe as ações do painel.
- **Ações rápidas** — Botões para sucesso crítico, falha crítica, queda e morte (com incremento controlado e persistência).
- **Modal de dano/cura** — Edição de totais e máximos com validação numérica estrita.
- **Modal de anotação custom** — Texto livre sanitizado (UTF-8, controle de tamanho e caracteres).
- **Desfazer** — Pilha por mensagem do painel, com reversão em memória e tratamento de falha ao salvar.
- **Limpar dados** — Zera todas as estatísticas e a anotação do **jogador selecionado** (com undo e persistência).
- **Editar jogador** — Modal único (limite do Discord: 5 campos) com N20, N1, bloco de quatro números (dano/cura), quedas/mortes e anotação custom.
- **Persistência** — Upsert por jogador em SQLite; carregamento na subida conforme a lista configurada.

## Arquitetura e padrões (Go)

O repositório segue o **layout canônico** de projetos Go: `cmd/` para o binário, `internal/` para código não importável por outros módulos, e documentação em `docs/`.

| Camada | Pacote | Papel |
|--------|--------|--------|
| Entrada | `cmd/grimoire` | Configuração, sinal de encerramento, delegação ao `app`. |
| Composição | `internal/app` | Liga Discord, SQLite, logging do repositório e handlers. |
| Adaptador Discord | `internal/bot` | Slash, botões, modals, renderização; depende só da **porta** `player.Repository`. |
| Domínio | `internal/domain/player` | Agregado `Player`, regras e validação, interface `Repository`. |
| Infraestrutura | `internal/storage` | Implementação SQLite do `Repository`. |
| Configuração | `internal/config` | Variáveis de ambiente e validação de nomes de jogadores. |

**DDD (prático):** o núcleo é o agregado **Player** e a port **Repository**; o bot e o SQLite não se referenciam — o acoplamento é invertido em `app`. Isso mantém o domínio livre de `discordgo` e facilita trocar persistência ou testar o domínio isoladamente.

**Ports and adapters:** `LoggingPlayerRepository` decora o repositório concreto sem alterar o contrato, útil para observabilidade sem vazar detalhes de infraestrutura para o domínio.

**Segurança e dados:** SQL parametrizado; validação e sanitização antes e na persistência; limites numéricos e de texto no domínio. Detalhes de fluxo estão em [docs/arquitetura.md](docs/arquitetura.md).

## Testes e TDD

Hoje o projeto combina testes de **domínio** (`internal/domain/player`), **bot** (render, undo, IDs de modal) e **storage** (round-trip SQLite em memória). Isso cobre regras centrais e persistência, mas ainda não é TDD “puro” em todo o fluxo.

**Direções recomendadas para TDD mais forte:**

- **Casos de uso explícitos** — Extrair operações “aplicar ação ao jogador e persistir” em funções pequenas testáveis sem sessão Discord, e depois integrar o handler.
- **Contratos do `Repository`** — Testes de integração com `:memory:` que já existem podem ganhar cenários de erro (validação, limites).
- **Handlers** — Testes com `discordgo` são pesados; priorizar extrair parsing/validação de modal para o domínio ou pacote `internal/bot` puro e testar com tabelas.

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
| `GRIMOIRE_PLAYERS` | Nomes separados por vírgula (opcional; há lista padrão no código). Nomes inválidos impedem a inicialização. |

O carregamento de `.env` sobe pelos diretórios pais a partir do diretório de trabalho atual.

## Executar

```bash
go run ./cmd/grimoire
```

## Testes

```bash
go test ./...
```

## Como evoluir o código (escala e qualidade)

Sugestões alinhadas a **Go**, **DDD** e **testabilidade**, sem obrigar uma reescrita grande:

1. **Casos de uso (application layer)** — Introduzir um pacote fino `internal/application` (ou `usecase`) com métodos do tipo `RecordCriticalSuccess`, `ApplyModalStats`, recebendo `Repository` e retornando erros de domínio. Handlers do Discord ficam só com IO; o domínio deixa de ser “empurrado” pelo handler linha a linha.

2. **Value objects** — Onde fizer sentido, tipos como `PlayerName` ou `StatValue` encapsulam validação uma vez e evitam `int` soltos em toda parte.

3. **Migrações de schema** — Hoje a tabela é criada inline; para evolução segura, ferramentas ou migrações versionadas (arquivos SQL) evitam divergência entre ambientes.

4. **Observabilidade** — Métricas (contadores de save, latência) e correlação por `interaction` ID; manter logs sem dados sensíveis (já alinhado ao uso de `slog`).

5. **Concorrência** — O mutex no bot cobre o painel; se o escopo crescer (vários canais, sharding), considerar filas ou limites por servidor para não serializar tudo num único lock global.

6. **CI** — `go test ./...`, `go vet`, e opcionalmente `staticcheck` ou `golangci-lint` em pipeline para travar regressões cedo (apoio direto ao TDD).

Para o diagrama de dependências e fluxo de interações, veja [docs/arquitetura.md](docs/arquitetura.md).
