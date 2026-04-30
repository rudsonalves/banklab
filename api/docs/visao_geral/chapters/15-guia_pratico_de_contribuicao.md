# Guia Prático para Novas Contribuições

Este capítulo final complementa a visão geral da aplicação com uma orientação prática para quem vai começar a modificar a API. O objetivo é transformar o entendimento arquitetural acumulado nos capítulos anteriores em critérios concretos de trabalho.

## Antes de alterar o código

Antes de iniciar uma implementação, vale responder algumas perguntas simples.

### 1. Qual problema esta mudança resolve?

Antes de abrir arquivos, é importante definir se a mudança é:

- nova funcionalidade;
- correção de bug;
- ajuste de contrato HTTP;
- mudança de modelagem;
- ajuste de persistência;
- melhoria de arquitetura;
- reforço de teste ou documentação.

Essa distinção ajuda a identificar desde cedo quais partes do sistema tendem a ser afetadas.

### 2. Qual módulo é dono da mudança?

A próxima pergunta deve ser:

```text
esta responsabilidade pertence a auth, account, customer ou admin?
```

Se a mudança parecer atravessar mais de um módulo, vale distinguir:

- qual módulo é dono da regra principal;
- quais módulos apenas participam do fluxo.

Essa distinção evita que o código comece a crescer em locais convenientes, porém conceitualmente incorretos.

### 3. Em que camada a mudança deve acontecer?

Também é importante identificar se a alteração pertence a:

- `delivery`;
- `application`;
- `domain`;
- `infrastructure`.

Uma heurística útil é:

- se o problema é request, response, rota ou status code, comece por `delivery`;
- se o problema é coordenação do fluxo, comece por `application`;
- se o problema é regra de negócio ou contrato, comece por `domain`;
- se o problema é query, transação ou persistência, comece por `infrastructure`.

## Como começar a explorar uma feature existente

Quando a mudança for sobre um fluxo já implementado, uma boa estratégia é localizar a feature de fora para dentro.

### Passo 1: começar pelo endpoint

Identifique a rota correspondente na documentação REST ou no `main.go`.

Isso ajuda a descobrir:

- qual handler recebe a chamada;
- qual módulo expõe o fluxo;
- quais middlewares estão envolvidos.

### Passo 2: localizar o handler

No handler, observe:

- como a entrada é lida;
- quais validações de transporte acontecem;
- qual caso de uso é chamado;
- como erros e respostas são mapeados.

### Passo 3: localizar o caso de uso

Na `application`, observe:

- qual é a entrada esperada;
- quais pré-condições são validadas;
- quais contratos são usados;
- se existe transação;
- quais regras dependem do contexto autenticado.

### Passo 4: localizar regras e contratos do domínio

No `domain`, observe:

- quais entidades participam;
- quais erros são relevantes;
- quais invariantes o fluxo depende de preservar;
- quais interfaces estruturam a operação.

### Passo 5: localizar a implementação concreta

Na `infrastructure`, observe:

- quais queries SQL são executadas;
- como o dado é persistido ou consultado;
- se existe lock de linha;
- como erros do banco são tratados.

Esse percurso ajuda a entender a feature como sistema, e não apenas como um conjunto de arquivos desconexos.

## Como decidir onde colocar uma nova regra

Uma das decisões mais frequentes no dia a dia é identificar onde uma nova regra deve morar.

### Regra de protocolo ou contrato externo

Se a regra diz respeito a:

- JSON malformado;
- parâmetro obrigatório de rota;
- query param inesperado;
- status code da resposta;
- envelope HTTP;

ela tende a pertencer à `delivery`.

### Regra de coordenação do fluxo

Se a regra diz respeito a:

- ordem das etapas;
- necessidade de transação;
- combinação entre múltiplos repositórios;
- uso de contexto autenticado;
- orquestração entre módulos;

ela tende a pertencer à `application`.

### Regra de domínio

Se a regra diz respeito a:

- validade da operação;
- invariantes da entidade;
- restrições de negócio;
- impossibilidade lógica de determinado comportamento;

ela tende a pertencer ao `domain`.

### Regra de persistência

Se a regra diz respeito a:

- forma de consultar dados;
- forma de persistir;
- lock de linha;
- ordenação de leitura;
- mapeamento entre banco e entidade;

ela tende a pertencer à `infrastructure`.

## Como avaliar impacto de uma mudança

