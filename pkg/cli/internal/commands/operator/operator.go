package operator

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/zapr"
	"github.com/solo-io/bumblebee/pkg/cli/internal/options"
	"github.com/solo-io/bumblebee/pkg/operator"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type operatorOptions struct {
	general *options.GeneralOptions
}

func Command(opts *options.GeneralOptions) *cobra.Command {
	// operatorOptions := &operatorOptions{
	// 	general: opts,
	// }
	cmd := &cobra.Command{
		Use: "operator",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := buildContext(cmd.Context(), false)
			if err != nil {
				return err
			}
			return operator.Start(ctx)
		},
		SilenceUsage: true,
	}
	return cmd
}

func buildContext(ctx context.Context, debug bool) (context.Context, error) {
	ctx, cancel := context.WithCancel(ctx)
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopper
		fmt.Println("got sigterm or interrupt")
		cancel()
	}()
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't create zap logger: '%w'", err)
	}

	// controller-runtime
	zapLogger := zapr.NewLogger(logger)
	log.SetLogger(zapLogger)
	klog.SetLogger(zapLogger)

	contextutils.SetFallbackLogger(logger.Sugar())

	return ctx, nil
}
