![Simple HTTP Redirect Middleware for Traefik](./.assets/banner.png)

# Simple HTTP Redirect Middleware

This is a simple middleware for [Traefik](https://traefik.io/) that allows you to set up
client redirection rules for HTTP routers. You can define multiple rules per router.
Rules are evaluated in the order of configuration; the first rule that matches is executed.
If no rule matches, the request is forwarded to the next middleware or the service.

If you're certain that at least one rule will match, you might use `noop@internal` as service
for the router.

# Configuration

## Guide

The rules are a list under the key `redirections`. Each rule consists of two required and one
optional value:

| Field name | Description                             | Default value   | Example           |
|:-----------|:----------------------------------------|:----------------|:------------------|
| `from`     | Regular expression for the incoming URL | _(required)_    | `/old/(.+)`       |
| `to`       | Replacement rule for the target URL     | _(required)_    | `/new/${1}`       |
| `code`     | HTTP status code for the redirect       | `307`           | `302`             |

## Details

### The _required_ `from` and `to` fields

Under the hood, `from` is a Go regular expression, and `to` is a replacement string for the
[`regexp.ReplaceAllString`](https://pkg.go.dev/regexp#Regexp.ReplaceAllString) method. This
allows the administrator to use matched elements from the incoming URL.

### The _optional_ `code` field

This is the HTTP status code that is sent alongside with the `Location` header to let the client
execute the redirection. Please refer to the
[MDN Guide on Redirection](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Redirections)
to find out the right status code for your specifiv scenario.

## Example

The following example redirects any URL from `old.example.com` to `new.example.com` as client-side
permanent redirection.

### Static configuration

```yaml
experimental:
  plugins:
    TmHttpRedirectPlugin:
      moduleName: "github.com/edelbluth/tm_http_redirect"
      version: "v0.2.1"
```

### Dynamic configuration

```yaml
http:
  middlewares:
    OldWebsiteHttpRedirectMiddleware:
      plugin:
        TmHttpRedirectPlugin:
          redirections:
            - from: "/(.*)"  # Redirect any request...
              to: "https://new.example.com/${1}"

  routers:
    OldWebsiteRedirectionRouterHttps:
      rule: "Host(`old.example.com`)"
      entryPoints:
        - "websecure"
      middlewares:
        - "OldWebsiteHttpRedirectMiddleware@file"
      service: "noop@internal"
      tls: {}
    OldWebsiteRedirectionRouterHttp:
      rule: "Host(`old.example.com`)"
      entryPoints:
        - "web"
      middlewares:
        - "OldWebsiteHttpRedirectMiddleware@file"
      service: "noop@internal"
```
