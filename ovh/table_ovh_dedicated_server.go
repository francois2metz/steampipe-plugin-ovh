package ovh

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOvhDedicatedServer(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_dedicated_server",
		Description: "OVH Dedicated Server inventory with hardware and configuration details.",
		List: &plugin.ListConfig{
			Hydrate: listDedicatedServers,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getDedicatedServer,
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getDedicatedServer,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name/hostname of the dedicated server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_id",
				Description: "The unique server ID.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDedicatedServer,
				Transform:   transform.FromField("ServerId"),
			},
			{
				Name:        "ip",
				Description: "The main IP address of the server.",
				Type:        proto.ColumnType_IPADDR,
				Hydrate:     getDedicatedServer,
				Transform:   transform.FromField("Ip"),
			},
			{
				Name:        "reverse",
				Description: "The reverse DNS hostname for the IP.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "state",
				Description: "The current state of the server (ok, error, etc.).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "power_state",
				Description: "The power state of the server (poweron, poweroff, etc.).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "monitoring",
				Description: "Whether monitoring is enabled for this server.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "os",
				Description: "The operating system installed on the server.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "datacenter",
				Description: "The datacenter where the server is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "region",
				Description: "The region where the server is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "availability_zone",
				Description: "The availability zone of the server.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "rack",
				Description: "The physical rack location of the server.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "commercial_range",
				Description: "The commercial range/model of the server.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "link_speed",
				Description: "The network link speed in Mbps.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDedicatedServer,
				Transform:   transform.FromField("LinkSpeed"),
			},
			{
				Name:        "support_level",
				Description: "The support level for this server.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "professional_use",
				Description: "Whether the server is for professional use.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "no_intervention",
				Description: "Whether interventions are disabled on this server.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "boot_id",
				Description: "The boot configuration ID.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getDedicatedServer,
				Transform:   transform.FromField("BootId"),
			},
			{
				Name:        "boot_script",
				Description: "The boot script configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "root_device",
				Description: "The root device configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "rescue_ssh_key",
				Description: "SSH key for rescue mode.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "rescue_mail",
				Description: "Email address for rescue notifications.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "new_upgrade_system",
				Description: "Whether the new upgrade system is enabled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "efi_bootloader_path",
				Description: "The EFI bootloader path if configured.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "iam_display_name",
				Description: "The IAM display name for the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Iam.DisplayName"),
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "iam_id",
				Description: "The IAM ID for the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Iam.Id"),
				Hydrate:     getDedicatedServer,
			},
			{
				Name:        "iam_urn",
				Description: "The IAM URN for the server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Iam.Urn"),
				Hydrate:     getDedicatedServer,
			},
		},
	}
}

//// LIST FUNCTION

func listDedicatedServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_dedicated_server.listDedicatedServers", "connection_error", err)
		return nil, err
	}

	var serverNames []string
	if err := client.Get("/dedicated/server", &serverNames); err != nil {
		plugin.Logger(ctx).Error("ovh_dedicated_server.listDedicatedServers", "api_error", err)
		return nil, err
	}

	for _, serverName := range serverNames {
		d.StreamListItem(ctx, DedicatedServer{Name: serverName})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDedicatedServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_dedicated_server.getDedicatedServer", "connection_error", err)
		return nil, err
	}

	var serverName string
	if h.Item != nil {
		server := h.Item.(DedicatedServer)
		serverName = server.Name
	} else {
		serverName = d.EqualsQuals["name"].GetStringValue()
	}

	var server DedicatedServer
	if err := client.Get(fmt.Sprintf("/dedicated/server/%s", serverName), &server); err != nil {
		plugin.Logger(ctx).Error("ovh_dedicated_server.getDedicatedServer", "api_error", err)
		return nil, err
	}

	return server, nil
}

//// STRUCTS

type DedicatedServer struct {
	Name              string  `json:"name"`
	ServerId          int     `json:"serverId"`
	Ip                string  `json:"ip"`
	Reverse           string  `json:"reverse"`
	State             string  `json:"state"`
	PowerState        string  `json:"powerState"`
	Monitoring        bool    `json:"monitoring"`
	Os                string  `json:"os"`
	Datacenter        string  `json:"datacenter"`
	Region            string  `json:"region"`
	AvailabilityZone  string  `json:"availabilityZone"`
	Rack              string  `json:"rack"`
	CommercialRange   string  `json:"commercialRange"`
	LinkSpeed         int     `json:"linkSpeed"`
	SupportLevel      string  `json:"supportLevel"`
	ProfessionalUse   bool    `json:"professionalUse"`
	NoIntervention    bool    `json:"noIntervention"`
	BootId            int     `json:"bootId"`
	BootScript        *string `json:"bootScript"`
	RootDevice        *string `json:"rootDevice"`
	RescueSshKey      *string `json:"rescueSshKey"`
	RescueMail        *string `json:"rescueMail"`
	NewUpgradeSystem  bool    `json:"newUpgradeSystem"`
	EfiBootloaderPath *string `json:"efiBootloaderPath"`
	Iam               IAM     `json:"iam"`
}

type IAM struct {
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
	Urn         string `json:"urn"`
}
