# Soul Hunter

Esse projeto é o backend da página do jogo Soul Hunter. É um projeto simples que faz a conexão com banco de dados e que tem um chat em tempo real.

## Tecnologias

- Go
- Redis
- FireBase
- WebSocket

## Estrutura do projeto
- cmd/
  - api/
    - main.go
- internal/
  -  database/
  -  middleware/
  -  player/
  -  websocket/

## Como rodar localmente
### 1. Clone o repositório
~~~
git clone git@github.com:JoaoGeraldoS/Backend-TCC.git
cd Backend-TCC
~~~
### 2. Instale as dependências
~~~
go mod tidy
~~~
### 3. Conecte com o banco de dados
Ter um banco já configurado no firebase com uma coleção "Ordem"
Nessa coleção tem que ter os atributos "usuario", "pontos", "level".

Arquivo json de conexão com o nome de "service-account.json".

### 4. Subir o Redis pro uso do chat
~~~
docker run -d --name redis -p 6379:6379 redis
~~~

### 5. Rode a aplição
~~~
./cmd/api/main.go
~~~
## Endpoints 
Exemplos
~~~
/chat
POST /login
GET /ranking
GET /filter

~~~
---
## Padronização de erros
| Código | Significado |
|-----|-----|
| 400 | BadRequest |
| 404 | NotFound |
| 500 | InternalServerError |

## Exemplos de payload
Ler Usuarios 
``` Json
{
    "total_points": 22297,
    "total_players": 7,
    "players": [
        {
            "nick": "Bixovino",
            "points": 7263,
            "level": 0
        },
        {
            "nick": "T-Pose",
            "points": 6300,
            "level": 0
        },
    ]
}
```
Resposta do chat
``` json
{
  "user": "Bixovino",
  "message": "Teste",
  "time": "14:35"
}
