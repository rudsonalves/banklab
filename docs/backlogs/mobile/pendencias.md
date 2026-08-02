Ficaram estas pendências reais:

1. **Como detectar reinstalação/limpeza de dados**
   - Decidir o mecanismo do marker local fora do Keychain/secure storage.
   - Confirmar onde salvar esse marker no Flutter.
   - Regra proposta: se marker sumiu, gerar nova `installation_id`, mesmo se o secure storage ainda devolver uma antiga.

2. **Backup/restore**
   - Decidir se o marker e/ou `installation_id` podem ser restaurados por backup do sistema.
   - Para a semântica desejada, restore não deveria ressuscitar uma instalação antiga sem querer.

3. **Comportamento quando storage falha**
   - O backlog diz que no rollout inicial a requisição pode seguir sem header.
   - Falta decidir se isso continua aceitável mesmo com a API exigindo `X-Installation-Id`.
   - Na prática: se leitura/escrita falhar, o app mostra erro de login/rede? tenta de novo? segue sem header?

4. **Concorrência**
   - Definir que múltiplas chamadas simultâneas ao serviço compartilham a mesma operação de read-or-create.
   - Aqui é mais implementação, mas a regra precisa ficar explícita.

5. **Resposta de login restrito**
   - Confirmar UX quando a API retornar `restricted_access_token`.
   - O usuário vai para uma tela de senha transacional para registrar a instalação?
   - Pode cancelar? Ao cancelar, volta para login e descarta tudo?

6. **`installation_limit_reached`**
   - Definir a mensagem/estado final.
   - Decidir se haverá botão para “Entendi”, “Voltar ao login”, ou algo mais.
   - Confirmar que não haverá step-up nesse caso.

7. **Senha transacional ausente/bloqueada**
   - Decidir para onde o mobile direciona quando a API disser `TRANSACTION_PASSWORD_NOT_SET` ou `TRANSACTION_PASSWORD_LOCKED` no fluxo de nova instalação.
   - Como é nova instalação ainda não operacional, talvez não dê para abrir a área logada para configurar a senha.

8. **Expiração/cancelamento do token restrito**
   - Definir o que mostrar quando o `restricted_access_token` expirar antes do registro.
   - Provável regra: descartar fluxo e pedir login novamente.

9. **Gerenciamento de instalações**
   - Experiência de listar instalações.
   - Experiência de revogar uma instalação.
   - Como identificar/mostrar a instalação atual.
   - O que fazer se uma instalação revogada derrubar a sessão atual.

10. **Telemetria segura**
   - O que registrar quando storage/geração/interceptor falhar.
   - Sem logar UUID completo, tokens ou senha.

11. **Escopo do primeiro corte**
   - Decidir se o primeiro pacote mobile entrega só identidade + header + login normal, ou se já inclui o fluxo completo de registro de nova instalação.