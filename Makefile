install:
	go build -o  ~/.steampipe/plugins/hub.steampipe.io/plugins/francois2metz/ovh@latest/steampipe-plugin-ovh.plugin *.go

test:
	go test -v ./ovh/...

update_fixtures:
	@curl -s -X GET "https://eu.api.ovh.com/v1/cloud/project/$(OVH_CLOUD_PROJECT_ID)/region" -H "accept: application/json" -H "authorization: Bearer $(OVH_BEARER_TOKEN)" > fixtures/ovh_regions.json
	@for region in $(cat fixtures/ovh_regions.json | jq -r .[]) ; do \
		curl -s -X GET  "https://eu.api.ovh.com/v1/cloud/project/$(OVH_CLOUD_PROJECT_ID)/region/$region" -H "accept: application/json" -H "authorization: Bearer $(OVH_BEARER_TOKEN)" > fixtures/ovh_region_$region.json ; \
	done
