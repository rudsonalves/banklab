# Mobile: Pré-onboarding com verificação de contato

## Problema

A API mudou o fluxo de cadastro para exigir verificação prévia de e-mail e
telefone antes de criar o usuário. O endpoint `POST /auth/register` não aceita
mais apenas `name`, `email`, `cpf` e `password`.

O novo contrato exige:

```json
{
  "email": "user@example.com",
  "phone": "+5511999999999",
  "password": "P@ssword123",
  "name": "Maria Silva",
  "birth_date": "1990-01-15",
  "cpf": "12345678901",
  "email_verification_token": "token-confirmado-email",
  "phone_verification_token": "token-confirmado-phone"
}
```

Além disso, o login pode retornar `CONTACT_NOT_VERIFIED` quando e-mail ou
telefone ainda não estiverem verificados.

O app mobile ainda usa o fluxo antigo de cadastro e não possui chamadas para os
endpoints de verificação de contato.

## Objetivo

Atualizar o app mobile para suportar o novo fluxo de pré-onboarding da API:

1. Solicitar verificação de e-mail.
2. Confirmar código de e-mail.
3. Solicitar verificação de telefone.
4. Confirmar código de telefone.
5. Criar cadastro com telefone, data de nascimento, CPF e tokens confirmados.
6. Tratar `CONTACT_NOT_VERIFIED` no login com mensagem específica.

## Contexto atual

- `RegisterRequestDto` ainda envia apenas `name`, `email`, `password` e `cpf`.
- `RegisterPage` possui campos de nome, e-mail, CPF e senha.
- `AuthApi` possui `register`, `login` e `getProfile`, mas não possui chamadas
  para contact verification.
- `AuthRepository` expõe `register`, mas não expõe operações de verificação de
  contato.
- `dio_error_mapper.dart` trata `ACCOUNT_APPROVAL_REQUIRED`, mas não trata
  `CONTACT_NOT_VERIFIED`.
- CPF continua existindo no contrato mobile de cadastro e perfil; a mudança de
  armazenamento do CPF é interna da API.

## Fora de escopo

- Não remover CPF das telas mobile.
- Não alterar o fluxo de transferência por CPF/documento.
- Não implementar envio real de SMS/e-mail; a API atual retorna o código no
  response para ambiente de desenvolvimento.
- Não implementar recuperação de senha.
- Não implementar edição de telefone/e-mail após cadastro.
- Não alterar o fluxo de aprovação administrativa.
- Não criar um fluxo de onboarding completo pós-cadastro; isso pertence ao
  backlog de onboarding.

## Contratos da API

### Solicitar verificação

`POST /auth/contact-verifications`

Headers:

```text
X-App-Token: <app_token>
```

Body:

```json
{
  "channel": "email",
  "target": "user@example.com"
}
```

`channel` aceita `email` ou `phone`.

Resposta esperada:

```json
{
  "data": {
    "verification_id": "uuid",
    "channel": "email",
    "target": "user@example.com",
    "token": "123456",
    "expires_at": "2026-05-18T12:10:00Z"
  },
  "error": null
}
```

### Confirmar verificação

`POST /auth/contact-verifications/confirm`

Headers:

```text
X-App-Token: <app_token>
```

Body:

```json
{
  "verification_id": "uuid",
  "token": "123456"
}
```

Resposta esperada:

```json
{
  "data": {
    "verification_token": "token-confirmado",
    "channel": "email",
    "target": "user@example.com",
    "verified_at": "2026-05-18T12:03:00Z"
  },
  "error": null
}
```

## Comportamento de produto

O cadastro deve guiar o usuário por uma jornada clara:

1. Preencher dados pessoais e credenciais.
2. Informar e-mail e telefone.
3. Solicitar código para e-mail.
4. Confirmar código de e-mail.
5. Solicitar código para telefone.
6. Confirmar código de telefone.
7. Enviar cadastro final.

O fluxo deve ser implementado em múltiplas etapas. A tela de cadastro atual deve
ser redesenhada para acomodar essa jornada, em vez de apenas adicionar novos
campos ao formulário existente.

O `RegisterViewmodel` continuará sendo o responsável por gerenciar o registro e
orquestrar o estado das etapas.

## Decisões fechadas

- O cadastro será dividido em múltiplas etapas.
- O `RegisterViewmodel` continuará gerenciando o fluxo de registro.
- O telefone será nacional brasileiro, no formato visual `(27) 99999-9999`.
- A data de nascimento deve ser coletada com date picker.
- Enquanto não houver envio real de e-mail/SMS, o token retornado pela API deve
  ser impresso no terminal/log de depuração da aplicação.
- A mensagem de `CONTACT_NOT_VERIFIED` deve ser específica por canal, usando os
  detalhes `email_verified` e `phone_verified` quando disponíveis.
