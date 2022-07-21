package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"hellerwach.com/go/hnotes/config"
	"hellerwach.com/go/hnotes/server"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hnotes",
	Short: "A simple HTTP file server that converts Markdown to HTML before serving",
	Long: `HNotes is an HTTP file server that converts Markdown to HTML before serving
it. It also provides directory viewing functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			logrus.Fatalf("Could not get port from command-line arguments: %v\n", err)
		}
		if viper.GetInt("port") != 0 {
			port = viper.GetInt("port")
		}
		server.Run(port)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(config.MustRead)
	rootCmd.PersistentFlags().IntP("port", "p", 4545, "port on which the server will run")

	rootCmd.AddCommand(newCmd)
}
