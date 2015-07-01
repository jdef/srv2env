## Building

    $ go get github.com/jdef/srv2env

To generate a static build:

    $ env CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
    -ldflags '-extld ld -extldflags -static' -a -x .

## Usage

    $ ./srv2env {fully-qualified-service-domain-name} {executable} [ {args...} ]

## Examples

    $ ./srv2env _sip._udp.sip2sip.info sh -c 'env'|egrep -e '_HOST|_PORT'
    _SIP_HOST0=proxy.sipthor.net
    _SIP_PORT0=5060

    $ ./srv2env _sip._udp.qxip.net sh -c 'env'|egrep -e 'cname|_HOST|_PORT|_ENDP'
    _SIP_UDP_ENDPOINT1=sip://sbc.qxip.net:5060
    _SIP_UDP_ENDPOINT0=sip://sbc.qxip.net:5060
    _SIP_UDP_PORT1=5060
    _SIP_UDP_PORT0=5060
    _SIP_UDP_ENDPOINTS=sip://sbc.qxip.net:5060,sbc.qxip.net:5060
    _SIP_UDP_HOST0=sbc.qxip.net
    _SIP_UDP_HOST1=sbc.qxip.net
