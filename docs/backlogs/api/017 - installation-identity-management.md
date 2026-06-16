# Backlog API 017: Installation Identity Management

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Implementar a gestão das instalações associadas ao usuário autenticado:
listagem e revogação lógica.

## 3. Escopo

### `GET /security/installations`

- Listar instalações do usuário autenticado.
- Expor identificador público de gerenciamento separado de `installation_id`.
- Refletir estados `known` e `revoked`.
- Retornar apenas os metadados mínimos definidos para o MVP.

### `DELETE /security/installations/{installation_resource_id}`

- Não exigir step-up no MVP.
- Permitir revogar apenas outra instalação do mesmo usuário.
- Impedir revogação da instalação vinculada à sessão atual.
- Aplicar revogação lógica com `revoked_at`.
- Preservar histórico.
- Invalidar imediatamente refresh tokens da instalação revogada.
- Fazer access tokens já emitidos deixarem de valer na hora.

## 4. Fora de escopo

- Registro de nova instalação.
- Recuperação de acesso quando todas as instalações anteriores estiverem
  indisponíveis.
- Painel administrativo de instalações.

## 5. Dependências

- Backlog 013: domínio e repositórios.
- Backlog 014: vínculo de sessão com instalação.

## 6. Orientação para tasks

As tasks deste backlog devem separar listagem, autorização de propriedade,
bloqueio de remoção da instalação atual e invalidação imediata de acesso.

## 7. Referências

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
