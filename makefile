run-dev:
	PASSBOOK_ENV=DEV CGO_ENABLED=0 go run cmd/passbook-app/main.go
build-dev:
	CGO_ENABLED=0 GOOS=linux go build -o bin/passbook-app /cmd/passbook-app/main.go
docker-build-image:
	docker build -t passbook-app-backend -f Dockerfile.multistage .
docker-run-image:
	docker run -p 8080:8080 -d --env-file dev.env passbook-app-backend