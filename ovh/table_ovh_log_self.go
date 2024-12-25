package ovh

import (
	"context"
	"fmt"
	"time"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type Log struct {
	ID              int       `json:"logId"`
	Date            time.Time `json:"date"`
	Account         string    `json:"account"`
	IP              string 	  `json:"ip"`
	Method          string    `json:"method"`
	Route           string    `json:"route"`
	Path            string    `json:"path"`
}

func tableOvhLog() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_log_self",
		Description: "Logs of your account.",
		List: &plugin.ListConfig{
			Hydrate: listLog,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    getLog,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{Func: getLogInfo},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the log.",
			},
			{
				Name:        "date",
				Hydrate:     getLogInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date of the log.",
			},
			{
				Name:        "account",
				Hydrate:     getLogInfo,
				Type:        proto.ColumnType_STRING,
				Description: "User performing the action.",
			},
			{
				Name:        "ip",
				Hydrate:     getLogInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Origin IP of the action.",
			},
			{
				Name:        "method",
				Hydrate:     getLogInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Method requested.",
			},
			{
				Name:        "route",
				Hydrate:     getLogInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Route used for the action.",
			},
			{
				Name:        "path",
				Hydrate:     getLogInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Path used for the action with project and object IDs.",
			},
		},
	}
}

func getLogInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	log := h.Item.(Log)

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_logs_self.getLogInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/me/api/logs/self/%s", strconv.Itoa(log.ID)), &log)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_logs_self.getLogInfo", err)
		return nil, err
	}

	return log, nil
}

func listLog(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_logs_self.listLog", "connection_error", err)
		return nil, err
	}

	var logsId []int
	err = client.Get("/me/api/logs/self", &logsId)

	if err != nil {
		plugin.Logger(ctx).Error("ovh_logs_self.listLog", err)
		return nil, err
	}

	for _, logId := range logsId {
		var log Log
		log.ID = logId
		d.StreamListItem(ctx, log)
	}

	return nil, nil
}

func getLog(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	strId := d.Quals.ToEqualsQualValueMap()["id"].GetStringValue()
	var log Log
	intId, err := strconv.Atoi(strId)
	if err != nil {
		return nil, err
	}
	log.ID = intId
	return log, nil
}
