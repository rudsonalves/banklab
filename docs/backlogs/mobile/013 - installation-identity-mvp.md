# Backlog Mobile: Installation Identity MVP

## 1. Status

- Tipo: Research
- Área: Security
- Prioridade: High
- Estado: Discussão

Este backlog define as responsabilidades do aplicativo mobile para criar e
propagar uma identidade aleatória por instalação. Associação com usuário,
sessão, estados e revogação pertencem ao backlog API 010.

Ele ainda não autoriza implementação: as decisões de armazenamento,
reinstalação, concorrência e falha precisam ser fechadas antes da criação das
tasks mobile.

## 2. Objetivo

O mobile deve:

- gerar um UUID v4 aleatório para a instalação atual;
- persistir o valor entre execuções e atualizações do app;
- enviar o valor em `X-Installation-Id`;
- preservar o identificador em logout e troca de usuário;
- gerar uma nova identidade após reinstalação ou limpeza dos dados;
- não bloquear o uso básico do app caso o armazenamento falhe durante o
  rollout inicial;
- não tratar o identificador como segredo, autenticação ou prova de posse.

## 3. Contrato compartilhado

```http
X-Installation-Id: <UUID v4>
```

O formato canônico usa letras minúsculas e hífens:

```text
550e8400-e29b-41d4-a716-446655440000
```

O header deve ser enviado no login, no registro de nova instalação, refresh,
logout e em todas as chamadas autenticadas por access token.

Para o MVP, não existe fase tolerante sem esse header. O cliente deve adotar o
envio de `X-Installation-Id` desde o primeiro release da feature.

No primeiro login do usuário que nunca teve instalação associada, a API
cadastra automaticamente a primeira instalação. Essa regra representa o
primeiro evento cronológico e não limita o usuário a uma única instalação.

Quando já existir ou tiver existido outra instalação, inclusive revogada, o
login pode retornar
`installation_registration_required` com access token restrito e sem refresh
token. Nesse caso, o mobile deve solicitar step-up e chamar o endpoint de
registro da instalação.

O MVP permite até três instalações conhecidas por usuário. Quando o login
retornar `installation_limit_reached`, o mobile não deve solicitar a senha
transacional, pois não existe vaga disponível para confirmar a instalação.

O identificador representa somente a instalação atual do app. Ele não
identifica o aparelho físico e não deve ser derivado de IMEI, serial,
advertising ID, fingerprint ou atributos semelhantes.

## 4. Ciclo de vida

### Primeira execução

```text
App procura installation_id
  -> se não existir, gera UUID v4
  -> persiste localmente
  -> reutiliza o mesmo valor nas requisições seguintes
```

### Atualização do app

Uma atualização normal deve preservar o identificador.

### Logout e troca de usuário

O identificador pertence à instalação, não à sessão ou ao usuário. Portanto:

- logout não apaga o valor;
- limpeza das credenciais não apaga o valor;
- outro usuário no mesmo app reutiliza a identidade da instalação.

### Reinstalação ou limpeza de dados

Uma nova instalação ou limpeza dos dados deve gerar uma nova identidade. O app
não deve tentar reconstruir silenciosamente a identidade anterior por
fingerprint físico.

## 5. Persistência

Chave candidata:

```text
banklab.installation.id
```

O projeto já possui `LocalSecureStorage` sobre `flutter_secure_storage`. A
implementação deve decidir como distinguir atualização de reinstalação,
considerando que mecanismos como o Keychain do iOS podem preservar valores
após a remoção do app.

A identidade da instalação deve ter ciclo de vida separado dos tokens. O
`deleteAll()` usado para limpar sessão não pode apagar esse valor sem uma
decisão explícita.

## 6. Serviço e interceptor

O mobile deve possuir um serviço responsável por:

- executar `read-or-create` do UUID;
- evitar geração concorrente de valores diferentes durante o bootstrap;
- validar localmente o valor recuperado antes de reutilizá-lo;
- disponibilizar o identificador ao transporte HTTP;
- nunca registrar o valor completo em logs.

O esqueleto atual de `DeviceInterceptor` é apenas referência. A implementação
deve adotar terminologia de instalação, por exemplo:

```text
InstallationIdentityService
InstallationInterceptor
```

O interceptor deve adicionar somente:

```http
X-Installation-Id: <UUID v4>
```

Metadados como plataforma e versão do app devem seguir contrato separado, caso
a API decida solicitá-los. Para o MVP, o conjunto mínimo aceito para esse
contrato separado é `platform`, `app_version` e `app_build`. Eles não devem
ser codificados dentro do UUID.

## 7. Comportamento em falhas

No rollout inicial:

- falha de leitura ou escrita não deve gerar um UUID novo por requisição;
- uma tentativa de geração deve ser compartilhada entre requisições
  concorrentes;
- se não houver identidade estável disponível, a requisição segue sem o header;
- o erro deve ser observável sem registrar o identificador;
- a ausência do header não deve encerrar a sessão por decisão exclusiva do
  mobile.

O comportamento futuro para operações sensíveis será definido pela API.

## 8. Registro de nova instalação

### Conta nova

```text
Mobile gera e persiste installation_id
  -> envia credenciais + X-Installation-Id
  -> API valida o primeiro login
  -> API cadastra silenciosamente a primeira instalação
  -> mobile recebe e persiste os tokens da sessão
```

Não existe interação adicional para cadastrar a primeira instalação da conta.

### Conta existente em nova instalação

```text
Mobile gera e persiste installation_id
  -> envia credenciais + X-Installation-Id
  -> API valida as credenciais
  -> mobile recebe installation_registration_required
  -> mantém restricted_access_token durante o fluxo
  -> solicita step-up para POST /security/installations
  -> usuário informa senha transacional
  -> mobile recebe step_up_token
  -> chama POST /security/installations
  -> API cadastra a instalação e cria a sessão operacional
  -> mobile substitui o token restrito pelos tokens operacionais
```

