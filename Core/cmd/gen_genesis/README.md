# genesis
- db:
```
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build .
```
- run:
```
go run . --config=config.yaml --genesis=genesis.json
```
- copy folder and move into folder node