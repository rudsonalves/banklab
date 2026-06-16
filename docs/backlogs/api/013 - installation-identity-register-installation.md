# Backlog API 013: Register installation on `POST /security/installations`

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Registrar explicitamente uma nova instalação após login com autorização
restrita e `step-up` bem-sucedido.

## 3. Contrato

```http
POST /security/installations
Authorization: Bearer <restricted_access_token>
X-Step-Up-Token: <step_up_token>
X-Installation-Id: <UUID v4>
```

## 4. Regras consolidadas

O endpoint deve:

1. validar a autorização restrita;
2. validar e consumir o `step_up_token`;
3. confirmar que o `X-Installation-Id` corresponde ao apresentado no login;
4. confirmar atomicamente que ainda existe vaga no limite de três;
5. criar a associação como `known`;
6. invalidar a autorização restrita;
7. criar a sessão operacional vinculada à instalação;
8. retornar `access_token` e `refresh_token` operacionais.

Depois do sucesso:

- não existe novo login intermediário;
- o `restricted_access_token` é descartado;
- o `access_token` operacional carrega `installation_id`;
- o `refresh_token` permanece opaco, mas vinculado à sessão da instalação.

## 5. Regras de negação

- `400 INVALID_INSTALLATION_ID` para header ausente ou malformado, quando
  obrigatório;
- bloquear quando a senha transacional estiver `not_set` ou `locked`;
- bloquear quando o limite de instalações tiver sido atingido;
- bloquear quando o `X-Installation-Id` divergir do autorizado no login.

## 6. Tasks derivadas

- validar grant restrito e `step_up_token`;
- criar instalação e sessão operacional em fluxo atômico;
- invalidar grant restrito após sucesso;
- emitir tokens normais vinculados à instalação;
- cobrir cenários de corrida no limite de instalações.

## 7. Referências

- [Installation Identity MVP umbrella](<010 - installation-identity-mvp.md>)
- [Step-up authorize](<012 - installation-identity-step-up-authorize.md>)
- [Shared infrastructure](<018 - installation-identity-shared-infrastructure.md>)
