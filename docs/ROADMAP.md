# Roadmap do BankLab

Este roadmap organiza a direção do BankLab para colaboração técnica. Ele começa explicando o que o projeto é, por que ele existe e quais caminhos fazem sentido para sua evolução.

Ele não é uma promessa fechada de entrega. É um guia para alinhar prioridades, abrir issues e escolher boas contribuições.

## O que é o BankLab

BankLab é um laboratório de engenharia de software aplicado a um domínio bancário simplificado, com ambição de evoluir para uma plataforma experimental de segurança, onboarding, pagamentos e gestão operacional.

O projeto implementa uma API em Go e um app mobile em Flutter para estudar, na prática, problemas que aparecem em sistemas financeiros: autenticação, autorização, onboarding, clientes, contas, saldo, ledger, transferência, extrato, consistência transacional e separação clara de responsabilidades.

Mesmo com escopo reduzido, o BankLab não deve ser entendido como um CRUD de cadastro. Ele está mais próximo de um **core bancário orientado a engenharia**, porque já lida com regras que exigem mais disciplina:

- identidade autenticada propagada para os fluxos;
- autorização por papel e ownership;
- aprovação administrativa antes do acesso completo;
- vínculo entre usuário e cliente de negócio;
- contas associadas ao cliente;
- movimentações financeiras registradas em ledger;
- saldo protegido por transações e regras de consistência;
- documentação e testes como parte da evolução.

O diferencial do projeto não é apenas a stack. A parte mais importante é o raciocínio arquitetural aplicado a um domínio financeiro: modularidade, invariantes, concorrência, idempotência, ledger, autorização e evolução incremental.

O projeto ainda não pretende ser um banco completo. Ele é uma base arquitetural e funcional para estudar um recorte específico do domínio financeiro com profundidade suficiente para gerar decisões reais de engenharia. A partir dessa base, a intenção é evoluir para temas mais ambiciosos, como Zero Trust Architecture, KYC, meios de pagamento locais, app web, painel administrativo e, quando fizer sentido, decomposição em serviços menores.

Parte importante dessa proposta é documentar também o processo de decisão. O BankLab deve mostrar não apenas o software implementado, mas o caminho de engenharia que levou a ele: discussões, alternativas, trade-offs, correções de rumo e decisões preservadas em backlog.

## Por que este projeto existe

O BankLab foi criado para praticar decisões reais de engenharia em um domínio onde não basta criar telas e endpoints. Em sistemas financeiros, regras de negócio, autorização, saldo e histórico de movimentações precisam ser tratados com clareza e consistência desde o início.

A motivação principal é praticar engenharia em um contexto onde algumas decisões importam de verdade:

- onde colocar uma regra de negócio;
- como separar transporte HTTP, caso de uso, domínio e persistência;
- como impedir que o cliente decida coisas que pertencem ao backend;
- como proteger saldo, ledger e transferência;
- como validar ownership de conta;
- como evoluir um projeto sem perder clareza arquitetural;
- como documentar decisões para que outras pessoas consigam contribuir.

Em outras palavras: o BankLab existe para ser um projeto de aprendizado sério, com cara de sistema real, mas sem a pretensão de operar dinheiro real.

Ele também serve como vitrine técnica. O objetivo é mostrar capacidade de construir, explicar e evoluir uma aplicação com backend, mobile, testes, documentação e raciocínio de produto financeiro. Para colaboradores, ele deve funcionar como um ambiente de prática, troca e crescimento.

## Onde o projeto quer chegar

O BankLab quer evoluir de um core bancário simplificado para um laboratório mais amplo de produto financeiro, segurança aplicada e colaboração técnica.

O caminho do projeto tem sete frentes principais:

- consolidar uma base bancária simplificada, segura e consistente;
- criar uma experiência mobile útil para validar fluxos reais de uso;
- desenvolver um onboarding progressivo, persistente e retomável;
- estruturar uma abordagem de Zero Trust Architecture, com sinais coletados no mobile e decisões controladas pelo backend;
- experimentar integrações externas de KYC e serviços auxiliares sempre mediadas pelo backend;
- abrir caminho para app web, painel administrativo e evolução arquitetural incremental;
- transformar a engenharia já construída em uma narrativa pública clara, capaz de atrair colaboradores alinhados com a visão do projeto.

