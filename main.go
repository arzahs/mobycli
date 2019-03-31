package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"github.com/spf13/cobra"
	"log"
	"os"
	"text/tabwriter"
)

func main() {
	if reexec.Init() {
		return
	}
	moby, err := NewMobyClients()
	if err != nil {
		fmt.Printf("Moby connection error: %s\n", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Use:   "mobycli",
		Short: "mobycli is simple client for moby",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please, use one of this commands: \"run\" \"stop\" \"ps\"")
		},
	}

	var runCmd = &cobra.Command{
		Use:   "run [image name]",
		Short: "run container",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			imageName := args[0]
			ctx := context.Background()
			log.Printf("Pulling image: %s\n", imageName)
			output, err := moby.PullImage(ctx, imageName)
			if err != nil {
				log.Fatalf("Pulling image %s has error %s\n", imageName, err)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.DiscardEmptyColumns)
			if err := PrintToWriter(output, w); err != nil {
				log.Fatalf("reader error %s", err)
			}
			w.Flush()
			log.Printf("Creating containter for image: %s", imageName)
			containerID, err := moby.CreateContainer(ctx, imageName)
			if err != nil {
				log.Fatalf("Creating container %s has error %s\n", imageName, err)
			}
			log.Printf("Running containter: %s", containerID)
			err = moby.Run(ctx, containerID)
			if err != nil {
				log.Fatalf("Running container %s has error %s\n", containerID, err)
			}
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop [container id]",
		Short: "stop container",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			containerID := args[0]
			ctx := context.Background()
			log.Printf("Stopping container %s", containerID)
			err = moby.Stop(ctx, containerID)
			if err != nil {
				log.Fatalf("Stopping container %s has error %s\n", containerID, err)
			}
		},
	}

	var listCmd = &cobra.Command{
		Use:   "ps",
		Short: "list containers",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			containers, err := moby.List(ctx)
			if err != nil {
				log.Fatalf("%s\n", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.DiscardEmptyColumns)
			fmt.Fprintf(w, "CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAMES\n")
			for _, c := range containers {
				fmt.Fprintf(w, "%s\t%s\t\"%s...\"\t%s\t%s\t%s\t%s\n",
					c.ID[:12], c.Image, c.Command[:12], UnixToStr(c.Created), c.Status, PortsToStr(c.Ports), c.Names[0])

			}
			w.Flush()
		},
	}

	rootCmd.AddCommand(runCmd, stopCmd, listCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution error %s\n", err)
	}
}
