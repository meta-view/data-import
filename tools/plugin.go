package tools

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"meta-view-service/services"
	"path"

	"github.com/coreos/go-semver/semver"
	"github.com/robertkrimen/otto"
)

// Plugin - a basic plugin structur
type Plugin struct {
	ID              string                `json:"ID"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	DownloadRequest string                `json:"download_request"`
	URLWhitelist    []string              `json:"url_whitelist"`
	Provider        string                `json:"provider"`
	Version         *semver.Version       `json:"version"`
	Path            string                `json:"path"`
	VM              *otto.Otto            `json:"-"`
	DB              *services.Database    `json:"-"`
	FS              *services.FileStorage `json:"-"`
}

// LoadPlugin - loads one plugin
func LoadPlugin(pluginPath string, db *services.Database, fs *services.FileStorage) (*Plugin, error) {
	infoFile := path.Join(pluginPath, "info.json")
	data, err := ioutil.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	checksum, err := services.GetSha1ChecksumOfFile(infoFile)
	if err != nil {
		return nil, err
	}
	var packageInfo map[string]interface{}
	err = json.Unmarshal(data, &packageInfo)
	log.Printf("loading %s from %s\n", packageInfo, pluginPath)
	return &Plugin{
		ID:              checksum,
		Name:            packageInfo["name"].(string),
		Description:     packageInfo["description"].(string),
		DownloadRequest: packageInfo["download_request"].(string),
		URLWhitelist:    parseList(packageInfo["url_whitelist"]),
		Path:            pluginPath,
		Version:         semver.New(packageInfo["version"].(string)),
		Provider:        packageInfo["provider"].(string),
		VM:              otto.New(),
		DB:              db,
		FS:              fs,
	}, nil
}

func parseList(list interface{}) []string {
	elements := make([]string, 0)
	vals, ok := list.([]interface{})
	if ok {
		for _, val := range vals {
			elements = append(elements, val.(string))
		}
	}
	return elements
}

// Detect - returns the percentage if a given payload is of the type of the plugin
func (plugin *Plugin) Detect(payloadPath string) (float64, error) {
	log.Printf("Detecting if %s is for [%s]\n", payloadPath, plugin.Name)
	output := 0.0

	err := LoadPluginExtenstions(plugin.VM)
	if err != nil {
		return output, err
	}

	err = plugin.loadFileTools(payloadPath)
	if err != nil {
		return output, err
	}

	defaultProfile, err := plugin.FS.SaveFile(path.Join("static", "images", "default_profile.png"))
	if err != nil {
		return output, err
	}

	err = plugin.VM.Set("_provider", plugin.Provider)
	if err != nil {
		return output, err
	}

	err = plugin.VM.Set("_version", plugin.Version)
	if err != nil {
		return output, err
	}

	err = plugin.VM.Set("_defaultProfile", defaultProfile)
	if err != nil {
		return output, err
	}

	err = plugin.VM.Set("_payloadPath", payloadPath)
	if err != nil {
		return output, err
	}

	script, err := ioutil.ReadFile(path.Join(plugin.Path, "detector.js"))
	if err != nil {
		return output, err
	}

	result, err := plugin.VM.Run(script)
	if err != nil {
		return output, err
	}

	value, err := result.ToFloat()
	if err != nil {
		return output, err
	}
	log.Printf("Detector result: %f", value)
	return math.Round(value*100) / 100, nil
}

// Import - imports the payload into a specific data structure
func (plugin *Plugin) Import(payloadPath string) error {
	log.Printf("Importing data of %s for [%s]\n", payloadPath, plugin.Name)

	err := LoadPluginExtenstions(plugin.VM)
	if err != nil {
		return err
	}

	err = plugin.loadFileTools(payloadPath)
	if err != nil {
		return err
	}

	err = plugin.loadDBTools()
	if err != nil {
		return err
	}

	err = plugin.VM.Set("_payloadPath", payloadPath)
	if err != nil {
		return err
	}

	script, err := ioutil.ReadFile(path.Join(plugin.Path, "importer.js"))
	if err != nil {
		return err
	}

	result, err := plugin.VM.Run(script)
	if err != nil {
		return err
	}

	log.Printf("Import Result: %s\n", result)
	return nil
}

// Present - queries and presents a given list of found elements
func (plugin *Plugin) Present(entry map[string]interface{}, render string) (string, error) {
	log.Printf("Presenting result id %s for [%s]", entry["id"], plugin.Name)
	output := ""

	script, err := ioutil.ReadFile(path.Join(plugin.Path, "presenter.js"))
	if err != nil {
		return output, err
	}

	err = LoadPluginExtenstions(plugin.VM)
	if err != nil {
		return output, err
	}
	plugin.VM.Set("render", render)
	plugin.VM.Set("entry", entry)
	result, err := plugin.VM.Run(script)
	if err != nil {
		return output, err
	}

	value, err := result.ToString()
	if err != nil {
		return output, err
	}

	return value, nil
}

func (plugin *Plugin) loadFileTools(payloadPath string) error {
	log.Printf("Installing file tools in path %s for plugin %s", payloadPath, plugin.Name)

	err := plugin.VM.Set("readDir", func() []string {
		return readFiles(payloadPath, "", false)
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("listFiles", func() []string {
		return readFiles(payloadPath, "", true)
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("saveFile", func(file string) string {
		path := path.Join(payloadPath, file)
		filePath, err := plugin.FS.SaveFile(path)
		if err != nil {
			log.Printf("error saving file %s\n", path)
			return ""
		}
		return filePath
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("getContent", func(file string) string {
		path := path.Join(payloadPath, file)
		content, err := services.GetFileContent(path)
		if err != nil {
			log.Printf("error %s reading content of %s\n", err, path)
			return ""
		}
		return content
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("getBase64", func(file string) string {
		path := path.Join(payloadPath, file)
		content, err := services.GetFileBase64(path)
		if err != nil {
			log.Printf("error %s reading content of %s\n", err, path)
			return ""
		}
		return content
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("getContentType", func(file string) string {
		path := path.Join(payloadPath, file)
		contentType, err := services.GetFileContentType(path)
		if err != nil {
			log.Printf("error %s reading contentType of %s\n", err, path)
			return ""
		}
		return contentType
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("getFileChecksum", func(file string) string {
		path := path.Join(payloadPath, file)
		checksum, err := services.GetSha1ChecksumOfFile(path)
		if err != nil {
			log.Printf("error %s calculating sha1 checksum of %s\n", err, path)
			return ""
		}
		return checksum
	})
	if err != nil {
		return err
	}

	return nil
}

func (plugin *Plugin) loadDBTools() error {
	log.Printf("Installing DB tools for plugin %s", plugin.Name)
	err := plugin.VM.Set("saveEntry", func(data map[string]interface{}) string {
		id, err := plugin.DB.SaveEntry(data)
		if err != nil {
			log.Printf("error %s saving data\n", err)
			return ""
		}
		return id
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("readEntry", func(query map[string]interface{}) map[string]interface{} {
		entries, err := plugin.DB.ReadEntries(query)
		if err != nil {
			log.Printf("error %s reading with query %s\n", err, query)
			empty := make(map[string]interface{}, 0)
			return empty
		}
		return entries
	})
	if err != nil {
		return err
	}

	return nil
}

func readFiles(parent string, child string, filesOnly bool) []string {
	folder := path.Join(parent, child)
	files := make([]string, 0)
	log.Printf("Reading folder %s\n", folder)
	folders, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Printf("error reading %s\n", folder)
		return files
	}
	for _, file := range folders {
		if file.IsDir() {
			if !filesOnly {
				files = append(files, path.Join(child, file.Name()))
			}
			subFolderFiles := readFiles(parent, path.Join(child, file.Name()), filesOnly)
			files = append(files, subFolderFiles...)
		} else {
			files = append(files, path.Join(child, file.Name()))
		}
	}
	return files
}
