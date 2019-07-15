package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lichangjiang/iperf-operator/cmd"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var rootCmd = &cobra.Command{
	Use:    "iperf-operator",
	Hidden: true,
}

func main() {
	klog.InitFlags(nil)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	addCommands()
	if err := rootCmd.Execute(); err != nil {
		klog.Errorf("iperf-operator error: %s\n", err.Error())
	}
	select {
	case <-signalChan:
	}
}

func addCommands() {
	rootCmd.AddCommand(cmd.OperatorCmd)
	rootCmd.AddCommand(cmd.DeployCmd)
}
