# Go + Vite

Server rendered HTML using Go, optimized frontend with Vite. 

```
pnpm vite dev
```

```
go run .
```

By default, the Go server expects the Vite server to be on port `5173`. Set `VITE_PORT` to set the dev server's port.

```
VITE_PORT=4000 go run .
```

## Production

```
pnpm tsx vite/build.tsx

go build

ENV="PROD" ./go-vite
```
s