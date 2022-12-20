gqlgen:
	go run github.com/99designs/gqlgen

migration:
	go install go-graph/db/migration
	migration	