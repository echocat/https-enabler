# https-enabler

Simple proxy to add HTTPS to endpoints that only support HTTP.

## Initial idea

This was originally designed to solve the problem that [Prometheus exporters normally has no HTTPS support](https://prometheus.io/docs/introduction/faq/#why-don't-the-prometheus-server-components-support-tls-or-authentication?-can-i-add-those?).

If you want to be able to host your services in different datacenters and use the Internet for service-to-service communication, or you cannot guarantee that the network your services communicate with is safe and you do not want the overhead of VPN tunnels you'll need HTTPS (TLS) for a secure communication.

https-enabler solves the leak of TLS support in some tools. 

## Get it

Download the current version from the [releases page](https://github.com/echocat/https-enabler/releases/latest). For older version see [archive page](https://github.com/echocat/https-enabler/releases).

Example:
```bash
sudo curl -SL https://github.com/echocat/https-enabler/releases/download/v0.1.0/https-enabler-linux-amd64 \
    > /usr/bin/https-enabler
sudo chmod +x /usr/bin/https-enabler
```

## Use it

### Usage

```
Usage: https-enabler <flags> [<wrapped tool to start> [<args to pass to tool>]]
Flags:
  -connect.address string
        Address to connect to and proxy content to 'listen.address'.
  -listen.address string
        Address to serve the HTTPS socket for access from the outsite world. (default ":9000")
  -listen.ca string
        Path to PEM file that contains the CAs that are trused for incoming client connections.
        If provided: Connecting clients must present a certificate signed by one of these CAs.
        If not provided: Expects that 'listen.cert' also contains CAs to trust.
  -listen.cert string
        Path to PEM file that contains the certificate (and optionally also the private key in PEM format)
        to create the HTTPS socket with.
        The whole certificate chain must be included.
  -listen.private-key string
        Path to PEM file that contains the private-key.
        If not provided: The private key should be contained in the 'listen.cert' PEM file.
```

### Examples

```bash
# Connecting to a local HTTP only web server on port 8080 and exposing it on every network interface on port 8443.
# Expects the my.server.com.pem file to contain the client privateKey, the client certificate and the CA chain.
https-enabler -listen.address=:8443 \
    -listen.cert=my.server.com.pem \
    -connect.address=localhost:8080

# Connecting to a local HTTP only web server on port 8080 and exposing it on every network interface on port 8443.
# Expects the my.server.com.pem file to contain the client certificate and the CA chain.
# ... And the my.server.com.key file to contain only the client privateKey.
https-enabler -listen.address=:8443 \
    -listen.cert=my.server.com.pem \
    -listen.private-key=my.server.com.key \
    -connect.address=localhost:8080

# Like the first example but also start the prometheus node_exporter and exposes it on port 8443 in a secure way. The
# node_exporter is bound to localhost only to prevent access from the network.
https-enabler -listen.address=:8443 \
    -listen.cert=my.server.com.pem \
    -connect.address=localhost:9100 \
    node_exporter -web.listen-address=localhost:9100

```

## Build it

### Precondition

For building https-enabler you need:

1. a compatible operating system (Linux, Windows or Mac OS X)
2. and a working [Java 8](http://www.oracle.com/technetwork/java/javase/downloads/index.html) installation.

There is no need for a working and installed Go installation (or anything else). The build system will download every dependency and build it if necessary.

> **Hint:** The Go runtime build by the build system will be placed under ``~/.go/sdk``.

### Run build process

On Linux and Mac OS X:
```bash
# Build binaries (includes test)
./gradlew build

# Run tests (but do not build binaries)
./gradlew test

# Build binaries and release it on GitHub
# Environment variable GITHUB_TOKEN is required
./gradlew build githubRelease
```

On Windows:
```bash
# Build binaries (includes test)
gradlew build

# Run tests (but do not build binaries)
gradlew test

# Build binaries and release it on GitHub
# Environment variable GITHUB_TOKEN is required
gradlew build githubRelease
```

### Build artifacts

* Compiled and linked binaries can be found under ``./build/out/https-enabler-*``

## Contributing

https-enabler is an open source project of [echocat](https://echocat.org).
So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/https-enabler)
by [fork us](https://github.com/echocat/https-enabler/fork).

If you commit code to this project you have to accept that this code will be released under the [license](#license) of this project.


## License

See [LICENSE](LICENSE) file.
