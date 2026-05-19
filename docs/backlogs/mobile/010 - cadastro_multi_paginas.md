# Mobile: Cadastro em múltiplas páginas

## Problema

A tela atual de cadastro concentra muitos estados em uma única `RegisterPage`, criando uma experiência pesada, pouco legível e difícil de manter.

O fluxo de pré-onboarding exige dados pessoais, senha, verificação de e-mail e
verificação de telefone antes de criar o usuário. Esse fluxo fica mais claro se
cada etapa tiver sua própria página, com navegação explícita e estado de
cadastro compartilhado.

## Objetivo

Redesenhar o cadastro mobile como uma jornada em páginas pequenas e separadas:

1. CPF.
2. Nome completo.
3. Data de nascimento.
4. E-mail.
5. Token de confirmação do e-mail.
6. Telefone.
7. Token de confirmação do telefone.
8. Senha.
9. Confirmação da senha.
10. Criação da conta do usuário.

O `RegisterViewmodel` continua sendo o orquestrador do registro e deve manter o
estado acumulado entre as páginas.

O `RegisterViewmodel` pode continuar registrado como `lazySingleton`, desde que
seu estado seja limpo ao concluir, cancelar ou reiniciar explicitamente o
cadastro. Assim ele só é carregado quando o onboarding de cadastro começa, mas
preserva dados ao navegar entre páginas.

A implementação deve começar com uma base específica para o cadastro de usuário,
sem criar ainda um gerenciador genérico de onboarding. O desenho deve, porém,
deixar clara a separação entre:

- estado do rascunho do cadastro;
- snapshot persistível em secure storage;
- serviço de persistência do rascunho;
- orquestração feita pelo `RegisterViewmodel`.

Esse formato permite evoluir no futuro para um gerenciador de onboarding mais
amplo quando houver mais etapas, como documentos, endereço e KYC, sem antecipar
uma abstração genérica agora.

## Fluxo proposto

### Página 1: CPF

Campos:

- CPF.

Comportamento:

- CPF deve aceitar entrada formatada e armazenar apenas números no estado.
- Ao avançar, validar formato básico do CPF.
- Chamar `POST /auth/cpf-check` com `X-App-Token` para verificar se o CPF já
  está cadastrado.
- Se `available = false`, bloquear o cadastro e orientar o usuário a usar login
  ou recuperação de acesso quando existir.
- Após um CPF válido, tentar recuperar rascunho local de onboarding associado ao
  hash do CPF.
- Se houver rascunho, hidratar o `RegisterViewmodel` com os dados persistidos e
  permitir continuar do ponto salvo.
- Não chamar API nesta página.

### Página 2: Nome completo

Campos:

- Nome completo.

Comportamento:

- Validar campo obrigatório.
- Persistir rascunho local após avanço.
- Não chamar API nesta página.

### Página 3: Data de nascimento

Campos:

- Data de nascimento.

Comportamento:

- Data de nascimento deve usar date picker.
- Validar campo obrigatório.
- Persistir rascunho local após avanço.
- Não chamar API nesta página.

### Página 4: E-mail

Campos:

- E-mail.

Comportamento:

- Validar formato básico de e-mail.
- Permitir solicitar código por e-mail.
- Chamar `POST /auth/contact-verifications` com `channel = email`.
- Enviar `X-App-Token`.
- Em ambiente de desenvolvimento, o token retornado pela API continua sendo
  impresso no log de depuração pelo `AuthApi`.
- Persistir `email` e `email_verification_id` no rascunho local.

### Página 5: Token de confirmação do e-mail

Campos:

- Token de confirmação do e-mail.

Comportamento:

- Permitir informar o token recebido.
- Chamar `POST /auth/contact-verifications/confirm`.
- Guardar o `verification_token` confirmado no estado do `RegisterViewmodel`.
- Persistir a confirmação do e-mail no rascunho local.
- Só permitir avançar depois de confirmar o e-mail.

### Página 6: Telefone

Campos:

- Telefone.

Comportamento:

- Telefone visual em formato nacional: `(27) 99999-9999`.
- Converter para formato aceito pela API antes da chamada, por exemplo
  `+5527999999999`.
- Chamar `POST /auth/contact-verifications` com `channel = phone`.
- Enviar `X-App-Token`.
- Em ambiente de desenvolvimento, o token retornado pela API continua sendo
  impresso no log de depuração pelo `AuthApi`.
- Persistir `phone` e `phone_verification_id` no rascunho local.

### Página 7: Token de confirmação do telefone

Campos:

- Token de confirmação do telefone.

Comportamento:

