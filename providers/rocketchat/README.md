# Rocket.Chat Setup

This provider implementation uses Rocket.Chat's [Incoming WebHook][webhook] integration.

To create a new Incoming WebHook:

1. Go to your admin page.
1. Go to Integrations.
1. Create a New Integration and select Incoming WebHook.
1. Name the integration "AWS News" or something similar to identify the service.
1. Select the channel where you will receive the alerts.
1. __Post as:__ `rocket.cat` or any other existing user.
1. Leave the Script Disabled.
1. Save the integration.
1. Copy the generated Webhook URL and paste it in to the [config.yaml](/config.yaml) under
  `providers.rocketchat.webhookURL`.
1. Enable the Rocket.Chat provider by setting `providers.rocketchat.enabled` to `true`.
1. Optionally override the default AWS Logo used on the Avatar for posts.

## Example Config

```yaml
providers:
  rocketchat:
    enabled: true
    webhookURL: "https://rocket.chat.server/hooks/{token}"
    iconURL: "{ override url to aws logo }"
```

[webhook]:https://rocket.chat/docs/administrator-guides/integrations/
