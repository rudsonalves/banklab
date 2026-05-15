# Discussão: evolução do onboarding

## 1. Contexto

Hoje o BankLab cadastra poucos dados do usuário. O fluxo atual cobre uma base mínima para autenticação e criação do customer, mas ainda não representa um onboarding mais próximo de uma aplicação financeira real.

A evolução do onboarding deve permitir coletar dados cadastrais, dados de contato, endereço, documentos pessoais e evidências documentais para KYC. Antes de implementar, é importante discutir o modelo com cuidado, porque essa decisão impacta:

- banco de dados;
- modelo de domínio;
- contratos da API;
- fluxo mobile;
- integração futura com serviços de terceiros;
- processo de aprovação do usuário;
- base futura para controles de segurança e Zero Trust Architecture.

Este documento registra a ideia inicial e separa a discussão em seções para que cada parte seja amadurecida antes da implementação.

A reestruturação cadastral prévia está descrita em [000 - pre-onboarding.md](<000 - pre-onboarding.md>). Este documento assume que essa fase preparatória será tratada antes do onboarding por checkpoints.

## 2. Objetivo da discussão

Definir um modelo de onboarding simples o bastante para ser implementado em etapas, mas organizado o suficiente para não bloquear evoluções futuras.

O objetivo não é implementar KYC completo agora. Neste primeiro momento, o KYC será tratado como captura e armazenamento de evidências enviadas pelo usuário. Integrações com provedores externos ficam para uma fase posterior.

## 3. Visão geral do onboarding

O onboarding será dividido em dois fluxos principais.

### Fluxo 1: Registro inicial e onboarding cadastral

Este fluxo representa a experiência inicial de cadastro percebida pelo usuário.

Do ponto de vista técnico, ele possui duas partes:

- registro inicial do usuário/customer, tratado no backlog de pre-onboarding;
- onboarding editável, composto por documentos, endereço e KYC.

Campos e informações do onboarding editável inicialmente considerados:

- documentos estruturados, como CPF, RG ou CNH;
- endereço principal;
- evidências documentais para KYC.

O endereço deve incluir:

- CEP;
- rua/logradouro;
- bairro;
- número;
- complemento;
- cidade;
- estado.

A consulta de CEP deve ser feita por meio do backend. O app mobile não deve acessar diretamente o serviço externo.

Os dados retornados pela consulta de CEP preenchem apenas parte do endereço. Campos como número e complemento continuam sendo informados pelo usuário.

### Fluxo 2: KYC documental

Este fluxo representa as evidências enviadas pelo usuário.

Arquivos inicialmente considerados:

- selfie;
- imagem do documento;
- comprovante de residência.

Neste primeiro momento, o sistema apenas recebe, armazena e organiza esses arquivos. Validação automática, consulta a terceiros e análise avançada ficam fora do primeiro incremento.

## 4. Premissas iniciais

- O backend deve ser a fonte de verdade sobre o estado do onboarding.
- O frontend pode persistir rascunhos locais quando fizer sentido, mas não decide sozinho que uma etapa foi concluída.
- Serviços de terceiros devem ser acessados pelo backend, não diretamente pelo mobile.
- A consulta de CEP pode usar um serviço gratuito no primeiro momento, apenas para ilustrar o padrão de integração via backend.
- Futuras integrações de KYC devem seguir a mesma ideia: o backend gerencia o provedor externo e expõe uma API própria para o cliente.
- Verificação de email e telefone deve ser pensada como integração externa, mas será simulada no primeiro momento.
- Enquanto serviços reais de email/SMS não forem definidos, os tokens de verificação poderão ser retornados nas respostas da API backend apenas para facilitar desenvolvimento e validação local.
- O primeiro desenho deve evitar internacionalização completa antes da hora, mas também não deve prender o projeto exclusivamente a CPF de forma difícil de evoluir.

## 5. Seções de decisão

As seções abaixo devem ser discutidas separadamente antes da implementação.

## 6. Seção A: Modelo de estado do onboarding

### Pergunta principal

Como o backend deve representar o progresso do usuário no onboarding?

### Direção inicial

O backend deve controlar o progresso do onboarding por meio de checkpoints.

Esse controle deve indicar:

- etapa corrente;
- etapas já concluídas;
- pendências da etapa atual;
- estado geral do onboarding.

Para o primeiro modelo, as etapas imaginadas são:

1. **Contato / registro inicial**
   - nome completo;
   - email;
   - telefone.

2. **Documentos**
   - CPF;
   - RG ou CNH.

