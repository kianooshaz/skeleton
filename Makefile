APP=rostam-stage
PORT=8080

.which-go:
	@which go > /dev/null || (echo "install go from https://go.dev/dl/" & exit 1)

.which-go-mod-upgrade:
	@which go-mod-upgrade > /dev/null || (echo "install go-mod-upgrade from https://github.com/oligot/go-mod-upgrade" & exit 1)

.which-golangci-lint:
	@which golangci-lint > /dev/null || (echo "install golangci-lint from https://github.com/golangci/golangci-lint" & exit 1)

.which-swag:
	@which swag > /dev/null || (echo "install swag from https://github.com/swaggo/swag" & exit 1)

.which-sqlc:
	@which sqlc > /dev/null || (echo "install sqlc from https://sqlc.dev" & exit 1)

.now:
	@date

upgrade: .now  .which-go-mod-upgrade
	go-mod-upgrade

lint: .now  .which-golangci-lint
	golangci-lint run

swag: .now .which-swag
	swag init --pd -g "../http.go"  --o "./docs/api/" --ot "go,json" --dir "./internal/transport/server/http/handler"

sqlc: .now .which-sqlc
	sqlc -f build/sqlc.yaml generate