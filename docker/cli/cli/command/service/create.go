package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/container"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/phayes/freeport"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newCreateCommand(dockerCli command.Cli) *cobra.Command {
	var runOpts container.RunOptions
	var copts *container.ContainerOptions
	execOpts := container.NewExecOptions()
	var restartOpts container.RestartOptions

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] IMAGE [COMMAND] [ARG...]",
		Short: "Create a new service", 
		Args:  cli.RequiresMinArgs(1), 
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(args)
			fmt.Println(len(args))
			////////////////////////////////////////////////////////////////////////////////////////////////////////
			//                                   NGINX, DOCKER RUN, NOT WORKING
			// ./build/docker service create hostip network nginxname containerport memory swap cpu image1 rep1 image2 rep2 ...
			// the firs image is multimedia
			// we only have one container of each
			hostIP := args[0]
			network := args[1]
			nginxName := args[2]
			memory := args[3]
			swap := args[4]
			cpu := args[5]
			fmt.Println("host ip ", hostIP)
			fmt.Println("network ", network)
			fmt.Println("nginxName ", nginxName)
			fmt.Println("args length: " + strconv.Itoa(len(args)))
			portList := []string{}
			nginxPort, err2 := freeport.GetFreePort()
			nginxPortStr := strconv.Itoa(nginxPort)
			// nginxPortStr = "80"
			fmt.Println("nginx port: " + nginxPortStr)
			if err2 != nil {
				log.Fatal(err2)
			}
			var nginxOpts []string
			for i := 6; i < len(args); i++ {
				fmt.Println("i: ", i)
				fmt.Println("image ", args[i])
				copts.Image = args[i]
				rep, _ := strconv.Atoi(args[i+1])
				fmt.Println("rep ", rep)
				i++
				containerPort := args[i+1]
				i++
				for j := 0; j < rep; j++ {
					fmt.Println("j: ", j)
					copts.Publish.Set(containerPort)
					copts.NetMode = network
					copts.Memory.Set(memory)
					copts.MemorySwap.Set(swap)
					copts.Cpus.Set(cpu)
					fmt.Println("net mode set")
					contName := nginxName + "_" + randomString(8)
					runOpts.CreateOption.Name = contName
					ports := strings.Split(containerPort, ":")
					portList = append(portList, ports[0])
					container.RunRun(dockerCli, cmd.Flags(), &runOpts, copts)
				}
			}
			fmt.Println("port list ", portList)

			copts.Args = nginxOpts
			fmt.Println("changing image")
			copts.Image = "sgholami/znn-nginx:latest"
			// copts.Image = "nginx"
			copts.Publish.Set(nginxPortStr + ":80") //hostport:containerport
			runOpts.CreateOption.Name = nginxName
			container.RunRun(dockerCli, cmd.Flags(), &runOpts, copts)

			execOpts.Container = nginxName

			servers_both := "server " + hostIP + ":" + portList[0] + " weight=99;server " + hostIP + ":" + portList[1] + " weight=1;"
			conf_multi := "echo 'upstream servers {" + servers_both + "} server {listen " + "80" + " default_server; access_log  /var/log/nginx/access.log  main; server_name _; location / { proxy_pass http://servers;}}' > /etc/nginx/conf.d/default.conf"

			execOpts.Command = []string{"sh", "-c", conf_multi}
			container.RunExec(dockerCli, execOpts)

			restartOpts.Containers = []string{nginxName}
			restartOpts.NSecondsChanged = cmd.Flags().Changed("time")
			container.RunRestart(dockerCli, &restartOpts)
			////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		},
	}
	flags := cmd.Flags()
	flags.SetInterspersed(false)

	// These are flags not stored in Config/HostConfig
	flags.BoolVarP(&runOpts.Detach, "detach", "d", true, "Run container in background and print container ID")
	flags.BoolVar(&runOpts.SigProxy, "sig-proxy", true, "Proxy received signals to the process")
	flags.StringVar(&runOpts.CreateOption.Name, "name", "", "Assign a name to the container")
	flags.StringVar(&runOpts.DetachKeys, "detach-keys", "", "Override the key sequence for detaching a container")

	// Add an explicit help that doesn't have a `-h` to prevent the conflict
	// with hostname
	flags.Bool("help", false, "Print usage")

	command.AddPlatformFlag(flags, &runOpts.CreateOption.Platform)
	command.AddTrustVerificationFlags(flags, &runOpts.CreateOption.Untrusted, dockerCli.ContentTrustEnabled())
	copts = container.AddFlags(flags)

	return cmd
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func runCreate(dockerCli command.Cli, flags *pflag.FlagSet, opts *serviceOptions) error {
	apiClient := dockerCli.Client()
	createOpts := types.ServiceCreateOptions{}

	ctx := context.Background()

	service, err := opts.ToService(ctx, apiClient, flags)
	if err != nil {
		return err
	}

	specifiedSecrets := opts.secrets.Value()
	if len(specifiedSecrets) > 0 {
		// parse and validate secrets
		secrets, err := ParseSecrets(apiClient, specifiedSecrets)
		if err != nil {
			return err
		}
		service.TaskTemplate.ContainerSpec.Secrets = secrets
	}

	specifiedConfigs := opts.configs.Value()
	if len(specifiedConfigs) > 0 {
		// parse and validate configs
		configs, err := ParseConfigs(apiClient, specifiedConfigs)
		if err != nil {
			return err
		}
		service.TaskTemplate.ContainerSpec.Configs = configs
	}

	if err := resolveServiceImageDigestContentTrust(dockerCli, &service); err != nil {
		return err
	}

	// only send auth if flag was set
	if opts.registryAuth {
		// Retrieve encoded auth token from the image reference
		encodedAuth, err := command.RetrieveAuthTokenFromImage(ctx, dockerCli, opts.image)
		if err != nil {
			return err
		}
		createOpts.EncodedRegistryAuth = encodedAuth
	}

	// query registry if flag disabling it was not set
	if !opts.noResolveImage && versions.GreaterThanOrEqualTo(apiClient.ClientVersion(), "1.30") {
		createOpts.QueryRegistry = true
	}

	response, err := apiClient.ServiceCreate(ctx, service, createOpts)
	if err != nil {
		return err
	}

	for _, warning := range response.Warnings {
		fmt.Fprintln(dockerCli.Err(), warning)
	}

	fmt.Fprintf(dockerCli.Out(), "%s\n", response.ID)

	if opts.detach || versions.LessThan(apiClient.ClientVersion(), "1.29") {
		return nil
	}

	return waitOnService(ctx, dockerCli, response.ID, opts.quiet)
}
