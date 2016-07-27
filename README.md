# https-enabler

Simple proxy to add HTTPS to endpoints that only supports HTTP.

## Initial idea

This was originally designed to solve the problem that [Prometheus exporters normally has no
 HTTPS support](https://prometheus.io/docs/introduction/faq/#why-don't-the-prometheus-server-components-support-tls-or-authentication?-can-i-add-those?).
 
But in a world wide network of services that should communicate directly with each other (without overhead of VPN tunnels that does not scale) we need HTTPS (TLS) for a secure communication.

https-enabler solve the leak of support in these tools. 

## Usage

```
Usage: https-enabler <flags> [<enclosed tool to start> [<args to pass to tool>]]
Flags:
  -connect.address string
        Address to connect to and proxy this content to 'listen.address'.
  -listen.address string
        Address to listen on to serve the HTTPS socket to access from the outsite world. (default ":9000")
  -listen.ca string
        Path to PEM file that conains the CAs that are trused for incoming client connections.
        If provided: Connecting clients should present a certificate signed by one of this CAs.
        If not provided: Expected that 'listen.cert' also contains CAs to trust.
  -listen.cert string
        Path to PEM file that contains the certificate (and optionally also the private key in PEM format)
        to create the HTTPS socket with.
        This should include the whole certificate chain.
  -listen.private-key string
        Path to PEM file that contains the private-key.
        If not provided: The private key should be contained also in 'listen.cert' PEM file.
```

## Examples

```bash
# Connect to a local HTTP only web server on 8080 and expose it on every network interface on 8443.
# We expect that my.server.com.pem contains privateKey, certificate and CA chain.
https-enabler -listen.address=:8443 \
    -listen.cert=my.server.com.pem \
    -connect.address=localhost:8080

# Connect to a local HTTP only web server on 8080 and expose it on every network interface on 8443.
# We expect that my.server.com.pem contains certificate and CA chain.
# ... And my.server.com.key contains only the privateKey.
https-enabler -listen.address=:8443 \
    -listen.cert=my.server.com.pem \
    -listen.private-key=my.server.com.key \
    -connect.address=localhost:8080

# Like the first example but also start the prometheus node_exporter and connect it node_exporter
# itself only to localhost:8080 so nobody from the outside world can access it without client certificate.
https-enabler -listen.address=:8443 \
    -listen.cert=my.server.com.pem \
    -connect.address=localhost:9100 \
    node_exporter -web.listen-address=localhost:9100

```

## Build

### Precondition

For building https-enabler there is only:

1. a compatible operating system (Linux, Windows or Mac OS X)
2. and a working [Java 8](http://www.oracle.com/technetwork/java/javase/downloads/index.html) installation required.

There is no need for a working and installed Go installation (or anything else). The build system will download every dependency and build it if necessary.

> **Hint:** The Go runtime build by the build system will be placed under ``~/.go/sdk``.

### Run

On Linux and Mac OS X:
```bash
# Build binaries only
./gradlew build

# Run tests (includes compile)
./gradlew test
```

On Windows:
```bash
# Build binaries only
gradlew build

# Run tests (includes compile)
gradlew test
```

### Build artifacts

* Compiled and lined binaries can be found under ``./build/out/https-enabler-*``

## Contributing

https-enabler is an open source project of [echocat](https://echocat.org).
So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/https-enabler)
by [fork us](https://github.com/echocat/https-enabler/fork).

If you commit code to this project you have to accept that this code will be released under the [license](#license) of this project.


## License

See [LICENSE](LICENSE) file.
