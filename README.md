github-es-mgmt
==============

This is a command line tool for setting the certificate and the key for the GitHub Enterprise Server.
It uses [REST API endpoints for managing GitHub Enterprise Server - GitHub Enterprise Server 3.15 Docs](https://docs.github.com/en/enterprise-server@3.15/rest/enterprise-admin/manage-ghes?apiVersion=2022-11-28).
It is written in [Go](https://go.dev/).

This project is open source but closed development.

## Usage

### Set and apply certificate and key

```
printf "%s\n%s\n" _YOUR_USERNAME_ _YOUR_PASSWORD_ \
  | github-es-mgmt certificate set --apply --endpoint https://your-github-es.example.jp:8443 --cert /path/to/your.crt --key /path/to/your.key
```
