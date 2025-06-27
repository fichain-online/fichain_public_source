# client
- build:
```
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build .
```
- run :
```
./client --config=config.yaml
```