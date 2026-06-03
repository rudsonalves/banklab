# Backlog: endpoint de sessão pós-login

## 1. Contexto

O mobile hoje monta o perfil do usuário após login combinando duas chamadas:

```text
GET /customers/me
GET /auth/me
```

No mobile, essa composição está concentrada em
`mobile/lib/data/services/apis/auth/auth_api.dart`, no método `getProfile`.

Esse desenho já funciona para montar `UserProfile`, mas começa a ficar limitado
para decisões de fluxo pós-login. O app passará a precisar de um snapshot único
para responder perguntas como:

- quem é o usuário autenticado?
- qual é o customer vinculado?
- o onboarding foi concluído?
- a conta está aprovada?
- existe conta operacional?
- o usuário possui senha transacional ativa?
- qual é o próximo passo obrigatório antes de liberar a Home?

Para o fluxo mobile de senha transacional, a backlog
`docs/backlogs/mobile/011 - cadastro-senha-transacional.md` definiu que, após
login, o app deve verificar se o usuário possui senha transacional ativa. Hoje
não existe um contrato limpo para isso.

## 2. Objetivo

Criar um endpoint autenticado de sessão pós-login:

```http
GET /auth/session
Authorization: Bearer <access_token>
```

Esse endpoint deve consolidar os dados necessários para o bootstrap do app
mobile após autenticação, evitando que o cliente precise juntar múltiplas
respostas e reimplementar regras de prontidão de usuário.

Decisão de contrato:

- `GET /auth/session` será o endpoint canônico de sessão para clientes
  autenticados;
- o primeiro consumidor será o mobile no fluxo pós-login;
- ele deve trazer, em uma única resposta, as informações que hoje o mobile
  precisa buscar em `GET /auth/me` e `GET /customers/me`;
- o método `getProfile` do mobile deve poder ser migrado para ler apenas
  `GET /auth/session`;
- campos de readiness, como `can_access_home` e `next_required_step`, fazem
  parte do estado de produto da sessão e podem ser usados por qualquer cliente,
  não apenas pelo mobile.
- neste primeiro corte, `next_required_step` pode não ser retornado; a API deve
  expor a prontidão principal por `can_access_home` e pelos status de
  readiness.

## 3. Decisão de contrato

`GET /auth/session` deve retornar um snapshot composto, com seções explícitas
para:

- identidade/autorização do usuário;
- dados cadastrais do customer, quando houver;
- estado de prontidão para navegação pós-login;
- próximo passo obrigatório, quando o usuário ainda não puder acessar a Home.

Formato inicial, respeitando o envelope padronizado da aplicação:

```json
{
  "data": {
    "user": {
      "id": "<user_id>",
      "email": "user@example.com",
      "phone": "+5527999999999",
      "role": "customer"
    },
    "customer": {
      "id": "<customer_id>",
      "name": "Maria Silva",
      "cpf": "12345678901",
      "birth_date": "1990-01-15",
      "created_at": "2026-05-29T10:00:00Z"
    },
    "readiness": {
      "onboarding_completed": true,
      "approved": true,
      "has_operational_account": true,
      "transaction_password_status": "active",
      "can_access_home": true
    }
  },
  "error": null
}
```

Na aplicação atual, um usuário customer só deve existir após o fluxo de criação
coletar os dados de user e customer. Portanto, para usuários customer,
`customer` não é opcional em uma sessão saudável.

Se um usuário customer autenticado não tiver customer associado, isso deve ser
tratado como problema de estado da conta, não como resposta normal de sessão.

Campos obrigatórios no primeiro corte:

- `user.id`;
- `user.email`;
- `user.phone`;
- `user.role`;
- `customer`, obrigatório para usuário customer;
- `customer.id`;
- `customer.name`;
- `customer.cpf`;
- `customer.birth_date`;
- `customer.created_at`;
- `readiness.onboarding_completed`;
- `readiness.approved`;
- `readiness.has_operational_account`;
- `readiness.transaction_password_status`;
- `readiness.can_access_home`.

`readiness.next_required_step` fica fora dos campos obrigatórios do primeiro
corte. Ele pode ser adicionado depois, quando a aplicação fechar a taxonomia de
próximos passos obrigatórios.

O primeiro corte deve cobrir todos os campos usados hoje para montar
`UserProfile` no mobile:

```text
user.id
user.email
user.phone
user.role
customer.id
customer.name
customer.cpf
customer.created_at
```

Assim, o mobile deixa de precisar chamar `GET /customers/me` e `GET /auth/me`
em sequência para montar o perfil após login.

Campos de readiness sugeridos:

```text
onboarding_completed
approved
has_operational_account
transaction_password_status
can_access_home
```

Decisões para o primeiro corte:

- `onboarding_completed` deve retornar `true`, porque hoje o customer é criado
  junto com o novo usuário no cadastro do sistema;
- o onboarding futuro será um fluxo independente;
- endereço fica fora do cadastro inicial e deve entrar no onboarding futuro,
  em tabela própria vinculada por `customer_id`;
- `approved` deve ser exposto para indicar o estado de aprovação do usuário ou
  da conta conforme a regra atual da API;
- `next_required_step` não precisa existir neste momento.

Valores sugeridos para `transaction_password_status`:

