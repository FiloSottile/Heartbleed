.PHONY: HBserver

HBserver:
	cd ../server && GOOS=linux GOARCH=amd64 go build -o ../deploy/HBserver

hosts.txt: inventory.ini
	ansible -i inventory.ini --list-hosts bleed | cut -c 5- > hosts.txt
