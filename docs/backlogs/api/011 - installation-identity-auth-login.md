# Backlog API 011: Installation Identity on `POST /auth/login`

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Tratar a identidade de instalação já no login, exigindo `X-Installation-Id`,
classificando o estado da instalação e decidindo entre:

- sessão operacional normal;
- bootstrap automático da primeira instalação;
- autorização restrita para registrar nova instalação;
- negação por instalação revogada;
- negação por limite de instalações atingido.

## 3. Contrato

O endpoint exige:

```http
POST /auth/login
X-Installation-Id: <UUID v4>
```

Regras consolidadas:

- o header é obrigatório desde o primeiro release da feature;
- formato inválido ou ausência, quando obrigatório, retorna
  `400 INVALID_INSTALLATION_ID`;
- instalação conhecida retorna sessão operacional normal;
- primeira instalação da conta é cadastrada automaticamente de forma atômica;
- instalação revogada é negada;
- nova instalação com vaga disponível retorna autorização restrita;
- nova instalação sem vaga retorna `installation_limit_reached`.

## 4. Respostas esperadas

### Instalação conhecida ou primeira instalação

Retorna sessão operacional normal com `access_token` e `refresh_token`.

### Nova instalação com vaga

```json
{
  "data": {
    "authentication_status": "installation_registration_required",
    "restricted_access_token": "<jwt>",
    "expires_in": 300
  },
  "error": null
}
```

### Limite atingido

```json
{
  "data": {
    "authentication_status": "installation_limit_reached",
    "max_installations": 3,
    "known_installations_count": 3,
    "installation_registered": false,
    "next_action": "revoke_existing_installation"
  },
  "error": null
}
```

## 5. Regras de negócio

- a primeira instalação é evento cronológico, não condição de existir uma
  única instalação ativa;
- instalações revogadas permanecem no histórico e contam como associação
  anterior;
- o bootstrap da primeira instalação deve ser atômico;
- o limite considera apenas instalações `known`.

## 6. Tasks derivadas

- validar `X-Installation-Id` no login;
- classificar instalação por estado e histórico;
- criar bootstrap atômico da primeira instalação;
- emitir `restricted_access_token` quando aplicável;
- retornar `installation_limit_reached` com o contrato final.

## 7. Referências

- [Installation Identity MVP umbrella](<010 - installation-identity-mvp.md>)
- [Installation registration](<013 - installation-identity-register-installation.md>)
- [Shared infrastructure](<018 - installation-identity-shared-infrastructure.md>)