3. **Endereço**
   - CEP;
   - logradouro;
   - bairro;
   - cidade;
   - estado;
   - número;
   - complemento.

4. **KYC**
   - selfie;
   - documento;
   - comprovante de residência, como conta de energia ou documento equivalente.

O frontend deve consultar o backend para saber qual etapa renderizar. O estado oficial do onboarding não deve ser inferido apenas por dados locais do app.

Para o usuário, registro e onboarding devem parecer um fluxo contínuo de cadastro. Para o backend, porém, contato pertence ao registro inicial e não deve ser tratado como etapa editável do onboarding.

### Decisões iniciais

- O onboarding será sequencial.
- Todas as etapas são obrigatórias.
- O backend deve separar consulta simples de status e consulta detalhada do checkpoint.
- `GET /onboarding/status` deve retornar apenas o status geral do onboarding.
- `GET /onboarding` deve retornar o checkpoint completo, com etapa atual e etapas concluídas.
- O app mobile/web deve enviar cadastros completos por etapa. O endpoint de cada etapa não deve aceitar atualização parcial.
- Uma etapa concluída não volta para pendente por alteração posterior dentro do fluxo normal.
- O usuário não pode avançar sem concluir a etapa atual.
- O KYC depende obrigatoriamente da conclusão das etapas anteriores.
- O backend não deve permitir pular etapas em ambiente de desenvolvimento.
- Contato certificado não pode ser alterado diretamente pelo fluxo comum de onboarding.
- O registro inicial do usuário/customer é tratado no backlog de pre-onboarding.
- Dados usados para criação do usuário/customer não entram no fluxo editável do onboarding.
- Contato deve aparecer como etapa concluída no checkpoint, mas não deve possuir endpoint de edição dentro do onboarding.
- Troca de senha, email ou telefone deve ser tratada em outro fluxo da aplicação, fora deste onboarding.
- A aprovação administrativa acontece somente após conclusão do onboarding.
- A aprovação administrativa inicial pode continuar usando o endpoint atual de aprovação por `user_id`.
- A aprovação administrativa deve mudar tanto o status do onboarding quanto o status do usuário.
- A conta bancária só deve ser criada após aprovação administrativa.
- Um painel administrativo para revisão/aprovação deve ser pensado em etapa futura.

### Representação inicial sugerida

Consulta simples de status:

```http
GET /onboarding/status
```

Resposta:

```json
{
  "data": {
    "status": "submitted"
  },
  "error": null
}
```

Esse endpoint deve ser usado em fluxos que só precisam decidir se o usuário pode acessar o app ou qual mensagem deve ser apresentada.

Se o status não for `approved`, o mobile não deve permitir acesso à área autenticada principal e deve apresentar uma mensagem adequada ao estado do onboarding.

Quando o usuário já tiver criado conta, mas ainda não tiver onboarding aprovado, ele deve conseguir autenticar para continuar o onboarding. Esse login não deve liberar acesso normal ao app.

Uma resposta possível para login bloqueado por onboarding ainda não aprovado:

```json
{
  "data": null,
  "error": {
    "code": "ONBOARDING_NOT_APPROVED",
    "message": "Onboarding not approved",
    "details": {
      "onboarding_status": "submitted"
    }
  }
}
```

Essa resposta informa o motivo do bloqueio, mas ainda precisa ser definido como o usuário autenticado continuará o onboarding sem receber uma sessão completa da área principal.

Consulta detalhada do checkpoint:

```http
GET /onboarding
```

Resposta:

```json
{
  "data": {
    "current_step": "address",
    "completed_steps": ["contact", "documents"]
  },
  "error": null
}
```

Estados gerais inicialmente considerados:

- `submitted`: onboarding concluído e aguardando aprovação;
- `approved`: cadastro aprovado;
- `rejected`: cadastro rejeitado;
- `revision`: cadastro exige revisão ou ajuste antes de aprovação.

O KYC pode exigir um status próprio por evidência ou por etapa no futuro. Essa necessidade deve ser discutida separadamente na seção de KYC documental, para não complicar o contrato inicial de status.

### Pontos para discutir

- O login de usuário com onboarding não aprovado deve emitir um token/sessão limitada para continuar o onboarding?
- Essa sessão limitada deve permitir apenas endpoints de onboarding e consulta de status?
- O erro `ONBOARDING_NOT_APPROVED` deve carregar também a etapa atual ou apenas o status geral?
- `GET /onboarding/status` ainda é necessário se o login já retornar o status no erro?
- Qual mensagem o mobile deve mostrar para cada status: `submitted`, `revision`, `rejected` e ausência de onboarding completo?

