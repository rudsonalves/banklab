# Backlog API: Installation Identity MVP

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

Este backlog define as responsabilidades da API no segundo sinal contextual da
evolução Zero Trust do BankLab. A geração, persistência local e propagação do
identificador pertencem ao backlog mobile 013.

A implementação deve respeitar a ordem de dependência definida neste documento:
primeiro a base compartilhada, depois os fluxos que consultam, registram,
revogam ou aplicam instalação em sessão.

## 2. Contexto

O primeiro incremento ZTA protege a transferência interna com senha
transacional e step-up token de uso único. Esse fluxo confirma intenção, mas
trata sessões válidas de ambientes diferentes da mesma forma.

O próximo incremento deve permitir que a API reconheça uma instalação do app
sem afirmar que reconhece o aparelho físico. O identificador da instalação não
é segredo, fator forte ou prova de posse.

O contrato compartilhado usa `X-Installation-Id`. O ciclo de vida desse valor
no cliente está documentado no backlog mobile correspondente.

## 3. Objetivo

Definir a parte da API para um MVP de identidade de instalação que:

- valide o identificador recebido em `X-Installation-Id`;
- associe instalações a usuários e sessões;
- permita consultar e revogar instalações;
- disponibilize a instalação como sinal contextual para políticas futuras;
- preserve o backend como fonte de verdade sobre o estado da associação.

Este backlog volta a ser a fonte única da API para este MVP. A implementação
deve seguir a ordem de dependência descrita aqui, começando pela base
compartilhada antes dos fluxos que dependem dela.

## 4. Princípios de segurança

- O identificador não é autenticação e não substitui JWT ou step-up.
- O identificador representa a instalação do app, não o aparelho físico.
- A API não deve inferir identidade física a partir desse valor.
- A API decide se uma instalação é conhecida, permitida ou revogada.
- O header pode ser copiado por um cliente comprometido e deve ser tratado como
  sinal fraco.
- Logs não devem expor tokens, senha transacional ou atributos excessivos do
  ambiente.
- Revogar uma instalação deve ter efeito definido sobre sessões existentes.
- A política deve distinguir instalação conhecida, nova, revogada e ausente.
- Reinstalar o app ou limpar seus dados gera uma nova identidade.

## 5. Modelo inicial para discussão

```text
app_installations
- id
- user_id
- installation_id
- status
- platform
- app_version
- app_build
- first_seen_at
- last_seen_at
- revoked_at
- created_at
- updated_at
```

Status candidatos:

```text
known
revoked
```

Esses estados descrevem apenas o cadastro:

- `known`: a instalação está associada ao usuário e pode participar da
  avaliação contextual;
- `revoked`: a associação foi revogada e não deve ser aceita como instalação
  válida para novas sessões ou operações.

`trusted` não é um terceiro status deste MVP. O fato de uma instalação ser
conhecida não prova posse do aparelho nem torna a requisição confiável.

Uma eventual confiança de dispositivo ou instalação deve ser modelada
separadamente, somente depois de definir:

- qual evidência concede confiança;
- quem ou qual fluxo pode concedê-la;
- por quanto tempo ela permanece válida;
- quais eventos removem ou reduzem essa confiança.

## 6. Contrato inicial

Decisão adotada para o MVP:

```http
X-Installation-Id: <UUID v4>
```

O valor identifica somente a instalação atual do app. Deve ser gerado
aleatoriamente pelo mobile, persistido entre execuções e atualizações e
substituído após reinstalação ou limpeza dos dados. O formato canônico usa
letras minúsculas e hífens.

Responsabilidades da API:

- aceitar apenas UUID v4 em formato canônico;
- rejeitar valores malformados com erro de contrato estável;
- não usar o header como autenticação ou autorização isolada;
- associar a instalação somente dentro de um contexto autenticado;
- não expor `installation_id` como identificador público de revogação.

O login decide se a instalação pode ser associada automaticamente ou se o
cliente precisa seguir o endpoint explícito de registro:

```http
POST /auth/login
X-Installation-Id: <UUID v4>
```

Após autenticar as credenciais, a API segue uma destas regras:

1. **Instalação já conhecida:** cria uma sessão vinculada à associação
   existente e retorna os tokens normalmente.
2. **Primeira instalação do usuário:** se o usuário nunca teve uma instalação
   associada, cria a primeira associação, vincula a sessão e retorna os tokens.
3. **Nova instalação com outra já cadastrada:** emite uma autorização de acesso
   restrito para permitir somente o step-up e o registro da instalação.
4. **Instalação revogada:** não reativa nem reassocia silenciosamente. O acesso
   deve ser negado.
5. **Limite atingido:** se o usuário já possuir três instalações `known`, a API
   não permite o registro de uma quarta instalação. Uma associação existente
   precisa ser revogada antes de continuar.

