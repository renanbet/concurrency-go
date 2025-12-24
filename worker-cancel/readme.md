**Problema de Concorrência**
**Cenário**

Você precisa processar uma lista de tarefas de forma concorrente usando goroutines e channels, garantindo:

1. Limite de concorrência (ex: no máximo N tarefas rodando ao mesmo tempo)
2. Coleta de resultados
3. Tratamento de erros
4. Encerramento limpo (sem goroutine vazando)

**Descrição do problema**

Você recebe uma lista de números inteiros.
Cada número representa uma "tarefa".

**Para cada número:**

- Simule um processamento que demora entre 100ms e 500ms
- Se o número for múltiplo de 7, a tarefa falha e retorna um erro
- Caso contrário, o resultado é o quadrado do número

**Requisitos obrigatórios**

Implemente um programa em Go que:

1. Processe os números concorrentemente
2. Tenha um limite máximo de workers (ex: 3 workers)
3. Use channels para:
  - Enviar tarefas
  - Receber resultados
4. Retorne:
  - Um slice com os resultados bem-sucedidos
  - Um slice com os erros
5. Garanta:
  - Nenhuma goroutine fique bloqueada
  - O programa finalize corretamente