## v0.11.0 [2025-03-01]

_What's new?_

- Update go to 1.23
- Update steampipe sdk to 5.11.3

## v0.10.1 [2024-12-30]

_What's new?_

- Fix `ovh_log_self` log messages and avoid int to string conversion.

## v0.10.0 [2024-12-29]

_What's new?_

- Add `ovh_log_self` table. Thanks [@romain-pix-cyber](https://github.com/romain-pix-cyber).
- Update SDK to v5.11.0

## v0.9.0 [2024-06-20]

_What's new?_

- Add `ovh_refund` and `ovh_refund_detail` tables. Thanks [@jdenoy](https://github.com/jdenoy).
- Update SDK to v5.10.1

## v0.8.0 [2024-04-27]

_What's new?_

- Renamed `ovh_cloud_storage` table to `ovh_cloud_storage_swift`. Thanks [@jdenoy](https://github.com/jdenoy).
- Add `ovh_cloud_storage_s3` table. Thanks [@jdenoy](https://github.com/jdenoy).
- Add `ovh_cloud_region` table.
- Update SDK to v5.9.0

## v0.7.0 [2024-01-08]

_What's new?_

- Add `ovh_cloud_volume_snapshot` table. Thanks [@jdenoy](https://github.com/jdenoy).
- Update SDK to v5.8.0

## v0.6.0 [2023-10-15]

_What's new?_

- Update SDK to v5.6.2
- Update go to 1.21

## v0.5.1 [2023-03-15]

_What's new?_

- Rename `start` and `end` columns of the `ovh_bill_detail` to `period_start` and `period_end`.

## v0.5.0 [2023-03-14]

_What's new?_

- Add `ovh_bill_detail` table. Thanks @emeric-martineau.

## v0.4.0 [2023-03-08]

_What's new?_

- Add `ovh_bill` table. Thanks @emeric-martineau.
- Update SDK to v5.2.0

## v0.3.0 [2023-02-17]

_What's new?_

- Add `ovh_cloud_ai_app`, `ovh_cloud_ai_job` and `ovh_cloud_ai_notebook` tables. Thanks @Benzhaomin.
- Update SDK to v5

## v0.2.0 [2022-09-01]

_What's new?_

- Update SDK to 4.1.5
- Update to go 1.19
- The default API key is commented

## v0.1.0 [2022-06-15]

_What's new?_

- Add table ovh_cloud_database table
- Add errors logs to all API calls
- Update steampipe sdk to 3.2.0

## v0.0.3 [2022-05-15]

_What's new?_

- Rename table ovh_cloud_sshkey table to ovh_cloud_ssh_key
- Documentation updates

## v0.0.2 [2022-05-06]

_What's new?_

- Update steampipe sdk to 3.1.0
- Update to go 1.18
- Build ARM64 binaries

## v0.0.1 [2022-04-02]

_What's new?_

- New tables added

  - ovh_cloud_data_job
  - ovh_cloud_flavor
  - ovh_cloud_image
  - ovh_cloud_instance
  - ovh_cloud_postgres
  - ovh_cloud_project
  - ovh_cloud_sshkey
  - ovh_cloud_storage
  - ovh_cloud_volume
