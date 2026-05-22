# Tasks: Pré-onboarding com verificação de contato

Estas tasks dividem o backlog mobile de pré-onboarding em passos executáveis.

O objetivo é adaptar o app ao novo contrato da API: antes de registrar o usuário,
o mobile deve verificar e-mail e telefone, coletar telefone e data de nascimento,
e enviar o cadastro final com os dois tokens confirmados.

## Task 1/10: Adicionar DTOs de verificação de contato

### Objetivo

Representar no mobile os contratos dos endpoints de verificação de contato.

### Escopo

- Criar DTO para solicitar verificação:
  - `channel`;
  - `target`.
- Criar DTO para resposta da solicitação:
  - `verification_id`;
  - `channel`;
  - `target`;
  - `token`;
  - `expires_at`.
- Criar DTO para confirmar verificação:
  - `verification_id`;
  - `token`.
- Criar DTO para resposta da confirmação:
  - `verification_token`;
  - `channel`;
  - `target`;
  - `verified_at`.
- Usar o envelope existente `data`/`error`.
- Manter os DTOs em `mobile/lib/data/services/auth/api/dtos`.

### Critérios de aceite

- DTO de solicitação serializa o body esperado pela API.
- DTO de confirmação serializa o body esperado pela API.
- DTOs de resposta parseiam payloads válidos.
- Campos de data são tratados de forma compatível com os padrões atuais do app.
- Testes cobrem serialização e parse.

### Depende de

- Nenhuma.

## Task 2/10: Adicionar métodos de contact verification no AuthApi

### Objetivo

Permitir que a camada de dados chame os endpoints de verificação de contato.

### Escopo

- Adicionar método para `POST /auth/contact-verifications`.
- Adicionar método para `POST /auth/contact-verifications/confirm`.
- Enviar `X-App-Token` em ambas as chamadas.
- Parsear resposta com `ApiEnvelope`.
- Retornar `Result` seguindo o padrão atual de `AuthApi`.
- Em ambiente de desenvolvimento, imprimir no terminal/log o `token` curto
  retornado pela solicitação de verificação.
- Não persistir tokens de sessão.

### Critérios de aceite

- `AuthApi` solicita verificação de e-mail.
- `AuthApi` solicita verificação de telefone.
- `AuthApi` confirma código de e-mail.
- `AuthApi` confirma código de telefone.
- Token curto retornado pela API aparece no log de depuração.
- Falhas da API retornam `Failure(AppError)`.
- Testes cobrem sucesso e erro nos dois endpoints.

### Depende de

- Task 1.

## Task 3/10: Atualizar RegisterRequestDto para o novo contrato

### Objetivo

Fazer o cadastro final enviar todos os campos obrigatórios pela API.

### Escopo

- Adicionar ao `RegisterRequestDto`:
  - `phone`;
  - `birthDate`;
  - `emailVerificationToken`;
  - `phoneVerificationToken`.
- Serializar `birthDate` como `birth_date` no formato `YYYY-MM-DD`.
- Serializar `emailVerificationToken` como `email_verification_token`.
- Serializar `phoneVerificationToken` como `phone_verification_token`.
- Continuar enviando:
  - `name`;
  - `email`;
  - `password`;
  - `cpf`.
- Manter CPF normalizado com 11 dígitos.
- Atualizar testes, mocks e fake repositories que constroem
  `RegisterRequestDto`.

### Critérios de aceite

- `RegisterRequestDto.toMap()` gera o JSON esperado pela API.
- `cpf` continua no payload.
- `phone` é enviado no formato definido para a integração.
- `birth_date` é enviado como data ISO curta.
- Testes existentes compilam com a nova assinatura.
- Testes cobrem o novo payload completo.

### Depende de

- Nenhuma.

## Task 4/10: Expor verificação de contato no AuthRepository

### Objetivo

Permitir que o `RegisterViewmodel` orquestre a verificação sem depender
diretamente do `AuthApi`.

