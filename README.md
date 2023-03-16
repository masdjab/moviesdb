## MoviesDB REST
This project is written as part of technical test for MaxSol.

## Setup for Testing
1. run `make copy-config`
2. adjust configs in `config.yml` based on your local setup
3. adjust configs in `docker/.env` (optional)
4. run `make docker.start`
5. run `make migrate`
6. run `make start` to start the server
7. run `make docker.stop` when you done testing

## Users Available for Test
1. username: admin, password: admin
2. username: someuser, password: admin
3. username: anotheruser, password: admin

## Endpoints
- GET /ping
  - no header/params required
  - response: "pong"
- POST /login
  - request params:
    - username, ex: "admin"
    - password, ex: "admin"
  - response:
    {
        "success": true,
        "error": "",
        "data": {
            "user_id": 2,
            "user_name": "someuser",
            "display_name": "Some User",
            "email": "someuser@moviesdb.com",
            "token": "MSAW5IFQDCY2X3ZD4GIUBVMW556UXWCDVNGMHRGHONT3ZEYUVBDQ====",
            "last_login": "0001-01-01T00:00:00Z"
        }
    }    
- GET /logout
  - header
    - token, ex: "MSAW5IFQDCY2X3ZD4GIUBVMW556UXWCDVNGMHRGHONT3ZEYUVBDQ===="
  - response:
    {
        "success": true,
        "error": "",
        "data": null
    }  
- GET /movies
  - header
    - token, ex: "MSAW5IFQDCY2X3ZD4GIUBVMW556UXWCDVNGMHRGHONT3ZEYUVBDQ===="
  - request params:
    - keywords (optional), ex: "rass par"
  - response:
    {
        "success": true,
        "error": "",
        "data": [
            {
                "id": 9,
                "title": "jurassic park",
                "description": "",
                "rating": 0,
                "image_url": "http://123.com/2.jpg",
                "created_at": "2023-03-16T01:59:55+07:00",
                "updated_at": "2023-03-16T01:59:55+07:00"
            }
        ]
    }
- GET /movies/{id}
  - header
    - token, ex: "MSAW5IFQDCY2X3ZD4GIUBVMW556UXWCDVNGMHRGHONT3ZEYUVBDQ===="
  - no request params needed
  - response:
    {
        "success": true,
        "error": "",
        "data": {
            "id": 11,
            "title": "beowolf",
            "description": "",
            "rating": 5.5,
            "image_url": "http://123.com/4.jpg",
            "created_at": "2023-03-16T02:00:18+07:00",
            "updated_at": "2023-03-16T02:00:18+07:00"
        }
    }
- POST /movies
  - header
    - token, ex: "MSAW5IFQDCY2X3ZD4GIUBVMW556UXWCDVNGMHRGHONT3ZEYUVBDQ===="
  - request body (json)
    {
      "title": "19 oktober",
      "description": "contoh deskripsi film",
      "image_url": "http://123.com/53.jpg"
    }
  - response:
    {
        "success": true,
        "error": "",
        "data": {
            "id": 14,
            "title": "19 oktober",
            "description": "contoh deskripsi film",
            "rating": 0,
            "image_url": "http://123.com/53.jpg",
            "created_at": "2023-03-16T05:28:49+07:00",
            "updated_at": "2023-03-16T05:28:49+07:00"
        }
    }
- PATCH /movies/{id}
  - header
    - token, ex: "MSAW5IFQDCY2X3ZD4GIUBVMW556UXWCDVNGMHRGHONT3ZEYUVBDQ===="
  - request body (json)
    {
      "title": "tragedi 19 oktober",
      "description": "contoh deskripsi film",
      "image_url": "http://123.com/53.jpg"
    }
  - response:
  {
      "success": true,
      "error": "",
      "data": {
          "id": 14,
          "title": "tragedi 19 oktober",
          "description": "contoh deskripsi film",
          "rating": 0,
          "image_url": "http://123.com/53.jpg",
          "created_at": "2023-03-16T05:28:49+07:00",
          "updated_at": "2023-03-16T05:30:51+07:00"
      }
  }
- DELETE /movies/{id}
  - header
    - token, ex: "MSAW5IFQDCY2X3ZD4GIUBVMW556UXWCDVNGMHRGHONT3ZEYUVBDQ===="
  - no params needed
  - response:
    {
        "success": true,
        "error": "",
        "data": null
    }  
- POST /movies/{id}/vote
  - header
    - token, ex: "MSAW5IFQDCY2X3ZD4GIUBVMW556UXWCDVNGMHRGHONT3ZEYUVBDQ===="
  - params:
    - score (0-10)
  - response:
    {
        "success": true,
        "error": "",
        "data": null
    }
- GET /goroutine-example
  - no header or params needed
  - response:
    Goroutine #1 completed in 427 ms
    Goroutine #3 completed in 445 ms
    Goroutine #2 completed in 467 ms
