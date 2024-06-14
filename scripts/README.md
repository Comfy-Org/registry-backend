# Scripts

## Create JWT Token

This script is used to create JWT Token representing an existing user which can then be used to invoke the API Server as an alternative authentication method to firebase auth.

```bash
export JWT_SECRET
go run ./scripts/create-jwt-token --user-id "<existing user id>"
```

Notes:

1. You need to use the same values as the [JWT_SECRET](./../docker-compose.yml#L20) environment variables defined for the API Server.
2. By default, the token will [expire in 30 days](./create-jwt-token/main.go#L14), which can be overriden using `--expiry` flag.

## Ban a Publisher

This script is used to invoke ban publisher API using jwt token as authentication.

```bash
go run ./scripts/ban-publisher \
    --base-url "<base url of the API>" \
    --token "<jwt token representing an admin user>" \
    --publisher-id "<publisher id to be banned>"
```

## Ban a Node

This script is used to invoke ban publisher API using jwt token as authentication.

```bash
go run ./scripts/ban-node \
    --base-url "<base url of the API>" \
    --token "<jwt token representing an admin user>" \
    --publisher-id "<publisher id of the node>"
    --node-id "<node id to be banned>"
```
