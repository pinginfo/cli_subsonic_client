compile:
	go build main.go
	mv main cli_subsonic

image:
	docker build -t cli_subsonic_client_build_image .

image_musl:
	docker build -f Dockerfile.musl -t cli_subsonic_client_build_image .

build: image
	docker run --rm -v $(PWD):/app --workdir /app cli_subsonic_client_build_image sh -c './module.sh; make'

build_musl: image_musl
	docker run --rm -v $(PWD):/app --workdir /app cli_subsonic_client_build_image sh -c './module.sh; make'

fmt: image
	docker run --rm -v $(PWD):/app --workdir /app cli_subsonic_client_build_image sh -c 'gofmt -w -s .'

bash: image
	docker run -it --rm -v $(PWD):/app --workdir /app cli_subsonic_client_build_image bash

install: build
	sudo cp cli_subsonic /usr/bin/

install_musl: build_musl
	sudo cp cli_subsonic /usr/bin

clean:
	rm -f cli_subsonic
