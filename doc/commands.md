#### Alles bauen
`go build -v ./...`

#### Rekursiv alle Unit-Tests starten
`go test -cover ./...`

#### Coverage
`go test -coverprofile=coverage.out ./...`
`go tool cover -html=coverage.out -o ./result.html`
