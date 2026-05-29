# Backlogs e decisões

Esta pasta organiza discussões, decisões e tarefas planejadas do BankLab.

A intenção é separar com clareza o que ainda está em discussão do que já foi resolvido ou implementado.

Os backlogs fazem parte da superfície colaborativa do projeto. Eles não são apenas listas de tarefas: também registram debates de arquitetura, alternativas consideradas, decisões tomadas e caminhos abandonados.

## Estrutura

```text
backlogs/
|-- api/          # discussões e backlogs ativos da API
|-- api/done/     # histórico resolvido ou implementado da API
|-- mobile/       # discussões e backlogs ativos do mobile
|-- mobile/done/  # histórico resolvido ou implementado do mobile
`-- discussion.md # notas amplas de discussão técnica
```

## Convenção

- Arquivos diretamente em `api/` ou `mobile/` representam assuntos ativos ou em planejamento.
- Arquivos dentro de `done/` representam backlogs já resolvidos, implementados ou substituídos por decisões mais novas.
- Backlogs concluídos devem permanecer no repositório como histórico de deliberação.

## Backlogs ativos

### API

- [007 - public-step-up-endpoint-contract.md](<api/007 - public-step-up-endpoint-contract.md>): corrigir o contrato público de step-up para usar método e path HTTP em vez de chave interna de policy.
- [007 - public-step-up-endpoint-contract_tasks.md](<api/007 - public-step-up-endpoint-contract_tasks.md>): tasks do contrato público de step-up por operação HTTP.

### Mobile

- [011 - senha-transacional-e-step-up.md](<mobile/011 - senha-transacional-e-step-up.md>): senha transacional e step-up no fluxo mobile de transferência interna.

## Histórico

Os diretórios `done/` ajudam novos colaboradores a entender decisões anteriores sem confundir esses arquivos com trabalho ainda aberto.

Manter esse histórico é importante porque o BankLab trata maturidade arquitetural como um processo incremental. Discussões, refinamentos e mudanças de direção devem continuar acessíveis quando forem relevantes para entender o estado atual do projeto.

Antes de implementar uma mudança relevante, consulte:

- o backlog ativo correspondente;
- o histórico em `done/`, quando existir;
- o [ROADMAP](../ROADMAP.md);
- o [CONTRIBUTING](../../CONTRIBUTING.md).

## O que preservar

Preserve nos backlogs:

- decisões de modelagem;
- trade-offs relevantes;
- motivos para rejeitar uma alternativa;
- mudanças de escopo;
- relação entre decisões antigas e decisões novas;
- contexto necessário para novos colaboradores entenderem uma implementação.

O objetivo não é guardar ruído. É preservar o raciocínio técnico que explica a evolução do projeto.
