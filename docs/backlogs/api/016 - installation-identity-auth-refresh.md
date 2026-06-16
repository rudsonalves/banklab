# Backlog API 016: Installation binding on `POST /auth/refresh`

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Vincular o refresh da sessão à mesma instalação que originou a autenticação
operacional.

## 3. Contrato

```http
POST /auth/refresh
X-Installation-Id: <UUID v4>
```

## 4. Regras consolidadas

- o header é obrigatório desde o primeiro release da feature;
- o refresh deve revalidar o mesmo `X-Installation-Id` da sessão em curso;
- sessões de instalação revogada não renovam tokens;
- divergência entre header e vínculo da sessão deve negar a operação;
- o novo `access_token` emitido continua carregando `installation_id`.

## 5. Erros

- `400 INVALID_INSTALLATION_ID` para header ausente ou malformado, quando
  obrigatório;
- `403 INSTALLATION_MISMATCH` para divergência entre header e sessão;
- negar refresh de instalação revogada.

## 6. Tasks derivadas

- validar header no refresh;
- validar vínculo entre sessão, usuário e instalação;
- impedir renovação de instalação revogada;
- cobrir rotação de refresh token mantendo o vínculo da instalação.

## 7. Referências

- [Installation Identity MVP umbrella](<010 - installation-identity-mvp.md>)
- [Session enforcement](<017 - installation-identity-session-enforcement.md>)
