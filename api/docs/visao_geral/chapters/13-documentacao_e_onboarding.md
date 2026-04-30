# Documentação, Onboarding e Uso desta Visão Geral

Este capítulo explica o papel da documentação dentro do projeto, como esta visão geral se encaixa no conjunto de materiais existentes e de que forma novos desenvolvedores podem usar esse conteúdo para acelerar entendimento e contribuir com mais segurança.

## Documentação como parte da arquitetura

Na API atual, a documentação não deve ser entendida como um anexo opcional ao código. Ela faz parte da forma como o sistema é comunicado, estudado e mantido.

Isso é especialmente importante porque a aplicação já possui:

- múltiplos módulos;
- separação em camadas;
- decisões de modelagem relevantes;
- fluxos financeiros com cuidado transacional;
- contrato HTTP consumido por outro cliente do monorepo.

Sem documentação, grande parte desse entendimento ficaria dispersa entre arquivos, testes, handlers, casos de uso, queries e decisões implícitas no código. A consequência seria maior custo de onboarding e maior risco de mudanças incoerentes.

## Papel desta visão geral

Esta pasta `api/docs/visao_geral/` existe para oferecer uma leitura orientada da aplicação no estágio atual.

Seu papel não é substituir:

- o código;
- a documentação REST detalhada;
- a documentação de banco;
- os fluxos específicos de caso de uso.

Seu papel é servir como camada de entendimento entre visão arquitetural ampla e implementação concreta.

Em outras palavras, esta visão geral ajuda um novo desenvolvedor a responder:

- o que a aplicação faz hoje?
- como ela está organizada?
- quais decisões importam para novas implementações?
- onde cada responsabilidade tende a morar?
- como os fluxos principais se relacionam?

## Relação com os demais documentos

O projeto já possui outros documentos técnicos relevantes em `api/docs/`.

Cada um deles cumpre um papel diferente.

### `ARCHITECTURE.md`

Esse documento oferece a visão arquitetural macro da API.

Ele é útil para entender:

- estilo arquitetural;
- módulos principais;
- direção de dependências;
- princípios gerais do sistema.

### `02-use_case_flows.md`

Esse documento descreve fluxos operacionais de caso de uso.

Ele é útil quando o objetivo é entender, de forma mais direta, a sequência de passos de uma operação específica.

### `06-implementation.md`

Esse documento descreve a implementação atual com foco mais técnico e estrutural.

Ele é útil para quem precisa relacionar arquitetura e código existente.

### `07-api-rest.md`

Esse documento descreve o contrato HTTP da API.

Ele é a principal referência para:

- rotas;
- payloads;
- autenticação;
- exemplos de request e response;
- códigos de erro.

### `08-auth_implementation.md`

Esse documento aprofunda aspectos da autenticação e da sessão.

Ele é particularmente útil quando a mudança envolve tokens, refresh, middleware de autenticação ou identidade do usuário.

### `09-database.md`

Esse documento detalha o modelo persistido.

Ele é importante quando a mudança envolve:

- tabelas;
- relações;
- significado de campos persistidos;
- impacto de evolução de schema.

## Como esta visão geral deve ser usada

Esta visão geral funciona melhor como documento de entrada e orientação.

Uma leitura recomendada para um novo desenvolvedor seria:

1. ler a introdução e os capítulos iniciais desta visão geral;
2. entender o estilo arquitetural e a organização dos módulos;
3. estudar os principais fluxos já implementados;
4. só então aprofundar nos documentos específicos da área onde pretende atuar;
5. por fim, navegar no código com base nesse mapa mental.

Esse uso em camadas é importante porque reduz a chance de o desenvolvedor entrar diretamente em arquivos específicos sem contexto suficiente para entender por que determinada regra está naquele lugar.

## Roteiro prático de onboarding

Para alguém que acabou de chegar ao projeto, um roteiro útil pode ser:

### 1. Entender o escopo

Ler:

- `00-visao_geral.md`
- `01-objetivo_e_escopo.md`
- `02-estilo_arquitetural.md`

Objetivo:

entender o que a aplicação é, o que ela resolve e como sua arquitetura foi pensada.

### 2. Entender o mapa do código

Ler:

- `03-organizacao_do_projeto.md`
- `04-anatomia_dos_modulos.md`
- `11-startup_e_wiring.md`

