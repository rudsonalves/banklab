# Backlog: Device Identity MVP

## 1. Status

- Tipo: Research
- Área: Security
- Prioridade: High
- Estado: Discussão

Este backlog inicia a discussão do segundo sinal contextual da evolução Zero
Trust do BankLab. Ele ainda não autoriza implementação: as decisões de
identidade, ciclo de vida e confiança do dispositivo devem ser fechadas antes
da criação das tasks.

## 2. Contexto

O primeiro incremento ZTA protege a transferência interna com senha
transacional e step-up token de uso único. Esse fluxo confirma intenção, mas
trata sessões válidas de ambientes diferentes da mesma forma.

O próximo incremento deve permitir que a API reconheça uma instalação do app
sem transformar o identificador do dispositivo em segredo, fator forte ou
prova de posse.

Já existe no mobile um esqueleto não ativo de `DeviceInterceptor`. Ele deve ser
tratado apenas como referência inicial; o contrato precisa ser definido pelo
backend antes de sua ativação.

## 3. Objetivo

Definir um MVP de identidade de dispositivo que:

- gere um identificador aleatório por instalação do app;
- persista esse identificador no armazenamento seguro do mobile;
- envie o identificador à API em um header explícito;
- associe dispositivos a usuários e sessões;
- permita consultar e revogar dispositivos;
- disponibilize o dispositivo como sinal para políticas futuras;
- preserve o backend como fonte de verdade sobre estado e confiança.

## 4. Princípios de segurança

- O identificador não é autenticação e não substitui JWT ou step-up.
- O valor não deve ser derivado de IMEI, serial, advertising ID ou outro
  identificador físico estável.
- O mobile não decide se o dispositivo é confiável.
- O header pode ser copiado por um cliente comprometido e deve ser tratado como
  sinal fraco.
- Logs não devem expor tokens, senha transacional ou atributos excessivos do
  dispositivo.
- Revogar um dispositivo deve ter efeito definido sobre sessões existentes.
- A política deve distinguir dispositivo conhecido, novo, revogado e ausente.

## 5. Modelo inicial para discussão

```text
devices
- id
- user_id
- installation_id
- status
- platform
- app_version
- first_seen_at
- last_seen_at
- revoked_at
- created_at
- updated_at
```

Status candidatos:

```text
known
revoked
```

`trusted` não deve ser adotado como booleano sem definir quem concede a
confiança, com base em qual evidência e por quanto tempo.

## 6. Contrato inicial para discussão

Header candidato:

```http
X-Device-Id: <installation UUID>
```

Endpoints candidatos:

```http
POST   /security/devices/register
GET    /security/devices
DELETE /security/devices/{device_id}
```

Questões de contrato:

- O registro acontece explicitamente após o login ou de forma idempotente no
  bootstrap da sessão?
- O login deve aceitar `X-Device-Id` antes de existir JWT?
- O identificador público do endpoint deve ser diferente do UUID gerado pelo
  cliente?
- Quais metadados são necessários no MVP além de plataforma e versão do app?
- A ausência do header bloqueia operações ou apenas reduz o nível de confiança?

## 7. Fluxos a definir

### Primeira instalação

```text
App gera installation_id
  -> persiste em secure storage
  -> autentica usuário
  -> registra ou associa instalação ao usuário
  -> API retorna estado do dispositivo
```

### Reinstalação ou limpeza de dados

Uma nova instalação gera uma nova identidade. O sistema não deve tentar
reconstruir silenciosamente a identidade anterior usando fingerprint físico.

### Revogação

Precisamos decidir se a revogação:

- encerra imediatamente todas as sessões associadas;
- impede apenas novas operações sensíveis;
- exige novo step-up em outro dispositivo para recuperação.

## 8. Relação com sessão e policy enforcement

O bootstrap de sessão deve ser considerado o ponto de retorno do estado do
dispositivo, mas o JWT não deve carregar uma confiança imutável durante toda a
sua validade.

Uma política futura poderá avaliar:

```text
Evaluate(user, session, device, operation)
  -> allow
  -> require_step_up
  -> deny
```

No primeiro corte, o dispositivo deve ser coletado e registrado antes de ser
obrigatório para transferências. Isso permite observar o contrato e os estados
sem introduzir bloqueios prematuros no fluxo financeiro.

## 9. Decisões necessárias antes das tasks

- [ ] Definir geração e persistência do `installation_id`.
- [ ] Definir header e formato do identificador.
- [ ] Definir vínculo entre usuário, dispositivo e sessão.
- [ ] Definir estados e transições do dispositivo.
- [ ] Definir comportamento de registro idempotente.
- [ ] Definir efeitos da revogação sobre refresh e access tokens.
- [ ] Definir dados retornados no bootstrap da sessão.
- [ ] Definir retenção, auditoria e minimização de metadados.
- [ ] Definir estratégia de compatibilidade para clientes sem o header.
- [ ] Definir em qual etapa o sinal passa a afetar operações sensíveis.

## 10. Fora de escopo

- biometria local como prova para o backend;
- attestation de plataforma;
- device fingerprinting;
- prova de vida;
- geolocalização;
- score antifraude;
- notificações push para aprovação;
- confiança permanente concedida somente pelo identificador;
- painel administrativo de dispositivos.

## 11. Critérios para encerrar a discussão

- Modelo de ameaça e limitações do sinal documentados.
- Contrato HTTP definido entre API e mobile.
- Ciclo de vida e revogação definidos.
- Relação com `user_sessions` definida.
- Estratégia de rollout sem bloqueio definida.
- Backlogs de implementação da API e do mobile derivados desta decisão.

## 12. Referências internas

- [ZTA MVP - fundação e decisões](<done/006 - zta-mvp-foundation.md>)
- [Auth session bootstrap](<done/009 - auth-session-bootstrap.md>)
- [Discussão de segurança transacional](../discussion.md)
- [Roadmap](../../ROADMAP.md)
