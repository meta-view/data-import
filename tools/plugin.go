package tools

import (
	"github.com/coreos/go-semver/semver"
	"github.com/robertkrimen/otto"
)

type Plugin struct {
	Provider *Provider
	Path     string
	Version  *semver.Version
}

func (plugin *Plugin) getAccountName() (string, error) {
	script := `
		// Sample xyzzy example
		(function(){
			var result           = '';
			var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
			var charactersLength = characters.length;
			for ( var i = 0; i < length; i++ ) {
			   result += characters.charAt(Math.floor(Math.random() * charactersLength));
			}
			return result;
		})();
	`
	vm := otto.New()
	result, err := vm.Run(script)
	if err != nil {
		return "", err
	}
	value, err := result.ToString()
	if err != nil {
		return "", err
	}
	return value, nil
}
