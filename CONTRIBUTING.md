# Guia de contribuição

Obrigado por considerar contribuir com o BankLab.

Este projeto é um laboratório local para estudar engenharia de software aplicada a um domínio financeiro simplificado. A ideia é manter um ambiente acolhedor para quem está aprendendo, mas também organizado o suficiente para que cada contribuição preserve clareza, consistência e qualidade técnica.

## Antes de começar

Leia primeiro:

- [README.md](README.md), para entender objetivo, escopo e como rodar o projeto;
- [docs/README.md](docs/README.md), para entender a organização da documentação;
- [docs/backlogs/README.md](docs/backlogs/README.md), para acompanhar discussões ativas e histórico de decisões;
- [api/README.md](api/README.md), se for mexer na API;
- [mobile/README.md](mobile/README.md), se for mexer no app Flutter;
- [api/docs/ARCHITECTURE.md](api/docs/ARCHITECTURE.md), se a mudança tocar arquitetura ou regras importantes.

Antes de implementar, tente responder:

- qual problema esta mudança resolve?
- ela é uma feature, melhoria, bug, pesquisa ou documentação?
- qual área do projeto é dona da mudança?
- existe impacto em contrato HTTP, banco de dados, saldo, ledger ou autenticação?

Mudanças financeiras precisam de cuidado extra. Se a alteração mexer em saldo, transações, extrato, concorrência ou persistência, abra uma issue ou descreva bem a decisão antes de enviar o pull request.

## Formas de contribuir

Você pode ajudar em várias frentes:

- **Backend Go**: casos de uso, handlers, domínio, persistência, autenticação e testes.
- **Mobile Flutter**: telas, navegação, estados de carregamento, erros, integração com API e experiência de uso.
- **Testes**: cobertura de regras, fluxos críticos, casos de erro e cenários de integração.
- **Documentação**: guias, exemplos de payload, diagramas, decisões técnicas e tutoriais de setup.
- **Produto financeiro**: refinamento de fluxos, regras de negócio, nomes e limites do domínio.
- **Infra e DevEx**: scripts, Makefile, Docker, CI, Postman e melhoria do ambiente local.

## Organização das issues

Toda issue deve deixar claro o tipo, a área e a prioridade do trabalho.

### Tipo

- **Feature**: adiciona uma nova funcionalidade.
- **Improvement**: melhora algo existente sem mudar a intenção principal.
- **Bug**: corrige comportamento incorreto.
- **Research**: investiga uma decisão antes da implementação.
- **Docs**: melhora documentação, exemplos ou guias.
- **Test**: adiciona ou melhora cobertura de testes.

### Área

- **Auth**: autenticação, JWT, refresh token, sessão e autorização.
- **Account**: contas, saldo, extrato e ciclo de vida da conta.
- **Ledger**: transações, consistência, transferências e núcleo financeiro.
- **Customer**: dados de cliente e identidade.
- **Admin**: aprovação, provisionamento e fluxos administrativos.
- **Mobile**: aplicativo Flutter.
- **Security**: endurecimento de segurança, tokens, roles e proteção de rotas.
- **Infra**: Docker, banco, migrations, scripts e configuração.
- **Docs**: documentação geral ou técnica.
- **DevEx**: experiência de desenvolvimento, automações e organização do projeto.

### Prioridade

- **High**: bloqueia evolução importante, corrige risco relevante ou afeta fluxo crítico.
- **Medium**: importante, mas não bloqueia o projeto.
- **Low**: melhoria pequena, ajuste futuro ou refinamento.

## Sugestão de issue

```md
## Contexto

Explique o problema ou oportunidade.

## Objetivo

Descreva o resultado esperado.

## Escopo

- O que entra na mudança
- O que fica fora da mudança

## Área

Exemplo: Ledger

## Tipo

Exemplo: Improvement

## Prioridade

Exemplo: Medium

## Critérios de aceite

- [ ] comportamento esperado 1
- [ ] comportamento esperado 2
- [ ] testes/documentação atualizados quando necessário
```

## Fluxo recomendado

1. Escolha uma issue ou proponha uma nova.
2. Comente na issue que pretende trabalhar nela.
3. Crie uma branch curta e descritiva.
4. Faça mudanças pequenas, focadas e revisáveis.
5. Rode os testes relevantes.
6. Abra um pull request explicando o que mudou e como foi validado.

Exemplos de branch:

