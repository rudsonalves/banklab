# Guia de documentação orientada a IA

Este documento explica como o BankLab organiza documentação que também serve como contexto para agentes de IA, assistentes de código e colaboradores que usam ferramentas automatizadas durante o desenvolvimento.

O objetivo não é substituir documentação para humanos. A ideia é tornar explícito quais arquivos ajudam a orientar decisões, navegação no projeto, padrões locais e contexto arquitetural.

## Por que este guia existe

O BankLab possui documentação em diferentes níveis:

- visão pública do projeto;
- guias de contribuição;
- documentação técnica da API e do mobile;
- backlogs e histórico de decisões;
- instruções locais para agentes de IA;
- relatórios de implementação.

Parte desses documentos é lida por pessoas. Outra parte também é útil para agentes de IA entenderem limites, padrões e decisões já tomadas.

Registrar essa organização ajuda a evitar respostas inconsistentes, refatorações fora de escopo e perda de contexto histórico.

## Tipos de documentação

### Documentação pública

Arquivos principais:

- [README.md](../README.md)
- [CONTRIBUTING.md](../CONTRIBUTING.md)
- [docs/README.md](README.md)
- [docs/ROADMAP.md](ROADMAP.md)

Função:

- explicar o que é o projeto;
- orientar novos colaboradores;
- apresentar visão, escopo e direção;
- indicar onde encontrar documentação técnica e backlogs.

Esses documentos devem ser escritos para humanos primeiro, mas também devem dar contexto suficiente para agentes de IA entenderem a intenção geral do projeto.

### Backlogs e decisões

Arquivos principais:

- [docs/backlogs/README.md](backlogs/README.md)
- [docs/backlogs/api/000 - pre-onboarding.md](<backlogs/api/000 - pre-onboarding.md>)
- [docs/backlogs/api/001 - onboarding.md](<backlogs/api/001 - onboarding.md>)
- diretórios `done/`

Função:

- registrar decisões em aberto;
- preservar discussões já resolvidas;
- documentar trade-offs;
- explicar por que um caminho foi escolhido;
- separar trabalho ativo de histórico.

Para agentes de IA, os backlogs são fonte importante de intenção. Antes de implementar mudanças ligadas a onboarding, identidade, customer, documentos ou fluxos bancários, esses arquivos devem ser consultados.

### Documentação técnica

Arquivos principais:

- [api/docs](../api/docs)
- [mobile/docs](../mobile/docs)
- [api/README.md](../api/README.md)
- [mobile/README.md](../mobile/README.md)

Função:

- descrever arquitetura atual;
- registrar contratos HTTP;
- explicar domínio, persistência, autenticação e fluxos;
- orientar setup e manutenção.

Esses documentos descrevem o estado técnico do sistema. Quando o código muda comportamento, contrato ou estrutura, a documentação técnica correspondente deve ser atualizada.

### Instruções para agentes

Arquivos encontrados no projeto:

- `.github/instructions/*.md`
- `mobile/AGENT.md`
- `mobile/lib/**/AGENT.md`

Função:

- orientar agentes de IA em áreas específicas do código;
- registrar padrões locais de arquitetura;
- explicar limites de responsabilidade por pasta;
- reduzir mudanças fora de escopo;
- manter consistência entre módulos.

Esses arquivos são mais operacionais e locais. Eles devem ser objetivos, específicos e alinhados com a documentação pública e técnica.

## Como decidir onde documentar

Use esta regra geral:

- **README/ROADMAP**: visão pública, direção e posicionamento.
- **CONTRIBUTING**: processo de contribuição, PRs, issues e cuidados gerais.
- **docs/backlogs**: decisões em discussão, trade-offs e histórico de deliberação.
- **api/docs e mobile/docs**: comportamento técnico já implementado ou contrato aceito.
- **AGENT.md e .github/instructions**: orientação local para agentes e navegação por área.

## Quando atualizar

Atualize a documentação quando uma mudança:

- altera contrato de API;
- muda schema de banco;
- altera regra de domínio;
- muda autenticação, autorização ou onboarding;
- cria novo fluxo mobile;
- muda organização de pastas ou responsabilidades;
- resolve uma discussão de backlog;
- invalida instrução usada por agentes de IA;
- afeta como novos colaboradores devem entender o projeto.

## Cuidados ao usar IA no projeto

Agentes de IA podem ajudar a navegar, resumir, implementar e revisar mudanças. Ainda assim, as decisões arquiteturais importantes devem permanecer explícitas em documentação versionada.

Ao usar IA:

- forneça o backlog ou documento de decisão relevante;
- confirme se a proposta respeita o roadmap e o escopo;
- peça que mudanças em contrato/schema atualizem documentação;
- evite aceitar refatorações amplas sem backlog ou decisão registrada;
- preserve decisões anteriores, mesmo quando uma abordagem for substituída;
- diferencie claramente implementação de discussão.

## Preservação de histórico

O BankLab trata maturidade arquitetural como processo incremental.

Por isso, discussões, ajustes, alternativas rejeitadas e decisões substituídas podem continuar no repositório quando ajudam a entender a evolução do projeto.

Backlogs concluídos devem ser movidos para `done/`, não apagados sem motivo.

## Relação entre IA e decisões humanas

Documentos orientados a IA existem para melhorar contexto, consistência e produtividade. Eles não substituem decisão humana.

Quando houver conflito entre sugestão automatizada e documentação do projeto, a documentação versionada deve prevalecer até que uma nova decisão seja registrada.
