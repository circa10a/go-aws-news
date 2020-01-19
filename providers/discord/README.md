# Discord Setup

See Discord's [Intro to Webhooks][webhooks] for information on how to create a webhook.

1. Name the webhook "AWS News", or something similar to distinguish the service using it.
1. Copy the webhook URL and paste it in to the [config.yaml](/config.yaml) under
  `providers.discord.webhookURL`. The webhook format will look similar to this:

    ```yaml
    "https://discordapp.com/api/webhooks/{webhook.id}/{webhook.token}"
    ```

1. Enable the Discord provider by setting `providers.discord.enabled` to `true`.

## Example Config

```yaml
providers:
  discord:
    enabled: true
    webhookURL: "https://discordapp.com/api/webhooks/{webhook.id}/{webhook.token}"
```

[webhooks]:https://support.discordapp.com/hc/en-us/articles/228383668-Intro-to-Webhooks