### Resultado esperado

- Definição dos estados ou checkpoints.
- Regras de transição entre etapas.
- Contrato inicial para consulta simples de status.
- Contrato inicial para consulta detalhada do checkpoint.
- Regras para submissão, aprovação, rejeição e eventual correção.
- Definição do contrato de login quando o onboarding ainda não estiver aprovado.

## 7. Seção B: Dados pessoais e contato

### Pergunta principal

Quais dados pertencem ao usuário de autenticação e quais pertencem ao customer?

### Direção inicial

Email e telefone devem ser tratados como contatos verificáveis.

No produto final, a verificação deve acontecer por serviços de terceiros, como envio de email ou SMS com token de confirmação. Como esses serviços podem ser pagos, exigir conta externa ou depender de cotas de desenvolvimento, a primeira versão deve usar uma solução simples para ambiente local: o backend gera o token de verificação e o retorna na própria resposta da API.

Essa abordagem é apenas um mecanismo de desenvolvimento. Ela permite validar o fluxo de verificação sem acoplar o projeto cedo demais a um provedor específico.

### Decisões iniciais

- As decisões de modelagem de `users`, `customers`, telefone, email e campos de verificação pertencem ao backlog de pre-onboarding.
- Para o onboarding, contato é considerado uma etapa já concluída após o registro inicial.
- Contato não possui endpoint de edição dentro do onboarding.
- Alteração futura de email ou telefone deve ser tratada fora deste fluxo.

### Pontos para discutir

- Como o checkpoint do onboarding deve representar a etapa `contact` após o registro inicial?
- O checkpoint deve sempre retornar `contact` em `completed_steps` para usuários já registrados?
- O mobile deve exibir contato como etapa concluída ou apenas iniciar visualmente em documentos?

### Resultado esperado

- Definição de como a etapa `contact` aparece no checkpoint.
- Definição de como o mobile apresenta a continuidade entre registro inicial e onboarding editável.

## 8. Seção C: Documento de identificação

### Pergunta principal

Como representar CPF, RG, CNH e possíveis documentos futuros sem complicar o modelo inicial?

### Decisões iniciais

- A modelagem de `customer_documents` pertence ao backlog de pre-onboarding.
- Este backlog trata de como o onboarding coleta e valida os dados estruturados dos documentos.
- CPF continua obrigatório para o contexto brasileiro inicial.
- RG ou CNH entra como documento adicional do onboarding.
- Imagens dos documentos pertencem ao KYC documental, não ao cadastro estruturado do documento.

### Pontos para discutir

- Como deixar caminho aberto para passaporte ou documento estrangeiro no futuro?
- Quais campos de RG/CNH são realmente necessários no primeiro incremento?

### Resultado esperado

- Campos documentais exigidos pelo onboarding.
- Definição do que é dado cadastral e do que é evidência KYC.

## 9. Seção D: Endereço principal

### Pergunta principal

Como modelar o endereço do customer neste primeiro incremento?

### Decisões iniciais

- A modelagem de `customer_addresses` pertence ao backlog de pre-onboarding.
- Este backlog trata de como o onboarding coleta e valida o endereço principal.
- O primeiro incremento do onboarding deve trabalhar com um endereço principal.
- O modelo preparado no pre-onboarding deve permitir evolução futura para múltiplos endereços.

### Pontos para discutir

- Quais campos mínimos entram no endereço?
- O que vem preenchido pela consulta de CEP e o que o usuário informa manualmente?
- Como lidar com complemento, número e bairro?

### Resultado esperado

- Campos obrigatórios e opcionais.
- Separação entre campos vindos da consulta de CEP e campos informados manualmente pelo usuário.

## 10. Seção E: Consulta de CEP via backend

### Pergunta principal

Como o BankLab deve integrar um serviço externo de CEP sem acoplar o mobile ao provedor?

### Pontos para discutir

- Qual endpoint o backend deve expor?
- A consulta pertence ao módulo de onboarding, customer ou a um serviço compartilhado de endereço?
- Qual serviço gratuito usar inicialmente?
- Como normalizar a resposta?
- Como tratar CEP inválido, timeout, indisponibilidade ou resposta incompleta?
- Deve haver cache?
- Como registrar falhas de integração sem expor detalhes ao cliente?
- Quais campos retornados pelo serviço de CEP devem preencher o formulário de endereço?
- Como deixar claro que número e complemento não vêm do serviço de CEP e devem ser informados pelo usuário?

