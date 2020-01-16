# api-template
api-template is MVC template to create micrososervices with postgres as backend database.


## Making new project from template
Run following command in shell. Replace `some\/new\/myproject` with desired project path.

```
curl -s https://raw.githubusercontent.com/samtech09/api-template/master/scripts/create.sh | bash -s "some\/new\/myproject"
```


## Technologies used
- [Chi router](https://github.com/go-chi/chi) for request routing
- [pgx](https://github.com/jackc/pgx) and [pgxpool](https://github.com/jackc/pgx/v4/pgxpool) for connection to PostgreSQL
- [zerolog](https://github.com/rs/zerolog) for logging
