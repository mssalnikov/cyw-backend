# install golang
https://golang.org/doc/install

# deploy bases
```
cd deploy
docker-compose up -d
```

# insatll migrations
```
go get -u github.com/pressly/goose/cmd/goose
```

# migrate
```
cd migrations
goose postgres "user=developer password=somepass dbname=cywdb sslmode=disable" up
```

# psql
```
docker exec -ti deploy_postgres_1 psql "postgresql://developer:somepass@localhost/cywdb"
```

# run application
```
go run main.go
```