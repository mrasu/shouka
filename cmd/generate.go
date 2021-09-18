package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/mrasu/shouka/configs"

	"github.com/mrasu/shouka/injections"

	"github.com/mrasu/shouka/cmd/ask"
	"github.com/mrasu/shouka/generators"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate files to setup",
	Long: `Generate files like tf or yml by your input.
With them, you can setup your own environment including AWS or GitHub Actions.
In the environment, You can evaluate a way for CI/CD.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generate()
	},
}

var jsonFilePath string
var outputPath string

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&jsonFilePath, "json", "j", "", "Path to json instead of answering questions interactively")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to output your answers as json")
}

func generate() error {
	printAsciiBanner()

	cnf, err := getConfig()
	if err != nil {
		return err
	}

	if outputPath != "" {
		if err := writeConfig(outputPath, cnf); err != nil {
			return err
		}
	}

	// if _, err := os.Stat(dir); err == nil {
	// 	panic("already exists")
	// }
	//
	// if err := os.MkdirAll(dir, 0755); err != nil {
	// 	panic(err)
	// }
	g := generators.NewGenerator(injections.EmbedFs, cnf)
	if err := g.Generate(); err != nil {
		return err
	}

	path, err := filepath.Abs(cnf.Directory)
	if err != nil {
		// Ignore error as getting absolute path is just for usability
		path = cnf.Directory
	}
	fmt.Printf("Successfully generated files at %s\n", path)
	return nil
}

const banner = `
         888                        888               
         888                        888               
         888                        888               
.d8888b  88888b.   .d88b.  888  888 888  888  8888b.  
88K      888 "88b d88""88b 888  888 888 .88P     "88b 
"Y8888b. 888  888 888  888 888  888 888888K  .d888888 
     X88 888  888 Y88..88P Y88b 888 888 "88b 888  888 
 88888P' 888  888  "Y88P"   "Y88888 888  888 "Y888888 
                                                      
                                                      `

func printAsciiBanner() {
	fmt.Println(banner)
}

func getConfig() (*configs.Config, error) {
	if jsonFilePath == "" {
		cnf, err := ask.AskConfig()
		if err != nil {
			return nil, err
		}
		return cnf, nil
	}

	if _, err := os.Stat(jsonFilePath); err != nil {
		return nil, errors.New(fmt.Sprintf("no file found: %s", jsonFilePath))
	}
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read: %s", jsonFilePath))
	}

	cnf := &configs.Config{}
	if err := json.Unmarshal(data, cnf); err != nil {
		return nil, err
	}

	return cnf, nil
}

func writeConfig(filename string, cnf *configs.Config) error {
	data, err := json.MarshalIndent(cnf, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal Config")
	}

	if err := ioutil.WriteFile(filename, data, 0664); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write Config to file: %s", filename))
	}

	return nil
}
