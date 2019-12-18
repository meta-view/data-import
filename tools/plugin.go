package tools

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"strconv"

	"github.com/coreos/go-semver/semver"
	"github.com/robertkrimen/otto"
)

// Plugin - a basic plugin structur
type Plugin struct {
	Provider *Provider
	Path     string
	Version  *semver.Version
	VM       *otto.Otto
}

// LoadPlugin - loads one plugin
func LoadPlugin(pluginPath string) (*Plugin, error) {
	data, err := ioutil.ReadFile(path.Join(pluginPath, "package.json"))
	if err != nil {
		return nil, err
	}
	var packageInfo map[string]string
	err = json.Unmarshal(data, &packageInfo)
	log.Printf("loading %s from %s\n", packageInfo, pluginPath)
	return &Plugin{
		Provider: &Provider{Name: packageInfo["name"], DownloadRequest: packageInfo["download_request"]},
		Path:     pluginPath,
		Version:  semver.New(packageInfo["version"]),
		VM:       otto.New(),
	}, nil
}

// Detect - returns the percentage if a given payload is of the type of the plugin
func (plugin *Plugin) Detect(payloadPath string) (float64, error) {
	log.Printf("Detect if %s is for [%s]", payloadPath, plugin.Provider.Name)
	output := 0.0
	plugin.VM.Set("payloadPath", payloadPath)
	script, err := ioutil.ReadFile(path.Join(plugin.Path, "detector.js"))
	if err != nil {
		return output, err
	}

	result, err := plugin.VM.Run(script)
	value, err := result.ToString()
	if err != nil {
		return output, err
	}

	output, err = strconv.ParseFloat(value, 64)
	if err != nil {
		return output, err
	}

	return output, nil
}