Com o tempo, o projeto deve evoluir de uma base funcional para um laboratório mais completo de segurança transacional, experiência mobile, canais web, gestão administrativa, pagamentos locais e colaboração open source.

## Eixos estratégicos

### Narrativa pública e colaboração

O BankLab precisa ser apresentado como um projeto de engenharia aplicada a sistemas financeiros, não apenas como uma API bancária.

Boa parte do valor técnico já existe no código e na documentação, mas precisa ser comunicada de forma mais acessível para atrair colaboradores e gerar conversa pública.

Esse eixo inclui:

- README forte, com visão, stack, arquitetura, quick start, roadmap e convite para colaboração;
- roadmap público com direção clara;
- issues pequenas, bem delimitadas e com critérios de aceite;
- templates de issue e pull request;
- screenshots do app mobile;
- posts técnicos explicando decisões arquiteturais do projeto;
- discussões públicas sobre temas como ledger, concorrência, idempotência, monólito modular e Zero Trust.

O objetivo não é criar marketing vazio. É transformar engenharia real em narrativa compreensível para pessoas que possam aprender, opinar e contribuir.

### Onboarding progressivo

O onboarding deve evoluir para um fluxo guiado por checkpoints controlados pelo backend.

A intenção é permitir que o usuário avance por etapas, tenha dados persistidos localmente no frontend quando fizer sentido e consiga continuar o processo de onde parou. O backend deve continuar sendo a fonte de verdade sobre estado, etapa atual, validações obrigatórias e liberação de conta.

Esse eixo inclui:

- checkpoints de onboarding definidos pelo backend;
- persistência local temporária no mobile;
- retomada de onboarding interrompido;
- validação gradual de dados;
- integração futura com KYC;
- aprovação ou bloqueio orientado por estado do processo.

### Zero Trust Architecture

Um dos objetivos centrais é estruturar uma abordagem de Zero Trust Architecture para o BankLab.

No contexto do projeto, isso significa que uma requisição autenticada não deve ser automaticamente tratada como confiável. O backend deve avaliar contexto, instalação do app, fator adicional, tipo de operação e sinais de risco antes de permitir ações sensíveis.

Esse eixo inclui:

- senha transacional para operações críticas;
- controle de instalações conhecidas do app;
- prova de vida para validação de identidade ou operação;
- TOTP para autenticação forte, especialmente em possível app web;
- coleta de sinais adicionais para decisão de confiança;
- políticas backend para decidir quando uma operação pode prosseguir, exigir reforço ou ser bloqueada.

### Integrações externas via backend

Serviços de terceiros devem ser integrados pelo backend, não diretamente pelos clientes.

Esse eixo pode incluir:

- KYC;
- validação documental;
- prova de vida;
- consulta ou simulação de serviços financeiros;
- provedores auxiliares de risco, identidade ou comunicação.

A regra central é manter o backend como camada de controle, auditoria, proteção de credenciais e normalização de contratos externos.

### Pagamentos locais

O projeto deve estudar meios de pagamento locais, mesmo que inicialmente mockados.

Esse eixo inclui:

- DOC/TED simulados;
- Pix mockado;
- fluxos de entrada e saída de dinheiro;
- estados de pagamento;
- recibos;
- idempotência;
- trilha de auditoria;
- separação entre movimentação interna e integração externa.

### Canais web e administração

Além do mobile, o BankLab deve abrir caminho para um app web e um painel administrativo.

Esse eixo inclui:

- app web para fluxos de cliente quando fizer sentido;
- autenticação forte para web;
- painel administrativo para gestão de usuários, clientes, contas e onboarding;
- aprovação e revisão operacional;
- visualização de eventos relevantes;
- auditoria de ações administrativas.

### Evolução arquitetural

O projeto deve começar simples e evoluir com critério.

A base atual em monólito modular é adequada para o estágio presente. A evolução para microserviços deve acontecer somente quando houver motivo claro, como isolamento de domínio, necessidade operacional, diferença de escala, fronteira de segurança ou integração externa que justifique separação.

Microserviços não são objetivo por si só. Eles são uma opção futura quando a arquitetura pedir.

