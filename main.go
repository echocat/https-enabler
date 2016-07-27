package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	name        = "https-enabler"
	version     = "devel"
	description = ""
)

var (
	listenAddress = flag.String("listen.address", ":9000", "Address to listen on to serve the HTTPS socket to access from the outsite world.")
	listenCert    = flag.String("listen.cert", "", "Path to PEM file that contains the certificate (and optionally also the private key in PEM format)\n"+
		"    \tto create the HTTPS socket with.\n"+
		"    \tThis should include the whole certificate chain.")
	listenPrivateKey = flag.String("listen.private-key", "", "Path to PEM file that contains the private-key.\n"+
		"    \tIf not provided: The private key should be contained also in 'listen.cert' PEM file.")
	listenCa = flag.String("listen.ca", "", "Path to PEM file that conains the CAs that are trused for incoming client connections.\n"+
		"    \tIf provided: Connecting clients should present a certificate signed by one of this CAs.\n"+
		"    \tIf not provided: Expected that 'listen.cert' also contains CAs to trust.")

	connectAddress = flag.String("connect.address", "", "Address to connect to and proxy this content to 'listen.address'.")

	enclosedCommand          = ""
	enclosedCommandArguments = []string{}

	flagsBuffer = &bytes.Buffer{}
)

func main() {
	parseUsage()
	if len(strings.TrimSpace(enclosedCommand)) > 0 {
		runEnclosedCommand(enclosedCommand, enclosedCommandArguments)
	}
	startHttpsServer(*listenAddress, *listenCert, *listenPrivateKey, *listenCa)
}

func parseUsage() {
	plainFlags := parseArguments(os.Args[1:])
	flags := flag.CommandLine
	flags.SetOutput(flagsBuffer)
	flags.Usage = func() {
		errorString := flagsBuffer.String()
		if len(errorString) > 0 {
			printUsage(strings.TrimSpace(errorString))
		} else {
			printUsage(nil)
		}
	}
	if len(os.Args) <= 1 {
		fail(nil)
	}

	flags.Parse(plainFlags)
	assertUsage()
}

func parseArguments(plains []string) []string {
	result := []string{}
	inFlags := true
	for _, plain := range plains {
		if inFlags {
			if strings.HasPrefix(plain, "-") {
				result = append(result, plain)
			} else {
				inFlags = false
				enclosedCommand = plain
			}
		} else {
			enclosedCommandArguments = append(enclosedCommandArguments, plain)
		}
	}

	return result
}

func assertUsage() {
	if len(strings.TrimSpace(*listenAddress)) == 0 {
		fail("Missing --listen.address")
	}
	if len(strings.TrimSpace(*listenCert)) == 0 {
		fail("Missing --listen.cert")
	}
	if len(strings.TrimSpace(*connectAddress)) == 0 {
		fail("Missing --connect.address")
	}
}

func fail(err interface{}) {
	printUsage(err)
	os.Exit(1)
}

func printUsage(err interface{}) {
	fmt.Fprintf(os.Stderr, "%v (version: %v, url: https://github.com/echocat/https-enabler)\n", name, version)
	if description != "" {
		fmt.Fprintf(os.Stderr, "%v\n", description)
	}
	fmt.Fprint(os.Stderr, "Author(s): Gregor Noczinski (gregor@noczinski.eu)\n")
	fmt.Fprint(os.Stderr, "\n")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
	}

	fmt.Fprintf(os.Stderr, "Usage: %v <flags> [<enclosed tool to start> [<args to pass to tool>]]\n", os.Args[0])
	fmt.Fprint(os.Stderr, "Flags:\n")
	flag.CommandLine.SetOutput(os.Stderr)
	flag.CommandLine.PrintDefaults()
	flag.CommandLine.SetOutput(flagsBuffer)
}