- Permitir informar o token recebido.
- Chamar `POST /auth/contact-verifications/confirm`.
- Guardar o `verification_token` confirmado no estado do `RegisterViewmodel`.
- Persistir a confirmação do telefone no rascunho local.
- Só permitir avançar depois de confirmar o telefone.

### Página 8: Senha

Campos:

- Senha.

Comportamento:

- Validar senha mínima conforme regra atual do mobile/API.
- Não persistir senha no rascunho local.
- Manter senha apenas em memória no `RegisterViewmodel`.

### Página 9: Confirmação da senha

Campos:

- Conferir senha.

Comportamento:

- Validar igualdade entre senha e confirmação.
- Não persistir confirmação de senha.
- Após validação, permitir criar a conta.

### Página 10: Criação concluída

Comportamento:

- Chamar `POST /auth/register` com todos os dados acumulados.
- Criar o usuário na API.
- Não autenticar o usuário automaticamente.
- Não persistir tokens de sessão no mobile.
- Limpar o estado em memória do cadastro.
- Remover o rascunho local do onboarding.
- Navegar para uma tela de sucesso com ação para ir ao login.

## Viabilidade da criação da conta no fim do fluxo

É possível criar a conta do usuário nesse ponto com o contrato atual da API.

O `POST /auth/register` já recebe todos os dados necessários:

```json
{
  "email": "user@example.com",
  "phone": "+5527999999999",
  "password": "P@ssword123",
  "name": "Maria Silva",
  "birth_date": "1990-01-15",
  "cpf": "12345678901",
  "email_verification_token": "token-confirmado-email",
  "phone_verification_token": "token-confirmado-phone"
}
```

Até o fim da confirmação da senha, o mobile já deve ter acumulado:

- `cpf` da página 1.
- `name` da página 2.
- `birth_date` da página 3.
- `email` e `email_verification_token` das páginas 4 e 5.
- `phone` e `phone_verification_token` das páginas 6 e 7.
- `password` das páginas 8 e 9.

Com esses dados, o `RegisterViewmodel` consegue montar o `RegisterRequestDto` e
executar o cadastro final sem endpoint adicional.

Após esse processo, o usuário terá uma conta criada na API, mas ainda não estará
conectado no app. O login continua sendo uma etapa separada.

## Persistência local do onboarding

Enquanto o cadastro não for concluído, o mobile deve persistir um rascunho local
do onboarding para permitir retomada pelo CPF.

O rascunho deve ser persistido em secure storage, usando o serviço seguro já
existente no app ou `FlutterSecureStorage` equivalente. Mesmo sem senha e sem
tokens confirmados, o rascunho contém dados pessoais sensíveis.

### Chave do rascunho

O CPF não deve ser usado em texto claro como chave de storage.

Usar uma chave derivada do CPF normalizado:

```text
onboarding_draft:{sha256(cpf_normalizado)}
```

O usuário informa o CPF ao retomar. O app normaliza o CPF, calcula a chave e
tenta recuperar o rascunho correspondente.

O valor persistido deve ser um JSON do rascunho.

O estado do rascunho deve ter dirty tracking simples para registrar quais campos
foram alterados desde a última persistência. Como o rascunho é pequeno e o
secure storage trabalha bem com chave/valor, a persistência pode salvar o JSON
inteiro quando houver campos alterados, em vez de aplicar patches parciais.

### Dados que podem ser persistidos

- CPF normalizado, se necessário para reidratar a tela.
- Nome completo.
- Data de nascimento.
- E-mail.
- Telefone.
- Etapa atual.
- `email_verification_id`.
- `phone_verification_id`.
- Status de e-mail verificado.
- Status de telefone verificado.

### Dados que não devem ser persistidos

- Senha.
- Confirmação de senha.
- `email_verification_token`.
- `phone_verification_token`.
- Tokens de sessão.

### Expiração e retomada

- O rascunho deve ter TTL de 24 horas.
- Ao carregar o rascunho, se `created_at` ou `updated_at` indicarem expiração,
  apagar o rascunho e iniciar o cadastro do zero para aquele CPF.
- Se os tokens de verificação expirarem, limpar os tokens confirmados e voltar o
  fluxo para a etapa de verificação correspondente.
- Se o app for encerrado e reaberto, tokens confirmados não devem ser
  recuperados do rascunho. O usuário deve refazer a confirmação de e-mail ou
  telefone conforme a etapa retomada.
- Se o cadastro for concluído com sucesso, remover o rascunho.
- Se o usuário cancelar explicitamente o cadastro, remover o rascunho.

## Requisitos de autenticação dos endpoints

