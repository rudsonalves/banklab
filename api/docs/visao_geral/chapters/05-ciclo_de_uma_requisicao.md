# Ciclo de uma Requisição

Este capítulo descreve o caminho percorrido por uma requisição dentro da API. O objetivo é ajudar novos desenvolvedores a entender como uma chamada HTTP atravessa as camadas da aplicação, onde cada responsabilidade é exercida e como o sistema transforma entrada externa em regra de negócio, persistência e resposta.

## Visão geral do fluxo

Uma requisição típica percorre o sistema na seguinte direção:

```text
HTTP Request
  -> Delivery handler
  -> Application use case
  -> Domain rules/contracts
  -> Infrastructure implementation
  -> Application result
  -> Delivery response
  -> HTTP Response
```

Esse fluxo representa a forma como a arquitetura se manifesta em tempo de execução.

A chamada chega pela borda HTTP, é adaptada para um caso de uso, usa regras e contratos do domínio, aciona implementações concretas de persistência e retorna à camada HTTP como resposta estruturada.

## A entrada: `HTTP Request`

O ponto de início é uma requisição feita por um cliente externo, normalmente a aplicação mobile.

Essa requisição contém:

- método HTTP;
- path;
- headers;
- query params;
- body JSON, quando aplicável.

Exemplo:

```http
GET /accounts
Authorization: Bearer <access_token>
```

Esse exemplo pede à API a lista de contas do usuário autenticado.

Em termos arquiteturais, a requisição ainda está no nível de protocolo. Ela ainda não foi convertida em intenção de negócio. Essa conversão começa na próxima etapa.

## A borda HTTP: `delivery`

A camada `delivery` recebe a chamada e a interpreta do ponto de vista do contrato REST.

Suas responsabilidades nessa etapa incluem:

- validar o método e a rota;
- ler path params;
- ler query params;
- decodificar body JSON;
- extrair o usuário autenticado do contexto;
- verificar erros básicos de transporte e formato;
- montar a entrada do caso de uso.

No exemplo de `GET /accounts`, o handler:

- confirma a presença de contexto autenticado;
- rejeita query params inesperados;
- constrói a entrada esperada pelo caso de uso de listagem.

Esse ponto é importante: o handler trabalha com adaptação HTTP, não com a regra central da listagem.

## O contexto autenticado no fluxo

Antes de o caso de uso ser chamado, a requisição já pode ter passado por middleware de autenticação.

Esse middleware valida o token JWT e popula o contexto da requisição com dados do usuário autenticado, como:

- `user_id`;
- `role`;
- `email`;
- `customer_id`.

Essas informações passam a acompanhar o fluxo interno da requisição. Isso permite que o sistema use identidade e ownership sem depender de dados enviados pelo cliente no corpo da operação.

No caso de `GET /accounts`, isso significa que a API descobre o `customer_id` a partir do usuário autenticado, e não de um parâmetro informado pelo cliente.

## A operação de negócio: `application`

Depois da adaptação HTTP, a requisição chega a um caso de uso na camada `application`.

O caso de uso representa a intenção operacional do sistema.

No exemplo, o caso de uso é a listagem de contas do usuário autenticado. A `application` passa a ser responsável por:

- validar pré-condições da operação;
- garantir que o usuário possui contexto suficiente;
- acionar o contrato de repositório adequado;
- organizar o resultado do fluxo.

Essa é a etapa em que a requisição deixa de ser apenas um evento HTTP e passa a ser tratada como ação do sistema.

Para `GET /accounts`, o caso de uso valida, por exemplo, se o usuário autenticado possui `customer_id`. Sem isso, a operação não pode determinar a quem pertencem as contas a serem listadas.

## O papel do `domain`

O domínio participa do fluxo fornecendo contratos, tipos, invariantes e erros.

Nem toda requisição aciona uma regra de entidade complexa, mas toda requisição relevante depende, de alguma forma, do modelo conceitual do sistema.

No caso da listagem de contas, o domínio define o contrato necessário para buscar contas por `customer_id`. Esse contrato expressa uma necessidade do negócio:

```text
listar contas pertencentes ao customer autenticado
```

O domínio também define erros e conceitos que a aplicação pode usar para comunicar falhas de forma consistente.

Em fluxos mais complexos, como depósito, saque e transferência, o domínio também fornece regras mais explícitas, como validade de valor, status de conta e saldo suficiente.

