package tools

import (
	"log"
	"meta-view-service/services"
	"strings"

	"github.com/robertkrimen/otto"
)

// LoadPluginExtenstions - adds JS extensions to otto vm runtime
func (plugin *Plugin) LoadPluginExtenstions() error {
	err := plugin.VM.Set("getChecksum", func(content string) string {
		return services.GetSha1Checksum(content)
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("log", func(message string) {
		log.Printf("[%s] %s", plugin.Name, message)
	})
	if err != nil {
		return err
	}

	err = loadStringExtensions(plugin.VM)
	if err != nil {
		return err
	}
	return nil
}

func loadStringExtensions(vm *otto.Otto) error {
	err := vm.Set("StringEndsWith", func(a string, b string) bool {
		return strings.HasSuffix(a, b)
	})
	if err != nil {
		return err
	}

	err = vm.Set("StringStartsWith", func(a string, b string) bool {
		return strings.HasPrefix(a, b)
	})
	if err != nil {
		return err
	}

	err = vm.Set("StringReplace", func(s string, old string, new string) string {
		return strings.TrimSpace(strings.Replace(s, old, new, -1))
	})
	if err != nil {
		return err
	}

	return nil
}
