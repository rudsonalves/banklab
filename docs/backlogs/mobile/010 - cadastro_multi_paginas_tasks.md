# Tasks: Cadastro em múltiplas páginas

Estas tasks dividem o backlog mobile de cadastro em múltiplas páginas em passos
executáveis.

O objetivo é substituir a `RegisterPage` única por um fluxo de onboarding de
cadastro com páginas pequenas, estado compartilhado no `RegisterViewmodel`,
rascunho local em secure storage por hash do CPF e criação da conta apenas ao
fim do processo.

A base inicial deve ser específica do cadastro de usuário. Não criar ainda um
gerenciador genérico de onboarding. O desenho deve separar estado do rascunho,
snapshot persistível e store seguro, deixando a evolução futura possível sem
antecipar abstrações.

## Task 1/12: Criar estado e snapshot do rascunho local do cadastro

### Objetivo

Representar o estado persistível do onboarding de cadastro, com dirty tracking
simples, sem incluir senha ou tokens confirmados.

### Escopo

- Criar `RegisterDraftState` ou nome equivalente para o estado em memória do
  rascunho.
- Criar `RegisterDraftSnapshot` ou nome equivalente para a forma persistível em
  JSON.
- Incluir campos persistíveis:
  - `cpf`;
  - `name`;
  - `birthDate`;
  - `email`;
  - `phone`;
  - `currentStep`;
  - `emailVerificationId`;
  - `phoneVerificationId`;
  - `isEmailVerified`;
  - `isPhoneVerified`;
  - `createdAt`;
  - `updatedAt`.
- Registrar campos alterados desde a última persistência com dirty tracking
  simples.
- Expor operação para marcar o estado como limpo após persistência bem-sucedida.
- Expor operação para hidratar o estado a partir de um snapshot.
- Expor operação para gerar snapshot persistível a partir do estado.
- Não incluir:
  - senha;
  - confirmação de senha;
  - `emailVerificationToken`;
  - `phoneVerificationToken`;
  - tokens de sessão.
- Serializar o rascunho como JSON.
- Parsear rascunhos inválidos como ausência de rascunho recuperável.

### Critérios de aceite

- O snapshot serializa e parseia todos os campos persistíveis.
- O estado registra campos alterados.
- O estado consegue limpar dirty tracking após persistência.
- O estado consegue hidratar a partir de snapshot válido.
- Senha e tokens confirmados não existem no snapshot persistido.
- `createdAt` e `updatedAt` são preservados.
- Testes cobrem serialização, parse, hidratação, dirty tracking e payload
  inválido.

### Depende de

- Nenhuma.

## Task 2/12: Criar store de persistência segura do rascunho

### Objetivo

Persistir e recuperar o rascunho do onboarding em secure storage usando chave
derivada do CPF.

### Escopo

- Criar store local para rascunho de cadastro.
- Usar secure storage já disponível no app ou abstração equivalente.
- Receber um snapshot persistível, não o `RegisterViewmodel`.
- Gerar chave no formato:

```text
onboarding_draft:{sha256(cpf_normalizado)}
```

- Normalizar CPF antes de calcular a chave.
- Salvar rascunho em JSON.
- Salvar o JSON inteiro do snapshot quando houver alteração, sem patch parcial
  por campo.
- Recuperar rascunho por CPF.
- Remover rascunho por CPF.
- Tratar ausência de rascunho sem erro.
- Tratar JSON corrompido removendo ou ignorando o rascunho.
- Retornar um resultado selado para o carregamento:
  - `RegisterDraftFound`;
  - `RegisterDraftNotFound`.
- Não aplicar TTL nesta camada.

### Critérios de aceite

- Rascunho é salvo em secure storage.
- CPF puro não é usado como chave.
- Store não conhece regras de UI nem depende do `RegisterViewmodel`.
- Store não conhece regra de TTL.
- Recuperação por CPF encontra o rascunho correto.
- Remoção apaga o rascunho correto.
- Testes cobrem save, get, delete, ausência e JSON inválido.

### Depende de

- Task 1.

## Task 3/12: Criar repository do rascunho com expiração de 24 horas

### Objetivo

Criar a camada de repository que encapsula regra de TTL e expiração do rascunho,
mantendo o store restrito a secure storage, hash e JSON.

### Escopo

