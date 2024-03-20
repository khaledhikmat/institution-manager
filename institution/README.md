## Setup

```bash
go mod init github.com/khaledhikmat/campaign-manager/country
go get -u github.com/dapr/go-sdk
```

## Uprade Go Version

```bash
go mod edit -go 1.22
go mod tidy
```

## Change Module Name

```bash
go mod edit -module <new_name>
```

But remember to refactor the import and do `go mod tidy`.