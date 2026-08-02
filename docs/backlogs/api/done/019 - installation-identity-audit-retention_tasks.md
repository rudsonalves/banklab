# Tasks de Installation Identity Audit, Retention and Operational Docs

Backlog pai: `019 - installation-identity-audit-retention.md`

Campos sugeridos por task:
- Status: Backlog | Em andamento | Concluida
- Dono:
- PR:
- Observacoes:

## Task 1/9: Definir politica de retencao de instalacoes revogadas

Status: Concluida

### Objetivo
Definir por quanto tempo instalacoes revogadas permanecem persistidas e quais
dados podem continuar armazenados apos a revogacao.

### Escopo
- Definir estado final esperado para uma instalacao revogada.
- Definir se registros revogados sao mantidos, anonimizados ou removidos.
- Documentar prazo de retencao e criterio operacional para limpeza.
- Explicitar que a retencao nao transforma `installation_id` em identidade forte.

### Criterios de aceite
- Existe uma politica documentada para instalacoes revogadas.
- A politica diferencia necessidade operacional, auditoria e minimizacao.
- A documentacao descreve claramente o que pode ser apagado ou preservado.
- Nenhum texto sugere confianca forte no `installation_id`.

### Depende de
- Backlogs 011 a 018 com contrato final estabilizado.

## Task 2/9: Definir politica de retencao de autorizacoes restritas

Status: Concluida

### Objetivo
Definir retencao e descarte de autorizacoes restritas consumidas, revogadas ou
expiradas.

### Escopo
- Definir prazo de retencao por estado da autorizacao.
- Definir quais campos podem permanecer apos consumo, revogacao ou expiracao.
- Definir comportamento esperado para tokens ou hashes associados.
- Registrar impacto operacional da limpeza em suporte e auditoria.

### Criterios de aceite
- A politica cobre autorizacoes consumidas, revogadas e expiradas.
- Tokens sensiveis nao permanecem recuperaveis alem do necessario.
- O comportamento de limpeza esta documentado para cada estado.
- A politica preserva somente dados compativeis com auditoria e minimizacao.

### Depende de
- Task 1.

## Task 3/9: Definir eventos obrigatorios de auditoria

Status: Concluida

### Objetivo
Definir os eventos do ciclo de instalacao e autorizacao restrita que devem ser
auditados no MVP.

### Escopo
- Listar eventos de criacao, ativacao, uso, revogacao e expiracao.
- Definir campos minimos de auditoria por evento.
- Definir eventos de falha relevantes para seguranca.
- Definir correlacao operacional sem expor segredos ou atributos excessivos.

### Criterios de aceite
- Existe uma lista fechada de eventos obrigatorios.
- Cada evento possui payload minimo definido.
- Eventos de falha relevantes estao cobertos.
- O contrato de auditoria nao inclui tokens, senha transacional ou fingerprint.

### Depende de
- Task 1.
- Task 2.

## Task 4/9: Validar logs seguros no fluxo de instalacao

Status: Concluida

### Objetivo
Garantir que logs dos fluxos de instalacao, revogacao e autorizacao restrita nao
exponham tokens, senha transacional ou atributos excessivos do ambiente.

### Escopo
- Revisar logs existentes nos handlers, servicos e repositorios envolvidos.
- Remover ou mascarar valores sensiveis quando necessario.
- Garantir mensagens operacionais suficientes sem dados sensiveis.
- Adicionar testes quando houver risco de regressao observavel.

### Criterios de aceite
- Logs nao expõem tokens, hashes sensiveis, senha transacional ou payload bruto.
- Dados de ambiente nao sao registrados em detalhe excessivo.
- Mensagens de erro continuam uteis para operacao e suporte.
- Testes ou revisoes cobrem os pontos alterados.

### Depende de
- Task 3.

## Task 5/9: Validar minimizacao do modelo persistido

Status: Concluida

### Objetivo
Confirmar que o MVP persiste apenas os atributos necessarios para o contrato de
instalacao, autorizacao restrita, auditoria e operacao.

### Escopo
- Revisar tabelas, entidades e DTOs relacionados a instalacoes.
- Revisar campos de autorizacoes restritas e historico operacional.
- Identificar atributos redundantes, excessivos ou ambiguos.
- Propor migracoes de remocao ou ajustes documentais quando necessario.

### Criterios de aceite
- Campos persistidos possuem justificativa funcional, operacional ou de auditoria.
- A revisao nao encontra atributos equivalentes a fingerprinting fora do escopo.
- Ajustes necessarios estao implementados ou registrados como decisao explicita.
- A documentacao final reflete o modelo persistido real.

