build:
	go build main.go
	mv main cli_subsonic

install: build
	sudo cp cli_subsonic /usr/bin/

clean:
	rm -f cli_subsonic
