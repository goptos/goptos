package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/goptos/goptos/codegen"
	"github.com/goptos/goptos/goesive"
	"github.com/goptos/goptos/project"
)

const cliVersion = "v0.1.4"

func main() {
	const msg = "expected 'version' or 'init' or 'genview' or 'build' or 'package' or 'serve'"

	var initCmd = flag.NewFlagSet("init", flag.ExitOnError)
	var initVersion = initCmd.String("version", "latest", "version")

	var genViewCmd = flag.NewFlagSet("genview", flag.ExitOnError)
	var genViewSrc = genViewCmd.String("src", ".", "source code directory")

	var buildCmd = flag.NewFlagSet("build", flag.ExitOnError)
	var buildSrc = buildCmd.String("src", "src", "source code directory")

	var packageCmd = flag.NewFlagSet("package", flag.ExitOnError)
	var packageDist = packageCmd.String("dist", "dist", "distribution directory")

	var serveCmd = flag.NewFlagSet("serve", flag.ExitOnError)
	var serveSrc = serveCmd.String("src", "src", "source code directory")
	var serveDist = serveCmd.String("dist", "dist", "distribution directory (to serve)")
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
		initCmd.Parse(os.Args[2:])
		project.Init("goptos", "app", *initVersion)
	case "genview":
		genViewCmd.Parse(os.Args[2:])
		codegen.View(*genViewSrc)
	case "build":
		buildCmd.Parse(os.Args[2:])
		goesive.Build(*buildSrc)
	case "package":
		packageCmd.Parse(os.Args[2:])
		goesive.Pack(*packageDist)
	case "serve":
		serveCmd.Parse(os.Args[2:])
		goesive.Serve(*serveSrc, *serveDist, *servePort)
	default:
		fmt.Println(msg)
		os.Exit(1)
	}
}
