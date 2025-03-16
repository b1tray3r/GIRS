install:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0
	go install github.com/air-verse/air@latest

image:
	cd docker; docker buildx build -t girs:test -f Dockerfile --no-cache ..

run:
	docker run -it --rm -p 8085:8085 --env-file ./.env girs:test
	