```text
active
not_set
locked
unknown
```

O backend deve ser a fonte da decisão de `can_access_home`. O mobile deve usar
esse campo e os status explícitos de readiness para roteamento, em vez de
recriar a policy completa com base em múltiplos sinais isolados.

## 4. Escopo de API

- Criar use case de sessão pós-login no módulo de auth ou em um módulo
  transversal apropriado.
- Reutilizar dados atualmente expostos por:
  - `GET /auth/me`;
  - `GET /customers/me`.
- Garantir que a resposta contenha todos os campos necessários para substituir
  a composição atual feita por `AuthApi.getProfile` no mobile.
- Consultar o estado da senha transacional no repositório de segurança.
- Consultar ou derivar o estado de aprovação para preencher
  `readiness.approved`.
- Retornar `readiness.onboarding_completed = true` no primeiro corte.
- Calcular `readiness.transaction_password_status`.
- Calcular `readiness.can_access_home`.
- Registrar a rota `GET /auth/session` com autenticação JWT.
- Manter o padrão de envelope `{ data, error }`.
- Atualizar documentação REST.
- Atualizar collection Postman se ela for mantida como artefato de contrato.

## 5. Regras iniciais de prontidão

Primeiro corte sugerido para usuário customer:

```text
onboarding_completed = true

se conta ainda não estiver aprovada:
  can_access_home = false

se não houver conta operacional:
  can_access_home = false

se não houver senha transacional ativa:
  can_access_home = false

caso contrário:
  can_access_home = true
```

Se o usuário autenticado tiver papel customer e não houver customer associado,
a API deve retornar erro de estado inválido da conta, como
`INVALID_USER_STATE` ou erro equivalente já padronizado.

As regras exatas podem ser ajustadas conforme o estado atual do domínio de
contas. O importante é que a decisão de `can_access_home` fique centralizada na
API.

## 6. Compatibilidade

`GET /auth/session` não precisa substituir imediatamente:

- `GET /auth/me`;
- `GET /customers/me`.

Esses endpoints podem continuar existindo para compatibilidade e usos
específicos.

No mobile, a migração pode ocorrer em etapa posterior:

1. API disponibiliza `GET /auth/session`.
2. Mobile passa a usar `GET /auth/session` no bootstrap pós-login.
3. `AuthApi.getProfile` deixa de compor `/customers/me` + `/auth/me`.
4. Endpoints antigos permanecem disponíveis até decisão explícita de
   depreciação.

## 7. Erros

Clientes devem depender de `error.code`.

Erros relevantes:

- `UNAUTHORIZED`: sessão ausente;
- `INVALID_TOKEN`: token inválido ou expirado;
- `FORBIDDEN`: usuário autenticado não está apto a consultar a sessão;
- `INVALID_USER_STATE`: vínculo inconsistente entre usuário e customer;
- `CUSTOMER_NOT_FOUND` ou erro equivalente quando o usuário customer existir,
  mas o customer não for encontrado;
- erro genérico mapeado para falhas inesperadas de infraestrutura.

Para usuário customer, ausência de customer associado deve ser erro. Para papéis
não-customer, a seção `customer` pode ser omitida ou retornar `null`, conforme
decisão específica de contrato para esses perfis.

O endpoint não deve retornar senha transacional, hash, pepper, token de step-up
ou qualquer material sensível.

## 8. Fora de escopo

- Remover `GET /auth/me`.
- Remover `GET /customers/me`.
- Criar senha transacional.
- Autorizar step-up.
- Executar transferência.
- Implementar recuperação/troca/reset de senha transacional.
- Resolver todos os gates futuros de KYC, termos ou dispositivo confiável além
  dos campos reservados no contrato.

## 9. Critérios de aceite

- `GET /auth/session` exige JWT válido.
- Resposta segue envelope `{ data, error }`.
- Resposta inclui seção `user`.
- Resposta inclui seção `customer` para usuário customer.
- Usuário customer sem customer associado retorna erro de estado inválido, não
  `customer: null`.
- Resposta contém todos os campos hoje usados pelo mobile para montar
  `UserProfile`.
- Resposta inclui seção `readiness`.
- `readiness.onboarding_completed` retorna `true` no primeiro corte.
- `readiness.approved` é exposto.
- `readiness.transaction_password_status` informa se a senha transacional está
  ativa ou ausente.
- `readiness.can_access_home` é calculado pela API.
- `readiness.next_required_step` não é obrigatório neste primeiro corte.
- O endpoint não retorna dados sensíveis de senha transacional ou step-up.
- Testes cobrem usuário com senha transacional ativa.
- Testes cobrem usuário sem senha transacional.
- Testes cobrem usuário customer sem customer associado retornando erro de
  estado inválido.
- Testes cobrem sessão inválida/ausente.
- Documentação REST é atualizada.

## 10. Referências

- `mobile/lib/data/services/apis/auth/auth_api.dart`
- `mobile/lib/domain/common/auth/models/user_profile.dart`
- `api/internal/auth/delivery/handler.go`
- `api/internal/customer/delivery/http.go`
- `api/internal/security/application/create_transaction_password.go`
- `docs/backlogs/mobile/011 - cadastro-senha-transacional.md`