A primeira instalação é um evento cronológico, não uma condição de existir
apenas uma instalação ativa. Depois desse bootstrap, o usuário pode possuir
várias instalações conhecidas.

A verificação de que o usuário nunca teve uma instalação associada e a criação
da primeira associação devem ser atômicas. Dois logins concorrentes não podem
cadastrar duas instalações como primeiro bootstrap.

Associações revogadas permanecem no histórico e continuam contando para
determinar que o usuário já teve sua primeira instalação. Portanto, revogar
todas as instalações não restaura a elegibilidade para o bootstrap automático.

### Limite de instalações

O MVP permite no máximo:

```text
3 instalações known por usuário
```

O limite considera somente associações com status `known`:

- instalações revogadas permanecem no histórico, mas não ocupam vaga;
- sessões diferentes da mesma instalação não criam novas vagas;
- registrar uma nova instalação deve reservar a vaga de forma atômica;
- registros concorrentes não podem ultrapassar o limite;
- uma quarta instalação só pode ser registrada após a revogação de uma das três
  instalações conhecidas.

Quando o limite já estiver atingido, `POST /auth/login` deve retornar:

```json
{
  "data": {
    "authentication_status": "installation_limit_reached",
    "max_installations": 3,
    "known_installations_count": 3,
    "installation_registered": false,
    "next_action": "revoke_existing_installation"
  },
  "error": null
}
```

Nesse cenário, a API não deve emitir autorização restrita nem abrir fluxo de
registro para a nova instalação. O login apenas informa que:

- a instalação atual ainda não está cadastrada;
- o usuário já atingiu o máximo de três instalações `known`.

A continuidade fica a cargo do mobile: orientar o usuário a acessar uma
instalação já cadastrada e, a partir dela, revogar outra instalação para
liberar vaga.

### Autorização restrita para nova instalação

A senha transacional não deve ser enviada para `POST /auth/login`. Quando as
credenciais estão corretas, mas a instalação é nova, o login retorna:

```json
{
  "data": {
    "authentication_status": "installation_registration_required",
    "restricted_access_token": "<jwt>",
    "expires_in": 300
  },
  "error": null
}
```

Essa autorização:

- identifica o usuário cujas credenciais já foram validadas;
- não representa uma sessão operacional completa;
- não emite refresh token;
- permite somente as rotas necessárias ao registro da instalação;
- não permite acesso a contas, saldo, extrato, transferências ou demais rotas
  autenticadas;
- expira sem criar associação caso o fluxo seja abandonado.

O MVP usará uma tabela própria, consistente com o desenho já adotado para
`step_up_tokens`, sem reutilizar `user_sessions`.

Claims mínimas do `restricted_access_token`:

```text
sub
jti
token_type = restricted_access
scope = installation.register
installation_id
iat
exp
```

Persistência sugerida:

```text
installation_registration_authorizations
- id
- jti
- user_id
- installation_id
- scope
- status
- expires_at
- consumed_at
- created_at
```

Status mínimos:

```text
active
consumed
revoked
```

Regras de consistência:

- `jti` deve ser único;
- deve existir no máximo uma autorização `active` por
  `(user_id, installation_id, scope)`;
- `expired` pode ser tratado como estado derivado de `expires_at`, sem exigir
  transição persistida;
- a autorização precisa ser invalidada ao concluir com sucesso o registro.

Aplicação no middleware:

- manter o middleware atual para sessão operacional;
- adicionar um middleware específico para acesso restrito;
- esse middleware deve validar o JWT, exigir
  `token_type = restricted_access` e o `scope` esperado;
- o contexto da requisição deve carregar ao menos `user_id`,
  `installation_id`, `jti` e `scope`;
- `POST /security/step-up/authorize` e `POST /security/installations` devem
  aceitar esse contexto restrito sem serem tratados como sessão operacional
  completa.

### Step-up para o registro

A senha transacional mantém sua responsabilidade atual: validar intenção e
emitir autorização curta para uma operação pública. Ela não registra a
instalação.

```http
POST /security/step-up/authorize
Authorization: Bearer <restricted-access-token>
```

```json
{
  "method": "POST",
  "path": "/security/installations",
  "transaction_password": "123456"
}
```

A resposta contém o `step_up_token` já adotado pelo projeto, curto, de uso
único e escopado para `POST /security/installations`.

Senha transacional ativa é pré-requisito para registrar uma nova instalação.
Se a senha estiver `not_set` ou `locked`, a API não deve emitir `step_up_token`
para `POST /security/installations` e o registro não pode prosseguir nesse
momento.

### Registro explícito

Instalações posteriores à primeira são cadastradas por:

