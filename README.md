# api-template
api-template is MVC template to create micrososervices with postgres / mssql as backend database. Redis and mongodb client also added.


## Making new project from template
Run following command in shell. Replace `some\/new\/myproject` with desired project path.

```
curl -s https://raw.githubusercontent.com/samtech09/api-template/master/scripts/create.sh | bash -s "some\/new\/myproject"
```


## Technologies used
- [Chi router](https://github.com/go-chi/chi) for request routing
- [dbtools/mssql](https://github.com/samtech09/dbtools) for SQL-Server database handling
- [dbtools/pgsql](https://github.com/samtech09/dbtools) for PostgreSQL database handling
- [dbtools/mango](https://github.com/samtech09/dbtools) for working with MongoDB
- [redicache](https://github.com/samtech09/redicache) for connection to Redis server
- [gosql](https://github.com/samtech09/gosql) for generating SQLs with design-time preview/details
