package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goptos/cli/goptos/codegen"
	"github.com/goptos/cli/goptos/goesive"
)

func main() {
	var genViewCmd = flag.NewFlagSet("genview", flag.ExitOnError)
	var genViewSrc = genViewCmd.String("src", ".", "source code directory")

	var packageCmd = flag.NewFlagSet("package", flag.ExitOnError)
	var packageDist = packageCmd.String("dist", "dist", "directory to serve")

	var serveCmd = flag.NewFlagSet("serve", flag.ExitOnError)
	var serveDist = serveCmd.String("dist", "dist", "directory to serve")
	var servePort = serveCmd.String("port", "8080", "port to listen on")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("expected 'genview' or 'package' or 'serve'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "genview":
		log.Printf("generating\n")
		genViewCmd.Parse(os.Args[2:])
		codegen.View(*genViewSrc)
	case "package":
		log.Printf("packaging\n")
		packageCmd.Parse(os.Args[2:])
		goesive.Build(*packageDist)
	case "serve":
		log.Printf("serving\n")
		serveCmd.Parse(os.Args[2:])
		goesive.Serve(*serveDist, *servePort)
	default:
		fmt.Println("expected 'genview' or 'package' or 'serve'")
		os.Exit(1)
	}
}