Objetivo:

entender onde as coisas estão no repositório e como se conectam em runtime.

### 3. Entender o comportamento externo

Ler:

- `06-superficie_rest.md`
- `07-autenticacao_e_contexto.md`
- `07-api-rest.md`

Objetivo:

entender como o cliente usa a API e como identidade e autorização influenciam os fluxos.

### 4. Entender os fluxos mais importantes

Ler:

- `08-principais_fluxos.md`
- `02-use_case_flows.md`

Objetivo:

entender como o sistema se comporta nas operações mais relevantes.

### 5. Entender modelagem, persistência e testes

Ler:

- `09-decisoes_de_modelagem.md`
- `10-persistencia_e_consistencia.md`
- `12-testes_e_validacao.md`
- `09-database.md`

Objetivo:

entender o que precisa ser preservado quando a mudança toca domínio, banco ou fluxo financeiro.

## Como usar a documentação ao implementar uma feature

Uma forma prática de usar a documentação no dia a dia é partir da pergunta concreta da tarefa.

### Se a dúvida for sobre contrato externo

Consultar:

- `06-superficie_rest.md`
- `07-api-rest.md`

### Se a dúvida for sobre onde implementar a lógica

Consultar:

- `02-estilo_arquitetural.md`
- `03-organizacao_do_projeto.md`
- `04-anatomia_dos_modulos.md`

### Se a dúvida for sobre identidade, ownership ou autorização

Consultar:

- `07-autenticacao_e_contexto.md`
- fluxos relevantes em `08-principais_fluxos.md`

### Se a dúvida for sobre modelagem

Consultar:

- `09-decisoes_de_modelagem.md`
- `09-database.md`

### Se a dúvida for sobre consistência e persistência

Consultar:

- `10-persistencia_e_consistencia.md`
- `06-implementation.md`

### Se a dúvida for sobre cobertura de testes

Consultar:

- `12-testes_e_validacao.md`

Esse uso orientado por pergunta torna a documentação mais útil do que uma leitura linear obrigatória em toda situação.

## Documentação viva

Um princípio importante deste projeto é tratar a documentação como material vivo.

Isso significa que, quando o sistema evolui de forma relevante, a documentação correspondente também deve evoluir.

Exemplos de mudança que normalmente exigem atualização documental:

- criação de novo endpoint;
- mudança de payload ou resposta;
- alteração de regra de ownership;
- mudança em modelagem de status;
- introdução de nova política de conta;
- alteração de fluxo transacional;
- mudança de schema relevante;
- reorganização de responsabilidade entre módulos.

Esse cuidado é importante porque, sem ele, a documentação deixa de servir como fonte confiável de entendimento.

## Documentação e colaboração com IA

Este projeto também se beneficia do uso de ferramentas de IA para estudo, apoio à implementação e geração de código.

Nesse contexto, documentação bem organizada tem valor adicional porque:

- reduz ambiguidade de intenção;
- melhora localização da responsabilidade correta;
- facilita instruções mais específicas;
- ajuda a IA e os desenvolvedores a partirem do mesmo modelo mental;
- diminui a chance de mudanças serem implementadas no módulo errado.

Ou seja, documentação viva não ajuda apenas humanos. Ela melhora a qualidade da colaboração técnica como um todo.

## O que novos desenvolvedores devem preservar

Ao usar e evoluir a documentação, é importante preservar alguns princípios:

- evitar duplicação contraditória entre documentos;
- manter o texto alinhado ao comportamento implementado;
- atualizar a documentação quando a mudança altera entendimento estrutural do sistema;
- preferir clareza conceitual a excesso de detalhe operacional no lugar errado;
- usar esta visão geral como material de orientação, não como substituto da leitura do código e dos documentos especializados.

## Síntese

Esta visão geral foi criada para funcionar como porta de entrada para o entendimento da API.

Ela se conecta aos demais documentos técnicos do projeto e ajuda novos desenvolvedores a construir um mapa mental do sistema antes de mergulhar em arquivos específicos. Seu valor está em reduzir custo de onboarding, melhorar a qualidade das mudanças e preservar coerência arquitetural ao longo da evolução do projeto.

Em conjunto com a documentação especializada e com o código, esta visão geral forma uma base de trabalho para desenvolvimento contínuo e colaboração mais segura dentro da aplicação.
