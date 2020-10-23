# silly todo

A small (and silly) todo app built while evaluating a simple yet powerful web stack;  
[Svelte](https://svelte.dev/) with [Materialize](https://materializecss.com/) + [Go](https://golang.org/) + [Postgres](https://www.postgresql.org/) (with Docker and Compose)


## prerequisites
* (run): docker and docker-compose
* (dev): go 1.14.x, node 12.16.x, docker and docker-compose  
## run
```docker-compose up``` and navigate to localhost:80 with your favourite browser


This small app is made up of 3 services: a Postgres database, a http API (Go) and a frontend SPA (Svelte + Materialize). When started with docker-compose a migration container will run migrations and then exit. If you run these services outside of docker there is a migration (bash) script prepared in *./scripts*.

## test
in ./api run ```go test ./...```