### Depende de
- Task 1.
- Task 2.
- Task 3.

## Task 6/9: Implementar limpeza operacional de registros expirados

Status: Concluida

### Objetivo
Implementar ou ajustar o mecanismo operacional que remove registros expirados ou
fora da politica de retencao.

### Escopo
- Definir alvo inicial de limpeza conforme politicas aprovadas.
- Implementar migration, job, comando ou rotina ja adotada pelo projeto.
- Garantir que a limpeza seja idempotente.
- Documentar como executar e verificar a limpeza em ambiente operacional.

### Criterios de aceite
- Existe mecanismo versionado para limpeza dos registros cobertos.
- A limpeza respeita as politicas das Tasks 1 e 2.
- A execucao repetida nao causa erro nem remove dados fora do escopo.
- Ha instrucao operacional para verificacao da rotina.

### Depende de
- Task 1.
- Task 2.

## Task 7/9: Documentar efeitos de revogacao sobre sessoes e tokens

Status: Concluida

### Objetivo
Documentar o efeito de revogacao de instalacao e autorizacao restrita sobre
sessoes, tokens ativos e novas operacoes autenticadas.

### Escopo
- Descrever comportamento esperado para sessoes existentes.
- Descrever comportamento esperado para novos tokens ou refresh.
- Descrever efeito sobre operacoes que exigem instalacao ativa.
- Explicitar limites do MVP quando a revogacao nao invalida artefatos ja emitidos.

### Criterios de aceite
- A documentacao cobre sessoes existentes, tokens e novas operacoes.
- Casos de revogacao ficam claros para API, suporte e operacao.
- Limites do MVP estao descritos sem criar promessa de seguranca inexistente.
- O texto esta coerente com o comportamento implementado nos backlogs 011 a 018.

### Depende de
- Backlogs 011 a 018 com contrato final estabilizado.

## Task 8/9: Atualizar documentacao tecnica e operacional final

Status: Concluida

### Objetivo
Consolidar a documentacao tecnica e operacional do MVP de installation identity
apos as decisoes de retencao, auditoria e revogacao.

### Escopo
- Atualizar documentos de arquitetura, operacao e modelo de ameaca aplicaveis.
- Incluir politicas de retencao e limpeza.
- Incluir eventos de auditoria e limites de logging.
- Incluir orientacoes de suporte sem expor dados sensiveis.

### Criterios de aceite
- Documentos relevantes foram atualizados de forma consistente.
- Politicas e operacao estao descritas em linguagem acionavel.
- Nao ha divergencia entre documentacao e comportamento implementado.
- O texto preserva a premissa de identidade fraca do `installation_id`.

### Depende de
- Task 1.
- Task 2.
- Task 3.
- Task 6.
- Task 7.

## Task 9/9: Validar suite e documentacao final do MVP

Status: Concluida

### Objetivo
Executar a validacao final do backlog 019 e garantir que o MVP nao ficou
documentado como um mecanismo de identidade forte.

### Escopo
- Executar testes automatizados relevantes.
- Revisar documentacao final em busca de linguagem ambigua.
- Confirmar ausencia de tokens, senha transacional e atributos excessivos em logs.
- Registrar pendencias explicitas, se houver.

### Criterios de aceite
- Testes relevantes passam ou falhas ficam documentadas com causa.
- A documentacao nao sugere confianca forte no `installation_id`.
- Logs e auditoria seguem as decisoes das tasks anteriores.
- O backlog 019 pode ser marcado como pronto para implementacao ou encerramento.

### Depende de
- Task 4.
- Task 5.
- Task 8.

## Resultado implementado

- Politica de retencao registrada em `api/docs/08-auth_implementation.md`,
  `api/docs/09-database.md` e no backlog pai.
- Limpeza operacional de autorizacoes restritas adicionada por
  `cleanup_installation_registration_authorizations()`, agendada via `pg_cron`
  na migration `000015_installation_authorizations_cleanup`.
- Repositorio Postgres concreto passou a expor `CleanupExpired` para execucao
  testavel da mesma politica sem ampliar a interface usada pelos use cases.
- Teste de integracao cobre remocao de autorizacoes antigas e preservacao de
  autorizacoes recentes.
- Documentacao REST explicita efeitos de revogacao, retencao historica e que
  `installation_id` e apenas sinal contextual fraco.
