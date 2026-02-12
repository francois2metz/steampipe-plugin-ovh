install:
	go build -o  ~/.steampipe/plugins/hub.steampipe.io/plugins/francois2metz/ovh@latest/steampipe-plugin-ovh.plugin *.go

build:
	go build -o steampipe-plugin-ovh.plugin *.go