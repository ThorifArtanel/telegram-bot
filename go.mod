module github.com/digitalocean/sample-golang

go 1.13

replace api => ./api

require (
	api v0.0.0-00010101000000-000000000000
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.2 // indirect
)