## Princípios de evolução

- **Consistência antes de volume**: funcionalidades financeiras precisam preservar ledger, saldo e rastreabilidade.
- **Produto antes de endpoint**: uma rota só deve existir se representar uma capacidade clara do sistema.
- **Pequenos passos revisáveis**: mudanças menores ajudam novos colaboradores a entrar com mais segurança.
- **Documentação como parte da entrega**: decisões importantes devem ficar registradas.
- **Comunicação clara**: a documentação pode existir em português ou inglês, conforme o contexto, mas deve permanecer compreensível e alinhada às decisões do projeto.

## Agora

Prioridades para estabilizar a base atual e melhorar a entrada de colaboradores.

### Vitrine pública do projeto

- Manter README, CONTRIBUTING e roadmap alinhados com o estado real do projeto.
- Melhorar a apresentação visual do projeto com diagrama, banner ou screenshots.
- Preparar material curto mostrando o app mobile.
- Criar templates de issue e pull request.
- Organizar futuras issues de colaboração com escopo pequeno e critérios de aceite.
- Selecionar temas técnicos para uma primeira série de posts públicos.
- Manter `docs/` e `docs/backlogs/` como documentação pública do processo de decisão.

### Documentação e onboarding

- Revisar guias de setup da API e do mobile.
- Atualizar exemplos de payload para endpoints principais.
- Melhorar documentação de autenticação, aprovação de usuário e uso de `X-App-Token`.
- Registrar decisões de onboarding antes da implementação.
- Reescrever discussões importantes quando o texto original não estiver claro para colaboradores.
- Separar backlogs ativos de backlogs concluídos usando diretórios `done/`.

### API

- Consolidar a separação entre rotas de cliente e rotas operacionais.
- Garantir que criação de conta seja tratada como provisionamento administrativo.
- Revisar a posição de endpoints diretos de depósito e saque.
- Fortalecer testes de transferência interna, saldo insuficiente e autorização.
- Atualizar documentação REST sempre que o contrato HTTP mudar.

### Mobile

- Melhorar mensagens de erro no login.
- Tratar corretamente o estado de conta aguardando aprovação.
- Evoluir a jornada de extrato e detalhes de transação.
- Revisar estados de loading, vazio e falha nas telas principais.
- Garantir que o app não dependa de rotas operacionais que não sejam produto para cliente final.

## Próximo

Prioridades para tornar o BankLab mais completo como laboratório de produto financeiro.

### Experiência de conta

- Melhorar a tela de extrato.
- Permitir abrir detalhes de transações a partir do extrato.
- Revisar o recibo de transferência.
- Melhorar o fluxo de seleção de destinatário em transferências internas.
- Refinar textos e estados visuais para cenários de erro financeiro.

### Onboarding

- Desenhar checkpoints controlados pelo backend.
- Definir estados possíveis do onboarding.
- Permitir retomada de fluxo interrompido.
- Separar dados temporários do frontend e estado oficial do backend.
- Preparar pontos de integração futura com KYC.

### Zero Trust inicial

- Modelar senha transacional como credencial separada da senha de login.
- Definir quais operações exigem senha transacional.
- Modelar cadastro de instalação do app.
- Adotar `X-Installation-Id` e definir os sinais mínimos da instalação.
- Avaliar identificação e confiança do aparelho físico somente em uma evolução
  futura, caso exista necessidade clara.
- Criar primeira política backend para avaliar operações sensíveis.
- Documentar o fluxo antes de implementar etapas mais avançadas.

### Pagamentos simulados

- Desenhar fluxo mockado de DOC/TED.
- Desenhar fluxo mockado de Pix.
- Separar pagamentos externos simulados de transferências internas.
- Definir estados, recibos e registros de auditoria.

### Qualidade e testes

- Expandir testes de aplicação e delivery na API.
- Cobrir cenários de concorrência em operações financeiras críticas.
- Criar testes mobile para fluxos de autenticação, conta e transferência.
- Adicionar validações automatizadas para contratos importantes.
- Avaliar pipeline de CI para rodar testes de API e mobile.

### DevEx

