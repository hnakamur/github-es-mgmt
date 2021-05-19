github-es-mgmt
==============

[GitHub Enterprise Server Management Console API](https://docs.github.com/en/enterprise-server@3.0/rest/reference/enterprise-admin#management-console) client and CLI written in Go.

This project is open source but closed development.

## Usage

Set management console password to the environment variable `MGMT_PASSWORD`.

```
export MGMT_PASSWORD=_your_password_here_
```

### Set certificate

```
github-es-mgmt set-cert -endpoint https://your-github-es.example.jp:8443 -cert /path/to/your.crt -key /path/to/your.key
```

### Get maintenance status

```
github-es-mgmt get-maintenance -endpoint https://your-github-es.example.jp:8443
```

### Enable maintenance mode

```
github-es-mgmt set-maintenance -endpoint https://your-github-es.example.jp:8443 -enabled=true -when now
```

### Disable maintenance mode

```
github-es-mgmt set-maintenance -endpoint https://your-github-es.example.jp:8443 -enabled=false -when now
```

