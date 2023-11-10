compile:
	go build main.go
	mv main cli_subsonic

image:
	podman build -t cli_subsonic_client_build_image .

build: image
	podman run --rm -v $(PWD):/app --workdir /app cli_subsonic_client_build_image sh -c './module.sh; make'

fmt: image
	podman run --rm -v $(PWD):/app --workdir /app cli_subsonic_client_build_image sh -c 'gofmt -w -s .'

bash: image
	podman run -it --rm -v $(PWD):/app --workdir /app cli_subsonic_client_build_image bash

install: build
	sudo cp cli_subsonic /usr/bin/

clean:
	rm -f cli_subsonic
