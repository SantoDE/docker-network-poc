package main

import (
	"context"
	"fmt"

	"github.com/Pallinder/sillyname-go"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	cresp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "traefik/whoami",
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, cresp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	cont, err := cli.ContainerInspect(ctx, cresp.ID)

	if err != nil {
		panic(err)
	}

	fmt.Println("Created container with names and networks")
	fmt.Println(cont.Name)
	fmt.Println(cont.NetworkSettings.Networks)

	sn := sillyname.GenerateStupidName()

	nwresp, err := cli.NetworkCreate(ctx, sn, types.NetworkCreate{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created network with name %s and id %s \n", sn, nwresp.ID)

	if err := cli.NetworkConnect(ctx, nwresp.ID, cresp.ID, &network.EndpointSettings{}); err != nil {
		panic(err)
	}

	cont, err = cli.ContainerInspect(ctx, cresp.ID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Attached new network to container. New networks on container %s \n", cont.Name)
	fmt.Println(cont.NetworkSettings)

	if err := cli.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{Force: true} ); err != nil {
		panic(err)
	}

	if err := cli.NetworkRemove(ctx, nwresp.ID); err != nil {
		panic(err)
	}
}
