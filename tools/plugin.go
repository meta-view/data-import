package tools

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"meta-view-service/services"
	"net/http"
	"os"
	"path"

	"github.com/coreos/go-semver/semver"
	"github.com/robertkrimen/otto"
)

// Plugin - a basic plugin structur
type Plugin struct {
	Provider *Provider
	Path     string
	Version  *semver.Version
	VM       *otto.Otto
	DB       *services.Database
}

// Object - a Javascript Object
type Object struct {
}

// LoadPlugin - loads one plugin
func LoadPlugin(pluginPath string, db *services.Database) (*Plugin, error) {
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
		DB:       db,
	}, nil
}

// Detect - returns the percentage if a given payload is of the type of the plugin
func (plugin *Plugin) Detect(payloadPath string) (float64, error) {
	log.Printf("Detect if %s is for [%s]", payloadPath, plugin.Provider.Name)
	output := 0.0

	err := plugin.loadTools(payloadPath)
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
	log.Printf("result: %f", value)
	return value, nil
}

// Import - imports the payload into a specific data structure
func (plugin *Plugin) Import(payloadPath string) error {
	log.Printf("Import data of %s for [%s]", payloadPath, plugin.Provider.Name)

	err := plugin.loadTools(payloadPath)
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

func (plugin *Plugin) loadTools(payloadPath string) error {
	LoadPluginExtenstions(plugin.VM)
	log.Printf("installing tools for path %s", payloadPath)
	plugin.VM.Set("_provider", plugin.Provider.Name)

	err := plugin.VM.Set("getFiles", func() []string {
		return readFiles(payloadPath, "")
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("getContent", func(file string) string {
		path := path.Join(payloadPath, file)
		content, err := getFileContent(path)
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
		contentType, err := getFileContentType(path)
		if err != nil {
			log.Printf("error %s reading contentType of %s\n", err, path)
			return ""
		}
		return contentType
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("getSha1Checksum", func(file string) string {
		path := path.Join(payloadPath, file)
		checksum, err := getSha1Checksum(path)
		if err != nil {
			log.Printf("error %s calculating sha1 checksum of %s\n", err, path)
			return ""
		}
		return checksum
	})
	if err != nil {
		return err
	}

	err = plugin.VM.Set("saveEntry", func(data map[string]interface{}) string {
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

	return nil
}

func readFiles(parent string, child string) []string {
	folder := path.Join(parent, child)
	files := make([]string, 0)
	log.Printf("reading folder %s\n", folder)
	folders, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Printf("error reading %s\n", folder)
		return files
	}
	for _, file := range folders {
		if file.IsDir() {
			subFolderFiles := readFiles(parent, path.Join(child, file.Name()))
			files = append(files, subFolderFiles...)
		} else {
			files = append(files, path.Join(child, file.Name()))
		}
	}
	return files
}

func getFileContent(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(content))
	return encoded, nil
}

func getFileContentType(file string) (string, error) {

	f, err := os.Open(file)
	if err != nil {
		log.Printf("error opening %s\n", file)
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 512)

	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func getSha1Checksum(file string) (string, error) {

	f, err := os.Open(file)
	if err != nil {
		log.Printf("error opening %s\n", file)
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