- Melhorar comandos do Makefile quando houver fricção no setup.
- Atualizar coleção Bruno ou alternativa equivalente.
- Criar uma lista de tarefas `good first issue`.
- Padronizar templates de issue e pull request.
- Documentar decisões arquiteturais relevantes de forma clara, em português ou inglês conforme o contexto.

### Comunicação técnica

- Preparar posts curtos sobre decisões reais do projeto.
- Explicar por que saldo bancário exige consistência forte.
- Mostrar como o projeto usa ledger para registrar movimentações.
- Explicar a escolha por monólito modular antes de microserviços.
- Discutir por que JWT sozinho não representa uma arquitetura Zero Trust.
- Usar diagramas, exemplos pequenos e screenshots para tornar o projeto mais acessível.

## Futuro

Ideias maiores, que precisam de discussão antes de implementação.

### KYC e serviços externos

- Integração com provedores de KYC.
- Validação documental.
- Prova de vida por provedor externo ou mock controlado.
- Normalização de respostas externas no backend.
- Auditoria das decisões tomadas com base nesses serviços.

### Zero Trust avançado

- Prova de vida para validação de dispositivo ou operação.
- TOTP para autenticação forte em app web.
- Step-up authentication para ações críticas.
- Motor de políticas para decisão contextual.
- Sinais adicionais de risco e confiança.
- Modelo de decisão por contexto da operação.
- Regras mais fortes para autorização e ownership.

### Superfícies administrativas

- Backoffice simples para aprovação de usuários.
- Provisionamento administrativo de contas.
- Visão operacional de clientes e contas.
- Auditoria de ações administrativas.
- Gestão de checkpoints de onboarding.
- Revisão de eventos sensíveis.

### Produto financeiro

- Modelagem futura de cash-in e cash-out.
- DOC/TED mockados.
- Pix mockado.
- Estudo futuro de integração real, somente após a modelagem mockada estar madura.
- Limites transacionais.
- Categorização ou enriquecimento de extrato.
- Simulações de bloqueio, encerramento ou alteração de status de conta.

### App web

- App web para fluxos de cliente.
- Login com autenticação forte.
- Integração com TOTP.
- Retomada de onboarding pelo navegador.
- Visualização de conta, extrato e recibos.

### Evolução para serviços

- Identificar fronteiras reais de domínio antes de decompor.
- Separar serviços somente quando houver necessidade clara.
- Avaliar candidatos como auth, onboarding, ledger, payments e admin.
- Preservar consistência e contratos explícitos durante qualquer decomposição.

### Observabilidade e operação

- Logs estruturados.
- Métricas básicas.
- Rastreamento de erros.
- Health checks.
- Melhor visibilidade para fluxos críticos.

### Plataforma educacional e referência técnica

- Organizar uma série "construindo um core bancário".
- Transformar decisões técnicas em artigos ou posts mais longos.
- Criar material de apoio para novos colaboradores.
- Usar o projeto como portfólio arquitetural.
- Abrir discussões no GitHub para decisões maiores.

## Primeiras contribuições

Esta seção ainda será refinada.

Como o BankLab tem uma direção ambiciosa, as primeiras contribuições precisam ser escolhidas com cuidado. A intenção é abrir tarefas pequenas, mas que apontem para os eixos reais do projeto: onboarding, segurança, experiência mobile, documentação técnica, testes e organização da API.

Antes de publicar uma lista definitiva de `good first issues`, o projeto precisa transformar os eixos deste roadmap em tarefas menores, com escopo claro e critérios de aceite.

## Fora do roadmap imediato

Por enquanto, o projeto não deve priorizar:

- integração real com bancos ou provedores financeiros;
- uso de dados bancários reais;
- antifraude complexo;
- múltiplas moedas;
- conciliação bancária externa;
- migração prematura para arquitetura distribuída;
- grandes refatorações sem problema concreto.

## Como propor mudanças no roadmap

Abra uma issue do tipo **Research** quando a proposta:

- altera regras financeiras;
- muda contrato público da API;
- exige migration;
- impacta mobile e API ao mesmo tempo;
- adiciona uma nova superfície de produto;
- muda premissas de segurança.

Para contribuições menores, abra uma issue de **Feature**, **Improvement**, **Bug**, **Docs** ou **Test**, seguindo o guia em [CONTRIBUTING.md](../CONTRIBUTING.md).