- A tela de cadastro deve ser redesenhada para o novo fluxo.

## Epic 1: DTOs e modelos de contact verification

### Objetivo

Representar os novos contratos de request/response da API no mobile.

### Escopo

- Criar DTO para solicitar verificação de contato.
- Criar DTO para resposta da solicitação.
- Criar DTO para confirmar verificação.
- Criar DTO para resposta de confirmação.
- Preservar parse por envelope `data`/`error`.
- Representar `verification_id`, `token`, `verification_token`, `channel`,
  `target`, `expires_at` e `verified_at`.

### Critérios de aceite

- DTOs serializam os bodies esperados pela API.
- DTOs parseiam as respostas de sucesso.
- Falhas de parse retornam erro de parsing como nos demais serviços.
- Testes cobrem serialização e parse dos novos DTOs.

## Epic 2: AuthApi com endpoints de verificação

### Objetivo

Adicionar as chamadas HTTP necessárias ao pré-onboarding.

### Escopo

- Adicionar método para `POST /auth/contact-verifications`.
- Adicionar método para `POST /auth/contact-verifications/confirm`.
- Enviar `X-App-Token`.
- Usar `ApiEnvelope` para parse.
- Mapear erro de API para `AppError`, preservando mensagem e detalhes quando
  existirem.

### Critérios de aceite

- `AuthApi` consegue solicitar verificação de e-mail.
- `AuthApi` consegue solicitar verificação de telefone.
- `AuthApi` consegue confirmar código de e-mail.
- `AuthApi` consegue confirmar código de telefone.
- Testes cobrem sucesso e erro de API nos dois endpoints.

## Epic 3: Registro com novo contrato

### Objetivo

Atualizar o DTO e a chamada de registro para o novo contrato da API.

### Escopo

- Atualizar `RegisterRequestDto` para incluir:
  - `phone`;
  - `birthDate`;
  - `emailVerificationToken`;
  - `phoneVerificationToken`.
- Enviar `birth_date` no formato `YYYY-MM-DD`.
- Continuar enviando `cpf` normalizado com 11 dígitos.
- Manter `email`, `password` e `name`.
- Atualizar testes e mocks que constroem `RegisterRequestDto`.

### Critérios de aceite

- `POST /auth/register` envia todos os campos obrigatórios.
- Registro sem telefone não é permitido pelo app.
- Registro sem data de nascimento não é permitido pelo app.
- Registro sem tokens confirmados não é permitido pelo app.
- CPF continua sendo enviado ao backend.
- Testes de DTO e repository são atualizados.

## Epic 4: Repositório de auth

### Objetivo

Expor operações de verificação de contato para a camada de UI/viewmodel.

### Escopo

- Adicionar métodos ao `AuthRepository` para:
  - solicitar verificação de contato;
  - confirmar verificação de contato.
- Implementar os métodos em `AuthRepositoryImpl`.
- Manter `register` como operação final do cadastro.
- Não persistir tokens de auth durante o pré-onboarding.
- Não alterar estado de login durante o cadastro.

### Critérios de aceite

- ViewModels conseguem solicitar e confirmar verificação sem depender
  diretamente de `AuthApi`.
- Falhas são propagadas como `AppError`.
- Cadastro final continua retornando `Unit`.
- Nenhum token de sessão é salvo antes do login.

## Epic 5: UI e ViewModel de cadastro

### Objetivo

Permitir que o usuário complete a verificação de e-mail e telefone antes do
cadastro final.

### Escopo

- Incluir campo de telefone.
- Incluir campo de data de nascimento.
- Usar máscara de telefone nacional no formato `(27) 99999-9999`.
- Enviar telefone para API em formato aceito pelo backend.
- Usar date picker para selecionar data de nascimento.
- Incluir controles para solicitar código de e-mail.
- Incluir campo para código de e-mail.
- Incluir controles para confirmar código de e-mail.
- Incluir controles para solicitar código de telefone.
- Incluir campo para código de telefone.
- Incluir controles para confirmar código de telefone.
- Bloquear cadastro final enquanto e-mail e telefone não estiverem confirmados.
- Exibir feedback de sucesso/falha por etapa.
- Evitar perda dos dados digitados ao alternar etapas.
- Imprimir no terminal/log de depuração o token curto retornado pela API ao
  solicitar verificação de e-mail ou telefone.

### Critérios de aceite

- Usuário consegue solicitar código de e-mail.
- Usuário consegue confirmar código de e-mail.
- Usuário consegue solicitar código de telefone.
- Usuário consegue confirmar código de telefone.
- Botão de cadastro final só fica disponível após os dois tokens confirmados.
- Cadastro final envia telefone, nascimento e tokens.
- Mensagens de erro são específicas por etapa quando possível.
- Data de nascimento é escolhida por date picker.
- Telefone é digitado no formato nacional brasileiro.
- Token curto retornado pela API é registrado em log de depuração.
- Layout não quebra em telas pequenas.

