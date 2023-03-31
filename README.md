# My Cooking Codex - API

## Environment Variables

| Name                 | Description                                         | Default   |
| :------------------- | :-------------------------------------------------- | :-------- |
| HOST                 | Host to listen on                                   | 127.0.0.1 |
| PORT                 | Port to bind to                                     | 8000      |
| DATA_PATH            | Where app data will be stored                       |           |
| DB__URI              | Database URI                                        |           |
| DB__TYPE             | The type of database (sqlite, mysql, postgres)      |           |
| JWT_SECRET           | base64 encoded secret for JWT authentication tokens |           |
| STATIC_PATH          | Serve static files at / (e.g. the frontend)         | -         |
| CORS_ORIGINS         | List of origins that may access the API             | *         |
| OPTIMIZED_IMAGE_SIZE | Max image size to shrink uploaded image to          | 2000      |
| MAX_UPLOAD_SIZE      | The max possible upload size                        | 4M        |

### DB__URI

```
# sqlite
./data/db.sqlite

# mysql
user:pass@tcp(127.0.0.1:3306)/my_cooking_codex?charset=utf8mb4&parseTime=True&loc=Local

# postgres
host=localhost user=user password=password dbname=my_cooking_codex port=9920 sslmode=disable TimeZone=Europe/London
```

## Without Docker
### Requirements
- Database (SQLite, MySQL, PostgreSQL)
- libvips installed
- go >= 1.20

### Build
Run these commands:

```
go build
```

Copy built binary `./api` and run.
