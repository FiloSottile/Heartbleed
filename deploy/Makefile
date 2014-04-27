GO ?= go

HBserver:
	GOOS=linux GOARCH=amd64 $(GO) build -o HBserver github.com/FiloSottile/Heartbleed/server

hosts.txt: inventory.ini
	ansible -i inventory.ini --list-hosts bleed | cut -c 5- > hosts.txt
