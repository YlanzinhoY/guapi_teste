Contexto Geral do Projeto
Você deve criar uma API RESTful utilizando Go, Echo, e sqlc para gerenciar salas de chat, participantes, mensagens, e curtidas. Além disso, deve implementar uma funcionalidade de WebSocket para notificação de novas mensagens e curtidas.

## Parte 1: Criação de Sala e Participantes
    Objetivo: Implementar rotas para criar salas de chat e adicionar participantes a essas salas.
    Requisitos:
    Uma rota para criar uma sala.
    Uma rota para adicionar participantes a uma sala.
    Cada sala deve ter um nome único.
    Participantes devem ser associados a uma sala existente.

### Parte 2: Mensagens e Validação de Exclusão de Sala
    Objetivo: Implementar rotas para criar e gerenciar mensagens dentro das salas, e validar a exclusão de uma sala.
    Requisitos:
    Uma rota para enviar mensagens em uma sala específica.
    A exclusão de uma sala só deve ser permitida se não houver nenhuma mensagem associada a essa sala.
    As mensagens devem ter um conteúdo e ser associadas a um participante da sala.

#### Parte 3: Curtidas nas Mensagens
    Objetivo: Implementar uma funcionalidade para gerenciar curtidas em mensagens.
    Requisitos:
    Uma rota que permita acrescentar curtidas em uma mensagem utilizando o método HTTP PATCH.
    Uma rota que permita remover curtidas em uma mensagem utilizando o método HTTP DELETE.
    O número de curtidas de uma mensagem nunca deve ser menor que zero.

### Parte 4: Implementação de WebSocket para Notificações
    Objetivo: Implementar uma funcionalidade de WebSocket para notificação em tempo real.
    Requisitos:
    Os usuários devem poder se inscrever em uma sala para receber notificações.
    Uma notificação deve ser enviada para todos os inscritos sempre que uma nova mensagem for enviada na sala.
    Uma notificação deve ser enviada para todos os inscritos sempre que uma mensagem receber uma curtida ou tiver uma curtida removida.

### Parte 5: Considerações Técnicas
    Banco de Dados: Use PostgreSQL para armazenar as informações.
    Migrações: Utilize migrações para criar e gerenciar as tabelas necessárias no banco de dados.
    SQLC: Utilize o sqlc para gerar o código de acesso ao banco de dados.
    Injeção de Dependência: Certifique-se de que a aplicação use injeção de dependência.
    Documentação: A API deve ser documentada utilizando Swagger.
    Testes Unitários: Testes unitários são um diferencial e devem ser incluídos se possível.
    Esse é o conjunto de requisitos do teste. O candidato deve se preocupar em criar uma solução que seja funcional, seguindo as boas práticas de desenvolvimento, e que atenda a todos os pontos descritos.