# to run go project locally
run-dev:
	export PASSBOOK_ENV=dev; go run cmd/passbook-app/main.go
run-prod:
	export PASSBOOK_ENV=prod;export GIN_MODE=release; go run cmd/passbook-app/main.go