## A implementação concreta: `infrastructure`

Depois que a `application` decide o que precisa fazer, ela delega à infraestrutura a execução concreta das operações técnicas.

No projeto atual, essa etapa costuma significar:

- executar queries SQL;
- consultar ou atualizar PostgreSQL;
- iniciar ou participar de transações;
- aplicar locks de banco;
- mapear linhas e registros para estruturas da aplicação.

No exemplo de `GET /accounts`, a infraestrutura implementa a busca por contas em PostgreSQL, filtrando por `customer_id` e ordenando os resultados de forma previsível.

Conceitualmente, a operação executada é equivalente a:

```sql
SELECT ...
FROM accounts
WHERE customer_id = $1
ORDER BY created_at ASC, id ASC
```

Essa implementação concreta continua escondida da `delivery` e do cliente. O contrato da API não depende de saber como a consulta foi implementada.

## O retorno do resultado

Depois de executar a operação concreta, a infraestrutura devolve o resultado à `application`.

A camada de aplicação então:

- interpreta o retorno;
- trata falhas conhecidas;
- organiza a saída do caso de uso;
- devolve esse resultado à camada `delivery`.

No caso de listagem de contas, isso significa retornar uma coleção de contas já apta a ser serializada como resposta HTTP.

## A resposta HTTP: `delivery` novamente

Com o resultado do caso de uso em mãos, a `delivery` volta a atuar.

Agora seu papel é:

- escolher o status code adequado;
- converter o resultado em payload HTTP;
- aplicar o envelope padrão da API;
- serializar a resposta.

Em caso de sucesso, a resposta segue o formato:

```json
{
  "data": [
    {
      "id": "fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b",
      "customer_id": "6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3",
      "number": "10000001",
      "branch": "0001",
      "status": "active"
    }
  ],
  "error": null
}
```

Em caso de falha, a `delivery` traduz o erro para o status code e o corpo esperados pelo contrato da API.

## Exemplo completo: `GET /accounts`

O fluxo completo pode ser resumido assim:

```text
GET /accounts
  -> middleware valida JWT e popula contexto autenticado
  -> account/delivery.ListAccounts
  -> account/application/account.ListAccounts
  -> account/domain.AccountRepository.ListByCustomerID
  -> account/infrastructure.PostgresRepository.ListByCustomerID
  -> account/application retorna resultado
  -> account/delivery monta envelope JSON
  -> HTTP response
```

Esse exemplo é útil porque mostra um fluxo relativamente simples, mas suficientemente completo para ilustrar:

- autenticação;
- adaptação HTTP;
- uso de caso de uso;
- dependência de contrato de domínio;
- persistência concreta;
- resposta padronizada.

## Como esse fluxo aparece em operações mais complexas

O mesmo modelo geral vale para operações mais sensíveis, como depósito, saque, transferência e aprovação de usuário.

A diferença é que, nesses casos, a `application` assume responsabilidades adicionais, como:

- abertura de transação;
- coordenação entre múltiplos repositórios;
- aplicação de locks;
- validação de regras de domínio mais ricas;
- garantia de atomicidade.

Mesmo assim, a estrutura do fluxo permanece a mesma: entrada HTTP, adaptação, caso de uso, uso de domínio, execução concreta e resposta.

## Por que esse entendimento é importante

Compreender o ciclo de uma requisição ajuda a responder perguntas frequentes no desenvolvimento:

- onde devo validar este dado?
- esta falha deve virar `400`, `403` ou `422`?
- esta lógica fica no handler ou no caso de uso?
- esta regra deve ser modelada no domínio?
- esta mudança depende de query SQL ou de política de aplicação?

Sem essa visão, é comum que novas implementações misturem responsabilidades e acabem colocando lógica importante na camada errada.

## Síntese

O ciclo de uma requisição traduz a arquitetura em execução real.

A chamada entra por HTTP, é interpretada pela `delivery`, executada como caso de uso na `application`, apoiada por regras e contratos do `domain`, concretizada pela `infrastructure` e devolvida ao cliente como resposta padronizada.

Esse caminho é uma das bases do sistema. Ele ajuda a manter separação de responsabilidades, melhora a testabilidade e reduz o acoplamento entre protocolo, regra de negócio e persistência. Para novos desenvolvedores, entender esse fluxo é essencial para modificar a aplicação com segurança e coerência arquitetural.
