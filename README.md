# lighthouse

Akin to a lighthouse alone at sea, this application is a standalone monitoring system that watches over the
Kubernetes cluster it runs in and reports potential issues.

| Environment variable  | Description                                        | Default value |
|-----------------------|----------------------------------------------------|---------------|
| `WEBHOOK_URL`         | Discord webhook URL where alerts will be sent to   | `""` Required |
| `INTERVAL_IN_MINUTES` | Number of minutes between each check               | `10`          |
| `DEBUG`               | Whether to enable debugging logs                   | `false`       |
| `ENVIRONMENT`         | Set to `dev` if you want to run lighthouse locally | `""`          |


## Debugging
To enable debugging logs, you may set the `DEBUG` environment variable to `true`
