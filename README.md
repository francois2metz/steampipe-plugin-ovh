# OVH plugin for Steampipe

Use SQL to query infrastructure from [OVH][].

- **[Get started →](docs/index.md)**
- Documentation: [Table definitions & examples](docs/tables)

## Quick start

Install the plugin with [Steampipe][]:

    steampipe plugin install francois2metz/ovh

## Development

To build the plugin and install it in your `.steampipe` directory

    make

Copy the default config file:

    cp config/ovh.spc ~/.steampipe/config/ovh.spc

## License

Apache 2

[steampipe]: https://steampipe.io
[ovh]: https://www.ovhcloud.com/fr/
