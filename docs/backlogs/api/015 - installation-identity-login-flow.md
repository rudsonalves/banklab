# Backlog API 015: Installation Identity Login Flow

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Implementar a decisão de identidade de instalação durante o `POST /auth/login`
somente depois da base compartilhada estar operacional.

## 3. Escopo

- Classificar a instalação recebida no login.
- Permitir sessão operacional normal para instalação `known`.
- Executar bootstrap automático da primeira instalação.
- Bloquear instalação `revoked`.
- Retornar `installation_limit_reached` quando o limite estiver atingido.
- Emitir autorização restrita quando a instalação for nova e houver vaga.
- Garantir atomicidade no bootstrap da primeira instalação.
- Garantir atomicidade na decisão de limite.

## 4. Fora de escopo

- Criar tabelas.
- Criar repositórios compartilhados.
- Implementar `POST /security/installations`.
- Consumir step-up.
- Listar ou revogar instalações.

## 5. Dependências

- Backlog 011: contrato de header.
- Backlog 013: domínio e repositórios.
- Backlog 014: sessão, JWT e autorização restrita.

## 6. Orientação para tasks

As tasks deste backlog devem seguir os ramos de decisão do login, mas só depois
da infraestrutura real existir. A classificação pode ser testada isoladamente,
enquanto o handler deve proteger os contratos HTTP finais.

## 7. Referências

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
