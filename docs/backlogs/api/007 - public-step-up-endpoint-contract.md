# Backlog: contrato público de step-up por operação HTTP

## 1. Contexto

O MVP ZTA da API foi concluído no backlog `006`. O fluxo atual de autorização
de step-up usa o campo público `endpoint_key`:

```json
{
  "endpoint_key": "internal_transfer.create",
  "transaction_password": "123456"
}
```

Isso funciona tecnicamente, mas vaza para clientes mobile uma chave de política
interna da API. O app não deveria conhecer nomes como
`internal_transfer.create`. O cliente conhece a superfície pública que deseja
acessar: método HTTP e path.

## 2. Objetivo

Alterar o contrato público de autorização de step-up para que o cliente solicite
permissão para uma operação HTTP pública, por exemplo:

```json
{
  "method": "POST",
  "path": "/accounts/internal-transfers",
  "transaction_password": "123456"
}
```

A API deve traduzir internamente essa operação pública para a chave de policy
interna correspondente, como `internal_transfer.create`.

## 3. Decisão de contrato

### Contrato atual a substituir

```json
{
  "endpoint_key": "internal_transfer.create",
  "transaction_password": "123456"
}
```

### Novo contrato público

```json
{
  "method": "POST",
  "path": "/accounts/internal-transfers",
  "transaction_password": "123456"
}
```

Regras:

- `method` deve usar o verbo HTTP público, em maiúsculas.
- `path` deve usar o path público documentado da API.
- `path` deve ser apenas o path, sem scheme, host, query string ou fragmento.
- Para rotas parametrizadas futuras, o cliente deve enviar o template público
  documentado, como `/accounts/{id}/withdraw`, não o valor concreto do recurso.
- O cliente não envia nem conhece `endpoint_key`.
- A API resolve `method + path` para uma policy interna.
- O step-up token emitido continua vinculado à policy interna usada pelo
  enforcement.
- O enforcement de `POST /accounts/internal-transfers` continua exigindo token
  válido para a policy interna da transferência.

## 4. Escopo de API

- Atualizar o request de `POST /security/step-up/authorize`.
- Introduzir um resolvedor interno de operação pública para policy interna.
- Mapear inicialmente:
  - `POST /accounts/internal-transfers` -> `internal_transfer.create`.
- Manter a whitelist de operações permitidas no backend.
- Atualizar validação de payload:
  - exigir `method`;
  - exigir `path`;
  - exigir `transaction_password`;
  - rejeitar `endpoint_key` como campo público quando o contrato novo estiver
    ativo.
- Validar que `path`:
  - começa com `/`;
  - não contém scheme ou host;
  - não contém query string;
  - não contém fragmento;
  - usa templates públicos quando a rota tiver parâmetros, por exemplo
    `/accounts/{id}/withdraw`.
- Atualizar DTOs de delivery.
- Atualizar use case de autorização para receber a operação pública ou o
  resultado resolvido, sem expor regra HTTP ao domínio de token.
- Preservar o claim interno `endpoint_key` no JWT, se essa continuar sendo a
  melhor representação de enforcement. `endpoint_key` deixa de ser input
  público, mas pode continuar existindo como claim interno assinado.
- Atualizar testes de delivery, application e policy.
- Atualizar toda a documentação afetada, incluindo REST, ZTA e referências de
  consumo mobile.

## 5. Compatibilidade

Decisão inicial para o projeto:

- não manter compatibilidade com `endpoint_key` no contrato público mobile;
- atualizar documentação e testes para o novo contrato;
- se for necessário período de transição no futuro, isso deve ser uma decisão
  explícita.

## 6. Erros

`STEP_UP_ENDPOINT_NOT_ALLOWED` deve continuar sendo retornado quando a operação
pública solicitada não estiver na whitelist de step-up.

Exemplos:

- método não suportado para o path;
- path público inexistente;
- path existente, mas não elegível a step-up;
- combinação `method + path` sem policy configurada.

`INVALID_REQUEST` deve cobrir JSON inválido ou campos inesperados.

`INVALID_DATA` deve cobrir método ou path ausentes/malformados, conforme o
padrão atual da API.

## 7. Critérios de aceite

- Mobile não precisa conhecer `internal_transfer.create`.
- `POST /security/step-up/authorize` aceita `method`, `path` e
  `transaction_password`.
- `POST /security/step-up/authorize` rejeita operação pública fora da whitelist.
- Token emitido para `POST /accounts/internal-transfers` é aceito pelo
  enforcement da transferência interna.
- Token emitido para outra operação pública, quando houver, não é aceito na
  transferência interna.
- Testes cobrem resolução `method + path -> endpoint_key`.
- Testes cobrem rejeição de `path` com host, query string, fragmento ou valor
  concreto em rota parametrizada quando o template público for exigido.
- Testes cobrem payload antigo com `endpoint_key` sendo rejeitado ou removido
  do contrato.
- Documentação REST não orienta clientes a enviar `endpoint_key`.
- Documentação REST e mobile usam `path`, não `endpoint`, no contrato público
  de autorização de step-up.

## 8. Fora de escopo

- Criar novas policies além de transferência interna.
- Vincular o token ao payload detalhado da operação.
- Mudar o header de enforcement `X-Step-Up-Token`.
- Alterar a semântica de consumo único por `jti`.
- Alterar a validade de 120 segundos do token.

## 9. Referências

- `docs/backlogs/api/done/006b - step-up-token.md`
- `docs/backlogs/api/done/006c - internal-transfer-step-up-enforcement.md`
- `docs/backlogs/api/done/006d - zta-contracts-and-docs.md`
- `docs/backlogs/mobile/011 - senha-transacional-e-step-up.md`
- `api/docs/07-api-rest.md`
