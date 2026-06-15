# Backlog API: Installation Identity MVP

## 1. Status

- Tipo: Research
- Área: Security
- Prioridade: High
- Estado: Discussão

Este backlog define as responsabilidades da API no segundo sinal contextual da
evolução Zero Trust do BankLab. A geração, persistência local e propagação do
identificador pertencem ao backlog mobile 013.

Ele ainda não autoriza implementação: vínculo, estados, revogação e rollout
devem ser fechados antes da criação das tasks da API.

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
   deve ser negado ou direcionado a um fluxo de recuperação ainda a definir.
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
    "max_installations": 3
  },
  "error": null
}
```

O contrato final ainda deve decidir se a autorização restrita poderá acessar o
gerenciamento de instalações para liberar uma vaga.

### Autorização restrita para nova instalação

A senha transacional não deve ser enviada para `POST /auth/login`. Quando as
credenciais estão corretas, mas a instalação é nova, o login retorna:

```json
{
  "data": {
    "authentication_status": "installation_registration_required",
    "access_token": "<restricted-access-token>",
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

O mecanismo final para representar e aplicar a restrição no token, persistência
e middleware ainda precisa ser definido.

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

Futuramente, prova de vida poderá autorizar o mesmo endpoint sem alterar a
responsabilidade do registro.

O `X-Installation-Id` também deve acompanhar:

- `POST /auth/refresh`;
- logout, quando exposto pela API;
- todas as requisições autenticadas por access token.

Para requisições autenticadas, o identificador declarado deve corresponder à
instalação vinculada à sessão. A API não deve reassociar silenciosamente uma
sessão existente quando receber outro identificador; a troca exige novo login.

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

Embora o contrato HTTP use `DELETE`, a revogação é lógica:

- a linha não é removida fisicamente;
- o status passa para `revoked`;
- `revoked_at` registra o momento da revogação;
- a associação permanece disponível para auditoria e histórico;
- uma instalação revogada não volta a `known` por novo login.

Questões de contrato ainda abertas:

- Quais metadados da instalação são necessários no MVP além de plataforma e
  versão do app?
- Como representar e aplicar a autorização restrita no JWT, persistência e
  middleware?
- A autorização restrita pode listar e revogar instalações para liberar uma
  vaga?
- `DELETE /security/installations/{installation_resource_id}` também exige
  step-up?
- Como registrar uma nova instalação quando a senha transacional estiver
  `not_set` ou `locked`?
- Como recuperar acesso quando todas as instalações anteriores estiverem
  indisponíveis ou revogadas?
- A ausência do header bloqueia login, refresh e requisições autenticadas desde
  o primeiro release ou somente após o período de rollout?
- Qual erro deve ser retornado quando o header divergir da instalação vinculada
  à sessão?

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

Precisamos decidir se a revogação de uma instalação:

- encerra imediatamente todas as sessões associadas;
- impede apenas novas operações sensíveis;
- exige novo step-up em outra instalação para recuperação.

## 8. Relação com sessão e policy enforcement

O bootstrap de sessão deve ser considerado o ponto de retorno do estado da
instalação, mas o JWT não deve carregar uma confiança imutável durante toda a
sua validade.

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

## 9. Decisões necessárias antes das tasks

- [x] Definir header e formato aceito: `X-Installation-Id` com UUID v4
  canônico.
- [ ] Definir erro retornado para identificador malformado.
- [ ] Definir vínculo entre usuário, instalação e sessão.
- [ ] Definir estados e transições da instalação.
- [x] Definir bootstrap: a primeira instalação do usuário é associada
  atomicamente durante o login.
- [x] Definir histórico: instalações revogadas não são removidas e continuam
  contando como associações anteriores.
- [x] Definir limite: no máximo três instalações `known` por usuário;
  revogadas não ocupam vaga.
- [ ] Definir contrato HTTP final para `installation_limit_reached`.
- [ ] Definir acesso da autorização restrita ao gerenciamento para liberar uma
  vaga.
- [x] Definir novas instalações: login emite autorização restrita; senha
  transacional emite step-up para `POST /security/installations`; o endpoint
  executa o registro.
- [x] Definir endpoints de registro, listagem e revogação.
- [ ] Definir representação, TTL, persistência e middleware da autorização
  restrita.
- [x] Definir `POST /security/installations` como operação elegível para
  step-up.
- [ ] Implementar a operação na whitelist de step-up.
- [ ] Definir recuperação para senha transacional ausente ou bloqueada.
- [ ] Definir validação do vínculo em refresh e requisições autenticadas.
- [ ] Definir efeitos da revogação sobre refresh e access tokens.
- [ ] Definir dados retornados no bootstrap da sessão.
- [ ] Definir retenção, auditoria e minimização de metadados.
- [ ] Definir estratégia de compatibilidade para clientes sem o header.
- [ ] Definir em qual etapa o sinal passa a afetar operações sensíveis.

## 10. Fora de escopo

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

## 11. Critérios para encerrar a discussão

- Modelo de ameaça e limitações do sinal documentados.
- Contrato HTTP definido entre API e mobile.
- Ciclo de vida e revogação definidos.
- Limite e comportamento concorrente definidos.
- Relação com `user_sessions` definida.
- Estratégia de rollout sem bloqueio definida.
- Backlogs de implementação da API e do mobile derivados desta decisão.

## 12. Referências internas

- [ZTA MVP - fundação e decisões](<done/006 - zta-mvp-foundation.md>)
- [Auth session bootstrap](<done/009 - auth-session-bootstrap.md>)
- [Installation Identity MVP mobile](<../mobile/013 - installation-identity-mvp.md>)
- [Discussão de segurança transacional](../discussion.md)
- [Roadmap](../../ROADMAP.md)