## Epic 6: Login com `CONTACT_NOT_VERIFIED`

### Objetivo

Tratar login bloqueado por contato não verificado com mensagem específica.

### Escopo

- Adicionar `AppErrorCode.contactNotVerified`.
- Mapear backend code `CONTACT_NOT_VERIFIED` no `dio_error_mapper.dart`.
- Preservar `error.details.email_verified` e `error.details.phone_verified`
  quando disponíveis.
- Atualizar `LoginPage` para mostrar mensagem específica.
- Atualizar `ShortLoginPage` para mostrar mensagem específica.
- Manter comportamento de `INVALID_CREDENTIALS` inalterado.
- Manter comportamento de `ACCOUNT_APPROVAL_REQUIRED` inalterado.

### Mensagem sugerida

```text
Confirme seu e-mail e telefone antes de entrar.
```

Se os detalhes indicarem apenas um canal pendente, a UI pode especializar:

```text
Confirme seu telefone antes de entrar.
```

```text
Confirme seu e-mail antes de entrar.
```

### Critérios de aceite

- Login completo exibe mensagem específica para `CONTACT_NOT_VERIFIED`.
- Short login exibe mensagem específica para `CONTACT_NOT_VERIFIED`.
- Mensagem diferencia e-mail pendente, telefone pendente ou ambos pendentes.
- Login não navega para home nesse erro.
- Tokens não são persistidos nesse erro.
- `INVALID_CREDENTIALS` continua usando o comportamento atual.
- `ACCOUNT_APPROVAL_REQUIRED` continua usando o comportamento atual.

## Epic 7: ApiError com details

### Objetivo

Permitir que o app use `error.details` retornado pela API.

### Escopo

- Adicionar `details` em `ApiError`.
- Parsear `details` em `ApiError.fromMap`.
- Garantir compatibilidade com erros sem `details`.
- Usar os detalhes no mapeamento de `CONTACT_NOT_VERIFIED` quando aplicável.

### Critérios de aceite

- Erros com `details` são parseados.
- Erros sem `details` continuam funcionando.
- `CONTACT_NOT_VERIFIED` preserva os flags `email_verified` e
  `phone_verified`.

## Epic 8: Testes e validação

### Escopo

- Testar DTOs de contact verification.
- Testar DTO atualizado de register.
- Testar `AuthApi` nos novos endpoints.
- Testar `AuthRepositoryImpl` nas operações de verificação.
- Testar mapeamento de `CONTACT_NOT_VERIFIED`.
- Testar feedback de login completo.
- Testar feedback de short login.
- Testar comportamento da tela de cadastro com tokens pendentes e confirmados.

### Critérios de aceite

- Testes de unidade dos DTOs passam.
- Testes de `AuthApi` passam.
- Testes de repository passam.
- Testes de mapper HTTP passam.
- Testes de UI afetados passam.
- `flutter test` passa ou testes não executados são listados com motivo.

## Arquivos provavelmente afetados

- `mobile/lib/data/services/auth/api/auth_api.dart`
- `mobile/lib/data/services/auth/api/dtos/register_request_dto.dart`
- `mobile/lib/data/services/auth/api/dtos/*contact_verification*.dart`
- `mobile/lib/data/repositories/auth/auth_repository.dart`
- `mobile/lib/data/repositories/auth/auth_repository_impl.dart`
- `mobile/lib/core/services/client_http/dio/dio_error_mapper.dart`
- `mobile/lib/core/result/errors/app_error_code.dart`
- `mobile/lib/data/services/apis/core/api_error.dart`
- `mobile/lib/ui/pages/auth/register/register_page.dart`
- `mobile/lib/ui/pages/auth/register/viewmodel/register_viewmodel.dart`
- `mobile/lib/ui/pages/auth/login/login_page.dart`
- `mobile/lib/ui/pages/auth/short_login/short_login_page.dart`
- `mobile/test/core/services/client_http/dio/dio_error_mapper_test.dart`
- `mobile/test/ui/pages/auth/login_feedback_behavior_test.dart`
- `mobile/test/data/repositories/auth/auth_repository_impl_test.dart`

## Critérios gerais de aceite

- Cadastro mobile funciona com o novo contrato da API.
- Usuário não consegue concluir cadastro sem e-mail e telefone verificados.
- Login trata `CONTACT_NOT_VERIFIED` com mensagem clara.
- Fluxos de login existentes continuam funcionando.
- CPF continua sendo enviado no cadastro e exibido onde o contrato atual exige.
- Lookup de transferência por CPF/documento não é alterado.
- Testes afetados são atualizados.
