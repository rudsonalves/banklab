# Autenticação e Contexto do Usuário

Este capítulo descreve como a API identifica o usuário autenticado e como essa identidade é propagada para os demais fluxos do sistema. O objetivo é mostrar que autenticação, nesta aplicação, não se resume à emissão de tokens. Ela também estabelece o contexto a partir do qual os módulos de negócio tomam decisões de ownership, autorização e execução.

## Papel da autenticação na aplicação

A autenticação é o mecanismo que permite à API responder à pergunta:

```text
quem está fazendo esta requisição?
```

Na aplicação atual, essa resposta não é importante apenas para permitir acesso a rotas protegidas. Ela também é usada para:

- determinar a identidade do usuário;
- identificar o `customer_id` vinculado a esse usuário;
- conhecer o `role` associado à sessão;
- aplicar regras de ownership;
- autorizar operações administrativas;
- evitar que o cliente informe dados sensíveis de escopo no payload.

Portanto, a autenticação participa diretamente do modelo de segurança e da modelagem do domínio.

## Fluxo geral de autenticação

O fluxo de autenticação da API possui três momentos principais:

1. entrada inicial na aplicação;
2. estabelecimento da sessão autenticada;
3. reutilização do contexto autenticado nas rotas protegidas.

### Entrada inicial

O processo começa normalmente em:

```text
POST /auth/register
POST /auth/login
```

No registro, o sistema cria um usuário e o customer associado.

No login, a API valida credenciais e retorna os tokens necessários para o uso posterior das rotas protegidas.

### Sessão autenticada

Depois do login, o cliente recebe:

- `access_token`;
- `refresh_token`;
- `user_id`;
- `email`;
- `role`;
- `customer_id`.

O `access_token` é o JWT usado nas chamadas autenticadas.

O `refresh_token` é usado para renovação da sessão sem exigir novo login.

### Reutilização do contexto

Nas chamadas seguintes, o cliente envia o JWT no header:

```http
Authorization: Bearer <access_token>
```

O backend valida esse token e, a partir dele, reconstrói o contexto autenticado da requisição.

Esse contexto passa então a orientar a execução de operações nos módulos `auth`, `customer`, `account` e `admin`.

## App Token e JWT

No estágio atual, a API utiliza dois mecanismos distintos de proteção em momentos diferentes do fluxo.

### App Token

As rotas de entrada inicial utilizam App Token:

```text
POST /auth/register
POST /auth/login
```

Esse mecanismo protege o acesso inicial à API mesmo antes de existir uma sessão autenticada do usuário.

### JWT

As rotas protegidas utilizam JWT Bearer Token.

Esse token representa a identidade autenticada do usuário e é exigido em operações como:

- `POST /auth/refresh`;
- `GET /auth/me`;
- `GET /customers/me`;
- `GET /accounts`;
- `POST /accounts`;
- `GET /accounts/{id}/balance`;
- `POST /accounts/{id}/deposit`;
- `POST /accounts/{id}/withdraw`;
- `GET /accounts/internal-transfers/recipients`;
- `POST /accounts/internal-transfers`;
- `GET /accounts/{id}/statement`;
- `POST /admin/users/{id}/approve`.

## Conteúdo do contexto autenticado

Depois da validação do JWT, a requisição passa a carregar um contexto autenticado interno.

Esse contexto inclui informações como:

- `user_id`;
- `email`;
- `role`;
- `customer_id`.

Esses dados são importantes porque os módulos da aplicação passam a operar com base neles.

O contexto autenticado permite, por exemplo:

- identificar qual usuário está consultando dados;
- saber se o usuário possui permissão administrativa;
- derivar a quem pertencem as contas do fluxo atual;
- decidir se a operação pode ser executada sem depender de identificadores enviados pelo cliente.

## Contexto autenticado como fonte de verdade

Uma das decisões arquiteturais mais importantes da API é tratar o contexto autenticado como fonte principal de identidade para as operações protegidas.

Isso significa que, em vários fluxos, o backend não confia em dados sensíveis enviados explicitamente pelo cliente para determinar o escopo da operação.

Por exemplo:

- `GET /accounts` não recebe `customer_id`;
- `POST /accounts` não recebe `customer_id`;
- `GET /customers/me` não recebe `customer_id`.

Em vez disso, a API deriva esse valor a partir do usuário autenticado.

Essa escolha reduz o risco de:

- manipulação indevida de escopo pelo cliente;
- erro de integração no frontend;
- confusão entre identidade autenticada e entidade de negócio.

Também reforça o papel do backend como guardião da autorização e do ownership.

