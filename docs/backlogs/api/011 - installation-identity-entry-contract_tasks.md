# Tasks do contrato de entrada da identidade de instalacao

Backlog pai:

- `011 - installation-identity-entry-contract.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Area: API
- Tipo: Contrato/Seguranca/Auth

## Task 1/6: Centralizar constante do header de instalacao

Status: Concluída

### Objetivo

Definir uma constante compartilhada para o header `X-Installation-Id`.

### Escopo

- Adicionar `InstallationID` no pacote compartilhado de headers, ou equivalente.
- Reutilizar o padrao ja existente para headers como `X-App-Token` e
  `X-Step-Up-Token`.
- Evitar strings literais duplicadas em handlers, middlewares e testes.

### Criterios de aceite

- O nome publico do header fica centralizado.
- `POST /auth/login` usa a constante compartilhada.
- Testes usam a mesma constante quando fizer sentido.

### Depende de

- Nenhuma dependencia.

## Task 2/6: Definir erro `INVALID_INSTALLATION_ID`

Status: Concluída

### Objetivo

Criar o erro de contrato retornado quando o header estiver ausente ou mal
formatado.

### Escopo

- Adicionar sentinel error para instalacao invalida.
- Adicionar codigo `INVALID_INSTALLATION_ID`.
- Registrar o erro no mapper compartilhado.
- Retornar HTTP 400.
- Usar mensagem estavel:
  `X-Installation-Id must be a canonical UUID v4.`

### Criterios de aceite

- Ausencia do header retorna `400 INVALID_INSTALLATION_ID`.
- Formato invalido retorna `400 INVALID_INSTALLATION_ID`.
- O envelope segue o padrao `{ data, error }`.
- O erro nao e confundido com `INVALID_REQUEST`.

### Depende de

- Task 1.

## Task 3/6: Implementar validacao canonica de UUID v4

Status: Concluída

### Objetivo

Validar o valor recebido em `X-Installation-Id` antes de chamar o use case de
login.

### Escopo

- Remover espacos laterais antes da validacao.
- Rejeitar valor vazio.
- Rejeitar UUID invalido.
- Rejeitar UUID que nao seja versao 4.
- Rejeitar representacao nao canonica.
- Aceitar somente formato canonico com letras minusculas e hifens.

### Criterios de aceite

- UUID v4 canonico e aceito.
- UUID v1/v3/v5 e rejeitado.
- UUID sem hifens e rejeitado.
- UUID em uppercase e rejeitado.
- Valor com lixo ou espacos internos e rejeitado.

### Depende de

- Task 2.

## Task 4/6: Propagar `installation_id` para a camada de aplicacao

Status: Concluída

### Objetivo

Garantir que o login receba o `installation_id` validado sem ainda interpretar
se a instalacao e conhecida, nova ou revogada.

### Escopo

- Adicionar `InstallationID` ao input do use case de login.
- Converter o header validado para o tipo interno usado pela aplicacao.
- Chamar o use case com o valor validado.
- Nao alterar tokens, sessao, refresh ou regras de elegibilidade existentes.

### Criterios de aceite

- Handler passa o UUID validado para `LoginUserInput`.
- Use case continua funcionando sem consultar instalacoes.
- Nenhuma classificacao de instalacao e executada neste backlog.
- Nenhum token restrito e emitido neste backlog.

### Depende de

- Task 3.

## Task 5/6: Cobrir testes de handler e use case

Status: Concluída

### Objetivo

Proteger o contrato minimo com testes automatizados.

### Escopo

- Testar login com header valido.
- Testar ausencia de `X-Installation-Id`.
- Testar formato invalido.
- Testar UUID nao v4.
- Testar representacao nao canonica.
- Testar que o use case nao e chamado quando o header e invalido.
- Testar que o valor validado chega ao input do use case.

### Criterios de aceite

- Casos invalidos retornam `400 INVALID_INSTALLATION_ID`.
- Casos invalidos nao executam login.
- Caso valido propaga o UUID esperado.
- Testes nao criam repositorio de instalacao falso.

### Depende de

- Task 4.

## Task 6/6: Atualizar documentacao do contrato minimo

Status: Concluída

### Objetivo

Registrar que `POST /auth/login` exige `X-Installation-Id` desde este corte.

### Escopo

- Atualizar documentacao REST afetada.
- Documentar formato esperado do header.
- Documentar resposta de erro `INVALID_INSTALLATION_ID`.
- Deixar explicito que refresh e rotas autenticadas ainda nao exigem o header
  neste backlog.

### Criterios de aceite

- Documentacao de login lista `X-Installation-Id`.
- Documentacao mostra erro 400 para ausencia/formato invalido.
- Documentacao nao antecipa classificacao, bootstrap, token restrito ou
  enforcement de sessao.

### Depende de

- Task 5.