```text
feature/internal-transfer-tests
fix/mobile-login-error
docs/api-payload-examples
improvement/account-statement-empty-state
```

## Padrão de commits

Prefira commits pequenos e com mensagem objetiva.

Exemplos:

```text
docs: melhora guia de contribuição
test: cobre transferência com saldo insuficiente
fix: corrige estado de erro no login mobile
feat: adiciona recibo de transferência no mobile
```

Não é obrigatório seguir Conventional Commits com rigidez absoluta, mas mensagens nesse formato ajudam a entender o histórico.

## Rodando o projeto

Na raiz do repositório:

```bash
make help
make setup
make run
```

Para subir apenas o Docker:

```bash
make docker-up
```

Para resetar o ambiente local:

```bash
make reset
```

Use `make reset` com cuidado, pois ele recria o banco local.

## Testes

Antes de abrir um pull request, rode os testes relacionados à sua mudança.

API:

```bash
make api-tests
```

Mobile:

```bash
make mobile-tests
```

Todos os testes:

```bash
make tests
```

Se você não conseguiu rodar algum teste, informe isso no pull request e explique o motivo.

## Cuidados na API

A API segue a ideia de monólito modular com camadas bem definidas.

Use este mapa como referência:

```text
delivery       # HTTP, request, response, status code e parsing
application    # casos de uso, orquestração e transações
domain         # entidades, invariantes, erros e contratos
infrastructure # banco, queries, repositórios e integrações concretas
```

Boas práticas:

- coloque regras de negócio no domínio quando elas forem invariantes do sistema;
- coloque coordenação de fluxo na aplicação;
- coloque detalhes HTTP na camada de delivery;
- coloque SQL, locks e mapeamento de banco na infrastructure;
- atualize migrations quando alterar o banco;
- atualize documentação REST quando mudar rota, payload, status code ou erro.

Mudanças em ledger, saldo e transferência devem preservar:

- registro auditável de movimentações;
- atomicidade da operação;
- consistência em caso de erro;
- autorização e ownership da conta;
- testes para sucesso e falha.

## Cuidados no mobile

Ao contribuir no Flutter:

- mantenha a separação entre `core`, `data`, `domain` e `ui`;
- preserve os contratos da API;
- trate estados de carregamento, erro e vazio;
- evite acoplar tela diretamente a detalhes de transporte quando já existir camada apropriada;
- rode os testes do escopo alterado.

Se a mudança mexer em chamadas HTTP, confira também se a API local está alinhada com o app.

## Pull requests

Um bom pull request deve responder:

- o que mudou?
- por que mudou?
- como testar?
- existe impacto em API, mobile, banco, contrato ou documentação?
- quais testes foram rodados?

Modelo sugerido:

```md
## Resumo

- item 1
- item 2

## Como validar

- comando ou fluxo manual usado

## Impactos

- API:
- Mobile:
- Banco:
- Documentação:

## Testes

- [ ] make api-tests
- [ ] make mobile-tests
- [ ] make tests
```

## Documentação

Atualize a documentação quando a mudança:

- cria ou altera endpoint;
- muda request, response, status code ou erro;
- altera comportamento de autenticação;
- muda regra de negócio;
- cria migration;
- afeta setup local;
- muda fluxo visível no mobile.

Backlogs em `docs/backlogs/` também fazem parte da documentação do projeto. Use essa pasta para registrar discussões, decisões e recortes de implementação antes de mudanças relevantes. Quando um backlog for resolvido ou substituído por decisão mais nova, mova-o para o diretório `done/` correspondente, preservando o histórico.

## Comunicação

Este projeto valoriza colaboração clara e respeitosa. Perguntas são bem-vindas, inclusive perguntas básicas. A documentação pode estar em português ou inglês, conforme o contexto; o importante é manter clareza e vontade de melhorar o projeto junto.

Se estiver em dúvida entre implementar direto ou abrir uma discussão, abra uma issue de **Research** ou descreva a dúvida no pull request.

## Fora do escopo de contribuição

Por enquanto, evite contribuições que adicionem:

- integração real com instituições financeiras;
- uso de dados bancários reais;
- coleta de dados sensíveis desnecessários;
- Pix, TED, antifraude ou conciliação externa sem discussão prévia;
- grandes refatorações sem issue e sem plano claro.

BankLab é um projeto educacional e experimental. Ele não deve ser usado como sistema bancário real.
