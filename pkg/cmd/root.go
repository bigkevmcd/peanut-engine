package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/argoproj/gitops-engine/pkg/utils/errors"
	"github.com/argoproj/pkg/kube/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/bigkevmcd/peanut-engine/pkg/engine"
)

const (
	repoURLFlag = "repo-url"
	branchFlag  = "branch"
	pathFlag    = "path"
	portFlag    = "port"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func makeRootCmd() *cobra.Command {
	var (
		clientConfig clientcmd.ClientConfig
		cfg          engine.Config
		port         int
	)
	cmd := cobra.Command{
		Use: "peanut-engine",
		RunE: func(cmd *cobra.Command, args []string) error {
			resync := make(chan bool)
			config, err := clientConfig.ClientConfig()
			errors.CheckErrorWithCode(err, errors.ErrorCommandSpecific)
			if cfg.Namespace == "" {
				cfg.Namespace, _, err = clientConfig.Namespace()
				errors.CheckErrorWithCode(err, errors.ErrorCommandSpecific)
			}

			http.HandleFunc("/api/v1/sync", func(writer http.ResponseWriter, request *http.Request) {
				log.Infof("Synchronization triggered by API call")
				resync <- true
			})

			go func() {
				errors.CheckErrorWithCode(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", viper.GetInt(portFlag)), nil), errors.ErrorCommandSpecific)
			}()

			engine.PeanutSync(config, cfg, resync)
			return nil
		},
	}

	clientConfig = cli.AddKubectlFlagsToCmd(&cmd)

	cmd.Flags().StringVar(&cfg.Git.RepoURL, repoURLFlag, "", "Repository to deploy")
	logIfError(cmd.MarkFlagRequired(repoURLFlag))

	cmd.Flags().StringVar(&cfg.Git.Branch, branchFlag, "", "Branch to checkout")
	logIfError(cmd.MarkFlagRequired(branchFlag))

	cmd.Flags().StringVar(&cfg.Git.Path, pathFlag, "", "Path within the Repository to deploy")
	logIfError(cmd.MarkFlagRequired(pathFlag))

	cmd.Flags().DurationVar(&cfg.Resync, "resync-seconds", time.Second*300, "Resync duration")
	cmd.Flags().IntVar(&port, portFlag, 9001, "Port number.")
	cmd.Flags().BoolVar(&cfg.Prune, "prune", true, "Enables resource pruning.")
	logIfError(viper.BindPFlag(portFlag, cmd.Flags().Lookup(portFlag)))
	cmd.Flags().BoolVar(&cfg.Namespaced, "namespaced", false, "Switches agent into namespaced mode.")
	cmd.Flags().StringVar(&cfg.Namespace, "default-namespace", "",
		"The namespace that should be used if resource namespace is not specified."+
			"By default resources are installed into the same namespace where peanut-engine is installed.")
	return &cmd
}

func initConfig() {
	viper.AutomaticEnv()
}

// Execute is the main entry point into this component.
func Execute() {
	if err := makeRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

func logIfError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
