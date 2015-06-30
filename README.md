## Building

    $ go get github.com/jdef/srv2env

To generate a static build:

    $ env CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
    -ldflags '-extld ld -extldflags -static' -a -x .

## Usage

    $ ./srv2env {fully-qualified-service-domain-name} {executable} [ {args...} ]

    $ ./srv2env _sip._udp.sip2sip.info sh -c 'env'|egrep -e '_HOST|_PORT'
    _SIP_HOST0=proxy.sipthor.net.
    _SIP_PORT0=5060

