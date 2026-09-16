package envs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/algrvvv/ali/v2/aliases"
	"github.com/algrvvv/ali/v2/logger"
	"github.com/algrvvv/ali/v2/utils"
	"github.com/spf13/viper"
)

const color = "\033[1;34m"

func Print(search string, aliEnvFileFlag string, asTable bool) {
	logger.SaveDebugf("print envs")
	fmt.Println("Available Envs:")

	env := Load(nil, aliEnvFileFlag)

	als := aliases.Load(viper.GetViper())
	for alias, entry := range als {
		for name, value := range entry.Env {
			key := fmt.Sprintf("%s (%s)", name, alias)
			env[key] = value
		}
	}

	if asTable {
		tablePrint(env, search)
		return
	}

	fullPrint(env, search)
}

func fullPrint(envs map[string]any, search string) {
	var count int
	for name, value := range envs {
		if !searchInEnvs(search, name) {
			continue
		}

		prefix := "  └── "
		if count != len(envs)-1 {
			prefix = "  ├── "
		}
		count++

		name = "$" + strings.ToUpper(name)
		name = utils.Colors["red"] + name + utils.Colors["reset"]

		re := regexp.MustCompile(`\(([^)]+)\)`)
		name = re.ReplaceAllStringFunc(name, func(alias string) string {
			return color + strings.ToLower(alias) + utils.Colors["reset"]
		})

		fmt.Printf("%s%s -> %v\n", prefix, name, value)
	}
}

func tablePrint(envs map[string]any, search string) {
	fmt.Printf("+%s+%s+\n", strings.Repeat("-", 30), strings.Repeat("-", 42))
	fmt.Printf("| Env%s| Command%s|\n",
		strings.Repeat(" ", 32-len(" alias")),
		strings.Repeat(" ", 42-len(" command")),
	)
	fmt.Printf("+%s+%s+\n", strings.Repeat("-", 30), strings.Repeat("-", 42))

	for name, value := range envs {
		if !searchInEnvs(search, name) {
			continue
		}

		name = "$" + strings.ToUpper(name)
		re := regexp.MustCompile(`\(([^)]+)\)`)
		name = re.ReplaceAllStringFunc(name, func(alias string) string {
			return strings.ToLower(alias)
		})

		fmt.Printf("| %s%-28s%s | %-40s |\n",
			utils.Colors["red"],
			utils.TruncateString(name, 28),
			utils.Colors["reset"],
			utils.TruncateString(fmt.Sprintf("%v", value), 40),
		)

		fmt.Printf("+%s+%s+\n", strings.Repeat("-", 30), strings.Repeat("-", 42))
	}
}

func searchInEnvs(search, envName string) bool {
	if search == "" {
		return true
	}

	if strings.Contains(envName, search) {
		return true
	}

	return false
}
