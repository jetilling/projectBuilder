package appBuilder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jetilling/projectBuilder/configVars"
)

func createProjectDirectory(data AppDetails) {
	cmd := exec.Command("/bin/sh", "./buildScripts/initialBuildScripts/createProjectDirectory.sh", data.Name, data.UniqueID)
	runBashScript(cmd)
}

func downloadLaravel(projectFolderString string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/initialBuildScripts/downloadLaravel.sh", projectFolderString)
	runBashScript(cmd)
}

func renameLaravelApp(projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/initialBuildScripts/renameLaravelApp.sh", projectFolderString, appName)
	runBashScript(cmd)
}

func copyEnvironmentFile(projectFilePath, projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/initialBuildScripts/copyEnvironmentFile.sh", projectFolderString, appName)
	runBashScript(cmd)

	inputEnvironmentFile := fmt.Sprintf("%s/.env.template", projectFilePath)
	outputEnvironmentFile := fmt.Sprintf("%s/.env", projectFilePath)
	findAndReplace(inputEnvironmentFile, outputEnvironmentFile, "{{project_name}}", strings.ToLower(appName))
}

func copyDockerFiles(projectFilePath, projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/initialBuildScripts/copyDockerFiles.sh", projectFolderString, appName)
	runBashScript(cmd)

	inputDockerFile := fmt.Sprintf("%s/docker-compose_template.yml", projectFilePath)
	outputDockerFile := fmt.Sprintf("%s/docker-compose.yml", projectFilePath)
	findAndReplace(inputDockerFile, outputDockerFile, "{{project_name}}", strings.ToLower(appName))

	inputTestDBFile := fmt.Sprintf("%s/create-testing-db_template.sql", projectFilePath)
	outputTestDBFile := fmt.Sprintf("%s/create-testing-db.sql", projectFilePath)
	findAndReplace(inputTestDBFile, outputTestDBFile, "{{project_name}}", strings.ToLower(appName))
}

func removeTemplateFiles(projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/removeTemplateFiles.sh", projectFolderString, appName)
	runBashScript(cmd)
}

func addReact(projectFilePath, projectFolderString, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/reactScripts/addReact.sh", projectFolderString, appName)
	runBashScript(cmd)

	inputMainFile := fmt.Sprintf("%s/resources/js/components/main.js.template", projectFilePath)
	outputMainFile := fmt.Sprintf("%s/resources/js/components/main.js", projectFilePath)
	findAndReplace(inputMainFile, outputMainFile, "{{project_name}}", strings.ToLower(appName))

	inputIndexFile := fmt.Sprintf("%s/resources/js/index.js.template", projectFilePath)
	outputIndexFile := fmt.Sprintf("%s/resources/js/index.js", projectFilePath)
	findAndReplace(inputIndexFile, outputIndexFile, "{{project_name}}", strings.ToLower(appName))

	inputReducerFile := fmt.Sprintf("%s/resources/js/reducer.js.template", projectFilePath)
	outputReducerFile := fmt.Sprintf("%s/resources/js/reducer.js", projectFilePath)
	findAndReplace(inputReducerFile, outputReducerFile, "{{project_name}}", strings.ToLower(appName))

	cmd = exec.Command("/bin/sh", "./buildScripts/reactScripts/cleanUpReactTemplates.sh", projectFolderString, appName)
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
