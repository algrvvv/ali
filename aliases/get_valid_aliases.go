package aliases

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	color      = "\033[1;34m"
	resetColor = "\033[0m"
)

func GetValidAliases(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	aliases, ok := viper.Get("aliases").(map[string]any)
	if !ok {
		fmt.Println("failed to get all aliases")
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	var res []string
	for alias, command := range aliases {
		res = append(res, fmt.Sprintf("%s%s%s\t%s", color, alias, resetColor, command))
	}

	return res, cobra.ShellCompDirectiveNoFileComp
}
