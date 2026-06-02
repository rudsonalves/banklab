# Backlog: pepper da senha transacional

## 1. Contexto

A senha transacional do MVP ZTA é persistida no banco apenas como hash bcrypt.
O bcrypt já gera um salt aleatório por hash e embute esse salt no valor salvo,
portanto o problema não é ausência de salt.

Mesmo assim, a senha transacional atual é um PIN numérico de 6 dígitos. Esse
formato tem um espaço de busca pequeno. Em caso de vazamento do banco, um
atacante poderia tentar validar PINs offline contra os hashes bcrypt. O custo do
bcrypt reduz a velocidade desse ataque, mas não elimina o risco estrutural de
um segredo curto.

Este backlog adiciona um pepper global, lido do ambiente da API, para fortalecer
o material de entrada antes do bcrypt. O pepper não fica no banco e deve ser
mantido como segredo operacional da aplicação.

## 2. Objetivo

Adicionar `TRANSACTION_PASSWORD_PEPPER` como segredo obrigatório da API e usá-lo
na geração e comparação da senha transacional.

O fluxo deve:

- manter o bcrypt como mecanismo de hash persistido;
- preservar o salt automático do bcrypt;
- aplicar o pepper antes do bcrypt;
- manter o pepper fora do banco;
- carregar o pepper via `.env`/variável de ambiente;
- falhar no startup se o pepper não estiver configurado;
- manter a senha transacional como segredo nunca retornado ao cliente.

## 3. Decisão técnica

O pepper deve ser aplicado com HMAC-SHA256 antes do bcrypt, em vez de simples
concatenação de strings.

Modelo recomendado:

```text
mac = HMAC-SHA256(key = TRANSACTION_PASSWORD_PEPPER, message = PIN)
peppered = base64(mac)
hash = bcrypt(peppered)
```

Na validação:

```text
mac = HMAC-SHA256(key = TRANSACTION_PASSWORD_PEPPER, message = PIN)
peppered = base64(mac)
bcrypt.compare(hash, peppered)
```

Essa abordagem evita ambiguidade de concatenação, separa claramente segredo e
mensagem, e permite que o bcrypt continue cuidando de salt e fator de custo. O
resultado do HMAC deve ser codificado em base64 antes de ser passado ao bcrypt.

## 4. Configuração

Adicionar a variável:

```env
TRANSACTION_PASSWORD_PEPPER=<segredo aleatório>
```

Recomendação para gerar o segredo:

```bash
openssl rand -base64 32
```

Regras:

- exigir tamanho mínimo de 32 caracteres para o valor configurado;
- não reutilizar `JWT_SECRET`;
- não reutilizar `APP_TOKEN`;
- usar valor diferente por ambiente;
- não versionar o valor real no Git;
- documentar exemplos apenas com placeholder ou valor dev descartável;
- tratar a ausência da variável como erro fatal no startup da API.

## 5. Escopo de API

- Atualizar `bootstrap.Config` para incluir `TransactionPasswordPepper`.
- Atualizar `LoadConfig` para exigir `TRANSACTION_PASSWORD_PEPPER`.
- Atualizar o wiring em `cmd/api/main.go`.
- Atualizar `BcryptTransactionPasswordHasher` para receber o pepper.
- Aplicar HMAC-SHA256 antes de `bcrypt.GenerateFromPassword`.
- Aplicar o mesmo pre-processamento antes de `bcrypt.CompareHashAndPassword`.
- Adicionar testes unitários do hasher com pepper.
- Atualizar testes afetados pelo novo construtor.
- Atualizar documentação de setup local e variáveis de ambiente.
- Atualizar scripts/templates de inicialização de `.env`, se existirem valores
  padrão para `JWT_SECRET` e `APP_TOKEN`.

## 6. Compatibilidade e migração

Decisão inicial para o projeto:

- não criar migração automática de hashes já existentes;
- aceitar que hashes antigos sem pepper deixam de validar após a mudança;
- como o app segue em desenvolvimento, sem dados reais, e o banco atual está
  vazio, não há massa legada para preservar neste momento;
- em ambiente local/dev, se houver senha transacional criada antes da mudança, o
  usuário pode recriá-la;
- se houver necessidade futura de ambiente com usuários reais, introduzir
  estratégia explícita de versão de hash antes de ativar a mudança.

Motivo:

O projeto ainda está em fase de laboratório/MVP. Manter compatibilidade dupla
sem requisito real aumentaria complexidade de segurança e de testes. Para
produção, a estratégia recomendada seria versionar o algoritmo/hash ou permitir
validação legada e rehash com pepper após sucesso.

## 7. Impacto no mobile

O mobile não deve conhecer o pepper e não deve alterar o payload da senha
transacional.

Contratos públicos permanecem:

- `POST /security/transaction-password`
- `POST /security/step-up/authorize`
- campo `transaction_password`

Esta alteração deve ser concluída antes de prosseguir com a implementação mobile
do fluxo de senha transacional e step-up, para evitar validar o app contra uma
API que ainda mudará a forma interna de persistência/validação da credencial.

## 8. Critérios de aceite

- API exige `TRANSACTION_PASSWORD_PEPPER` no startup.
- API rejeita `TRANSACTION_PASSWORD_PEPPER` vazio ou menor que 32 caracteres.
- Senha transacional nova é persistida como bcrypt de um valor com pepper.
- Comparação usa o mesmo pepper configurado.
- Hash gerado com um pepper não valida quando a API usa outro pepper.
- Testes cobrem sucesso, senha incorreta e pepper incorreto.
- Testes existentes de criação e autorização de senha transacional continuam
  passando após ajustes.
- Documentação de ambiente inclui `TRANSACTION_PASSWORD_PEPPER`.
- Mobile não recebe novo campo e não precisa alterar contrato.

## 9. Fora de escopo

- Alterar o formato público do PIN.
- Alterar tamanho da senha transacional.
- Trocar bcrypt por Argon2id.
- Criar rotação automática de pepper.
- Criar suporte a múltiplos peppers ativos.
- Criar versão de algoritmo/hash no banco.
- Alterar o contrato mobile de criação ou autorização de step-up.

## 10. Referências

- `api/internal/security/infrastructure/bcrypt_transaction_password_hasher.go`
- `api/internal/bootstrap/bootstrap.go`
- `api/cmd/api/main.go`
- `api/docs/00-getting_started.md`
- `docs/backlogs/mobile/011 - cadastro-senha-transacional.md`
- `docs/backlogs/mobile/012 - step-up-transferencia-interna.md`
