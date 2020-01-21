# SMTP Setup

1. Provide your email server and port settings in the [config.yaml](/config.yaml) under
  `providers.smtp`.
1. Comment out or leave `username` and `password` blank if using an unauthenticated smtp server.
1. Add the `from` address and 1 or more `to` addresses.
1. Optionally override the Email Subject.

## Example Config

```yaml
  smtp:
    enabled: true
    server: "smtp.gmail.com"
    port: "587"
    username: "mail@gmail.com"
    password: "abcdefghijklmno"
    subject: "AWS News"
    footer: "Brought to you by <a href='https://github.com/circa10a/go-aws-news'>go-aws-news</a>"
    from: "mail@gmail.com"
    to:
      - "some.email@mail.com"
      - "some.other@mail.com"
```
