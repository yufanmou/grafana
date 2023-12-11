package apiserver

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli"

	grafanaAPIServer "github.com/grafana/grafana/pkg/services/grafana-apiserver"
	"github.com/grafana/grafana/pkg/setting"
)

func newCommandStartExampleAPIServer(o *APIServerOptions, stopCh <-chan struct{}) *cobra.Command {
	devAcknowledgementNotice := "The apiserver command is in heavy development.  The entire setup is subject to change without notice"

	cmd := &cobra.Command{
		Use:   "apiserver [api group(s)]",
		Short: "Run the grafana apiserver",
		Long: "Run a standalone kubernetes based apiserver that can be aggregated by a root apiserver. " +
			devAcknowledgementNotice,
		Example: "grafana apiserver example.grafana.app",
		RunE: func(c *cobra.Command, args []string) error {
			cfg, _ := setting.NewCfgFromArgs(setting.CommandLineArgs{
				Config:   "conf/custom.ini",
				HomePath: "./",
			})

			// Parse builders for each group in the args
			builders, err := ParseAPIGroupArgs(cfg, args[1:])
			if err != nil {
				return err
			}

			if err := o.LoadAPIGroupBuilders(builders); err != nil {
				return err
			}

			// Finish the config (a noop for now)
			if err := o.Complete(); err != nil {
				return err
			}

			config, err := o.Config()
			if err != nil {
				return err
			}

			if err := o.RunAPIServer(config, stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	// Register standard k8s flags with the command line
	o.RecommendedOptions = options.NewRecommendedOptions(
		defaultEtcdPathPrefix,
		Codecs.LegacyCodec(), // the codec is passed to etcd and not used
	)
	o.RecommendedOptions.AddFlags(cmd.Flags())

	return cmd
}

func ParseAPIGroupArgs(cfg *setting.Cfg, args []string) ([]grafanaAPIServer.APIGroupBuilder, error) {
	builders := make([]grafanaAPIServer.APIGroupBuilder, 0)
	for _, g := range args {
		switch g {
		// case "example.grafana.app":
		// 	eb, err := initializeExampleAPIBuilder(cfg)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	builders = append(builders, eb)
		// case "playlist.grafana.app":
		// 	pb, err := initializePlaylistsAPIBuilder(cfg)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	builders = append(builders, pb)
		// case "snapshots.grafana.app":
		// 	sb, err := initializeSnapshotsAPIBuilder(cfg)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	builders = append(builders, sb)
		default:
			return nil, fmt.Errorf("unknown group: %s", g)
		}
	}

	if len(builders) < 1 {
		return nil, fmt.Errorf("expected group name(s) in the command line arguments")
	}

	return builders, nil
}

func RunCLI() int {
	stopCh := genericapiserver.SetupSignalHandler()

	options := newAPIServerOptions(os.Stdout, os.Stderr)
	cmd := newCommandStartExampleAPIServer(options, stopCh)

	return cli.Run(cmd)
}
