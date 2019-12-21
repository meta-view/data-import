package tools

import (
	"io/ioutil"
	"meta-view-service/services"
	"path"
)

// PluginCollection - a collection of all plugins on this system
type PluginCollection struct {
	Plugins map[string]*Plugin
}

// LoadPlugins - load all plugins of the current system
func LoadPlugins(pluginFolder string, db *services.Database) (map[string]*Plugin, error) {
	folders, err := ioutil.ReadDir(pluginFolder)
	if err != nil {
		return nil, err
	}
	plugins := make(map[string]*Plugin)
	for _, folder := range folders {
		if folder.IsDir() {
			plugin, err := LoadPlugin(path.Join(pluginFolder, folder.Name()), db)
			if err == nil {
				plugins[folder.Name()] = plugin
			}
		}
	}
	return plugins, nil
}
