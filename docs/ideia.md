# Ideia e propósito

O **Grimoire** é um bot para **Discord** pensado para mesas de RPG (por exemplo, D&D e sistemas parecidos). Ele funciona como um **painel de estatísticas da campanha**: natural 20, natural 1, dano e cura (totais e máximos por sessão ou por aventura), quedas, mortes e uma linha de **texto livre** por jogador para anotações rápidas.

## O que o jogador ou mestre faz no Discord

1. Usa o comando slash **`/grimoire`** para abrir uma mensagem com a tabela e os controles.
2. **Seleciona um jogador** no menu antes de usar botões ou modais (o foco fica associado àquela mensagem do painel).
3. Registra eventos com botões (N20, N1, Queda, Morte), abre modais para **editar dano/cura**, **custom**, **editar tudo** (counters, dano/cura, quedas/mortes e anotação) ou **limpa** os dados do jogador selecionado.

Os dados são **persistidos em SQLite**, então o painel pode ser reaberto depois de reiniciar o bot sem perder o histórico salvo.

## Para quem é útil

- Grupos que querem um **quadro único** no canal da mesa, visível para todos, sem planilhas externas.
- Campanhas longas em que vale **acompanhar números** e marcos (críticos, quedas, mortes) de forma simples.

## Configuração rápida (conceito)

- Token do bot no ambiente (`DISCORD_TOKEN`).
- Lista de nomes dos jogadores pode ser definida em variável de ambiente ou usar a lista padrão do projeto.
- Caminho do banco SQLite configurável; padrão é um arquivo local na pasta de execução.

Detalhes técnicos de pacotes e fluxo estão em [Arquitetura](arquitetura.md).
