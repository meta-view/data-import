package tools

import (
	"log"
	"strings"

	"github.com/robertkrimen/otto"
)

// LoadPluginExtenstions - adds JS extensions to otto vm runtime
func LoadPluginExtenstions(vm *otto.Otto) error {
	log.Println("installing JS Extensions")
	err := vm.Set("StringEndsWith", func(a string, b string) bool {
		return strings.HasSuffix(a, b)
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
