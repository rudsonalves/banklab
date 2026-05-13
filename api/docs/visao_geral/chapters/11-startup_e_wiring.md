# Startup, Wiring e Composição da Aplicação

Este capítulo descreve como a API é montada no momento em que o processo é iniciado. O objetivo é explicar onde as dependências concretas são instanciadas, como os módulos são conectados entre si e por que essa composição é centralizada em um ponto específico da aplicação.

## O que significa startup e wiring

No contexto desta API:

- **startup** é o processo de inicialização da aplicação;
- **wiring** é o ato de conectar dependências concretas entre si.

Ao longo do código, os casos de uso dependem de contratos, e os handlers dependem de casos de uso. No entanto, em algum ponto do sistema é necessário decidir quais implementações reais serão usadas para satisfazer esses contratos.

Essa decisão não deve ficar espalhada por vários módulos. Ela é concentrada no ponto de entrada da aplicação.

## Ponto central de composição

O principal ponto de composição da API está em:

```text
api/cmd/api/main.go
```

Esse arquivo funciona como raiz de composição do sistema.

É nele que a aplicação:

- inicializa infraestrutura compartilhada;
- cria implementações concretas;
- instancia casos de uso;
- instancia handlers;
- registra rotas e middlewares;
- sobe o servidor HTTP.

Ele não é o lugar onde a regra de negócio mora. Seu papel é montar o ambiente em que a regra de negócio será executada.

## Por que centralizar a composição

Centralizar o wiring em um único ponto traz algumas vantagens importantes.

### Clareza

Quando a composição está concentrada em `main.go`, fica mais fácil entender:

- quais implementações concretas existem;
- quais casos de uso dependem de quais componentes;
- quais handlers expõem quais fluxos;
- quais rotas estão ligadas a quais módulos.

### Controle

O sistema evita que dependências sejam instanciadas de forma oculta ou espalhada. Isso reduz acoplamento invisível e facilita alterações controladas.

### Coerência arquitetural

Os módulos não precisam “procurar” globalmente suas dependências nem se autoconfigurar. Eles permanecem focados em sua responsabilidade principal, enquanto o ponto de entrada assume a responsabilidade pela montagem.

## Etapas típicas do startup

O startup da aplicação segue uma sequência lógica.

## 1. Inicialização da infraestrutura base

O primeiro passo é preparar dependências técnicas compartilhadas, especialmente as que sustentam a persistência e integrações centrais.

No estágio atual, isso inclui principalmente:

- criação do pool de conexão PostgreSQL;
- preparação de componentes auxiliares de infraestrutura;
- inicialização de recursos necessários antes do atendimento das requisições.

Sem essa etapa, a aplicação não teria como instanciar repositórios concretos nem executar operações persistidas.

## 2. Instanciação das implementações concretas

Depois da infraestrutura base, a aplicação instancia implementações concretas dos módulos.

Isso inclui, por exemplo:

- repositórios de autenticação;
- repositórios de customer;
- repositórios de account;
- componentes ligados a token, hashing e políticas auxiliares.

Esse é o momento em que as abstrações encontradas na `application` e no `domain` recebem uma implementação real.

Enquanto os casos de uso conhecem contratos, o `main.go` conhece structs concretas.

## 3. Instanciação dos casos de uso

Com as dependências concretas disponíveis, a aplicação passa a construir os casos de uso.

Cada caso de uso recebe explicitamente aquilo de que precisa para funcionar.

Exemplos conceituais:

- login recebe componentes de autenticação, repositório de usuário e serviço de token;
- listagem de contas recebe repositório de contas;
- criação de conta recebe dependências relacionadas a customer, autorização e sequência de conta;
- aprovação de usuário recebe repositório de usuário, repositório de customer, repositório de account e política de branch;
- transferência recebe repositórios e recursos transacionais necessários para orquestrar débito, crédito e ledger.

Essa montagem explícita é importante porque deixa visível o contrato real de dependência de cada fluxo.

## 4. Instanciação dos handlers

Depois dos casos de uso, a aplicação instancia os handlers HTTP.

Cada handler agrupa os casos de uso que expõe por meio das rotas de seu módulo.

Por exemplo:

