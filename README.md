# My Cooking Codex - API

## Environment Variables

| Name                     | Description                                         | Default   |
| :----------------------- | :-------------------------------------------------- | :-------- |
| BIND__HOST               | Host to listen on                                   | 127.0.0.1 |
| BIND__PORT               | Port to bind to                                     | 8000      |
| DB__URI                  | Database URI                                        |           |
| DB__TYPE                 | The type of database (sqlite, postgres)             |           |
| DATA__RECIPE_IMAGES_BASE | Where recipe images will be stored                  |           |
| JWT_SECRET               | base64 encoded secret for JWT authentication tokens |           |
| STATIC_PATH              | Serve static files at / (e.g. the frontend)         | -         |
| CORS_ORIGINS             | List of origins that may access the API             | *         |
| OPTIMIZED_IMAGE_SIZE     | Max image size to shrink uploaded image to          | 2000      |
| MAX_UPLOAD_SIZE          | The max possible upload size                        | 4M        |

### DB__URI

```
# sqlite
./data/db.sqlite

# postgres
host=localhost user=user password=password dbname=my_cooking_codex port=9920 sslmode=disable TimeZone=Europe/London
```

## Without Docker
### Requirements
- Database (SQLite, PostgreSQL)
- libvips installed
- go >= 1.20

### Build
Run these commands:

```
go build
```

Copy built binary `./api` and run.