### Escopo

- Adicionar ao `AuthRepository` métodos para:
  - solicitar verificação de contato;
  - confirmar verificação de contato.
- Implementar os métodos em `AuthRepositoryImpl`.
- Delegar as chamadas para `AuthApi`.
- Manter `register` como operação final.
- Não alterar `currentUser`, `isLoggedIn` ou tokens salvos durante
  pré-onboarding.

### Critérios de aceite

- `RegisterViewmodel` consegue solicitar e confirmar contato via repository.
- Falhas são propagadas como `AppError`.
- Sucessos retornam os dados necessários para a próxima etapa.
- Nenhum token de login é salvo durante verificação ou cadastro.
- Testes de repository são atualizados.

### Depende de

- Task 2.

## Task 5/10: Adicionar suporte a details em ApiError

### Objetivo

Permitir que o app leia `error.details` retornado pela API.

### Escopo

- Adicionar campo `details` em `ApiError`.
- Parsear `details` em `ApiError.fromMap`.
- Preservar compatibilidade com erros sem `details`.
- Garantir que `ApiEnvelope` continue funcionando para todos os payloads
  atuais.

### Critérios de aceite

- `ApiError` parseia `code`, `message` e `details`.
- Erros sem `details` continuam válidos.
- Testes cobrem erro com details e erro sem details.

### Depende de

- Nenhuma.

## Task 6/10: Mapear CONTACT_NOT_VERIFIED no erro mobile

### Objetivo

Transformar o backend code `CONTACT_NOT_VERIFIED` em erro semântico no app.

### Escopo

- Adicionar `AppErrorCode.contactNotVerified`.
- Atualizar `dio_error_mapper.dart` para mapear
  `CONTACT_NOT_VERIFIED`.
- Preservar `details.email_verified` e `details.phone_verified`.
- Definir mensagem específica para:
  - e-mail pendente;
  - telefone pendente;
  - ambos pendentes.
- Manter comportamento de `INVALID_CREDENTIALS` inalterado.
- Manter comportamento de `ACCOUNT_APPROVAL_REQUIRED` inalterado.
- Não mapear `403` genérico para `contactNotVerified`.

### Critérios de aceite

- `CONTACT_NOT_VERIFIED` mapeia para
  `AppErrorCode.contactNotVerified`.
- Mensagem usa os detalhes por canal quando disponíveis.
- Sem detalhes, mensagem genérica é usada.
- `INVALID_CREDENTIALS` continua com o comportamento atual.
- `ACCOUNT_APPROVAL_REQUIRED` continua com o comportamento atual.
- Testes cobrem todos esses cenários.

### Depende de

- Task 5.

## Task 7/10: Tratar CONTACT_NOT_VERIFIED no login e short login

### Objetivo

Mostrar feedback claro quando login for bloqueado por contato não verificado.

### Escopo

- Atualizar `LoginPage`.
- Atualizar `ShortLoginPage`.
- Exibir mensagem específica para `AppErrorCode.contactNotVerified`.
- Usar `AppSnackbar` como nos fluxos atuais.
- Não navegar para home nesse erro.
- Não limpar a identidade lembrada no short login.
- Preservar mensagens e comportamento de:
  - credenciais inválidas;
  - aprovação pendente;
  - erros genéricos.

### Critérios de aceite

- Login completo mostra mensagem de contato não verificado.
- Short login mostra mensagem de contato não verificado.
- Mensagem diferencia e-mail pendente, telefone pendente ou ambos.
- Login permanece na tela atual.
- Short login mantém a identidade lembrada visível.
- Testes de UI cobrem o novo feedback.

### Depende de

- Task 6.

## Task 8/10: Redesenhar RegisterViewmodel para fluxo multi-etapas

### Objetivo

Fazer o `RegisterViewmodel` gerenciar a jornada de cadastro em múltiplas etapas.

### Escopo