Todos os endpoints usados nesse fluxo devem funcionar com `app_token`, sem JWT:

- `POST /auth/cpf-check`
- `POST /auth/contact-verifications`
- `POST /auth/contact-verifications/confirm`
- `POST /auth/register`

No estado atual da API, esses endpoints já estão registrados com middleware de
`X-App-Token` em `api/cmd/api/main.go`.

## Escopo mobile

- Substituir a `RegisterPage` única por páginas específicas de cadastro.
- Criar ou reorganizar widgets de formulário para cada etapa.
- Atualizar rotas de autenticação para suportar as páginas do fluxo.
- Manter um único `RegisterViewmodel` compartilhado durante a jornada.
- Registrar ou manter o `RegisterViewmodel` como `lazySingleton` para preservar
  estado entre páginas do cadastro.
- Criar uma base simples e específica para o cadastro, separando estado do
  rascunho, snapshot persistível e store seguro.
- Usar dirty tracking simples no estado do rascunho para evitar persistências
  desnecessárias.
- Não criar ainda um gerenciador genérico de onboarding.
- Evitar que voltar uma página perca os dados já informados.
- Bloquear avanço quando a etapa atual estiver inválida.
- Consultar disponibilidade do CPF antes de avançar da primeira página.
- Persistir rascunho local por hash do CPF enquanto o cadastro não for
  concluído.
- Persistir o rascunho em secure storage.
- Não persistir senha no rascunho.
- Após cadastro bem-sucedido, navegar para uma tela de sucesso e então permitir
  ir para login.
- Atualizar testes de fluxo, viewmodel e navegação.

## Escopo API

Não há alteração funcional obrigatória identificada para viabilizar este fluxo.

Possíveis ajustes de API já aplicados:

- Adicionar ou revisar testes garantindo que os endpoints do cadastro
  exigem `X-App-Token`.
- Adicionar `POST /auth/cpf-check` para consulta de disponibilidade do CPF com
  `X-App-Token`.
- Garantir que a documentação continue clara sobre `app_token` em
  CPF check, contact verification e register.

## Fora de escopo

- Envio real de e-mail.
- Envio real de SMS.
- Recuperação de senha.
- Cadastro parcial persistido no backend antes da verificação de telefone.
- Criação automática de conta bancária operacional.
- Alterar o fluxo de aprovação administrativa pós-cadastro.
- Remover o `RegisterViewmodel` como orquestrador do registro.
- Autenticar automaticamente o usuário após o cadastro.

## Critérios de aceite

- O cadastro não usa mais uma página única com todos os campos.
- Cada página apresenta apenas os campos da etapa correspondente.
- O estado do cadastro permanece consistente ao avançar e voltar.
- O rascunho do cadastro é persistido localmente por chave derivada do CPF.
- O rascunho é persistido em secure storage.
- A senha não é persistida no rascunho.
- Ao informar um CPF com rascunho existente, o fluxo pode ser retomado.
- CPF já cadastrado bloqueia o avanço do fluxo antes da coleta dos demais dados.
- E-mail só é considerado concluído após confirmação do token.
- Telefone só é considerado concluído após confirmação do token.
- A conta do usuário é criada ao final do fluxo, depois da confirmação da senha.
- Após criar a conta, o usuário não fica conectado automaticamente.
- `POST /auth/cpf-check`, `POST /auth/contact-verifications`,
  `POST /auth/contact-verifications/confirm` e `POST /auth/register` enviam
  `X-App-Token`.
- O app continua imprimindo o token retornado no log de depuração em ambiente de
  desenvolvimento.
- Testes cobrem o caminho feliz completo e falhas de validação por etapa.

## Decisões fechadas

- O fluxo será dividido em páginas pequenas, começando por CPF.
- O mesmo `RegisterViewmodel` gerencia todo o cadastro.
- O `RegisterViewmodel` pode ser `lazySingleton`.
- A base inicial será específica do cadastro de usuário.
- Não será criado um gerenciador genérico de onboarding neste incremento.
- O rascunho terá dirty tracking simples e salvará o JSON inteiro quando houver
  alterações.
- O onboarding de cadastro deve ter rascunho local persistido por CPF.
- O rascunho deve ser salvo em secure storage.
- A chave do rascunho deve usar hash do CPF normalizado, não CPF puro.
- O TTL do rascunho local é de 24 horas.
- Senha e confirmação de senha não devem ser persistidas.
- Tokens confirmados de e-mail e telefone não devem ser persistidos.
- Após criar a conta, o usuário não estará conectado.
- Após criar a conta, o app deve exibir uma tela de sucesso com ação para login.

## Decisões pendentes

- Nenhuma.
