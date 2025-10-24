package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type Iam struct {
	DisplayName string            `json:"displayName"`
	ID          string            `json:"id"`
	Tags        map[string]string `json:"tags"`
	URN         string            `json:"urn"`
}

type Ceph struct {
	ID             string   		 `json:"cephId"`
	CephMons       []string          `json:"cephMons,omitempty"`
	CephVersion    string            `json:"cephVersion"`
	CreateDate     string            `json:"createDate"`
	CrushTunables  string            `json:"crushTunables,omitempty"`
	Iam            Iam               `json:"iam,omitempty"`
	Label          string            `json:"label,omitempty"`
	Region         string            `json:"region"`
	ServiceName    string            `json:"serviceName"`
	Size           int               `json:"size"`
	State          string            `json:"state"`
	Status         string            `json:"status"`
	UpdateDate     string            `json:"updateDate"`
}

func tableOvhCeph() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_ceph",
		Description: "Cephs services.",
		List: &plugin.ListConfig{
			Hydrate: listCeph,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    getCeph,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{Func: getCephInfo},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the ceph.",
			},
			{
				Name:        "ceph_mons",
				Type:        proto.ColumnType_JSON,
				Description: "Ceph monitors.",
				Transform:   transform.FromField("CephMons"),
			},
			{
				Name:        "ceph_version",
				Type:        proto.ColumnType_STRING,
				Description: "Ceph version.",
				Transform:   transform.FromField("CephVersion"),
			},
			{
				Name:        "create_date",
				Type:        proto.ColumnType_STRING,
				Description: "Creation date.",
				Transform:   transform.FromField("CreateDate"),
			},
			{
				Name:        "update_date",
				Type:        proto.ColumnType_STRING,
				Description: "Last update date.",
				Transform:   transform.FromField("UpdateDate"),
			},
			{
				Name:        "service_name",
				Hydrate:     getCephInfo,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceName"),
				Description: "Ceph Name.",
			},
			{
				Name:        "region",
				Hydrate:     getCephInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Region where the Ceph cluster is located.",
				Transform:   transform.FromField("Region"),
			},
			{
				Name:        "size",
				Hydrate:     getCephInfo,
				Type:        proto.ColumnType_INT,
				Description: "Size of the Ceph cluster in TB.",
				Transform:   transform.FromField("Size"),
			},
			{
				Name:        "state",
				Hydrate:     getCephInfo,
				Type:        proto.ColumnType_STRING,
				Description: "State of the Ceph cluster.",
				Transform:   transform.FromField("State"),
			},
			{
				Name:        "status",
				Hydrate:     getCephInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Status of the Ceph cluster.",
				Transform:   transform.FromField("Status"),
			},

		},
	}
}

func getCephInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	ceph := h.Item.(Ceph)

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_ceph.getCephInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/dedicated/ceph/%s", ceph.ID), &ceph)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_ceph.getCephInfo", err)
		return nil, err
	}

	return ceph, nil
}

func listCeph(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_ceph.listCeph", "connection_error", err)
		return nil, err
	}

	var cephsId []string
	err = client.Get("/dedicated/ceph", &cephsId)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_ceph.listCeph", err)
		return nil, err
	}

	for _, cephId := range cephsId {
		var ceph Ceph
		ceph.ID = cephId
		d.StreamListItem(ctx, ceph)
	}

	return nil, nil
}

func getCeph(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetStringValue()
	var ceph Ceph
	ceph.ID = id
	return ceph, nil
}
