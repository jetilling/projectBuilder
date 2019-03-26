package appBuilder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/jetilling/projectBuilder/configVars"
)

func createProjectDirectory(data AppDetails) {
	cmd := exec.Command("/bin/sh", "./buildScripts/createProjectDirectory.sh", data.Name, data.UniqueID)
	runBashScript(cmd)
}

func downloadLaravel(projectFolderString string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/downloadLaravel.sh", projectFolderString)
	runBashScript(cmd)
}

func renameLaravelApp(projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/renameLaravelApp.sh", projectFolderString, appName)
	runBashScript(cmd)
}

func copyEnvironmentFile(projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/copyEnvironmentFile.sh", projectFolderString, appName)
	runBashScript(cmd)
}

func copyDockerFiles(projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/copyDockerFiles.sh", projectFolderString, appName)
	runBashScript(cmd)
}

func removeTemplateFiles(projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/removeTemplateFiles.sh", projectFolderString, appName)
	runBashScript(cmd)
}

func pushToGithub(projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/pushToGithub.sh", projectFolderString, appName, configVars.Config.GITHUB_PASS)
	runBashScript(cmd)
}

func runBashScript(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	}
	fmt.Printf("%s \n", out.String())
}

func findAndReplace(inputFilePath, outputFilePath, itemToFind, replacementValue string) {
	input, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.ReplaceAll(input, []byte(itemToFind), []byte(replacementValue))

	if err = ioutil.WriteFile(outputFilePath, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readFile(path string) {
	fmt.Println("Reading in: ", path)
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(dat))
}
