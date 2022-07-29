package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		// Command line arguments
		port := viper.GetInt("port")
		if p, err := cmd.Flags().GetInt("port"); err != nil {
			logrus.Fatalf("Could not get port from command-line arguments: %v\n", err)
		} else {
			port = p
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
	// Extensions
	// Just comparing with "new" is not good, but it is the only
	// subcommand planned to be in hnotes.
	if len(os.Args) >= 2 && !strings.HasPrefix(os.Args[1], "-") && os.Args[1] != "new" {
		ext := filepath.Join(config.DirPath, "extensions", os.Args[1])
		cmd := exec.Command(ext)
		if len(os.Args) >= 3 {
			cmd = exec.Command(ext, os.Args[2:]...)
		}

		// Setup output and input
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr

		// Run command
		if err := cmd.Run(); err != nil {
			logrus.Fatalln(err.Error())
		}
		os.Exit(0)
	}

	// Flags
	cobra.OnInitialize(config.MustRead)
	rootCmd.PersistentFlags().IntP("port", "p", 4545, "port on which the server will run")

	// Sub commands
	rootCmd.AddCommand(newCmd)
}
