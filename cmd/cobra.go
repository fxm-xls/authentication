package cmd

import (
	"bigrule/cmd/authentication"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"bigrule/common/global"
	"bigrule/pkg/format"
)

var rootCmd = &cobra.Command{
	Use:          "bigrule",
	Short:        "bigrule",
	SilenceUsage: true,
	Long:         `bigrule`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(format.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用 ` + format.Green(global.ProjectName+" "+global.Version) + ` 可以使用 ` + format.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s\n", usageStr)
}

func init() {
	rootCmd.AddCommand(authentication.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
