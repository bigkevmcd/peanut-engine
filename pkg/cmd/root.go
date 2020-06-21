package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/argoproj/gitops-engine/pkg/utils/errors"
	"github.com/argoproj/pkg/kube/cli"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/pkg/signals"

	"github.com/bigkevmcd/peanut-engine/pkg/engine"
	"github.com/bigkevmcd/peanut-engine/pkg/metrics"
)

const (
	repoURLFlag          = "repo-url"
	branchFlag           = "branch"
	pathFlag             = "path"
	portFlag             = "port"
	resyncFlag           = "resync"
	pruneFlag            = "prune"
	namespacedFlag       = "namespaced"
	defaultNamespaceFlag = "default-namespace"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func makeRootCmd() *cobra.Command {
	var (
		clientConfig clientcmd.ClientConfig
		cfg          engine.PeanutConfig
		gitCfg       engine.GitConfig
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

			http.Handle("/metrics", promhttp.Handler())
			http.HandleFunc("/api/v1/sync", func(writer http.ResponseWriter, request *http.Request) {
				log.Println("Synchronization triggered by API call")
				resync <- true
			})

			go func() {
				errors.CheckErrorWithCode(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", viper.GetInt(portFlag)), nil), errors.ErrorCommandSpecific)
			}()

			peanutRepo := engine.NewRepository(gitCfg)
			dir, err := ioutil.TempDir("", "peanut")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)
			log.Printf("Cloning to %s", dir)
			err = peanutRepo.Clone(dir)
			if err != nil {
				return fmt.Errorf("failed to clone repository: %w", err)
			}

			return engine.StartPeanutSync(config, cfg, peanutRepo, metrics.New("peanut", nil), resync, signals.SetupSignalHandler())
		},
	}

	clientConfig = cli.AddKubectlFlagsToCmd(&cmd)

	cmd.Flags().StringVar(&gitCfg.RepoURL, repoURLFlag, "", "Repository to deploy")
	logIfError(cmd.MarkFlagRequired(repoURLFlag))

	cmd.Flags().StringVar(&gitCfg.Branch, branchFlag, "", "Branch to checkout")
	logIfError(cmd.MarkFlagRequired(branchFlag))

	cmd.Flags().StringVar(&gitCfg.Path, pathFlag, "", "Path within the Repository to deploy")
	logIfError(cmd.MarkFlagRequired(pathFlag))

	cmd.Flags().DurationVar(&cfg.Resync, resyncFlag, time.Minute*5, "Resync frequency")
	cmd.Flags().IntVar(&port, portFlag, 8080, "Port number.")
	cmd.Flags().BoolVar(&cfg.Prune, pruneFlag, true, "Enables resource pruning.")
	logIfError(viper.BindPFlag(portFlag, cmd.Flags().Lookup(portFlag)))
	cmd.Flags().BoolVar(&cfg.Namespaced, namespacedFlag, false, "Switches agent into namespaced mode.")
	cmd.Flags().StringVar(&cfg.Namespace, defaultNamespaceFlag, "",
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
