package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

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

	cname, addrs, err := net.LookupSRV("", "", fqdn)
	if err != nil {
		if dnserr, ok := err.(*net.DNSError); ok {
			fmt.Fprintln(os.Stderr, "error:", dnserr.Err)
			os.Exit(1)
		}
		panic(fmt.Sprintf("srv lookup failed: %#+v", err))
	}

	s := strings.SplitN(cname, ".", 3)
	prefix := strings.ToUpper(strings.Replace(s[0], "-", "_", -1))
	scheme := ""
	if strings.HasPrefix(s[0], "_") {
		scheme = s[0][1:] + "://"
	}
	if strings.HasPrefix(s[1], "_") {
		prefix += strings.ToUpper(strings.Replace(s[1], "-", "_", -1))
	}

	envvar := os.Environ()
	addons := make([]string, 4*len(addrs))
	eps := make([]string, len(addrs))

	for i, a := range addrs {
		j := i * 4
		ii := strconv.Itoa(i)
		port := strconv.Itoa(int(a.Port))
		host := a.Target
		if strings.HasSuffix(host, ".") {
			host = host[:len(host)-1]
		}
		hostport := net.JoinHostPort(host, port)

		addons[j] = prefix + "_HOST" + ii + "=" + host
		addons[j+1] = prefix + "_PORT" + ii + "=" + port
		addons[j+2] = prefix + "_ADDR" + ii + "=" + hostport
		addons[j+3] = prefix + "_ENDPOINT" + ii + "=" + scheme + hostport
		eps[i] = hostport
	}

	envvar = append(envvar, addons...)
	envvar = append(envvar, prefix+"_ENDPOINTS="+scheme+strings.Join(eps, ","))

	err = syscall.Exec(binary, os.Args[2:], envvar)
	if err != nil {
		panic(fmt.Sprintf("failed to exec %q: %v", binary, err))
	}
}
