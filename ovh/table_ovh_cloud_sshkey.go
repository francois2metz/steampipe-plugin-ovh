package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableOvhCloudSshkey() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_sshkey",
		Description: "Get SSH Keys.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listSshkey,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getSshkey,
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project_id"),
				Description: "Project ID.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "SSH Key ID.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "SSH Key name.",
			},
			{
				Name:        "public_key",
				Type:        proto.ColumnType_STRING,
				Description: "SSH public key.",
			},
		},
	}
}

type Sshkey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

func listSshkey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	var sshkeys []Sshkey
	err = client.Get(fmt.Sprintf("/cloud/project/%s/sshkey", projectId), &sshkeys)
	if err != nil {
		return nil, err
	}
	for _, sshkey := range sshkeys {
		d.StreamListItem(ctx, sshkey)
	}
	return nil, nil
}

func getSshkey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	id := d.KeyColumnQuals["id"].GetStringValue()
	var sshkey Sshkey
	err = client.Get(fmt.Sprintf("/cloud/project/%s/sshkey/%s", projectId, id), &sshkey)
	if err != nil {
		return nil, err
	}
	return sshkey, nil
}