Antes de concluir uma alteração, vale mapear o que ela afeta.

### Contrato externo

Pergunte:

- muda rota?
- muda payload?
- muda código de erro?
- muda status code?

Se sim, a documentação REST e os testes de delivery provavelmente precisam ser atualizados.

### Contexto autenticado e ownership

Pergunte:

- a feature depende de `role`?
- depende de `customer_id`?
- depende de ownership de conta ou customer?

Se sim, é importante revisar também efeitos sobre autorização e contexto autenticado.

### Persistência e consistência

Pergunte:

- a mudança mexe com banco?
- cria ou altera tabela?
- muda transação?
- mexe com saldo?
- mexe com ledger?

Se sim, a mudança precisa considerar consistência, migração e cobertura adequada de testes.

## Como pensar testes ao implementar

Ao escrever ou alterar código, vale decidir desde cedo onde a mudança precisa ser validada.

Uma regra simples pode orientar:

- regra de negócio: teste de domínio;
- coordenação de fluxo: teste de aplicação;
- contrato HTTP: teste de delivery;
- SQL, lock ou integração concreta: teste de integração.

Isso evita dois extremos comuns:

- deixar a mudança sem teste relevante;
- testar em um nível tão distante que a validação perde precisão.

## Como pensar documentação ao implementar

A documentação também deve ser tratada como parte da mudança.

Algumas alterações que normalmente exigem revisão documental:

- novo endpoint;
- mudança de request ou response;
- novo fluxo importante;
- nova decisão de modelagem;
- mudança de regra de ownership;
- alteração relevante em persistência;
- mudança de responsabilidade entre módulos.

Na prática, uma boa pergunta é:

```text
se outro desenvolvedor ler a documentação depois dessa mudança, ele entenderá o sistema corretamente?
```

Se a resposta for “não”, a documentação correspondente precisa ser ajustada.

## Erros comuns a evitar

Ao contribuir no projeto, alguns erros tendem a causar erosão arquitetural mais rapidamente.

### Colocar regra de negócio no handler

Isso costuma acontecer quando a mudança parece pequena, mas acaba deslocando decisão importante para a camada HTTP.

### Resolver tudo diretamente no repositório

Quando a lógica do fluxo é empurrada para SQL ou para a implementação concreta, a aplicação perde clareza e testabilidade.

### Usar `shared/` como depósito genérico

Código compartilhado deve ir para `shared/` apenas quando realmente for transversal. Caso contrário, é melhor mantê-lo no módulo dono da responsabilidade.

### Receber do cliente dados que o backend já conhece

Sempre que possível, o sistema deve continuar usando o contexto autenticado como fonte de identidade e escopo.

### Tratar operações financeiras como simples atualização de registro

Mudanças em saldo, transferência e histórico exigem cuidado com transação, ledger, concorrência e idempotência.

## O que fazer quando a mudança parece atravessar tudo

Algumas features realmente tocam vários pontos do sistema. Nesses casos, o melhor caminho costuma ser:

1. identificar o dono da regra principal;
2. mapear módulos participantes;
3. definir o fluxo na `application`;
4. preservar a propriedade das regras de cada módulo;
5. atualizar testes e documentação junto com a implementação.

Quando uma mudança parece “sem dono”, isso costuma ser sinal de que ainda falta modelar melhor a responsabilidade antes de sair implementando.

## Fluxo recomendado para uma nova contribuição

Um fluxo prático de trabalho pode ser:

1. entender a tarefa e o módulo dono;
2. reler o capítulo ou documento mais próximo da mudança;
3. localizar o endpoint e o caso de uso correspondente;
4. mapear impacto em domínio, persistência, testes e documentação;
5. implementar no nível correto;
6. validar com testes adequados;
7. atualizar a documentação correspondente.

Esse processo tende a gerar mudanças menores, mais claras e menos acopladas.

## Síntese

Contribuir bem para esta API exige mais do que saber escrever código em Go. Exige reconhecer onde cada responsabilidade vive, como identidade e ownership influenciam os fluxos e por que operações financeiras exigem cuidado especial com consistência e persistência.

Este guia prático existe para transformar a visão arquitetural do sistema em critério de implementação. A ideia principal é simples: antes de alterar o código, entender o significado da mudança no domínio ajuda a colocá-la no lugar certo, testá-la no nível certo e documentá-la da forma certa.