Quando as credenciais forem válidas, mas a nova instalação precisar ser
registrada, o mobile deve:

1. manter o `restricted_access_token` separado da sessão operacional;
2. solicitar step-up para `POST /security/installations`;
3. informar a senha transacional somente em
   `POST /security/step-up/authorize`;
4. enviar o mesmo `X-Installation-Id` usado no login;
5. chamar `POST /security/installations` com `X-Step-Up-Token`;
6. persistir os tokens operacionais retornados pelo registro;
7. descartar token restrito, step-up token e senha em sucesso, erro definitivo
   ou cancelamento.

Depois do cadastro bem-sucedido, o fluxo segue como uma sessão autenticada
normal. Não existe novo login intermediário: o mobile passa a usar os tokens
operacionais retornados por `POST /security/installations`.

O mobile não deve:

- enviar a senha transacional no payload de login;
- enviar a senha transacional para o endpoint de registro;
- usar o token restrito fora do fluxo permitido;
- registrar senha, tokens ou identificador completo em logs.

No futuro, prova de vida poderá autorizar o mesmo endpoint de registro sem
alterar o papel do `X-Installation-Id`.

Senha transacional ativa é pré-requisito para esse fluxo. Se a API indicar que
ela está `not_set` ou `locked`, o mobile não deve tentar registrar a nova
instalação e deve direcionar o usuário para regularizar esse pré-requisito
antes de repetir o login nessa instalação.

### Limite atingido

Quando o login retornar `installation_limit_reached`, o mobile deve:

- informar que o limite de três instalações foi atingido;
- considerar `known_installations_count` e `max_installations` para compor a
  mensagem ou o estado exibido;
- informar que a instalação atual ainda não está cadastrada;
- não persistir tokens de sessão;
- não iniciar o step-up de registro enquanto não houver vaga;
- interpretar `next_action = revoke_existing_installation` como orientação
  estruturada de continuidade;
- orientar o usuário a acessar uma instalação já cadastrada e, a partir dela,
  revogar outra instalação para liberar vaga.

Instalações revogadas permanecem no histórico da API, mas não ocupam uma das
três vagas.

### Gerenciamento

A API disponibilizará:

```http
GET    /security/installations
DELETE /security/installations/{installation_resource_id}
```

O consumo mobile para listar e revogar instalações será detalhado depois que
for definida a experiência de gerenciamento.

Para o MVP, a revogação de instalação:

- não exige step-up;
- deve ocultar ou desabilitar a remoção da instalação atualmente em uso;
- deve orientar logout para encerrar a instalação atual, em vez de tentar
  revogá-la pelo gerenciamento.
- deve considerar que uma instalação revogada perde a sessão imediatamente,
  exigindo tratamento de retorno para acesso encerrado.

## 9. Decisões necessárias antes das tasks

- [x] Definir geração: UUID v4 aleatório por instalação.
- [x] Definir header: `X-Installation-Id`.
- [x] Preservar o identificador em logout e troca de usuário.
- [x] Gerar nova identidade após reinstalação ou limpeza dos dados.
- [ ] Definir chave e mecanismo final de persistência.
- [ ] Definir estratégia de detecção de reinstalação no iOS e Android.
- [ ] Definir exclusão de backup e restauração quando aplicável.
- [ ] Definir sincronização do `read-or-create`.
- [ ] Definir integração do interceptor com login, refresh e cliente principal.
- [ ] Definir tratamento do resultado
  `installation_registration_required`.
- [x] Definir pós-cadastro: `POST /security/installations` conclui o bootstrap
  da sessão operacional e retorna os tokens normais.
- [ ] Definir contrato e UX final do resultado `installation_limit_reached`.
- [x] Definir diretriz para `installation_limit_reached`: orientar revogação
  em instalação já cadastrada, sem step-up de registro nessa tentativa.
- [ ] Definir armazenamento temporário e descarte do access token restrito.
- [ ] Integrar step-up para `POST /security/installations`.
- [ ] Implementar `POST /security/installations` e troca pelos tokens
  operacionais.
- [x] Reconhecer os contratos de listagem e revogação definidos pela API.
- [ ] Definir a experiência mobile de gerenciamento das instalações.
- [x] Definir regra de revogação: sem step-up e sem remover a instalação em
  uso.
- [ ] Definir tela e estados do registro com senha transacional.
- [ ] Definir cancelamento e expiração da autorização restrita.
- [x] Definir regra para senha transacional ausente ou bloqueada: nova
  instalação não é registrada; é pré-requisito.
- [ ] Definir telemetria segura para falhas de armazenamento.
- [ ] Definir testes de ciclo de vida e concorrência.

## 10. Fora de escopo

- identificação do aparelho físico;
- attestation de plataforma;
- device fingerprinting;
- biometria como prova para o backend;
- decisão local de confiança;
- associação, listagem ou revogação de instalações na API;
- políticas para operações sensíveis;
- geolocalização e score antifraude.

## 11. Critérios para encerrar a discussão

- Estratégia de persistência e reinstalação definida por plataforma.
- Ciclo de vida separado das credenciais documentado.
- Integração HTTP definida para login, refresh e chamadas autenticadas.
- Fluxo de registro de nova instalação definido.
- Comportamento de falha e concorrência definido.
- Contrato alinhado com o backlog API 010.
- Tasks de implementação mobile derivadas.

## 12. Referências internas

- [Installation Identity MVP API](<../api/010 - installation-identity-mvp.md>)
- [Arquitetura mobile](../../../mobile/docs/ARCHITECTURE.md)
- [Roadmap](../../ROADMAP.md)
