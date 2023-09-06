package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var (
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgGreen)
)

func makeBoilerplate(ctx *cli.Context) (err error) {
	targetFolder := normalizeDir(ctx.String(folderFlag))
	projectName := ctx.String(projectFlag)
	moduleName := ctx.String(moduleFlag)
	appName := ctx.String(appFlag)
	noApiSpec := ctx.Bool(apiFlag)

	if len(projectName) <= 0 {
		for {
			fmt.Println("\nWhat's the name of the GH project? [Please enter a name in the form <gh_profile>/<repo>]")
			fmt.Scanln(&projectName)
			projectName = normalizeDir(projectName)

			if err := validateGhProject(projectName); err != nil {
				yellow.Printf("%s\n", err)
				continue
			}
			break
		}
	}

	if len(targetFolder) <= 0 {
		defaultFolder, _ := filepath.Abs(".")
		defaultFolder = filepath.Join(defaultFolder, projectName)
		fmt.Printf("\nWhat's the folder where I'm going to create the local repository? [Press ENTER to use %s]\n", defaultFolder)
		fmt.Scanln(&targetFolder)
		targetFolder = cleanAndExpandPath(strings.Trim(targetFolder, " \n"))
		if err != nil {
			return err
		}
		if len(targetFolder) <= 0 {
			targetFolder = defaultFolder
		}
	}
	targetFolder, _ = filepath.Abs(targetFolder)

	if err := makeDirectoryIfNotExists(targetFolder); err != nil {
		return fmt.Errorf("failed to create folder %s: %s", targetFolder, err)
	}

	defer func() {
		if err != nil {
			var delete bool
			red.Printf("\nSomething went wrong: %s\n", err)
			for {
				var reply string
				red.Println("\nDo you want to remove all local files and folders created so far? [yes/no]")
				fmt.Scanln(&reply)
				reply = normalizeYesNo(reply)

				if reply != "yes" && reply != "y" && reply != "no" && reply != "n" {
					red.Println("invalid answer")
					continue
				}

				delete = reply == "yes" || reply == "y"
				break
			}

			if delete {
				os.RemoveAll(targetFolder)
			}
		}
	}()

	if len(moduleName) <= 0 {
		defaultModule := fmt.Sprintf("github.com/%s", projectName)
		fmt.Printf("\nWhat's the name of the Go module? [Press ENTER to use github.com/%s]\n", projectName)
		fmt.Scanln(&moduleName)

		if len(moduleName) <= 0 {
			moduleName = defaultModule
		}
	}

	if len(appName) <= 0 {
		defaultApp := strings.Split(projectName, "/")[1]
		fmt.Printf("\nWhat's the name of the app? [Press ENTER to use %s]\n", defaultApp)
		fmt.Scanln(&appName)

		if len(appName) <= 0 {
			appName = defaultApp
		}
	}

	if !noApiSpec {
		var reply string
		for {
			fmt.Println("\nDo you need an api-spec folder and relative buf initialization? [yes/no]")
			fmt.Scanln(&reply)
			reply = normalizeYesNo(reply)

			if reply != "yes" && reply != "y" && reply != "no" && reply != "n" {
				yellow.Println("invalid answer")
				continue
			}

			noApiSpec = reply == "no" || reply == "n"
			break
		}
	}

	fmt.Println("\nI'm creating the local repository...")
	if err := makeScaffolding(targetFolder, projectName, moduleName, appName, noApiSpec); err != nil {
		return fmt.Errorf("i couldn't create the local repository for the following reason: %s", err)
	}
	fmt.Println("Done!")

	fmt.Printf("\nBefore moving forward, please make sure to create the remote repo https://gihtub.com/%s. [Press ENTER to continue]\n", projectName)
	fmt.Scanln()

	fmt.Println("\nAlright! I'm initializing the project...")
	if err := initProject(targetFolder, projectName, moduleName, noApiSpec); err != nil {
		return fmt.Errorf("i couldn't initialize the project for the following reason: %s", err)
	}
	fmt.Println("\nDone!")

	green.Println("\nYour project has been kick-started")
	green.Printf("\nLocal repo: %s\n", targetFolder)
	green.Printf("Remote repo: https://github.com/%s\n", projectName)

	return nil
}

func normalizeDir(dir string) string {
	return cleanAndExpandPath(strings.Trim(dir, " "))
}

func normalizeYesNo(reply string) string {
	return strings.ToLower(strings.Trim(reply, " "))
}

func validateGhProject(name string) error {
	if len(name) <= 0 {
		return fmt.Errorf("missing project name")
	}
	s := strings.Split(name, "/")
	if len(s) != 2 {
		return fmt.Errorf("invalid project name")
	}
	if len(s[0]) <= 0 {
		return fmt.Errorf("missing gh profile name")
	}
	if len(s[1]) <= 0 {
		return fmt.Errorf("missing gh project name")
	}
	return nil
}
