# Yammer Setup

See Yammers's [REST API Docs](https://developer.yammer.com/docs/messages-json-post) for information on how to post a message.

1. [Create a yammer app](https://developer.yammer.com/docs/getting-started)
1. [Get a bearer token](https://developer.yammer.com/docs/test-token)
1. [Get your group id](https://support.office.com/en-us/article/how-do-i-find-a-yammer-group-s-feedid-b0e49b2c-ca30-4025-b3bc-7bd764c3e2ec) to post under a specific topic in yammer.
1. Enable the Yammer provider by setting `providers.yammer.enabled` to `true`.

## Example Config

```yaml
providers:
  yammer:
    enabled: true
    APIURL: https://www.yammer.com/api/v1/messages.json
    groupID: 123456
    token: somesecrettext
```
