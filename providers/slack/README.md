# Slack Setup

1. In Slack, click "Add more apps" on the left sidebar.
1. In the search field, type "webhook".
1. Select __Incoming WebHooks__.
1. Click __Add to Slack__.
1. Select the channel to post to.
1. Click __Add Incoming WebHooks integration__.
1. Copy the generated webhook URL and paste it in to the [config.yaml](/config.yaml) under
  `providers.slack.webhookURL`. The webhook format will look similar to this:

## Example Config

```yaml
providers:
  slack:
    enabled: true
    iconURL: "https://cdn.iconscout.com/icon/free/png-256/aws-1869025-1583149.png"
    webhookURL: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
```

__Example post:__

![Slack post example](https://i.imgur.com/iw2SCJZ.png)


[webhooks]:https://api.slack.com/messaging/webhooks
