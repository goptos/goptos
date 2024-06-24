package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/goptos/cli/codegen"
	"github.com/goptos/cli/goesive"
)

func main() {
	var genViewCmd = flag.NewFlagSet("genview", flag.ExitOnError)
	var genViewSrc = genViewCmd.String("src", ".", "source code directory")

	var buildCmd = flag.NewFlagSet("build", flag.ExitOnError)
	var buildDist = buildCmd.String("dist", "dist", "directory to serve")

	var serveCmd = flag.NewFlagSet("serve", flag.ExitOnError)
	var serveDist = serveCmd.String("dist", "dist", "directory to serve")
	var servePort = serveCmd.String("port", "8080", "port to listen on")

	flag.Parse()

	if len(os.Args) < 1 {
		fmt.Println("expected 'genview' or 'build' or 'serve'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "genview":
		genViewCmd.Parse(os.Args[2:])
		codegen.View(*genViewSrc)
	case "build":
		buildCmd.Parse(os.Args[2:])
		goesive.Build(*buildDist)
	case "serve":
		serveCmd.Parse(os.Args[2:])
		goesive.Serve(*serveDist, *servePort)
	default:
		fmt.Println("expected one of 'genview' or 'build' or 'serve'")
		os.Exit(1)
	}
}
