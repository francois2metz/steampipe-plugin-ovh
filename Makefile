.PHONY: help

help: # show this message
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

install: # install the plugin as is
	go build -o  ~/.steampipe/plugins/hub.steampipe.io/plugins/francois2metz/ovh@latest/steampipe-plugin-ovh.plugin *.go


uninstall: # for testing purposes : uninstalls & backups config
	steampipe plugin uninstall francois2metz/ovh
	mv ~/.steampipe/config/ovh.spc ~/.steampipe/config/ovh.spc.old || echo "config file not found"

apply_config:
	mv ~/.steampipe/config/ovh.spc.old ~/.steampipe/config/ovh.spc

dev: install apply_config # for testing purposes : installs and apply config

reload: uninstall install apply_config # for testing purposes : uninstalls, saves config, re install and apply config