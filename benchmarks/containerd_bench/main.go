package main

import (
	"context"
	"fmt"
	"log"
	"time"

	v2 "github.com/containerd/cgroups/v3/cgroup2/stats"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/typeurl/v2"
)

const (
	CONTAINER_NAME = "longrunningapp"
	FIFO_PATH      = "/home/alex/Desktop/dev/framework/benchmarks/containerd_bench/fifo"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	client, err := containerd.New("/run/containerd/containerd.sock")
	check(err)
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "default")

	image, err := client.Pull(ctx, "docker.io/alexcogojocaru/long-running-app:latest", containerd.WithPullUnpack)
	check(err)

	container, err := client.LoadContainer(ctx, CONTAINER_NAME)
	if err != nil {
		// create container with a new RW root filesystem
		container, err = client.NewContainer(
			ctx,
			CONTAINER_NAME,
			containerd.WithNewSnapshot(fmt.Sprintf("%s-rootfs", CONTAINER_NAME), image),
			containerd.WithNewSpec(oci.WithImageConfig(image)),
		)
		check(err)

		fmt.Printf("Created container %s\n", CONTAINER_NAME)
	} else {
		fmt.Printf("Using existing container %s\n", CONTAINER_NAME)
	}

	task, err := container.NewTask(
		ctx,
		cio.NewCreator(cio.WithFIFODir(FIFO_PATH)),
	)
	check(err)
	defer task.Delete(ctx)

	pid := task.Pid()
	fmt.Printf("pid=%d\n", pid)

	err = task.Start(ctx)
	check(err)

	go func() {
		time.Sleep(2 * time.Second)

		for {
			metric, err := task.Metrics(ctx)
			check(err)

			var data v2.Metrics
			err = typeurl.UnmarshalTo(metric.Data, &data)
			check(err)

			fmt.Println(data.Memory.Usage)
			time.Sleep(1 * time.Second)
		}
	}()

	status, err := task.Wait(ctx)
	check(err)
	<-status
}
