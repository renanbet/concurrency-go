**Desafio: Circuit Breaker concorrente em Go**

- Você vai implementar um Circuit Breaker thread-safe, integrado ao seu pipeline concorrente (workers, retry, rate limit).

**Conceito**

- O Circuit Breaker protege um sistema quando um serviço começa a falhar.
- Estados clássicos:
  - CLOSED → OPEN → HALF-OPEN → CLOSED

**Regras do desafio**

Estados:

🔒 CLOSED
- Requisições fluem normalmente
- Erros são contabilizados
- Ao atingir o limite → OPEN

🚫 OPEN
- Todas as chamadas falham imediatamente
- Nenhuma chamada ao handler
- Após um tempo (cooldown) → HALF-OPEN

⚠️ HALF-OPEN
- Permite N chamadas de teste
- Se todas tiverem sucesso → CLOSED
- Se qualquer falhar → OPEN

**Configurações mínimas**

- Você deve suportar:
- FailureThreshold (ex: 5 erros)
- ResetTimeout (ex: 3s)
- HalfOpenMaxCalls (ex: 2)

