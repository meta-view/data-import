package tools

import (
	"io/ioutil"
	"log"
	"meta-view-service/services"
	"path"
)

// PluginCollection - a collection of all plugins on this system
type PluginCollection struct {
	Plugins map[string]*Plugin
}

// LoadPlugins - load all plugins of the current system
func LoadPlugins(pluginFolder string, db *services.Database, fs *services.FileStorage) (map[string]*Plugin, error) {
	folders, err := ioutil.ReadDir(pluginFolder)
	if err != nil {
		return nil, err
	}
	plugins := make(map[string]*Plugin)
	for _, folder := range folders {
		if folder.IsDir() {
			plugin, err := LoadPlugin(path.Join(pluginFolder, folder.Name()), db, fs)
			if err == nil {
				plugins[folder.Name()] = plugin
			}
			id, err := registerPlugin(plugin, db)
			if err != nil {
				log.Printf("error: %s", err)
			}
			log.Printf("registering plugin ID: %s ", id)
		}
	}
	return plugins, nil
}

func registerPlugin(plugin *Plugin, db *services.Database) (string, error) {
	data := make(map[string]interface{})
	data["id"] = plugin.ID
	data["owner"] = "system"
	data["table"] = "providers"
	data["provider"] = plugin.Provider
	data["name"] = plugin.Provider
	data["content-type"] = "application/json"

	data["content"] = plugin
	result, err := db.SaveEntry(data)
	if err != nil {
		return "", err
	}
	return result, nil
}
