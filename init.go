package appix

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Travix-International/appix/config"
)

type InitCommand struct {
	*Command
	appPath string // Path to where the app will be placed after the scaffold has taken place
}

func (cmd *InitCommand) Register(app *kingpin.Application, config config.Config) {
	command := app.Command("init", "Scaffold a new application into the specified folder").
		Action(func(parseContext *kingpin.ParseContext) error {
			// grab the absolute path
			appPathRelative := cmd.appPath
			appPathAbsolute, err := filepath.Abs(appPathRelative)
			if err != nil {
				log.Printf("Failed to obtain absolute path for %s\n%v", appPathRelative, err)
				return err
			}

			// tell the user what we're planning on doing
			log.Print("Initializing new application")
			if cmd.Verbose {
				log.Printf("Specified appPath to be %s", appPathRelative)
				log.Printf("Absolute appPath is %s", appPathAbsolute)
			}

			// First we'll check to see if the directory is empty. It's purely for safety purposes, to ensure we don't overwrite
			// anyting special. The command line handling has already validated that the folder actually exists
			isEmptyPath, err := isEmptyPath(appPathAbsolute)
			if !isEmptyPath || err != nil {
				log.Printf("The specified appPath '%s' does not appear to be an empty directory\n%v", appPathRelative, err)
				return err
			}

			// Scaffold
			err = scaffoldNewApp(appPathAbsolute, cmd.Verbose)
			if err != nil {
				return err
			}

			log.Print("All done!")
			return nil
		}).
		Alias("i")

	command.Arg("appPath", "Path to an empty folder. (default: current folder)").
		Default(".").
		ExistingDirVar(&cmd.appPath)
}

func isEmptyPath(appPath string) (bool, error) {
	// See http://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty

	// Open the directory, which must not fail
	f, err := os.Open(appPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// See if there's anything in the directory at all
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}
