package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goptos/goptos/codegen"
	"github.com/goptos/goptos/goesive"
	"github.com/goptos/goptos/project"
)

func main() {
	const cliVersion = "v0.1.0"
	const msg = "expected 'version' or 'init' or 'genview' or 'package' or 'serve'"

	var initCmd = flag.NewFlagSet("init", flag.ExitOnError)
	var initVersion = initCmd.String("version", "latest", "version")

	var genViewCmd = flag.NewFlagSet("genview", flag.ExitOnError)
	var genViewSrc = genViewCmd.String("src", ".", "source code directory")

	var packageCmd = flag.NewFlagSet("package", flag.ExitOnError)
	var packageDist = packageCmd.String("dist", "dist", "directory to serve")

	var serveCmd = flag.NewFlagSet("serve", flag.ExitOnError)
	var serveDist = serveCmd.String("dist", "dist", "directory to serve")
	var servePort = serveCmd.String("port", "8080", "port to listen on")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println(msg)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("goptos cli version %s\n", cliVersion)
	case "init":
		log.Printf("initialising\n")
		initCmd.Parse(os.Args[2:])
		project.Init("goptos", "app", *initVersion)
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
		fmt.Println(msg)
		os.Exit(1)
	}
}
