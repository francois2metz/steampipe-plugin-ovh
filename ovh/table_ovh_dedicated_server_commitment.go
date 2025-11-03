package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// DedicatedServerCommitment represents a dedicated server and its commitment/renewal information
type DedicatedServerCommitment struct {
	ServerName  string     `json:"serverName"`
	ServiceId   int        `json:"serviceId"`
	Status      string     `json:"status"`
	RenewMode   string     `json:"renewMode"`
	RenewPeriod int        `json:"renewPeriod"`
	EngagedUpTo *time.Time `json:"engagedUpTo"`
	Expiration  *time.Time `json:"expiration"`
	Creation    *time.Time `json:"creation"`
	AccountId   string     `json:"accountId,omitempty"`
	Region      string     `json:"region,omitempty"`
}

func tableOvhDedicatedServerCommitment() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_dedicated_server_commitment",
		Description: "Lists all dedicated servers and their commitment/renewal information.",
		List: &plugin.ListConfig{
			Hydrate: listDedicatedServerCommitments,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("server_name"),
			Hydrate:    getDedicatedServerCommitment,
		},
		Columns: []*plugin.Column{
			{
				Name:        "server_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerName"),
				Description: "Dedicated server name (serviceName).",
			},
			{
				Name:        "service_id",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServiceId"),
				Description: "Internal OVH service identifier.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Service status.",
			},
			{
				Name:        "renew_mode",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RenewMode"),
				Description: "Renewal mode (manual, automatic, etc.).",
			},
			{
				Name:        "renew_period",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("RenewPeriod"),
				Description: "Renewal period in months.",
			},
			{
				Name:        "engaged_up_to",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EngagedUpTo"),
				Description: "Commitment end date (engagement expiration).",
			},
			{
				Name:        "expiration",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Service expiration date.",
			},
			{
				Name:        "creation",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Service creation date.",
			},
			{
				Name:        "account_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountId"),
				Description: "OVH account ID (if available from context).",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Description: "API region (e.g. ovh-eu).",
			},
		},
	}
}

func listDedicatedServerCommitments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_dedicated_server_commitment.listDedicatedServerCommitments", "connection_error", err)
		return nil, err
	}

	// Get list of all dedicated servers
	var servers []string
	err = client.Get("/dedicated/server", &servers)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_dedicated_server_commitment.listDedicatedServerCommitments", "api_error", err)
		return nil, err
	}

	// For each server, get its service info (commitment data)
	for _, serverName := range servers {
		commitment, err := getServerCommitmentInfo(ctx, client, serverName, d)
		if err != nil {
			plugin.Logger(ctx).Error("ovh_dedicated_server_commitment.listDedicatedServerCommitments", "server_info_error", err, "server_name", serverName)
			continue // Skip this server and continue with others
		}

		d.StreamListItem(ctx, commitment)
	}

	return nil, nil
}

func getDedicatedServerCommitment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	serverName := quals["server_name"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_dedicated_server_commitment.getDedicatedServerCommitment", "connection_error", err)
		return nil, err
	}

	commitment, err := getServerCommitmentInfo(ctx, client, serverName, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_dedicated_server_commitment.getDedicatedServerCommitment", "api_error", err)
		return nil, err
	}

	return commitment, nil
}

// getServerCommitmentInfo fetches service info for a specific dedicated server
func getServerCommitmentInfo(ctx context.Context, client *ovh.Client, serverName string, d *plugin.QueryData) (*DedicatedServerCommitment, error) {
	type ServiceInfoResponse struct {
		Renew struct {
			Automatic bool `json:"automatic"`
			Period    int  `json:"period"`
		} `json:"renew"`
		RenewalType string  `json:"renewalType"`
		EngagedUpTo *string `json:"engagedUpTo"`
		ServiceId   int     `json:"serviceId"`
		Status      string  `json:"status"`
		Expiration  *string `json:"expiration"`
		Creation    *string `json:"creation"`
	}

	var info ServiceInfoResponse
	err := client.Get(fmt.Sprintf("/dedicated/server/%s/serviceInfos", serverName), &info)
	if err != nil {
		return nil, fmt.Errorf("failed to get service info for server %s: %w", serverName, err)
	}

	// Determine renew mode from the API response
	renewMode := info.RenewalType
	if renewMode == "" {
		if info.Renew.Automatic {
			renewMode = "automatic"
		} else {
			renewMode = "manual"
		}
	}

	commitment := &DedicatedServerCommitment{
		ServerName:  serverName,
		ServiceId:   info.ServiceId,
		Status:      info.Status,
		RenewMode:   renewMode,
		RenewPeriod: info.Renew.Period,
	}

	// Parse timestamps
	if info.EngagedUpTo != nil && *info.EngagedUpTo != "" {
		if engagedUpTo, err := time.Parse("2006-01-02", *info.EngagedUpTo); err == nil {
			commitment.EngagedUpTo = &engagedUpTo
		}
	}

	if info.Expiration != nil && *info.Expiration != "" {
		if expiration, err := time.Parse("2006-01-02", *info.Expiration); err == nil {
			commitment.Expiration = &expiration
		}
	}

	if info.Creation != nil && *info.Creation != "" {
		if creation, err := time.Parse("2006-01-02", *info.Creation); err == nil {
			commitment.Creation = &creation
		}
	}

	// Set region if available from connection config
	ovhConfig := GetConfig(d.Connection)
	if ovhConfig.Endpoint != nil {
		commitment.Region = *ovhConfig.Endpoint
	}

	return commitment, nil
}
