package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// To generate a static build:
//    env CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
//    -ldflags '-extld ld -extldflags -static' -a -x .
//
// Sample execution:
//    ./srv2env _sip._udp.sip2sip.info sh -c 'env'|egrep -e '_HOST|_PORT'
//    _SIP_HOST0=proxy.sipthor.net.
//    _SIP_PORT0=5060
func main() {
	if len(os.Args) < 2 {
		panic("missing required fqdn parameter")
	}
	if len(os.Args) < 3 {
		panic("missing required command line arguments")
	}

	fqdn := os.Args[1]
	cmd := os.Args[2]
	binary, err := exec.LookPath(cmd)
	if err != nil {
		panic(fmt.Sprintf("failed to locate command binary %q", cmd))
	}

	s := strings.SplitN(fqdn, ".", 3)
	prefix := strings.ToUpper(strings.Replace(s[0], "-", "_", -1))

	_, addrs, err := net.LookupSRV("", "", fqdn)
	if err != nil {
		panic(fmt.Sprintf("srv lookup failed: %v", err))
	}

	envvar := os.Environ()

	for i, a := range addrs {
		envvar = append(envvar, fmt.Sprintf("%s_HOST%d=%s\n", prefix, i, a.Target))
		envvar = append(envvar, fmt.Sprintf("%s_PORT%d=%d\n", prefix, i, a.Port))
	}

	err = syscall.Exec(binary, os.Args[2:], envvar)
	if err != nil {
		panic(fmt.Sprintf("failed to exec %q: %v", binary, err))
	}
}
