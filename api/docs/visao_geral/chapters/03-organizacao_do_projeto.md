# Organização do Projeto e Mapa de Módulos

Este capítulo apresenta a organização estrutural da API no repositório. O objetivo é ajudar novos desenvolvedores a localizar rapidamente as partes relevantes do sistema e entender como o código foi distribuído entre ponto de entrada, módulos de negócio, infraestrutura compartilhada, migrações e documentação.

## Visão geral da estrutura

No nível mais alto, a API está organizada em torno de alguns diretórios principais:

```shell
api
├── cmd
│   └── api
├── docs
├── internal
│   ├── account
│   ├── admin
│   ├── auth
│   ├── bootstrap
│   ├── customer
│   ├── database
│   └── shared
└── migrations
```

Cada uma dessas áreas tem um papel diferente no ciclo de vida da aplicação.

## `cmd/`: ponto de entrada da aplicação

O diretório `cmd/` contém os executáveis da aplicação. No estágio atual, o principal ponto de entrada da API está em:

```text
api/cmd/api/
```

É nesse diretório que fica o `main.go`, responsável por:

- inicializar o processo;
- criar dependências concretas;
- montar casos de uso e handlers;
- registrar middlewares e rotas;
- subir o servidor HTTP.

Esse diretório não é o lugar da regra de negócio. Ele é a raiz de composição do sistema. Seu papel é conectar as peças da aplicação.

## `internal/`: implementação principal da API

O diretório `internal/` concentra o núcleo implementado da aplicação.

Em Go, `internal` também possui um significado prático importante: ele restringe a importação desses pacotes ao próprio módulo, evitando que esse código seja usado como biblioteca externa por outros projetos. Isso reforça a ideia de que essa estrutura pertence à implementação interna da API.

Dentro de `internal/`, o sistema está organizado principalmente por módulos de negócio:

- `auth`;
- `account`;
- `customer`;
- `admin`.

Além desses módulos, existem diretórios de suporte compartilhado:

- `database`;
- `shared`;
- `bootstrap`.

## Módulos de negócio

### `internal/auth`

O módulo `auth` concentra autenticação, sessão e identidade do usuário.

Nele ficam responsabilidades como:

- registro de usuário;
- login;
- refresh token;
- leitura do usuário autenticado atual;
- integração com JWT;
- persistência de sessão.

Esse módulo é responsável por estabelecer a identidade do chamador e disponibilizar o contexto autenticado usado pelas demais áreas da aplicação.

### `internal/account`

O módulo `account` concentra o domínio de conta bancária e operações financeiras.

Nele ficam responsabilidades como:

- criação de conta;
- listagem de contas;
- consulta de saldo;
- depósito;
- saque;
- transferência;
- extrato;
- regras de status de conta;
- ledger de transações.

Esse é o módulo mais sensível em termos de consistência e concorrência, porque lida diretamente com saldo e histórico financeiro.

### `internal/customer`

O módulo `customer` representa o cliente de negócio vinculado ao usuário autenticado.

Nele ficam responsabilidades como:

- criação do customer no fluxo de registro;
- leitura do customer associado ao usuário atual;
- contratos e erros relacionados a dados de cliente.

Esse módulo é importante porque separa a identidade de autenticação da entidade de cliente usada no domínio.

### `internal/admin`

O módulo `admin` concentra operações restritas a administradores.

No estágio atual, sua principal responsabilidade implementada é:

- aprovação de usuário.

Esse módulo coordena fluxos administrativos, mas não se torna automaticamente dono das regras de outros módulos. Uma operação pode ser administrativa por autorização de acesso e ainda assim depender de regras pertencentes a `account`, `customer` ou `auth`.

## Módulos compartilhados

### `internal/database`

O diretório `database` concentra recursos compartilhados de acesso ao banco.

Seu papel inclui:

- criação do pool de conexão PostgreSQL;
- utilitários relacionados a transação;
- abstrações de apoio ao uso do banco.

Esse diretório existe para evitar duplicação de infraestrutura comum entre módulos.

### `internal/shared`

O diretório `shared` contém componentes realmente transversais ao sistema.

No estágio atual, isso inclui principalmente elementos reutilizados por mais de um módulo sem pertencer claramente a apenas um deles, como contexto autenticado compartilhado.

Esse diretório deve ser usado com cuidado. Ele é útil para recursos transversais genuínos, mas não deve se tornar um local genérico para código sem dono claro.

### `internal/bootstrap`

O diretório `bootstrap` concentra rotinas de inicialização da aplicação.

Seu papel é preparar aspectos necessários antes do servidor começar a atender requisições, mantendo essas responsabilidades separadas da lógica de negócio.

