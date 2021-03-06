package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var cfgFile string

var (
	// QriRepoPath is the path to the QRI repository
	QriRepoPath string
	// IpfsFsPath is the path to the IPFS repo
	IpfsFsPath string
)

// global pagination variables
var (
	pageNum  int
	pageSize int
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "qri",
	Short: "qri GDVCS CLI",
	Long: `
qri (pronounced "query") is a global dataset version control system 
on the distributed web.

https://qri.io

Feedback, questions, bug reports, and contributions are welcome!
https://github.com/qri-io/qri/issues`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		printErr(err)
		os.Exit(-1)
	}
}

func init() {
	flag.Parse()
	cobra.OnInitialize(initializeCLI)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $QRI_PATH/config.json)")
	RootCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "c", false, "disable colorized output")
}

// initializeCLI sets up the CLI, reading in config file and ENV variables if set.
func initializeCLI() {
	home := userHomeDir()

	QriRepoPath = os.Getenv("QRI_PATH")
	if QriRepoPath == "" {
		QriRepoPath = filepath.Join(home, ".qri")
	}
	// TODO - this is stupid
	QriRepoPath = strings.Replace(QriRepoPath, "~", home, 1)

	IpfsFsPath = os.Getenv("IPFS_PATH")
	if IpfsFsPath == "" {
		IpfsFsPath = filepath.Join(home, ".ipfs")
	}
	IpfsFsPath = strings.Replace(IpfsFsPath, "~", home, 1)

	setNoColor()
	loadConfig()
	return
}
