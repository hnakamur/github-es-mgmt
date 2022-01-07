github-es-mgmt
==============

[GitHub Enterprise Server Management Console API](https://docs.github.com/en/enterprise-server@3.0/rest/reference/enterprise-admin#management-console) client and CLI written in Go.

This project is open source but closed development.

## Usage

You can set the management console password to the environment variable `MGMT_PASSWORD`.
Or you can input the password at the prompt `Enter Management Console password: `.

```
export MGMT_PASSWORD=_your_password_here_
```

### Set certificate

```
github-es-mgmt certificate set -endpoint https://your-github-es.example.jp:8443 -cert /path/to/your.crt -key /path/to/your.key
```

### Get maintenance status

```
github-es-mgmt maintenance status -endpoint https://your-github-es.example.jp:8443
```

### Enable maintenance mode

```
github-es-mgmt maintenance enable -endpoint https://your-github-es.example.jp:8443 -when now
```

### Disable maintenance mode

```
github-es-mgmt maintenance disable -endpoint https://your-github-es.example.jp:8443 -when now
```

