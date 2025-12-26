**Refatore seu código (worker-cancel) para:**

**Refatore processTask para:**

1. Tentar executar a tarefa até 3 vezes
2. Usar backoff exponencial com jitter
3. Respeitar:
- ctx.Done()
- timeout por tentativa (250ms)
- Não retry para erro de negócio (task % 7 == 0)
4. Retornar:
- Sucesso se alguma tentativa passar
- Último erro caso todas falhem