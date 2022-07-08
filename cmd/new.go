package cmd

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"hellerwach.com/go/hnotes/server"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new note.",
	Long:  `Creates a new note from a template.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Fatalln("No name given")
		}
		for _, arg := range args {

			if err := new(arg); err != nil {
				logrus.Fatalf("Could not create note with name %q: %v\n", arg, err)
			}
		}
	},
}

func new(path string) error {
	mdTemplate, err := os.ReadFile(filepath.Join(server.ConfigPath, "templates/single.md"))
	if err != nil {
		return err
	}

	return os.WriteFile(path, mdTemplate, 0664)
}
