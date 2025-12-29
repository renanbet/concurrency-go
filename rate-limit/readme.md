**Contexto**

Você tem vários workers concorrentes processando tarefas, mas existe um limite global de taxa:

❗ No máximo N tarefas por segundo podem iniciar processamento, independente do número de goroutines.

📌 Requisitos obrigatórios
1. Limite global

- Exemplo: 5 tarefas por segundo
- Mesmo com 100 workers, não pode ultrapassar esse limite

2. Concorrência

- Vários workers lendo de um taskChan
- Todos devem respeitar o mesmo rate limiter

3. Bloqueante e cancelável

- Se o limite for atingido:
- a goroutine bloqueia aguardando permissão
- mas respeita context.Context

Se o contexto for cancelado:

- deve sair imediatamente

4. Teste

- Rate limit: 3 req/s
- 5 workers
- 20 tarefas
- Execução leva ~7 segundos
- Nenhum pico acima de 3 execuções iniciadas por segundo