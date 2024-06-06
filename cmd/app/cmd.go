package app

import (
	"context"
	"practice/infra/log"
	"practice/library/config"
	"practice/server"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	Env         = config.Dev
	serviceName = "library"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "library",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(false)
		},
	}

	cmd.AddCommand(newCmdApiserver())
	return cmd
}

func newCmdApiserver() *cobra.Command {
	return &cobra.Command{
		Use: "apiserver",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(true)
		},
	}
}

func runServer(isApiServer bool) error {
	log, err := log.New(serviceName)
	if err != nil {
		return err
	}

	undo := zap.ReplaceGlobals(log)
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
		undo()
		log.Sync()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		stopCh := server.SetupSignalHandler()
		<-stopCh
		cancel()
	}()

	s, err := NewServer(isApiServer)
	if err != nil {
		return err
	}

	s.Run(ctx)
	return nil
}