- Criar contrato `RegisterDraftRepository` ou nome equivalente.
- Criar implementação `RegisterDraftRepositoryImpl` ou nome equivalente.
- O repository deve depender do `RegisterDraftStore`.
- Expor operações:
  - carregar por CPF;
  - salvar snapshot;
  - remover por CPF.
- Aplicar TTL de 24 horas ao carregar rascunho.
- Ao carregar rascunho, considerar `updatedAt` como referência principal.
- Se o rascunho estiver expirado, chamar o store para removê-lo e retornar
  `RegisterDraftNotFound`.
- Se o rascunho estiver dentro do TTL, retornar `RegisterDraftFound`.
- Ao salvar, receber snapshot persistível e delegar ao store.
- O repository deve ser a dependência usada pelo `RegisterViewmodel` nas tasks
  seguintes, não o store diretamente.
- Permitir injetar `DateTime Function()` nos testes para controlar o relógio.

### Critérios de aceite

- Rascunho com menos de 24 horas é recuperado.
- Rascunho com mais de 24 horas é removido e não é usado.
- Store permanece sem regra de TTL.
- Repository retorna `RegisterDraftFound` para rascunho válido.
- Repository retorna `RegisterDraftNotFound` para rascunho ausente ou expirado.
- Repository remove rascunho expirado.
- Testes cobrem rascunho válido, expirado, ausente e falha do store.

### Depende de

- Task 2.

## Task 4/12: Ajustar RegisterViewmodel para fluxo em páginas pequenas

### Objetivo

Fazer o `RegisterViewmodel` representar o fluxo granular de cadastro e manter o
estado entre páginas.

### Escopo

- Atualizar as etapas do cadastro para:
  - CPF;
  - nome completo;
  - data de nascimento;
  - e-mail;
  - token do e-mail;
  - telefone;
  - token do telefone;
  - senha;
  - confirmação da senha;
  - sucesso.
- Manter dados em memória:
  - CPF;
  - nome;
  - data de nascimento;
  - e-mail;
  - telefone;
  - senha;
  - confirmação de senha;
  - verification ids;
  - verification tokens confirmados.
- Delegar o estado persistível para `RegisterDraftState` ou nome equivalente,
  em vez de espalhar todos os campos persistíveis diretamente no viewmodel.
- Manter senha, confirmação de senha e tokens confirmados fora do snapshot
  persistível.
- Expor métodos de atualização por etapa.
- Expor validações por etapa.
- Expor navegação para próxima etapa e etapa anterior.
- Implementar `reset()` para limpar estado ao concluir, cancelar ou reiniciar.
- Preservar `RegisterViewmodel` como orquestrador do cadastro.

### Critérios de aceite

- ViewModel representa todas as etapas do novo fluxo.
- ViewModel orquestra o fluxo usando o estado de rascunho específico do
  cadastro.
- Cada etapa possui validação própria.
- Voltar e avançar não perde dados em memória.
- `reset()` remove dados sensíveis em memória, incluindo senha e tokens
  confirmados.
- Testes cobrem transições principais e validações.

### Depende de

- Nenhuma.

## Task 5/12: Integrar RegisterViewmodel com rascunho local

### Objetivo

Permitir retomada do cadastro por CPF e persistência incremental do progresso.

### Escopo

- Injetar o serviço de rascunho no `RegisterViewmodel`.
- Ao validar CPF, tentar carregar rascunho por CPF.
- Se houver rascunho válido, hidratar o `RegisterDraftState` e refletir a etapa
  correta no `RegisterViewmodel`.
- Persistir rascunho ao avançar nas etapas persistíveis.
- Persistir apenas quando o estado do rascunho tiver campos alterados.
- Não persistir senha, confirmação de senha ou tokens confirmados.
- Persistir flags de contato verificado quando aplicável.
- Ao reabrir app ou recarregar rascunho, não recuperar tokens confirmados.
- Se contato estiver marcado como verificado sem token confirmado em memória,
  direcionar o usuário para refazer a confirmação antes do register.
- Remover rascunho ao concluir cadastro.
- Remover rascunho ao cancelar explicitamente cadastro.

### Critérios de aceite

- CPF com rascunho válido hidrata o fluxo.
- CPF sem rascunho inicia fluxo limpo.
- Rascunho expirado não é usado.
- Persistência não é chamada quando não há alterações no rascunho.
- Senha e tokens confirmados nunca são gravados.
- Cadastro concluído apaga o rascunho.
- Testes cobrem retomada, persistência incremental, expiração e conclusão.