- Modelar estado de etapas do cadastro.
- Guardar os dados preenchidos entre etapas:
  - nome;
  - CPF;
  - data de nascimento;
  - e-mail;
  - telefone;
  - senha;
  - verification id de e-mail;
  - verification token de e-mail;
  - verification id de telefone;
  - verification token de telefone.
- Expor comandos para:
  - solicitar código de e-mail;
  - confirmar código de e-mail;
  - solicitar código de telefone;
  - confirmar código de telefone;
  - registrar usuário.
- Bloquear avanço quando dados obrigatórios da etapa estiverem inválidos.
- Não perder dados quando o usuário volta para uma etapa anterior.
- Manter `RegisterViewmodel` como orquestrador principal do registro.

### Critérios de aceite

- ViewModel representa claramente a etapa atual.
- ViewModel sabe se e-mail foi confirmado.
- ViewModel sabe se telefone foi confirmado.
- Cadastro final só é executado com os dois tokens confirmados.
- Falhas por etapa são expostas para a UI.
- Testes cobrem transições principais e bloqueios.

### Depende de

- Task 3.
- Task 4.

## Task 9/10: Redesenhar RegisterPage para cadastro em múltiplas etapas

### Objetivo

Atualizar a experiência de cadastro para o novo fluxo de pré-onboarding.

### Escopo

- Redesenhar a tela de cadastro.
- Dividir o cadastro em etapas visíveis.
- Coletar dados pessoais:
  - nome;
  - CPF;
  - data de nascimento via date picker.
- Coletar contatos:
  - e-mail;
  - telefone nacional no formato visual `(27) 99999-9999`.
- Permitir solicitar e confirmar código de e-mail.
- Permitir solicitar e confirmar código de telefone.
- Exibir feedback de sucesso e falha por etapa.
- Mostrar loading por ação em andamento.
- Bloquear cadastro final até ambos os contatos estarem confirmados.
- Enviar telefone em formato aceito pela API.
- Manter layout responsivo em telas pequenas.

### Critérios de aceite

- Usuário consegue avançar pelas etapas do cadastro.
- Date picker é usado para data de nascimento.
- Telefone é digitado com máscara nacional.
- Código de e-mail pode ser solicitado e confirmado.
- Código de telefone pode ser solicitado e confirmado.
- Botão de cadastro final só habilita após confirmação dos dois contatos.
- Cadastro final chama o `RegisterViewmodel` com todos os dados obrigatórios.
- Tela não quebra em dispositivos pequenos.
- Testes de widget cobrem o fluxo principal.

### Depende de

- Task 8.

## Task 10/10: Validar fluxo completo e atualizar testes

### Objetivo

Garantir que o app mobile permanece saudável com o novo cadastro e o novo erro
de login.

### Escopo

- Atualizar testes de DTO.
- Atualizar testes de `AuthApi`.
- Atualizar testes de `AuthRepositoryImpl`.
- Atualizar testes de erro HTTP.
- Atualizar testes de feedback de login.
- Adicionar testes para o novo fluxo de cadastro.
- Rodar `dart format` nos arquivos Dart alterados.
- Rodar testes focados de auth.
- Rodar `flutter analyze`.
- Rodar `flutter test`, se viável.
- Documentar qualquer teste não executado e o motivo.

### Critérios de aceite

- DTOs novos e alterados têm cobertura.
- `AuthApi` tem cobertura para contact verification.
- `AuthRepositoryImpl` tem cobertura para os novos métodos.
- `CONTACT_NOT_VERIFIED` tem cobertura no mapper e na UI.
- Fluxo principal de cadastro tem cobertura.
- `dart format` foi executado nos arquivos alterados.
- `flutter analyze` passa.
- `flutter test` passa ou a impossibilidade é registrada.

### Depende de

- Task 1.
- Task 2.
- Task 3.
- Task 4.
- Task 5.
- Task 6.
- Task 7.
- Task 8.
- Task 9.