```http
POST /security/installations
Authorization: Bearer <restricted-access-token>
X-Step-Up-Token: <step_up_token>
X-Installation-Id: <UUID v4>
```

O endpoint:

1. valida a autorização restrita;
2. valida e consome o step-up emitido para essa operação;
3. confirma que o `X-Installation-Id` corresponde ao apresentado no login;
4. confirma atomicamente que ainda existe vaga no limite de três;
5. cria a associação como `known`;
6. invalida a autorização restrita;
7. cria a sessão operacional vinculada à instalação;
8. retorna access token e refresh token operacionais.

Após o cadastro bem-sucedido, não existe novo login intermediário. O
`restricted_access_token` é invalidado e descartado, e a autenticação passa a
ser a sessão operacional normal.

O `access_token` operacional deve carregar o `installation_id` associado à
sessão. Esse claim permite que middleware e casos de uso validem o vínculo
entre:

- usuário autenticado;
- instalação declarada no `X-Installation-Id`;
- instalação vinculada à sessão criada após o cadastro.

O `refresh_token` permanece como token opaco de sessão, mas o refresh deve
revalidar o mesmo `X-Installation-Id` da sessão em curso antes de emitir novos
tokens.

Futuramente, prova de vida poderá autorizar o mesmo endpoint sem alterar a
responsabilidade do registro.

O `X-Installation-Id` também deve acompanhar:

- `POST /auth/refresh`;
- logout, quando exposto pela API;
- todas as requisições autenticadas por access token.

Para requisições autenticadas, o identificador declarado deve corresponder à
instalação vinculada à sessão e ao `installation_id` presente no
`access_token`. A API não deve reassociar silenciosamente uma sessão existente
quando receber outro identificador; a troca exige novo login.

Endpoints consolidados:

```http
POST   /security/installations
GET    /security/installations
DELETE /security/installations/{installation_resource_id}
```

- `POST` cadastra uma nova instalação após step-up;
- `GET` lista as instalações associadas ao usuário autenticado;
- `DELETE` revoga uma associação usando um identificador público gerado pela
  API, nunca o `installation_id` fornecido pelo cliente.

Para o MVP, `DELETE /security/installations/{installation_resource_id}` não
exige step-up. A operação só pode revogar outra instalação do mesmo usuário:

- a API deve impedir a revogação da instalação vinculada à sessão atual;
- a resposta de erro deve deixar claro que a instalação em uso não pode ser
  removida por esse fluxo;
- sair da instalação atual continua sendo responsabilidade do logout, não da
  revogação.

Embora o contrato HTTP use `DELETE`, a revogação é lógica:

- a linha não é removida fisicamente;
- o status passa para `revoked`;
- `revoked_at` registra o momento da revogação;
- a associação permanece disponível para auditoria e histórico;
- uma instalação revogada não volta a `known` por novo login.

Erros já definidos para o header e o vínculo da sessão:

- `X-Installation-Id` ausente ou malformado, quando obrigatório:

```http
400 Bad Request
```

```json
{
  "data": null,
  "error": {
    "code": "INVALID_INSTALLATION_ID",
    "message": "X-Installation-Id must be a canonical UUID v4."
  }
}
```

- `X-Installation-Id` válido, mas divergente da instalação vinculada à sessão:

```http
403 Forbidden
```

```json
{
  "data": null,
  "error": {
    "code": "INSTALLATION_MISMATCH",
    "message": "The provided installation does not match the authenticated session."
  }
}
```

Para o MVP, não existe fase de rollout tolerante sem header. O
`X-Installation-Id` é obrigatório desde o primeiro release da feature em:

- `POST /auth/login`;
- `POST /auth/refresh`;
- logout, quando exposto pela API;
- todas as requisições autenticadas por `access_token`.

## 7. Fluxos a definir

### Conta nova: primeira instalação

```text
Mobile gera installation_id antes do login
  -> envia credenciais + X-Installation-Id
  -> API valida as credenciais
  -> API confirma que a conta nunca teve instalação associada
  -> cadastra silenciosamente a primeira instalação
  -> cria a sessão vinculada à instalação
  -> retorna access e refresh tokens
```

O cadastro silencioso acontece como parte da conclusão do primeiro login. Não
existe tela ou confirmação adicional para a primeira instalação da conta.

### Conta existente: nova instalação

```text
Mobile gera installation_id antes do login
  -> envia credenciais + X-Installation-Id
  -> API valida as credenciais
  -> API identifica uma instalação ainda não cadastrada
  -> API verifica que existe vaga no limite de 3 instalações known
  -> emite access token restrito, sem refresh token
  -> mobile solicita step-up para POST /security/installations
  -> senha transacional autoriza o endpoint
  -> mobile recebe step_up_token
  -> chama POST /security/installations
  -> API valida e consome o step-up
  -> cadastra a instalação como known
  -> invalida a autorização restrita
  -> cria sessão operacional vinculada
  -> retorna access e refresh tokens operacionais
```

