.PHONY: HBserver

HBserver:
	cd ../server && GOOS=linux GOARCH=amd64 make HBserver
	mv ../server/HBserver .

hosts.txt: inventory.ini
	ansible -i inventory.ini --list-hosts bleed | cut -c 5- > hosts.txt