### Resultado esperado

- Contrato inicial do endpoint de consulta de CEP.
- Decisão de módulo/responsabilidade.
- Estratégia inicial para integração com serviço externo gratuito.

## 11. Seção F: KYC documental inicial

### Pergunta principal

Como receber e armazenar evidências de KYC sem integrar provedor externo neste primeiro momento?

### Pontos para discutir

- Quais arquivos são obrigatórios?
- Quais tipos de arquivo são aceitos?
- Qual limite de tamanho?
- Onde armazenar os arquivos inicialmente?
- O banco deve guardar apenas metadados e referência de storage?
- Quais status um arquivo pode ter?
- Como diferenciar selfie, documento e comprovante de residência?
- Deve existir revisão manual no primeiro momento?

### Resultado esperado

- Modelo inicial para arquivos de KYC.
- Tipos de evidência aceitos.
- Estratégia de armazenamento.
- Estados iniciais de revisão.

## 12. Seção G: Integração futura com KYC externo

### Pergunta principal

Como preparar o modelo para provedores externos sem implementar integração real agora?

### Pontos para discutir

- Que informações precisam ser guardadas para futura integração?
- Como representar status externo de análise?
- Como manter o backend como intermediador do provedor?
- Como evitar que o mobile conheça detalhes do provedor?
- Como lidar com provedores diferentes no futuro?
- Como registrar auditoria das decisões recebidas?

### Resultado esperado

- Premissas para integração futura.
- Campos ou tabelas que não precisam entrar agora.
- Limites claros entre fase atual e fase futura.

## 13. Seção H: API e contratos

### Pergunta principal

Quais endpoints serão necessários para o primeiro incremento do onboarding?

### Pontos para discutir

- Endpoint para consultar status/checkpoint atual.
- Endpoint para enviar dados pessoais.
- Endpoint para enviar telefone.
- Endpoint para consultar CEP.
- Endpoint para salvar endereço.
- Endpoint para salvar dados de documento.
- Endpoint para upload de evidências KYC.
- Endpoint para finalizar etapa ou submeter para revisão.
- Como versionar ou evoluir contratos sem quebrar o mobile?

### Resultado esperado

- Lista inicial de endpoints.
- Formato geral dos requests e responses.
- Definição de autenticação exigida em cada etapa.

## 14. Seção I: Mobile e persistência local

### Pergunta principal

Como o app mobile deve participar do onboarding sem virar fonte de verdade do processo?

### Pontos para discutir

- Quais dados podem ser salvos localmente como rascunho?
- Quais dados não devem ser persistidos localmente?
- Como retomar um fluxo interrompido?
- Como o mobile decide qual tela renderizar a partir do status do backend?
- Como tratar perda de conexão durante o onboarding?
- Como lidar com upload interrompido?

### Resultado esperado

- Regras de persistência local.
- Contrato entre status backend e telas mobile.
- Estratégia de retomada de fluxo.

## 15. Seção J: Compatibilidade com o fluxo atual

### Pergunta principal

Como introduzir o onboarding por checkpoints sem quebrar o fluxo atual de registro, login e aprovação administrativa?

### Pontos para discutir

- O usuário recém-registrado entra em qual estado de onboarding?
- A aprovação administrativa atual muda ou permanece até o novo fluxo ficar pronto?
- Como o login deve se comportar para usuários ainda não aprovados?
- Como o mobile atual deve lidar com usuários em onboarding?
- Como manter testes e documentação alinhados durante a transição?

### Resultado esperado

- Plano de compatibilidade com o fluxo atual.
- Ordem sugerida de implementação.

## 16. Fora do escopo do primeiro incremento

- Integração real com provedores de KYC.
- Validação automática de documento.
- Prova de vida com provedor externo.
- Internacionalização completa.
- Suporte completo a múltiplos documentos por cliente.
- Suporte completo a múltiplos endereços por cliente.
- Painel administrativo completo para revisão de KYC.
- Redesenho total do fluxo de autenticação.
- Implementação de Zero Trust Architecture.

## 17. Resultado esperado da discussão

Ao final da discussão, o projeto deve ter:

- modelo de estado do onboarding;
- campos de documentos exigidos pelo onboarding;
- campos de endereço exigidos pelo onboarding;
- contrato de consulta de CEP via backend;
- modelo inicial de evidências KYC;
- lista inicial de endpoints;
- regras de persistência local no mobile;
- plano de compatibilidade com o fluxo atual;
- definição clara do que fica para fases futuras.
