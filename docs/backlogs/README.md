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

- [010 - installation-identity-mvp.md](<api/010 - installation-identity-mvp.md>):
  backlog guarda-chuva do MVP de identidade de instalação na API.
- [011 - installation-identity-auth-login.md](<api/011 - installation-identity-auth-login.md>):
  tratamento da instalação no `POST /auth/login`.
- [012 - installation-identity-step-up-authorize.md](<api/012 - installation-identity-step-up-authorize.md>):
  `step-up` para autorizar `POST /security/installations`.
- [013 - installation-identity-register-installation.md](<api/013 - installation-identity-register-installation.md>):
  registro explícito de nova instalação.
- [014 - installation-identity-list-installations.md](<api/014 - installation-identity-list-installations.md>):
  listagem de instalações cadastradas.
- [015 - installation-identity-revoke-installation.md](<api/015 - installation-identity-revoke-installation.md>):
  revogação de instalação e corte imediato de acesso.
- [016 - installation-identity-auth-refresh.md](<api/016 - installation-identity-auth-refresh.md>):
  vínculo de instalação no `POST /auth/refresh`.
- [017 - installation-identity-session-enforcement.md](<api/017 - installation-identity-session-enforcement.md>):
  enforcement do `X-Installation-Id` nas rotas autenticadas.
- [018 - installation-identity-shared-infrastructure.md](<api/018 - installation-identity-shared-infrastructure.md>):
  modelo de dados e infraestrutura compartilhada do MVP.

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
