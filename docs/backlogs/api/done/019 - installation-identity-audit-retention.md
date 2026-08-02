# Backlog API 019: Installation Identity Audit, Retention and Operational Docs

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Concluido

## 2. Objetivo

Fechar o MVP com política de retenção, auditoria, minimização e documentação
operacional coerentes com o modelo de ameaça.

## 3. Escopo

- Definir retenção de instalações revogadas.
- Definir retenção de autorizações restritas consumidas, revogadas ou
  expiradas.
- Definir quais eventos devem ser auditados.
- Garantir que logs não exponham tokens, senha transacional ou atributos
  excessivos do ambiente.
- Confirmar que o MVP não persiste atributos além do mínimo necessário.
- Documentar efeitos de revogação sobre sessões e tokens.
- Atualizar documentação técnica e operacional.
- Validar que a documentação final não sugere confiança forte no
  `installation_id`.

## 4. Fora de escopo

- Criar painel administrativo.
- Adicionar score antifraude.
- Adicionar geolocalização, biometria, attestation ou device fingerprinting.

## 5. Dependencias

- Backlogs 011 a 018 concluídos ou com contrato final estabilizado.

## 6. Preparacao para tasks

As tasks devem separar política de retenção, auditoria, limpeza operacional e
documentação final.

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)

## 8. Decisoes implementadas

- Instalacoes revogadas permanecem no historico. Elas nao ocupam slot `known`,
  mas continuam impedindo que o usuario volte a ser elegivel ao bootstrap de
  primeira instalacao.
- Autorizacoes restritas ativas expiradas ha mais de 24 horas sao removidas.
- Autorizacoes restritas consumidas sao removidas 24 horas apos `consumed_at`.
- Autorizacoes restritas revogadas sao removidas 24 horas apos `created_at`.
- A limpeza operacional e feita por `cleanup_installation_registration_authorizations()`,
  agendada via `pg_cron` para 03:30.
- A revogacao de instalacao invalida refresh sessions vinculadas a instalacao
  revogada. Access tokens ja emitidos permanecem limitados ao proprio TTL curto
  e ao enforcement de contexto nas rotas protegidas.
- Logs e auditoria devem usar nomes de evento, identificadores publicos,
  status, timestamps e codigos de erro. Tokens, senha transacional, hashes e
  payloads sensiveis nao devem ser registrados.
- O `installation_id` permanece documentado como sinal contextual fraco, nunca
  como identidade forte, fator de autenticacao ou prova de posse do aparelho.