### Depende de

- Task 2.
- Task 3.
- Task 4.

## Task 6/12: Garantir RegisterViewmodel como lazySingleton e ciclo de vida correto

### Objetivo

Preservar estado entre páginas do cadastro sem manter dados após o fim do fluxo.

### Escopo

- Revisar o registro atual do `RegisterViewmodel` em DI.
- Registrar ou manter como `lazySingleton`.
- Garantir que as páginas do fluxo usem a mesma instância.
- Chamar `reset()` ao concluir cadastro com sucesso.
- Chamar `reset()` ao cancelar cadastro explicitamente.
- Evitar recriação acidental do viewmodel ao navegar entre páginas.

### Critérios de aceite

- Todas as páginas do cadastro usam a mesma instância do `RegisterViewmodel`.
- O viewmodel só é instanciado ao entrar no fluxo de cadastro.
- Dados são preservados entre páginas.
- Dados são limpos ao concluir ou cancelar.
- Testes ou verificação de DI cobrem a instância compartilhada.

### Depende de

- Task 4.

## Task 7/12: Criar rotas do fluxo de cadastro

### Objetivo

Separar o cadastro em rotas/páginas dedicadas.

### Escopo

- Atualizar `AuthRoutes` com rotas para:
  - CPF;
  - nome;
  - data de nascimento;
  - e-mail;
  - token do e-mail;
  - telefone;
  - token do telefone;
  - senha;
  - confirmação da senha;
  - sucesso.
- Registrar as rotas em `auth_routes.dart`.
- Manter navegação por nomes de rota.
- Atualizar o botão de cadastro no login para iniciar pela rota de CPF.
- Evitar rotas órfãs para a `RegisterPage` antiga.

### Critérios de aceite

- Login abre o início do novo fluxo de cadastro.
- Cada página possui rota nomeada.
- Navegação entre etapas usa o `RegisterViewmodel`.
- A rota antiga de cadastro não deixa a página monolítica acessível.
- Testes de navegação são atualizados.

### Depende de

- Task 6.

## Task 8/12: Implementar páginas de dados pessoais

### Objetivo

Criar as primeiras páginas do fluxo: CPF, nome completo e data de nascimento.

### Escopo

- Criar página de CPF.
- Criar página de nome completo.
- Criar página de data de nascimento.
- Adicionar DTOs e métodos na camada de dados, se ainda não existirem, para
  `POST /auth/cpf-check`.
- Expor a checagem de CPF no `AuthRepository` ou serviço equivalente usado pelo
  `RegisterViewmodel`.
- Usar date picker para data de nascimento.
- Aplicar máscaras/normalização necessárias para CPF.
- Chamar `POST /auth/cpf-check` com `X-App-Token` antes de avançar da página de
  CPF.
- Bloquear avanço quando a API retornar CPF indisponível.
- Mostrar feedback de validação por etapa.
- Persistir rascunho ao avançar quando aplicável.
- Permitir voltar sem perder dados.

### Critérios de aceite

- Página de CPF valida e normaliza CPF.
- Página de CPF consulta disponibilidade na API.
- CPF já cadastrado bloqueia o avanço.
- Página de CPF tenta recuperar rascunho.
- Página de nome valida campo obrigatório.
- Página de data usa date picker.
- Avanço inválido mostra feedback claro.
- Dados permanecem ao voltar.
- Testes de widget cobrem caminho feliz e validações.

### Depende de

- Task 5.
- Task 7.

## Task 9/12: Implementar páginas de verificação de e-mail

### Objetivo

Separar entrada de e-mail e confirmação de token em páginas dedicadas.

### Escopo

- Criar página de e-mail.
- Criar página de token de e-mail.
- Validar formato básico de e-mail.
- Solicitar código com `POST /auth/contact-verifications`.
- Enviar `X-App-Token` pela camada já existente.
- Persistir `email` e `emailVerificationId`.
- Confirmar código com `POST /auth/contact-verifications/confirm`.
- Manter `emailVerificationToken` apenas em memória.
- Ao reabrir fluxo sem token em memória, exigir nova confirmação.
- Mostrar loading e feedback de erro/sucesso.

### Critérios de aceite

