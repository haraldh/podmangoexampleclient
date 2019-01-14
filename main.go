package main

import (
	"flag"
	"fmt"
	"github.com/haraldh/podmangoexampleclient/iopodman"
	"github.com/varlink/go/varlink"
	"os"
)

func help(name string) {
	fmt.Printf("Usage: %s <varlink address URL>\n", name)
	os.Exit(1)
}

func printError(methodname string, err error) {
	fmt.Fprintf(os.Stderr, "Error calling %s: ", methodname)
	switch e := err.(type) {
	case *iopodman.ImageNotFound:
		//error ImageNotFound (name: string)
		fmt.Fprintf(os.Stderr, "'%v' name='%s'\n", e, e.Name)

	case *iopodman.ContainerNotFound:
		//error ContainerNotFound (name: string)
		fmt.Fprintf(os.Stderr, "'%v' name='%s'\n", e, e.Name)

	case *iopodman.NoContainerRunning:
		//error NoContainerRunning ()
		fmt.Fprintf(os.Stderr, "'%v'\n", e)

	case *iopodman.PodNotFound:
		//error PodNotFound (name: string)
		fmt.Fprintf(os.Stderr, "'%v' name='%s'\n", e, e.Name)

	case *iopodman.PodContainerError:
		//error PodContainerError (podname: string, errors: []PodContainerErrorData)
		fmt.Fprintf(os.Stderr, "'%v' podname='%s' errors='%v'\n", e, e.Podname, e.Errors)

	case *iopodman.NoContainersInPod:
		//error NoContainersInPod (name: string)
		fmt.Fprintf(os.Stderr, "'%v' name='%s'\n", e, e.Name)

	case *iopodman.ErrorOccurred:
		//error ErrorOccurred (reason: string)
		fmt.Fprintf(os.Stderr, "'%v' reason='%s'\n", e, e.Reason)

	case *iopodman.RuntimeError:
		//error RuntimeError (reason: string)
		fmt.Fprintf(os.Stderr, "'%v' reason='%s'\n", e, e.Reason)

	case *varlink.InvalidParameter:
		fmt.Fprintf(os.Stderr, "'%v' parameter='%s'\n", e, e.Parameter)

	case *varlink.MethodNotFound:
		fmt.Fprintf(os.Stderr, "'%v' method='%s'\n", e, e.Method)

	case *varlink.MethodNotImplemented:
		fmt.Fprintf(os.Stderr, "'%v' method='%s'\n", e, e.Method)

	case *varlink.InterfaceNotFound:
		fmt.Fprintf(os.Stderr, "'%v' interface='%s'\n", e, e.Interface)

	case *varlink.Error:
		fmt.Fprintf(os.Stderr, "'%v' parameters='%v'\n", e, e.Parameters)

	default:
		fmt.Fprintf(os.Stderr, "%T - '%v'\n", err, err)
	}
}

func main() {
	var bridge bool
	flag.BoolVar(&bridge, "bridge", false, "Use bridge for connection")
	flag.Parse()
	var c *varlink.Connection
	var err error

	if bridge {
		c, err = varlink.NewBridge(flag.Arg(0))
	} else {
		c, err = varlink.NewConnection(flag.Arg(0))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to '%s': '%v'", flag.Arg(0), err)
		os.Exit(1)
	}

	info, err := iopodman.GetInfo().Call(c)

	if err != nil {
		printError("GetInfo()", err)
		os.Exit(1)
	}

	fmt.Printf("Info: %+v\n\n", info)

	fmt.Printf("Podman Version: %+v\n\n", info.Podman.Podman_version)

	containers, err := iopodman.ListContainers().Call(c)

	if err != nil {
		printError("ListContainers()", err)
		os.Exit(1)
	}

	for container := range containers {
		print(container)
	}

	mount, err := iopodman.MountContainer().Call(c, "foo")
	if err != nil {
		printError("MountContainer()", err)
	} else {
		print(mount)
	}

	// Test error handling
	{
		var in struct {
			test string
		}
		var out string
		in.test = "test"
		err = c.Call("io.unknowninterface.Strangemethod", in, out)
		if err != nil {
			printError("io.unknowninterface.Strangemethod()", err)
		}
	}
	{
		var in struct {
			test string
		}
		var out string
		in.test = "test"
		err = c.Call("io.podman.UnknownMethod", in, out)
		if err != nil {
			printError("io.podman.UnknownMethod()", err)
		}
	}
}
