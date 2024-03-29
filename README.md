# lighthouse
Akin to a lighthouse lost at sea, this application is a standalone monitoring system that watches over the
Kubernetes cluster it runs in and reports potential issues.

![Discord alert](.github/assets/discord-alert.png)

| Environment variable  | Description                                        | Default value |
|-----------------------|----------------------------------------------------|---------------|
| `WEBHOOK_URL`         | Discord webhook URL where alerts will be sent to   | `""` Required |
| `INTERVAL_IN_MINUTES` | Number of minutes between each check               | `10`          |
| `DEBUG`               | Whether to enable debugging logs                   | `false`       |
| `ENVIRONMENT`         | Set to `dev` if you want to run lighthouse locally | `""`          |


## Installation
```console
helm repo add twin https://twin.github.io/helm-charts
helm repo update
helm install lighthouse twin/lighthouse -n kube-system
```


## Debugging
To enable debugging logs, you may set the `DEBUG` environment variable to `true`