- E-mail inválido bloqueia solicitação de código.
- Código de e-mail é solicitado via repository.
- Token curto continua aparecendo no log de depuração via `AuthApi`.
- Confirmação válida avança o fluxo.
- Token confirmado não é persistido.
- Testes cobrem solicitação, confirmação e falhas.

### Depende de

- Task 5.
- Task 7.

## Task 10/12: Implementar páginas de verificação de telefone

### Objetivo

Separar entrada de telefone e confirmação de token em páginas dedicadas.

### Escopo

- Criar página de telefone.
- Criar página de token de telefone.
- Usar formato visual nacional `(27) 99999-9999`.
- Converter telefone para formato de API, por exemplo `+5527999999999`.
- Solicitar código com `POST /auth/contact-verifications`.
- Enviar `X-App-Token` pela camada já existente.
- Persistir `phone` e `phoneVerificationId`.
- Confirmar código com `POST /auth/contact-verifications/confirm`.
- Manter `phoneVerificationToken` apenas em memória.
- Ao reabrir fluxo sem token em memória, exigir nova confirmação.
- Mostrar loading e feedback de erro/sucesso.

### Critérios de aceite

- Telefone inválido bloqueia solicitação de código.
- Telefone visual é convertido corretamente para API.
- Código de telefone é solicitado via repository.
- Token curto continua aparecendo no log de depuração via `AuthApi`.
- Confirmação válida avança o fluxo.
- Token confirmado não é persistido.
- Testes cobrem solicitação, confirmação, conversão e falhas.

### Depende de

- Task 5.
- Task 7.

## Task 11/12: Implementar páginas de senha, criação de conta e sucesso

### Objetivo

Finalizar o fluxo criando o usuário apenas após senha e confirmação válidas.

### Escopo

- Criar página de senha.
- Criar página de confirmação de senha.
- Criar página de sucesso.
- Validar regra mínima de senha.
- Validar igualdade entre senha e confirmação.
- Não persistir senha nem confirmação.
- Chamar `RegisterViewmodel.register()` somente com:
  - dados pessoais válidos;
  - e-mail confirmado em memória;
  - telefone confirmado em memória;
  - senha válida.
- `register` deve chamar `POST /auth/register` com `X-App-Token` pela camada já
  existente.
- Após sucesso:
  - remover rascunho local;
  - limpar estado em memória;
  - não salvar tokens de sessão;
  - navegar para tela de sucesso.
- Na tela de sucesso, oferecer ação para ir ao login.

### Critérios de aceite

- Senha inválida bloqueia avanço.
- Confirmação diferente bloqueia criação da conta.
- Cadastro só é chamado com tokens confirmados em memória.
- Usuário criado não fica autenticado no app.
- Rascunho é removido após sucesso.
- Estado em memória é limpo após sucesso.
- Tela de sucesso navega para login.
- Testes cobrem sucesso, erro de register e senha inválida.

### Depende de

- Task 5.
- Task 7.
- Task 9.
- Task 10.

## Task 12/12: Remover RegisterPage monolítica e validar fluxo completo

### Objetivo

Concluir a migração do cadastro antigo para o novo fluxo e garantir saúde dos
testes.

### Escopo

- Remover ou desativar a `RegisterPage` monolítica antiga.
- Remover testes antigos que dependem do formulário único ou adaptá-los ao novo
  fluxo.
- Atualizar imports, rotas e referências quebradas.
- Atualizar testes de:
  - endpoint mobile de checagem de CPF, quando implementado na camada de dados;
  - modelo de rascunho;
  - serviço de rascunho;
  - `RegisterViewmodel`;
  - páginas do fluxo;
  - navegação;
  - `AuthApi` quando necessário.
- Rodar `dart format` nos arquivos Dart alterados.
- Rodar `flutter analyze`.
- Rodar testes focados de cadastro/auth.
- Rodar `flutter test`, se viável.
- Documentar qualquer teste não executado e o motivo.

### Critérios de aceite

- A página monolítica antiga não é mais usada no app.
- O novo fluxo cobre o caminho feliz completo.
- O novo fluxo cobre falhas de validação por etapa.
- Persistência segura do rascunho tem cobertura.
- Retomada por CPF tem cobertura.
- Expiração de 24 horas tem cobertura.
- Senha e tokens confirmados não são persistidos.
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
- Task 10.
- Task 11.
