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

func LoadPlugin(path string) (*Plugin, error) {
	return &Plugin{
		Provider: &Provider{Name: "Test", DownloadRequest: "/provider/test"},
		Path:     path,
		Version:  semver.New("0.0.1"),
	}, nil
}

func (plugin *Plugin) GetAccountName() (string, error) {
	script := `
		// Sample xyzzy example
		(function(){
			var length = 12;
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
