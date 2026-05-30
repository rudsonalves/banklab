# Tasks do pepper da senha transacional

Backlog pai:

- `008 - transaction-password-pepper.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Area: API
- Tipo: Seguranca/ZTA/Configuracao

## Task 1/6: Adicionar configuracao obrigatoria do pepper

Status: Concluída

### Objetivo

Adicionar `TRANSACTION_PASSWORD_PEPPER` ao carregamento de configuracao da API
com fail-fast no startup.

### Escopo

- Incluir `TransactionPasswordPepper` em `bootstrap.Config`.
- Ler `TRANSACTION_PASSWORD_PEPPER` em `LoadConfig`.
- Encerrar startup com erro fatal se a variavel estiver ausente.
- Encerrar startup com erro fatal se a variavel tiver menos de 32 caracteres.
- Nao reutilizar `JWT_SECRET` ou `APP_TOKEN`.

### Criterios de aceite

- API nao inicia sem `TRANSACTION_PASSWORD_PEPPER`.
- API nao inicia com `TRANSACTION_PASSWORD_PEPPER` menor que 32 caracteres.
- API inicia quando `APP_TOKEN`, `JWT_SECRET` e
  `TRANSACTION_PASSWORD_PEPPER` estao configurados.
- Nome da variavel fica consistente entre codigo e documentacao.

### Depende de

- Nenhuma dependencia.

## Task 2/6: Aplicar pepper no hasher da senha transacional

Status: Concluída

### Objetivo

Atualizar o hasher da senha transacional para aplicar HMAC-SHA256 com pepper
antes do bcrypt.

### Escopo

- Alterar `NewBcryptTransactionPasswordHasher` para receber o pepper.
- Validar que o pepper nao seja vazio no construtor ou no wiring.
- Implementar pre-processamento:
  - `HMAC-SHA256(key=pepper, message=transaction_password)`.
- Codificar o resultado do HMAC em base64 antes de passar ao bcrypt.
- Usar o resultado do HMAC em:
  - `Hash`;
  - `Compare`.
- Manter o bcrypt como hash persistido.
- Nao persistir pepper, PIN ou HMAC intermediario.

### Criterios de aceite

- Hash de senha transacional usa entrada derivada por HMAC.
- Entrada derivada do HMAC e codificada em base64 antes do bcrypt.
- `Compare` valida corretamente com o mesmo pepper.
- `Compare` falha com senha incorreta.
- `Compare` falha quando o hash foi gerado com outro pepper.
- O hash salvo continua sendo um hash bcrypt valido.

### Depende de

- Task 1.

## Task 3/6: Atualizar wiring e testes afetados

Status: Concluída

### Objetivo

Passar o pepper configurado para o hasher e ajustar os testes que constroem o
hasher diretamente.

### Escopo

- Atualizar `cmd/api/main.go`.
- Atualizar testes unitarios do hasher.
- Atualizar testes de application/delivery/integracao que usam
  `NewBcryptTransactionPasswordHasher`.
- Garantir que mocks de hasher nao precisem conhecer o pepper.
- Rodar a suite Go da API.

### Criterios de aceite

- Codigo compila.
- Testes de senha transacional passam.
- Testes de step-up passam.
- `go test ./...` passa na API.

### Depende de

- Task 2.

## Task 4/6: Atualizar setup local e documentacao

Status: Concluída

### Objetivo

Documentar o novo segredo operacional e garantir que ambientes locais consigam
iniciar a API com a nova variavel.

### Escopo

- Atualizar `api/docs/00-getting_started.md`.
- Atualizar `api/README.md` e READMEs relacionados quando citarem variaveis de
  ambiente.
- Atualizar scripts/templates de `.env`, especialmente
  `infra/scripts/ensure-env-files.sh`, se aplicavel.
- Documentar recomendacao:
  - `openssl rand -base64 32`.
- Deixar claro que o valor real nao deve ser commitado.
- Deixar claro que trocar o pepper invalida hashes existentes sem estrategia de
  migracao.

### Criterios de aceite

- Documentacao lista `TRANSACTION_PASSWORD_PEPPER`.
- Setup local novo cria ou orienta criar a variavel.
- Documentacao recomenda segredo aleatorio separado de `JWT_SECRET`.
- Documentacao nao expõe segredo real de producao.

### Depende de

- Task 1.

## Task 5/6: Validar fechamento antes do mobile

Status: Concluída

### Objetivo

Confirmar que a API esta pronta para servir como base do fluxo mobile de senha
transacional e step-up.

### Escopo

- Rodar `go test ./...` na API.
- Revisar:
  - configuracao;
  - hasher;
  - criacao de senha transacional;
  - autorizacao de step-up;
  - documentacao de ambiente.
- Confirmar que o contrato mobile nao mudou.
- Registrar no backlog ou changelog o fechamento da alteracao.

### Criterios de aceite

- API passa nos testes.
- Pepper esta ativo para novos hashes.
- Validacao com pepper incorreto falha em teste.
- Documentacao esta alinhada.
- Implementacao mobile pode prosseguir sem novo contrato publico.

### Depende de

- Task 3.
- Task 4.

## Task 6/6: Verificar alinhamento final antes do mobile

Status: Concluída

### Objetivo

Fechar o backlog garantindo que implementação, testes e documentação estão
alinhados para liberar a continuidade do fluxo mobile.

### Escopo

- Conferir:
  - configuração obrigatória de `TRANSACTION_PASSWORD_PEPPER`;
  - hasher com HMAC-SHA256 + base64 antes do bcrypt;
  - wiring da API com injeção do pepper;
  - criação de senha transacional;
  - autorização de step-up;
  - documentação de setup e variáveis de ambiente.
- Rodar `go test ./...` na API.
- Registrar no backlog e no changelog o fechamento da alteração.

### Critérios de aceite

- `go test ./...` passa.
- Pepper está ativo para novos hashes.
- Validação com pepper incorreto falha em teste.
- Contrato mobile segue inalterado:
  - `POST /security/transaction-password`;
  - `POST /security/step-up/authorize`;
  - campo `transaction_password`.
- Não há divergência conhecida entre código, testes e documentação.

### Depende de

- Task 5.

### Registro de alinhamento final

Data de fechamento: 2026-05-30

Verificações executadas:

- `go test ./...` executado na API com sucesso.
- `TRANSACTION_PASSWORD_PEPPER` obrigatório em `LoadConfig`, com validação de:
  - ausência;
  - tamanho mínimo de 32 caracteres;
  - não reutilização com `APP_TOKEN` e `JWT_SECRET`.
- `BcryptTransactionPasswordHasher` atualizado com:
  - HMAC-SHA256 com key=pepper;
  - codificação base64 da saída do HMAC;
  - uso do valor derivado em `Hash` e `Compare`.
- Testes de infraestrutura cobrem:
  - sucesso com pepper correto;
  - falha com senha incorreta;
  - falha com pepper incorreto.
- Wiring da API atualizado para injetar `config.TransactionPasswordPepper` no
  hasher de senha transacional.
- Documentação e setup local atualizados para incluir
  `TRANSACTION_PASSWORD_PEPPER`.
- Contrato mobile não mudou, mantendo:
  - `POST /security/transaction-password`;
  - `POST /security/step-up/authorize`;
  - campo `transaction_password`.

Conclusão:

- Implementação da API está alinhada com código, testes e documentação, pronta
  para suportar a continuidade do fluxo mobile sem alteração de contrato
  público.
