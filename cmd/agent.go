package cmd

import (
	"context"
	"fmt"
	"github.com/uanid/fakenews-server/application/configs"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/uanid/fakenews-server/application/fnc_agent"
)

func init() {
	var agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Fake News Challenge Agent를 실행합니다",
	}

	var interval time.Duration
	var runOnce bool
	agentCmd.Flags().DurationVar(&interval, "interval", time.Duration(10)*time.Second, "The Agent Work Iteration Loop Interval")
	agentCmd.Flags().BoolVar(&runOnce, "once", false, "Flag for not Loop Iteration")

	agentCmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := configs.LoadConfig(configPath)
		if err != nil {
			return err
		}

		app, err := fnc_agent.NewApplication(cfg)
		if err != nil {
			return fmt.Errorf("AgentInitFailed: %s", err.Error())
		}

		ctx, cancel := context.WithCancel(context.Background())
		go handleSigterm(cancel)

		if runOnce {
			err = app.Start(ctx)
		} else {
			fmt.Printf("Start Agent interval=%.2fs\n", interval.Seconds())
			err = app.StartWithTicker(ctx, interval)
		}
		if err != nil {
			if !strings.Contains(err.Error(), context.Canceled.Error()) {
				return err
			}
			fmt.Printf("Error: %s\n", err.Error())
		}

		fmt.Println("Stop Agent")
		return nil
	}

	rootCmd.AddCommand(agentCmd)
}

func handleSigterm(cancel func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	sig := <-signals
	fmt.Printf("Received Signal '%s'. Terminating...\n", sig.String())
	cancel()
}
