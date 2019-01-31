// Copyright © 2019 suquiya
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

//Config storage config
type Config struct {
	License map[string]string `json:"license"`
	Author  map[string]string `json:"author"`
}

//Record write config c as json to a file specified by p
func Record(c *Config, p string) {
	j, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(p)
	if err != nil {
		fmt.Println("Cannot create config file. Specified license of author will not record.")
	}
	defer f.Close()
	f.Write(j)

}

//ReadConfigFile read config from file locating p
func ReadConfigFile(p string) *Config {
	f, err := os.Open(p)
	defer f.Close()
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	c := NewConfig()
	err = json.Unmarshal(b, c)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return c
}

//NewConfig crate new instance of config.
func NewConfig() *Config {
	return &Config{make(map[string]string), make(map[string]string)}
}

//SetDefValue set default vaule
func (c *Config) SetDefValue() {
	c.License["last"] = "mit"
	c.License["fix"] = ""
	c.Author["last"] = "COPYRIGHT HOLDER"
	c.Author["fix"] = ""
}

//GetLicenseValue get license value
func (c *Config) GetLicenseValue() string {
	if c.License["fix"] != "" {
		return c.License["fix"]
	}
	return c.License["last"]
}

func getDefaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return filepath.Join(home, ".suquiya_liquid.json")
}

// rootCmd represents the base command when called without any subcommands

func newRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "liquid",
		Short: "liquid is utility for license management in golang.",
		Long: `liquid is utility for license management in golang. liquid can add LICENSE to top of source code, and replace its LICENSE to another LICENSE
		`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
		RunE: Process,
	}

	rootCmd.AddCommand(newAddCmd())

	rootCmd.Flags().StringP("license", "l", "mit", "name of license (first default is mit, and after first use, config record what user choose and set it default)")
	rootCmd.Flags().StringP("author", "a", "COPYRIGHT HOLDER", "author(copyright holder) name for copyright (default is COPYTIGHT HOLDER)")
	rootCmd.Flags().StringP("customLicense", "c", "", "custom license path in the case that use custom license")
	rootCmd.Flags().String("config", "", "config file")

	return rootCmd
}

//Process is function of inner commandline-tools
func Process(cmd *cobra.Command, args []string) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		panic(err)
	}

	if exist, _ := IsExistFilePath(configPath); !exist {
		configPath = getDefaultConfigPath()
	}

	config := ReadConfigFile(configPath)
	if config == nil {
		config = NewConfig()
		config.SetDefValue()
	}

	l, err := cmd.Flags().GetString("license")
	if err != nil {
		panic(err)
	}
	license := ""
	if l == "" {

	}

	return nil
}

//IsExistFilePath is validate whether val is exist filepath or not.
func IsExistFilePath(val string) (bool, error) {
	absPath, err := filepath.Abs(val)
	if err != nil {
		return false, err
	}
	if is, _ := govalidator.IsFilePath(absPath); !is {
		return is, fmt.Errorf("IsNotFilePath")
	}

	fi, err := os.Stat(val)
	if err != nil {
		return false, err
	}

	if fi.IsDir() {
		return false, fmt.Errorf("%s is not file", val)
	}

	return true, err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.liquid.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".liquid" (without extension).
		viper.AddConfigPath(home)
		fmt.Println(home)
		viper.SetConfigName(".liquid")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
	}
}
*/