## Relação entre usuário, role e customer

O contexto autenticado conecta três dimensões importantes da aplicação:

### Usuário

O usuário representa a identidade autenticada no sistema. Ele é a entidade associada à credencial de acesso, ao login e à sessão.

### Role

O `role` representa o nível de permissão do usuário.

Ele é usado para responder perguntas como:

- este usuário pode executar operação administrativa?
- esta rota exige privilégio elevado?
- o chamador é apenas um usuário customer ou possui papel administrativo?

### Customer

O `customer_id` conecta a identidade autenticada ao domínio de negócio.

Esse vínculo é importante porque muitas operações não dependem apenas de “quem está autenticado”, mas também de “qual cliente de negócio está associado àquela identidade”.

No módulo `account`, por exemplo, esse dado é central para:

- listar contas do usuário;
- criar contas vinculadas ao customer correto;
- validar ownership sobre recursos financeiros.

## Autenticação e autorização não são a mesma coisa

Embora estejam relacionadas, autenticação e autorização cumprem papéis diferentes.

Autenticação responde:

```text
quem é o usuário?
```

Autorização responde:

```text
esse usuário pode executar esta ação?
```

Na API, a autenticação é feita por App Token e JWT. A autorização é aplicada nos fluxos da aplicação a partir do contexto autenticado.

Exemplos:

- um usuário pode estar autenticado, mas não ter `role` de administrador;
- um usuário pode estar autenticado, mas não ser dono da conta consultada;
- um usuário autenticado pode não possuir `customer_id` suficiente para determinado fluxo;
- uma operação pode exigir não apenas autenticação, mas também status válido ou ownership.

## Middleware e propagação do contexto

O contexto autenticado não é montado manualmente em cada handler. Ele é preparado pelo mecanismo de autenticação aplicado às rotas protegidas.

Na prática, o middleware:

- lê o token;
- valida sua assinatura e integridade;
- extrai os dados da identidade;
- injeta essas informações no contexto da requisição.

Depois disso, handlers e casos de uso podem consumir esse contexto sem precisar decodificar o token repetidamente.

Essa separação é importante porque:

- a validação técnica do token fica centralizada;
- os handlers permanecem mais limpos;
- os casos de uso recebem contexto já pronto para decidir ownership e autorização.

## O endpoint `GET /auth/me`

O endpoint:

```text
GET /auth/me
```

tem um papel especial nesse desenho.

Ele funciona como uma forma explícita de consultar o contexto autenticado atual. Isso é útil para:

- confirmação de sessão no cliente;
- reidratação de estado no mobile;
- validação de identidade após login ou refresh;
- leitura de `role` e `customer_id`.

Esse endpoint não cria identidade nem muda estado de sessão. Ele apenas retorna, em formato HTTP, a identidade autenticada atualmente reconhecida pela API.

## Refresh token e continuidade da sessão

O `refresh_token` permite renovar o `access_token` sem exigir novo login.

No projeto atual, o fluxo de refresh é protegido por rotação de token. Isso significa que:

- o refresh token apresentado é validado;
- a sessão correspondente é verificada;
- o token antigo é revogado;
- um novo par de credenciais é emitido.

Esse comportamento reforça o controle do backend sobre a continuidade da sessão e reduz a chance de reutilização silenciosa de tokens antigos.

Para o restante da aplicação, isso significa que o contexto autenticado continua sendo renovável sem exigir nova autenticação de credenciais a cada expiração do access token.

## Impacto direto no desenho das features

Entender o contexto autenticado é essencial para implementar novas funcionalidades corretamente.

Quando uma nova feature depende da identidade do usuário, o desenvolvedor deve primeiro perguntar:

- essa operação depende do usuário autenticado?
- ela depende de `role`?
- ela depende de `customer_id`?
- o cliente precisa realmente enviar esse identificador, ou o backend já o conhece pelo contexto?

Na arquitetura atual, a resposta tende a favorecer o contexto autenticado como origem dessas informações sempre que possível.

## Síntese

A autenticação na API estabelece mais do que acesso: ela estabelece contexto.

Por meio de App Token, JWT, middleware e propagação do usuário autenticado, o sistema passa a conhecer identidade, role e vínculo com customer. Esse contexto é então reutilizado para autorizar operações, determinar ownership e evitar dependência de dados sensíveis enviados pelo cliente.

Para novos desenvolvedores, esse é um dos conceitos mais importantes da aplicação. Muitas decisões de arquitetura, segurança e modelagem só fazem sentido quando se entende que o contexto autenticado é uma peça central no comportamento dos módulos de negócio.
