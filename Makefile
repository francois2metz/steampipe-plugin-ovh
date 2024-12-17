install:
	go build -o  ~/.steampipe/plugins/hub.steampipe.io/plugins/francois2metz/ovh@latest/steampipe-plugin-ovh.plugin *.go

uninstall:
	steampipe plugin uninstall francois2metz/ovh && mv ~/.steampipe/config/ovh.spc ~/.steampipe/config/ovh.spc.old

test:
	mv ~/.steampipe/config/ovh.spc.old ~/.steampipe/config/ovh.spc

dev: install test