package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"bufio"
        "os"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/container"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/phayes/freeport"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type version struct {
    image string
    replication int
		servers []string
}

func findThisVersion(name string, allVersions []version) version {
	var v version
	for _, element := range allVersions {
		if element.image == name {
			return element
		}
	}
	return v
}

func newCreateCommand(dockerCli command.Cli) *cobra.Command {
	//////////////////////////////////////////////////////////////////////////
	///////////////////////////// START OF CHANGES ///////////////////////////
	//////////////////////////////////////////////////////////////////////////

	var runOpts container.RunOptions
	var copts *container.ContainerOptions
	execOpts := container.NewExecOptions()
	var restartOpts container.RestartOptions

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] IMAGE [COMMAND] [ARG...]",
		Short: "Create a new service",
		Args:  cli.RequiresMinArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// The format of the command:
			// ./build/docker service create e env_var=value hostIP network nginxName port ruleSet image1 rep1 memory1 swap1 cpu1 ...

			// list of enviroment variables
			index := 0
			my_env := []string{}
			for true {
				if args[index] == "e" {
					index++
					my_env = append(my_env, args[index])
					index++
				} else {
					break
				}
			}

			hostIP := args[index] // Host IP address that the containers will be created on
			index++
			network := args[index] // The name of an Overlay network to connect containers together
			index++
			nginxName := args[index] // The name of the service (The NGINX container)
			index++
			containerPort := args[index] // The port that the containers of the service listen on and need to expose it
			index++
			ruleAddress := args[index] // The absolute address of the file containing the rules for the load balancer
			index++

			nginx_image := "sgholami/nginx-monitoring" // The NGINX image with monitoring

			// Assigning the NGINX a random port
			nginxPort, err2 := freeport.GetFreePort()
			nginxPortStr := strconv.Itoa(nginxPort)
			if err2 != nil {
				log.Fatal(err2)
			}

			// Checking if a specific service port is mentioned
			// if not use the random port
			for _, element := range my_env {
				env := strings.Split(element, "=")
				if env[0] == "SERVICE_PORT" {
					nginxPortStr = env[1]
				}
			}

			// Creating containers
			var containers_version []version
			servers := ""

			for i := index; i < len(args); i++ {

				var v version

				copts.Image = args[i]
				v.image = args[i]
				i++

				rep, _ := strconv.Atoi(args[i])
				v.replication = rep
				i++

				memory := args[i] // The memory limit of this version of containers, if do not want to specify it type "_"
				i++

				swap := args[i] // The swap memory limit of this version of containers, if do not want to specify it type "_"
				i++

				cpu := args[i] // The CPU limit of this version of containers, if do not want to specify it type "_"

				for j := 0; j < rep; j++ {
					temp, err3 := freeport.GetFreePort()
					tmpPort := strconv.Itoa(temp)
					if err3 != nil {
						log.Fatal(err3)
					}

					// Set the published port value
					copts.Publish.Clear()
					copts.Publish.Set(tmpPort + ":" + containerPort)

					v.servers = append(v.servers, hostIP + ":" + tmpPort)

					// Set the network value
					copts.NetMode = network

					// Set the memory limit value
					if memory != "_" {
						copts.Memory.Set(memory)
					}

					// Set the swap memory limit value
					if swap != "_" {
						copts.MemorySwap.Set(swap)
					}

					// Set the CPU limit value
					if cpu != "_" {
						copts.Cpus.Set(cpu)
					}

					// Set the enviroment variables value
					copts.Env.Clear()
					for _, element := range my_env {
						copts.Env.Set(element)
					}

					// Assign a name to the container
					contName := nginxName + "_" + randomString(8)
					runOpts.CreateOption.Name = contName
					container.RunRun(dockerCli, cmd.Flags(), &runOpts, copts)

					// Keeping the servers for NGINX config file
					servers += "server " + hostIP + ":" + tmpPort + "; "
				}
				containers_version = append(containers_version, v)
			}

			// Creating the rule set
			ruleSet := ""
			file, err := os.Open(ruleAddress)
			if err != nil {
					log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				// Reading each line of rules
				line := scanner.Text()
				serverWeight := ""
				// parts[0] is the condition which is the same in rule-set
				// part[1] is the servers which should be updated
				parts:= strings.Split(line, " , ")
				upstreams := strings.Split(parts[1], ";")
				for l, element := range upstreams {
					if l < len(upstreams) - 1 {
						e := strings.TrimSpace(element)
						// upstream_parts[0] is version, not used
						// upstream_parts[1] is the image name
						// upstream_parts[2] is the percentage
						upstreamParts := strings.Split(e, " ")
						cont_v := findThisVersion(upstreamParts[1], containers_version)

						temp := strings.Split(upstreamParts[2], "=")

						totalWeight, _ := strconv.Atoi(temp[1])

						weight := int(math.Ceil(float64(totalWeight) / float64(len(cont_v.servers))))

						for _, server := range cont_v.servers {
							serverWeight = serverWeight + " server " + server + " weight=" + strconv.Itoa(weight) + ";"
						}
					}
				}
				// Writing the rule
				ruleSet = ruleSet + parts[0] + " ," + serverWeight + "\n"
			}
			ruleSet = strings.TrimSpace(ruleSet)

			// NGINX configuration
			var nginxOpts []string
			copts.Args = nginxOpts
			copts.Image = nginx_image
			copts.Publish.Clear()
			copts.Publish.Set(nginxPortStr + ":80") // NGINX listens on port 80. hostport:containerport
			runOpts.CreateOption.Name = nginxName
			container.RunRun(dockerCli, cmd.Flags(), &runOpts, copts)
			execOpts.Container = nginxName

			// The content of the config file of the NGINX
			conf := "echo 'upstream servers {" + servers + "} server {listen 80 default_server; access_log  /var/log/nginx/access.log  main; server_name _; location / { proxy_pass http://servers;}}' > /etc/nginx/conf.d/default.conf, echo '" + ruleSet + "' > rules.txt"
			execOpts.Command = []string{"sh", "-c", conf}

			container.RunExec(dockerCli, execOpts)

			restartOpts.Containers = []string{nginxName}
			restartOpts.NSecondsChanged = cmd.Flags().Changed("time")
			container.RunRestart(dockerCli, &restartOpts)

			//////////////////////////////////////////////////////////////////////////
			///////////////////////////// END OF CHANGES /////////////////////////////
			//////////////////////////////////////////////////////////////////////////
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
