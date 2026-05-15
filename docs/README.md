# Documentação do BankLab

Esta pasta reúne documentos de apoio ao desenvolvimento, decisões técnicas, roadmap, backlogs e materiais de apresentação do BankLab.

No BankLab, a documentação não serve apenas para descrever o estado final do software. Ela também preserva o caminho de raciocínio que levou às decisões: discussões, ajustes, reversões, trade-offs e alternativas descartadas.

Essa escolha é intencional. O projeto evolui de forma incremental, e tornar esse processo visível ajuda novos colaboradores a entender não só **o que** foi decidido, mas também **por que** foi decidido.

## Estrutura

```text
docs/
|-- ROADMAP.md          # direção pública do projeto
|-- backlogs/           # discussões ativas e histórico de decisões por área
|-- disclosure/         # material de apresentação e divulgação
|-- mermaid-images/     # imagens e fontes geradas a partir de diagramas
|-- relatorio-*.md      # relatórios de implementação
```

## Documentos principais

- [ROADMAP.md](ROADMAP.md): visão de evolução do projeto.
- [AI_DOCUMENTATION_GUIDE.md](AI_DOCUMENTATION_GUIDE.md): organização da documentação usada por humanos e agentes de IA.
- [backlogs/README.md](backlogs/README.md): organização dos backlogs ativos e concluídos.
- [disclosure/apresentacao_banklab.md](disclosure/apresentacao_banklab.md): material de apresentação pública do projeto.

## Como usar

- Use `ROADMAP.md` para entender a direção geral.
- Use `backlogs/` para acompanhar discussões em aberto, decisões tomadas e tarefas planejadas.
- Use `backlogs/*/done/` para consultar histórico de decisões já resolvidas ou implementadas.
- Use `disclosure/` para materiais voltados à comunicação pública do projeto.

Os backlogs não são apenas notas internas. Eles fazem parte da documentação do processo de decisão e ajudam novos colaboradores a entender o caminho do projeto.

## Filosofia de documentação

O repositório documenta duas dimensões do BankLab:

- o software implementado;
- o processo de engenharia que moldou esse software.

Por isso, documentos de discussão e backlogs devem preservar o histórico das deliberações relevantes. Quando uma abordagem for substituída, ela pode sair do caminho ativo, mas não precisa desaparecer. O histórico ajuda a revelar maturidade arquitetural, evolução de pensamento e critérios usados para tomar decisões.