- o handler de `auth` reúne registro, login, refresh e `me`;
- o handler de `customer` reúne o fluxo de leitura do customer atual;
- o handler de `account` reúne listagem de contas, criação de conta, saldo, depósito, saque, transferência e extrato;
- o handler de `admin` reúne o fluxo de aprovação de usuário.

Isso reforça a ideia de que o handler é uma camada de adaptação HTTP para um conjunto coerente de casos de uso do módulo.

## 5. Registro de rotas e middlewares

Com os handlers montados, o `main.go` registra:

- rotas HTTP;
- métodos associados;
- proteção por autenticação;
- middlewares necessários.

É nessa etapa que o sistema define o contrato executável entre superfície REST e módulos internos.

Exemplos atuais incluem:

- `POST /auth/register`;
- `POST /auth/login`;
- `POST /auth/refresh`;
- `GET /auth/me`;
- `POST /admin/users/{id}/approve`;
- `GET /customers/me`;
- `GET /accounts`;
- `POST /accounts`;
- `GET /accounts/{id}/balance`;
- `POST /terminal/accounts/{id}/deposit`;
- `POST /terminal/accounts/{id}/withdraw`;
- `GET /accounts/internal-transfers/recipients`;
- `POST /accounts/internal-transfers`;
- `GET /accounts/{id}/statement`.

Também é nessa fase que rotas protegidas passam a receber o middleware responsável por validar JWT e popular o contexto autenticado.

## 6. Inicialização do servidor

Por fim, com infraestrutura, casos de uso, handlers e rotas devidamente conectados, o servidor HTTP é iniciado e passa a aceitar requisições.

Nesse momento, a arquitetura deixa de ser apenas estrutura de código e passa a operar como sistema em execução.

## O papel do `main.go` em relação às camadas

É importante entender que o `main.go` não quebra a arquitetura em camadas. Ele a conecta.

Ele não deve:

- carregar regra de negócio;
- implementar validação de domínio;
- decidir autorização de fluxos específicos;
- substituir casos de uso por lógica improvisada.

Seu papel é exclusivamente de composição.

Em termos práticos:

- a regra continua na `application` e no `domain`;
- a implementação técnica continua na `infrastructure`;
- a entrada HTTP continua na `delivery`;
- o `main.go` apenas conecta tudo.

## Por que isso melhora a manutenção

Para um novo desenvolvedor, esse modelo de composição explícita é muito útil.

Quando uma feature precisa ser rastreada de ponta a ponta, o `main.go` ajuda a responder:

- qual implementação concreta está sendo usada;
- como o handler daquele módulo foi montado;
- quais casos de uso um módulo expõe;
- quais dependências reais entram em cada fluxo;
- em que ponto uma nova rota precisa ser registrada.

Isso acelera a navegação no sistema e reduz a sensação de “mágica” na inicialização.

## Comparação implícita com DI automatizada

A aplicação utiliza composição explícita em Go, em vez de depender de um framework complexo de injeção de dependência.

Essa escolha combina bem com o tamanho e o estágio do projeto porque:

- mantém o wiring visível;
- reduz camadas adicionais de abstração;
- facilita debugging;
- torna a montagem mais previsível.

Isso não impede evoluções futuras, mas no momento favorece simplicidade e legibilidade.

## Cuidados práticos ao evoluir o wiring

Ao adicionar novos fluxos ou dependências, é importante preservar a clareza da composição.

Isso significa:

- adicionar novas dependências de forma explícita;
- evitar inicialização escondida dentro dos módulos;
- manter handlers coerentes com seus casos de uso;
- registrar novas rotas no ponto de composição correto;
- não deslocar regra de negócio para o startup apenas por conveniência.

Se o wiring começar a crescer, a solução ideal tende a ser organizar melhor a composição, e não enfraquecer a separação das camadas.

## Síntese

O startup da aplicação é o momento em que a arquitetura é transformada em sistema executável.

O `cmd/api/main.go` centraliza a composição das dependências concretas, instancia repositórios, casos de uso e handlers, registra middlewares e rotas e inicia o servidor HTTP. Essa centralização torna a estrutura mais legível, reduz acoplamento invisível e facilita o entendimento de como os módulos realmente se conectam em tempo de execução.

Para novos desenvolvedores, compreender esse ponto de composição é essencial. Ele mostra como as abstrações do sistema são materializadas no runtime e como uma funcionalidade deixa de ser apenas código espalhado em módulos para se tornar uma rota viva da aplicação.
