package plugins

import (
	"fmt"

	"github.com/algrvvv/ali/v2/utils"
)

func Show() {
	plugs, err := GetPlugins()
	if err != nil {
		utils.PrintError("failed to get plugins list", err)
		return
	}

	if len(plugs) == 0 {
		fmt.Println("plugins not found")
		return
	}

	fmt.Println("plugins list:")
	for _, plug := range plugs {
		fmt.Println(plug)
	}
}
