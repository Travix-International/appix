package appix

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func PrepareAppUpload(configAppPath string, skipTest bool) (appPath string, appName string, manifestPath string, err error) {
	if configAppPath == "" {
		configAppPath = "."
	}

	appPath, err = filepath.Abs(configAppPath)

	if err != nil {
		log.Printf("Invalid App path: %s\n", appPath)
		return "", "", "", err
	}

	manifestPath = appPath + "/app.manifest"

	if _, err = os.Stat(manifestPath); os.IsNotExist(err) {
		log.Printf("App manifest not found: %s\n", manifestPath)
		return "", "", "", err
	}

	type AppManifestSettings struct {
		Name  string         `json:"name"`
		Build []stages.Stage `json:"build"`
		Tests []stages.Stage `json:"tests"`
	}

	var manifestObject AppManifestSettings

	manifestData, err := ioutil.ReadFile(manifestPath)

	if err != nil {
		log.Println("Couldn't read the app.manifest")
		return "", "", "", err
	}

	err = json.Unmarshal(manifestData, &manifestObject)

	if err != nil {
		log.Println("Couldn't parse the app.manifest")
		return "", "", "", err
	}

	if manifestObject.Name == "" {
		log.Println("The name is missing from the app manifest")
		return "", "", "", errors.New("The name is missing from the app manifest")
	}

	if len(manifestObject.Tests) > 0 && skipTest == false {
		// create a pool of stages
		// wait for the return value of chan bool (a boolean)
		rTests := <-stages.CreateStagePool(manifestObject.Tests)
		if rTests {
			return "", "", "", errors.New("The tests are failing. Please check the output")
		}
	}

	if len(manifestObject.Build) > 0 {
		// create a pool of stages
		// wait for the return value of chan bool (a boolean)
		bTests := <-stages.CreateStagePool(manifestObject.Build)
		if bTests {
			return "", "", "", errors.New("The build is failing. Please check the ouput")
		}
	}

	appName = manifestObject.Name

	return appPath, appName, manifestPath, nil
}

func createZapPackage(appPath string, verbose bool) (string, error) {
	tempFolder, err := ioutil.TempDir("", "appix")

	if err != nil {
		log.Println("Could not create temp folder!")
		return "", err
	}

	zapFile := tempFolder + "/app.zap"

	if verbose {
		log.Println("Creating ZAP file: " + zapFile)
	}

	err = zipFolder(appPath, zapFile, func(path string) bool {
		ignored, ignoredFolder := IgnoreFilePath(path)
		if verbose && !ignoredFolder {
			if ignored {
				log.Printf("\tSkipping %s\n", path)
			} else {
				log.Printf("\tAdding %s\n", path)
			}
		}
		return !ignored
	})

	if err != nil {
		log.Println("Could not process App folder!")
		return "", err
	}

	return zapFile, err
}
