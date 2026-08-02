# Backlogs e decisões

Esta pasta organiza discussões, decisões e tarefas planejadas do BankLab.

A intenção é separar com clareza o que ainda está em discussão do que já foi resolvido ou implementado.

Os backlogs fazem parte da superfície colaborativa do projeto. Eles não são apenas listas de tarefas: também registram debates de arquitetura, alternativas consideradas, decisões tomadas e caminhos abandonados.

## Estrutura

```text
backlogs/
|-- api/          # discussões e backlogs ativos da API
|-- api/done/     # histórico resolvido ou implementado da API
|-- api/olds/     # backlogs substituídos por uma organização mais nova
|-- mobile/       # discussões e backlogs ativos do mobile
|-- mobile/done/  # histórico resolvido ou implementado do mobile
`-- discussion.md # notas amplas de discussão técnica
```

## Convenção

- Arquivos diretamente em `api/` ou `mobile/` representam assuntos ativos ou em planejamento.
- Arquivos dentro de `done/` representam backlogs já resolvidos, implementados ou substituídos por decisões mais novas.
- Arquivos dentro de `olds/` representam versões substituídas de backlogs ainda úteis como histórico de discussão.
- Backlogs concluídos devem permanecer no repositório como histórico de deliberação.

## Backlogs ativos

### API

- [010 - installation-identity-mvp.md](<api/010 - installation-identity-mvp.md>):
  backlog principal do MVP de identidade de instalação na API.
- [010 - split-installation-identity-by-dependency.md](<api/010 - split-installation-identity-by-dependency.md>):
  separação dos backlogs API por ordem de dependência técnica.
- [011 - installation-identity-entry-contract.md](<api/011 - installation-identity-entry-contract.md>):
  contrato mínimo de `X-Installation-Id` no login.
- [012 - installation-identity-domain-contracts.md](<api/012 - installation-identity-domain-contracts.md>):
  domínio, estados, erros e portas internas.
- [013 - installation-identity-database-schema.md](<api/013 - installation-identity-database-schema.md>):
  migrations, tabelas, constraints e índices.
- [014 - installation-identity-repositories.md](<api/014 - installation-identity-repositories.md>):
  implementações Postgres e operações atômicas.
- [015 - installation-identity-session-tokens-context.md](<api/015 - installation-identity-session-tokens-context.md>):
  sessão, claims, token restrito e contexto autenticado.
- [016 - installation-identity-login-usecases.md](<api/016 - installation-identity-login-usecases.md>):
  classificação, bootstrap, limite e autorização restrita no login.
- [017 - installation-identity-management-usecases.md](<api/017 - installation-identity-management-usecases.md>):
  registro explícito, listagem, revogação e efeitos em sessão.
- [018 - installation-identity-delivery-enforcement.md](<api/018 - installation-identity-delivery-enforcement.md>):
  handlers, middlewares, enforcement e contrato REST.
- [019 - installation-identity-audit-retention.md](<api/019 - installation-identity-audit-retention.md>):
  auditoria, retenção, minimização e documentação operacional.

### Mobile

- [013 - installation-identity-mvp.md](<mobile/013 - installation-identity-mvp.md>):
  geração do UUID, persistência local, ciclo de vida da instalação e envio de
  `X-Installation-Id`.

Os backlogs API 010 e Mobile 013 compartilham o contrato de identidade da
instalação. Identificação do aparelho físico fica fora deste MVP.

## Histórico

Os diretórios `done/` ajudam novos colaboradores a entender decisões anteriores sem confundir esses arquivos com trabalho ainda aberto.

As entregas mais recentes movidas para o histórico incluem:

- API 009: bootstrap unificado da sessão autenticada;
- Mobile 011: cadastro da senha transacional;
- Mobile 012: step-up da transferência interna.

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
