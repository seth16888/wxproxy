package cmd

import (
	"fmt"
	"os"

	"github.com/seth16888/wxproxy/internal/bootstrap"
	"github.com/seth16888/wxproxy/internal/di"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "conf", "c",
		"conf/conf.yaml", "--conf config file (default is conf/conf.yaml)")
}

var rootCmd = &cobra.Command{
	Use:   "wxtoken [command] [flags] [args]",
	Short: "A WX token server",
	Long:  `A WX token server`,
  PreRunE: func(cmd *cobra.Command, args []string) error {
    defer func() {
      if err:= recover(); err != nil {
        fmt.Println("error: ", err)
        os.Exit(1)
      }
    }()

    di.NewContainer(configFile)
    return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
    return bootstrap.StartApp()
  },
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
