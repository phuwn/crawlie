# crawlie

A web application that will extract large amounts of data from the Google search results page

## Local Setup

1. Download this repository
2. Install [asdf](https://asdf-vm.com/guide/getting-started.html)
3. Run `asdf plugin add golang` to set up asdf for go
4. Install go through `asdf install golang 1.21`
5. Use go in the repository `asdf local golang 1.21`
6. Download dependencies `go mod tidy`
7. Run The Database by `psql` command or use `docker-compose`

```sh
docker compose up -d
```

8. Generate & fill the config file `cp config.example.json config.json` and fill in the details
9. Run the server `go run cmd/server/*.go`
10. Run the crawler `go run cmd/crawler/*.go`
11. Test app if needed `GOARCH=amd64 go test ./...`
