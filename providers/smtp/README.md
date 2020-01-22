# SMTP Setup

1. Provide your email server and port settings in the [config.yaml](/config.yaml) under
  `providers.smtp`.
1. Comment out or leave `username` and `password` blank if using an unauthenticated smtp server.
1. Add the `from` address and 1 or more `to` addresses.
1. Optionally override the Email Subject.
1. Optionaly provide a `customTemplate` (path to file). See [Custom Template](#custom-template)
  for more info.

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
    # customTemplate: /path/to/email.html
```

## Default Template

The default email template looks like this:

![AWS News default email template](https://i.imgur.com/v6Ec7iK.png)

The footer can be customized in the [config.yaml](/config.yaml) `providers.smtp.footer`. For more
customization see [custom template](#custom-template).

## Custom Template

To provide a custom email html template, start with the `defaultTemplate` included in the code
and build from there:

```html
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
                      "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
  <head>
    <style>
    .footer {
      text-align: center;
    }
    .footer p {
      margin-top: 40px;
    }
    </style>
  </head>
  <body>
    <h3>AWS News for {{ .Date }}:</h3>
    {{ .News }}
    <div class="footer">
      <p>{{ .Footer }}</p>
    </div>
  </body>
</html>
```

There are 3 keys available to use:

|          |                                                             |
| -------- | ----------------------------------------------------------- |
| `Date`   | The date of yesterday's news.                               |
| `News`   | The AWS News in an unordered list of links.                 |
| `Footer` | Custom footer (overridable in [config.yaml](/config.yaml) ) |