## Anatomia interna de cada módulo

Cada módulo principal segue uma estrutura em camadas:

```text
internal/<module>/
|-- delivery/
|-- application/
|-- domain/
`-- infrastructure/
```

Essa organização permite que cada contexto de negócio contenha:

- sua entrada HTTP;
- seus casos de uso;
- suas regras e contratos;
- suas implementações técnicas.

No módulo `account`, a camada `application` ainda possui subdivisões por grupo de caso de uso:

```text
internal/account/application/
|-- account/
|-- transaction/
`-- statement/
```

Essa subdivisão existe porque o módulo de conta concentra mais fluxos do que os demais e se beneficia de uma organização interna mais detalhada.

## `migrations/`: evolução do schema

O diretório `migrations/` contém a evolução incremental do banco de dados.

Seu papel é registrar mudanças estruturais de schema de forma versionada, reproduzível e controlada. Isso é importante porque, em uma aplicação de domínio financeiro, a estrutura persistida tem impacto direto nas regras de negócio, na consistência e na auditabilidade do sistema.

Para um novo desenvolvedor, esse diretório é especialmente útil quando a mudança proposta envolve:

- criação de novas tabelas;
- alteração de colunas;
- índices;
- constraints;
- evolução do modelo de persistência.

## `db/`: schema base e artefatos de banco

O diretório `db/` concentra artefatos relacionados ao banco, incluindo o schema base usado como referência estrutural da aplicação.

Esse diretório complementa `migrations/`: enquanto as migrations mostram a evolução incremental, o schema base ajuda a entender a forma geral do banco.

## `docs/`: documentação técnica

O diretório `docs/` reúne a documentação técnica da API.

Ele contém materiais como:

- visão arquitetural;
- fluxos de caso de uso;
- documentação REST;
- documentação de banco;
- documentação de implementação;
- visão geral para onboarding.

Esse diretório é parte importante do trabalho de desenvolvimento. Ele não deve ser visto como apêndice opcional, mas como apoio real para entendimento e evolução do sistema.

## `README.md`: ponto de partida rápido

O `README.md` da API oferece uma visão inicial mais resumida:

- stack;
- arquitetura em alto nível;
- lista principal de rotas;
- setup local;
- comandos de teste;
- mapa resumido de diretórios;
- links para a documentação.

Ele funciona como a primeira leitura para quem está entrando no projeto e quer rapidamente subir o ambiente ou localizar a documentação principal.

## Como navegar no código a partir de uma feature

Uma boa forma de usar a organização do projeto é começar pela feature que se deseja entender ou alterar.

Por exemplo:

- se a mudança envolve login ou token, o ponto inicial tende a ser `internal/auth`;
- se envolve conta, saldo ou transferência, o ponto inicial tende a ser `internal/account`;
- se envolve perfil de cliente, o ponto inicial tende a ser `internal/customer`;
- se envolve aprovação ou acesso administrativo, o ponto inicial tende a ser `internal/admin`.

Dentro do módulo, a leitura pode seguir a camada adequada:

- `delivery`, se a dúvida é sobre rota, request ou response;
- `application`, se a dúvida é sobre fluxo do caso de uso;
- `domain`, se a dúvida é sobre regra ou contrato;
- `infrastructure`, se a dúvida é sobre banco, SQL ou persistência.

Esse padrão de navegação reduz o tempo de exploração e ajuda a evitar mudanças feitas no lugar errado.

## Consequência prática da organização atual

A organização do projeto foi pensada para responder uma necessidade comum de manutenção: localizar responsabilidade com rapidez.

Quando o código está estruturado por módulo e camada, fica mais fácil entender:

- onde implementar uma nova feature;
- onde procurar uma regra existente;
- onde corrigir um problema de persistência;
- onde alterar um contrato REST;
- onde validar se uma mudança afeta domínio ou apenas adaptação técnica.

Essa clareza é especialmente útil em um projeto que já passou por evolução arquitetural e que pode continuar crescendo com apoio de múltiplos colaboradores e ferramentas de IA.

## Síntese

A estrutura atual da API não é apenas uma convenção de pastas. Ela é a forma concreta como a arquitetura do sistema aparece no repositório.

`cmd/api` concentra o startup e o wiring. `internal/` concentra os módulos de negócio e os componentes compartilhados. `migrations/` e `db/` registram o modelo persistido. `docs/` e `README.md` apoiam entendimento e onboarding.

Para novos desenvolvedores, esse mapa é um ponto de partida essencial. Antes mesmo de aprofundar em regras específicas, entender a distribuição do projeto ajuda a localizar melhor cada responsabilidade e a trabalhar com mais segurança dentro da base de código.