Nesse fluxo, “após login” significa que as credenciais primárias foram
validadas. O token restrito existe apenas para concluir o registro e não libera
o restante da API.

### Limite atingido

```text
API recebe X-Installation-Id desconhecido
  -> autentica usuário
  -> encontra 3 instalações known
  -> retorna installation_limit_reached
  -> usuário precisa revogar uma instalação antes de continuar
```

### Revogação

Para o MVP, revogar uma instalação paralisa esse acesso imediatamente:

- todos os `access_token` associados à instalação revogada deixam de valer na
  hora;
- todos os `refresh_token` associados à instalação revogada devem ser
  invalidados no mesmo momento;
- a instalação revogada não pode concluir novas operações nem renovar sessão;
- a API deve tratar a revogação como corte imediato de acesso, não como bloqueio
  apenas em operações futuras.

## 8. Relação com sessão e policy enforcement

O bootstrap de sessão deve ser considerado o ponto de retorno do estado da
instalação, mas o JWT não deve carregar uma confiança imutável durante toda a
sua validade.

Para o MVP:

- o `access_token` operacional carrega `installation_id`;
- o middleware de sessão operacional deve extrair esse claim e compará-lo com o
  `X-Installation-Id` da requisição;
- `POST /auth/refresh` deve exigir o mesmo `X-Installation-Id` vinculado à
  sessão atual antes de renovar os tokens;
- divergência entre header e claim deve negar a operação e exigir novo login;
- a API não deve usar esse claim para inferir confiança de dispositivo, apenas
  vínculo de instalação.

Uma política futura poderá avaliar:

```text
Evaluate(user, session, installation, operation)
  -> allow
  -> require_step_up
  -> deny
```

No primeiro corte, a instalação deve ser coletada e registrada antes de ser
obrigatória para transferências. Isso permite observar o contrato e os estados
sem introduzir bloqueios prematuros no fluxo financeiro.

## 9. Ordem de implementação por dependência

A separação operacional deste MVP fica em um documento próprio para evitar que
o backlog principal misture decisão de produto, modelo de ameaça e divisão de
execução.

Ver:

- [010 - split-installation-identity-by-dependency.md](<010 - split-installation-identity-by-dependency.md>)

## 10. Estado atual de implementação

O trabalho deve ser retomado considerando que apenas o contrato mínimo de
entrada pode avançar antes da infraestrutura compartilhada.

Estado esperado neste ponto:

- `X-Installation-Id` validado no login;
- erro `INVALID_INSTALLATION_ID` definido;
- constantes compartilhadas de header disponíveis;
- demais fluxos dependentes aguardando base de dados, domínio, repositórios,
  sessão e JWT com vínculo de instalação.

Não considerar como concluídos, até que a base exista:

- bootstrap automático da primeira instalação;
- classificação operacional de instalação no login;
- emissão de `restricted_access_token`;
- registro explícito de instalação;
- listagem e revogação de instalações;
- refresh e middleware vinculados à instalação.

## 11. Backlogs derivados

Os backlogs derivados devem seguir a separação por dependência técnica definida
no documento de split.

- [010 - split-installation-identity-by-dependency.md](<010 - split-installation-identity-by-dependency.md>):
  separação dos backlogs por ordem de dependência técnica.

## 12. Fora de escopo

- biometria local como prova para o backend;
- attestation de plataforma;
- identificação do aparelho físico;
- confiança atribuída ao aparelho;
- device fingerprinting;
- prova de vida;
- geolocalização;
- score antifraude;
- notificações push para aprovação;
- confiança permanente concedida somente pelo identificador;
- correlação entre reinstalações no mesmo aparelho;
- painel administrativo de instalações.

## 13. Critérios para encerrar a discussão

- Modelo de ameaça e limitações do sinal documentados.
- Contrato HTTP definido entre API e mobile.
- Ciclo de vida e revogação definidos.
- Limite e comportamento concorrente definidos.
- Relação com `user_sessions` definida.
- Estratégia de obrigatoriedade do header definida.
- Ordem de implementação por dependência definida.
- Relação com o backlog mobile definida.

## 14. Referências internas

- [ZTA MVP - fundação e decisões](<done/006 - zta-mvp-foundation.md>)
- [Auth session bootstrap](<done/009 - auth-session-bootstrap.md>)
- [Installation Identity MVP mobile](<../mobile/013 - installation-identity-mvp.md>)
- [Discussão de segurança transacional](../discussion.md)
- [Roadmap](../../ROADMAP.md)
