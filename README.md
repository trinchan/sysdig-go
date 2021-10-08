# sysdig-go #

[![sysdig-go release](https://img.shields.io/github/v/release/trinchan/sysdig-go?sort=semver)](https://github.com/trinchan/sysdig-go/releases)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/trinchan/sysdig-go/github)
[![Test Status](https://github.com/trinchan/sysdig-go/workflows/tests/badge.svg)](https://github.com/trinchan/sysdig-go/actions?query=workflow%3Atests)
[![Test Coverage](https://codecov.io/gh/trinchan/sysdig-go/branch/main/graph/badge.svg)](https://codecov.io/gh/trinchan/sysdig-go)

`sysdig-go` is a Go client library for accessing the [Sysdig](https://docs.sysdig.com/en/docs/developer-tools) and [IBM Cloud Monitoring](https://cloud.ibm.com/apidocs/monitor) APIs.

## Installation ##

```bash
go get github.com/trinchan/sysdig-go
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/trinchan/sysdig-go/github"
```

and run `go get` without parameters.

## Usage ##

```go
import "github.com/trinchan/sysdig-go/sysdig"
```

Construct a new Sysdig client, then use the various services on the client to
access different parts of the API. For example:

```go
authenticator, err := accesstoken.Authenticator(accessToken)
if err != nil {
    // handle error
}
client, err := sysdig.NewClient(sysdig.WithAuthenticator(authenticator))
if err != nil {
    // handle error
}

// Get the current user
me, _, err := client.Users.Me(context.Background())
if err != nil {
    // handle error
}
fmt.Printf("Logged in as %s %s", me.User.FirstName, me.User.LastName)
```

Some API methods have optional parameters that can be passed. For example:

```go
opts := sysdig.EventOptions{
    Name: "Event Name",
    Description: "Event Description",
    Severity: sysdig.SeverityInfo,
}
event, _, err := client.Events.Create(context.Background(), opts)
if err != nil {
    // handle error
}
fmt.Printf("Created event: %s", event.Event.ID)
```

The services of a client divide the API into logical chunks and correspond roughly to
the structure of the Sysdig API paths.

NOTE: Using the [context](https://pkg.go.dev/context) package, one can easily
pass cancellation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

For more sample code snippets, head over to the
[example](https://github.com/trinchan/sysdig-go/tree/master/example) directory.

### Authentication ###

The sysdig-go library handles authentication through an `Authenticator` interface defined in the
[authentication](https://github.com/trinchan/sysdig-go/tree/master/sysdig/authentication) package.
When creating a new client, pass an `authentication.Authenticator` that can handle authentication for
you. There are two methods of authentication supported.

The `accesstoken` subpackage authenticates each request with the provided [Sysdig API Token](https://docs.sysdig.com/en/docs/administration/administration-settings/user-profile-and-password/retrieve-the-sysdig-api-token).
```go
func main() {
    authenticator, err := accesstoken.Authenticator(accessToken)
    if err != nil {
        // handle error
    }
    client, err := sysdig.NewClient(sysdig.WithAuthenticator(authenticator))
}
```
The `ibmiam` subpackage authenticates each request with an [IBM Cloud IAM Token](https://cloud.ibm.com/docs/monitoring?topic=monitoring-api_token). It automatically retrieves and refreshes an IAM Token using an IBM Cloud API Key.

```go
authenticator, err := ibmiam.Authenticator(
    apiKey,
    ibmiam.WithIBMInstanceID(instanceID),
)
if err != nil {
    panic(err)
}
client, err := sysdig.NewClient(
    sysdig.WithAuthenticator(authenticator),
    sysdig.WithIBMBaseURL(sysdig.RegionUSSouth, false),
)
```

See the [example](https://github.com/trinchan/sysdig-go/tree/master/example) directory for more authentication examples.

## Prometheus API ##

Sysdig offers a limited [Prometheus HTTP API](https://prometheus.io/docs/prometheus/latest/querying/api).

This API is exposed via the `PrometheusClient()` function of the Client. You can use this client to run PromQL queries against your Sysdig instance.

_Most_ functionality of the HTTP API is not available from Sysdig, but they appear to be offering more and more.
See the [Prometheus example](https://github.com/trinchan/sysdig-go/tree/master/example/prometheus).

## Client Options ##

### Debug Mode

Debug mode can be enabled by setting the following client options:

```go
sysdig.NewClient(
    sysdig.WithDebug(true),
	sysdig.WithLogger(log.Default()), // Or other Logger
)
```

Setting `Debug` mode will print out the request URLs and response body and headers, along with other debug information.

This is useful for debugging parse issues and during development.

### Compression ###

The Sysdig API (and this client) supports [gzip](https://docs.sysdig.com/en/docs/developer-tools/sysdig-rest-api-conventions/#encoding) to reduce the size of responses. This can be useful for large queries.

```go
sysdig.NewClient(
    sysdig.WithResponseCompression(true),
)
```

For other options, check the [documentation](https://pkg.go.dev/github.com/trinchan/sysdig-go/github).

## Versioning ##

`sysdig-go` is currently undergoing its initial development and is still incomplete. As new APIs are documented by Sysdig and IBM Cloud, new
APIs will be added or changed. Since `sysdig-go` is a client library, breaking changes in the upstream API may require updates to the client.
`sysdig-go` will follow semver as closely as possible to minimize breaking changes.

## FAQ ##

### "Can you add X API?"

Yes! Open an issue with the API path and as much information about it as you can for a better chance of it getting developed. Or better yet, submit a patch!. In the mean time,
 you can also use the `client.Do()` method to send a custom request.

### "The response for this API is wrong/broken!" ###

That's not a question! The client is incomplete as documentation for most of the Sysdig API has not been published. I have had to leave some types as `interface{}` until documentation is released or I receive a sample response. Submit an issue and include the (redacted) client logs with `Debug` mode enabled.

### "Is this an official client?" ###

Nope. The only official client I know is the [Python SDK](https://github.com/sysdiglabs/sysdig-sdk-python).


## Credits ##
- Sysdig's [Python SDK](https://github.com/sysdiglabs/sysdig-sdk-python) for API reference.
- Google's [Github Client](https://github.com/google/go-github) for client and repo design reference.

## License ##

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.
