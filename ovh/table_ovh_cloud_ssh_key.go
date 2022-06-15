package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableOvhCloudSshKey() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_ssh_key",
		Description: "An ssh key allows you to connect to an instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listSshKey,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getSshKey,
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

type SshKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

func listSshKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ssh_key.listSshKey", "connection_error", err)
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	var sshKeys []SshKey
	err = client.Get(fmt.Sprintf("/cloud/project/%s/sshkey", projectId), &sshKeys)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ssh_key.listSshKey", err)
		return nil, err
	}
	for _, sshKey := range sshKeys {
		d.StreamListItem(ctx, sshKey)
	}
	return nil, nil
}

func getSshKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ssh_key.getSshKey", "connection_error", err)
		return nil, err
	}
	projectId := d.KeyColumnQuals["project_id"].GetStringValue()
	id := d.KeyColumnQuals["id"].GetStringValue()
	var sshKey SshKey
	err = client.Get(fmt.Sprintf("/cloud/project/%s/sshkey/%s", projectId, id), &sshKey)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_ssh_key.getSshKey", err)
		return nil, err
	}
	return sshKey, nil
}
