# Backlog API 017: Session enforcement with `X-Installation-Id`

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Aplicar o vínculo da instalação nas requisições autenticadas por
`access_token`.

## 3. Regras consolidadas

- o `access_token` operacional carrega `installation_id`;
- o `X-Installation-Id` é obrigatório em logout, quando exposto, e em todas as
  rotas autenticadas por `access_token`;
- o middleware de sessão operacional deve comparar o header com o claim do
  token;
- a API não reassocia silenciosamente uma sessão a outra instalação;
- divergência exige novo login.

## 4. Erros

- `400 INVALID_INSTALLATION_ID` para header ausente ou malformado, quando
  obrigatório;
- `403 INSTALLATION_MISMATCH` para divergência entre header e instalação da
  sessão.

## 5. Escopo do MVP

No MVP, esse sinal afeta:

- login;
- refresh;
- vínculo de sessão autenticada;
- registro de nova instalação;
- gerenciamento de instalações.

Ainda não afeta:

- transferências;
- demais operações financeiras sensíveis.

## 6. Tasks derivadas

- estender o middleware de sessão operacional para validar `installation_id`;
- propagar o vínculo da instalação para o contexto autenticado;
- negar requisições com divergência entre header e claim;
- cobrir testes de middleware e rotas autenticadas.

## 7. Referências

- [Installation Identity MVP umbrella](<010 - installation-identity-mvp.md>)
- [Refresh binding](<016 - installation-identity-auth-refresh.md>)
- [Shared infrastructure](<018 - installation-identity-shared-infrastructure.md>)
