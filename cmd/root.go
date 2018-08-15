// Copyright © 2018 Patrick Nuckolls <nuckollsp at gmail>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourfin/binappender"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "binappend",
	Short: "Cross-platform distributed ffmpeg-based transcoding pipeline",
	Long: `Transcodebot is designed to simplify distributing ffmpeg transcoding to the background of computers with other jobs, e.g. various home computers.
This is the server CLI, which can be used to generate statically complied clients that work with extremely minimal setup, as well as serve and recieve files to transcode from clients.`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		forceSuperuserInit()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	settingsHelpString := fmt.Sprintf("The directory containing settings and state information.\n(Default: %s)", common.GetDefaultSettingsDir())
	rootCmd.PersistentFlags().StringVar(&settingsDirProxy, "settings-dir", "", settingsHelpString)
	rootCmd.PersistentFlags().BoolVar(&forceSuperuser, "force-su", false, "Force transcodebot to use superuser defaults")
	rootCmd.PersistentFlags().BoolVar(&forceNoSuperuser, "force-no-su", false, "Force transcodebot to use normal user defaults")
	rootCmd.PersistentFlags().BoolVar(&common.AlwaysPanic, "always-panic", false, "Always panic instead of normal error messages")
	rootCmd.PersistentFlags().MarkHidden("always-panic")
}

func forceSuperuserInit() {
	if forceSuperuser && forceNoSuperuser {
		fmt.Println(`Cannot force superuser and force no superuser
(both --force-su and --force-no-su flags present)`)
		os.Exit(1)
	} else if forceSuperuser {
		common.ForceSuperuser(true)
	} else if forceNoSuperuser {
		common.ForceSuperuser(false)
	}
}

// Reads in config file and sets settings dir
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match settings
	if settingsDirProxy == "" {
		settingsDirProxy = common.GetDefaultSettingsDir()
	}
	common.SetSettingsDir(settingsDirProxy)
	viper.AddConfigPath(common.SettingsDir())
	viper.SetConfigName("server-config.yaml")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		common.PrintVerbose("Using config file:", viper.ConfigFileUsed())
	}
}
