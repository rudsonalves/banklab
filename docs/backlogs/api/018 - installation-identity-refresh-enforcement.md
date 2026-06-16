# Backlog API 018: Installation Identity Refresh and Enforcement

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Aplicar o vínculo de instalação em refresh e requisições autenticadas depois
que sessão, tokens e gestão de revogação estiverem prontos.

## 3. Escopo

- Exigir `X-Installation-Id` em `POST /auth/refresh`.
- Negar refresh para instalação revogada.
- Negar refresh quando header, sessão e autenticação operacional divergirem.
- Exigir `X-Installation-Id` nas rotas autenticadas por `access_token`.
- Retornar `403 INSTALLATION_MISMATCH` quando houver divergência.
- Garantir que revogação de instalação tenha efeito imediato sobre refresh e
  access tokens.
- Manter o sinal de instalação fora das políticas financeiras sensíveis neste
  primeiro corte.

## 4. Fora de escopo

- Criar confiança de dispositivo.
- Fazer score de risco.
- Exigir step-up adicional para transferências com base apenas na instalação.
- Implementar recuperação de acesso.

## 5. Dependências

- Backlog 014: sessão, JWT e contexto autenticado.
- Backlog 017: revogação de instalação e invalidação de acesso.

## 6. Orientação para tasks

As tasks deste backlog devem ser criadas depois que revogação e sessão com
`installation_id` estiverem funcionais. A ordem deve começar pelo refresh e só
depois avançar para enforcement geral em rotas autenticadas.

## 7. Referências

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
