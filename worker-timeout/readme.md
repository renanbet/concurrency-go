**Refatore seu código (worker-cancel) para:**

1. Cada tarefa ter timeout de 250ms
2. Se a tarefa estourar:
  - Retornar erro context.DeadlineExceeded
4. Se o main cancelar:
  - Workers param imediatamente
5. Nenhuma goroutine vaza
6. Nenhum worker bloqueia ao enviar resultado

**Conceitos**

- Worker pool
- Fan-out / fan-in
- Cancelamento hierárquico
- Cancelamento global
- Timeout por tarefa
- Encerramento limpo
- Evitar goroutine leak
- Padrões idiomáticos Go
- Consumer lento ou inexistente
- Zero vazamento