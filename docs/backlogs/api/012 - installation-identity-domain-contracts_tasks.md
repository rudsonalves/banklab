# Tasks dos contratos de dominio da identidade de instalacao

Backlog pai:

- `012 - installation-identity-domain-contracts.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Area: API
- Tipo: Dominio/Contratos/Seguranca

## Task 1/8: Definir value object de identificador de instalacao

Status: Backlog

### Objetivo

Criar a representacao interna do `installation_id` recebido do cliente.

### Escopo

- Representar `installation_id` como UUID v4 canonico.
- Centralizar validacao pura de dominio, quando aplicavel.
- Diferenciar `installation_id` do identificador publico de gerenciamento.
- Manter a semantica de sinal fraco: o valor nao autentica e nao prova posse.

### Criterios de aceite

- `installation_id` tem validacao reutilizavel.
- Identificador publico de gerenciamento nao reutiliza o valor enviado pelo
  cliente.
- Testes cobrem valores validos e invalidos sem depender de banco ou HTTP.

### Depende de

- Backlog 011 para contrato de entrada do header.

## Task 2/8: Modelar instalacao de app e estados persistidos

Status: Backlog

### Objetivo

Definir a entidade de instalacao como associacao entre usuario e instalacao do
app.

### Escopo

- Modelar campos de dominio necessarios:
  - usuario;
  - `installation_id`;
  - identificador publico de gerenciamento;
  - status;
  - timestamps relevantes.
- Definir estados persistidos:
  - `known`;
  - `revoked`.
- Definir transicoes permitidas.
- Garantir que instalacao revogada nao volte a `known` por novo login.

### Criterios de aceite

- Estados persistidos ficam tipados no dominio.
- Transicao para `revoked` e representada.
- Nao existe estado `trusted` neste MVP.
- Testes cobrem criacao e revogacao logica em memoria.

### Depende de

- Task 1.

## Task 3/8: Definir classificacoes derivadas para login

Status: Backlog

### Objetivo

Criar a linguagem de decisao usada futuramente pelo login sem ainda ligar o
fluxo ao handler.

### Escopo

- Definir classificacoes:
  - `known`;
  - `first`;
  - `new`;
  - `revoked`;
  - `limit_reached`.
- Definir estrutura de decisao com dados auxiliares necessarios, como contagem
  de instalacoes conhecidas e limite maximo.
- Separar classificacao derivada dos estados persistidos.

### Criterios de aceite

- Classificacoes nao sao persistidas como status da instalacao.
- `first` representa ausencia historica de qualquer instalacao.
- `limit_reached` representa nova instalacao com tres instalacoes `known`.
- Testes cobrem a semantica das classificacoes em regras puras ou tabeladas.

### Depende de

- Task 2.

## Task 4/8: Definir erros de dominio e codigos de aplicacao

Status: Backlog

### Objetivo

Criar os erros necessarios para os proximos backlogs sem acoplar o dominio ao
delivery HTTP.

### Escopo

- Definir erro para instalacao divergente.
- Definir erro para instalacao revogada.
- Definir erro para limite de instalacoes atingido.
- Definir erro para autorizacao restrita ausente, invalida, expirada,
  consumida ou revogada.
- Definir quais erros precisam de codigo publico futuro.
- Nao mapear detalhes HTTP neste backlog, exceto quando o erro compartilhado
  ja existir por contrato.

### Criterios de aceite

- Erros sao comparaveis com `errors.Is` quando necessario.
- Erros carregam detalhes apenas quando isso for seguro.
- Nenhum erro expõe tokens, senha transacional ou hashes.

### Depende de

- Task 3.

## Task 5/8: Definir portas de leitura e classificacao de instalacoes

Status: Backlog

### Objetivo

Definir interfaces necessarias para descobrir o estado de uma instalacao sem
implementar banco.

### Escopo

- Porta para buscar instalacao por usuario e `installation_id`.
- Porta para contar instalacoes `known` por usuario.
- Porta para verificar se o usuario ja teve qualquer instalacao associada.
- Porta para buscar instalacao por identificador publico de gerenciamento.
- Definir contratos de retorno para ausencia, erro e estado invalido.

### Criterios de aceite

- Interfaces ficam no limite de dominio/aplicacao apropriado.
- Nenhuma interface importa pacote de Postgres ou HTTP.
- Contratos permitem implementar a classificacao do backlog 016.

### Depende de

- Task 4.

## Task 6/8: Definir portas de escrita e operacoes atomicas de instalacao

Status: Backlog

### Objetivo

Definir os contratos que permitirao implementar bootstrap, reserva de vaga e
revogacao sem expor detalhes de transacao ao delivery.

### Escopo

- Porta para bootstrap atomico da primeira instalacao.
- Porta para reservar vaga e criar instalacao `known` respeitando o limite.
- Porta para revogacao logica.
- Porta para invalidacao de sessoes por instalacao revogada, se ficar sob o
  mesmo limite de aplicacao.
- Definir comportamento esperado em conflito de concorrencia.

### Criterios de aceite

- Contratos deixam claro o que deve ser atomico.
- Contratos permitem impedir duas primeiras instalacoes concorrentes.
- Contratos permitem impedir mais de tres instalacoes `known`.
- Nenhuma implementacao Postgres e criada neste backlog.

### Depende de

- Task 5.

## Task 7/8: Modelar autorizacao restrita de registro

Status: Backlog

### Objetivo

Definir o dominio da autorizacao restrita que sera usada para registrar uma
nova instalacao depois do login.

### Escopo

- Modelar campos:
  - `jti`;
  - usuario;
  - `installation_id`;
  - escopo;
  - status;
  - expiracao;
  - consumo.
- Definir escopo `installation.register`.
- Definir estados:
  - `active`;
  - `consumed`;
  - `revoked`.
- Tratar expiracao como estado derivado de `expires_at`.
- Definir regras de consumo unico.

### Criterios de aceite

- Autorizacao restrita nao e sessao operacional.
- Autorizacao restrita nao implica refresh token.
- Consumo unico e representado no dominio.
- Testes cobrem ativa, expirada, consumida e revogada.

### Depende de

- Task 4.

## Task 8/8: Definir portas de autorizacao restrita

Status: Backlog

### Objetivo

Definir as interfaces que permitirao criar, validar, consumir e revogar
autorizacoes restritas nos proximos backlogs.

### Escopo

- Porta para criar autorizacao restrita.
- Porta para buscar por `jti`.
- Porta para consumir autorizacao ativa.
- Porta para revogar autorizacao.
- Porta para garantir no maximo uma autorizacao ativa por usuario,
  instalacao e escopo.
- Definir contratos de concorrencia e idempotencia quando necessario.

### Criterios de aceite

- Interfaces permitem emissao futura de `restricted_access_token`.
- Interfaces permitem validacao futura pelo middleware restrito.
- Interfaces permitem consumo atomico no registro da instalacao.
- Nenhuma migration ou implementacao Postgres e criada neste backlog.

### Depende de

- Task 7.
