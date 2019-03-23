package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/jetilling/projectBuilder/configVars"
)

type AppDetails struct {
	Name     string   `json:"name"`
	UniqueID string   `json:"uniqueID"`
	Models   []Models `json:"models"`
}

type Models struct {
	ModelName   string        `json:"modelName"`
	ColumnNames []ColumnNames `json:"columnNames"`
	ForeignKeys []ForeignKeys `json:"foreignKeys"`
}

type ColumnNames struct {
	ColumnName     string `json:"columnName"`
	ColumnDataType string `json:"columnDataType"`
	Nullable       bool   `json:"nullable"`
	Unique         bool   `json:"unique"`
}

type ForeignKeys struct {
	ForeignKeyName  string `json:"foreignKeyName"`
	ReferenceTable  string `json:"referenceTable"`
	ReferenceColumn string `json:"referenceColumn"`
}

func main() {

	// INITIALIZE CONFIG VARIABLES
	configVars.InitConfigVars()

	file, _ := ioutil.ReadFile("appDetails.json")

	data := AppDetails{}

	_ = json.Unmarshal([]byte(file), &data)

	projectFolderString := fmt.Sprintf("%s_%s", data.Name, data.UniqueID)

	createProjectDirectory(data)
	downloadLaravel(projectFolderString)
	renameLaravelApp(projectFolderString, data.Name)
	updateEnvironmentFile(projectFolderString, data.Name)
	pushToGithub(projectFolderString, data.Name)

	fmt.Println("Initial Project Built")

	readFile(fmt.Sprintf("/Users/jetilling/projects/%s/%s/app/User.php", projectFolderString, data.Name))

}

func createProjectDirectory(data AppDetails) {
	cmd := exec.Command("/bin/sh", "./buildScripts/createProjectDirectory.sh", data.Name, data.UniqueID)
	runBashScript(cmd)
}

func downloadLaravel(projectFolderString string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/downloadLaravel.sh", projectFolderString)
	runBashScript(cmd)
}

func renameLaravelApp(projectFolderString string, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/renameLaravelApp.sh", projectFolderString, appName)
	runBashScript(cmd)
}

func updateEnvironmentFile(projectFolderString string, appName string) {
	cmd := exec.Command("/bin/sh", "./buildScripts/updateEnvironmentFile.sh", projectFolderString, appName)
	runBashScript(cmd)
}

func pushToGithub(projectFolderString string, appName string) {
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

func readFile(path string) {
	fmt.Println("Reading in: ", path)
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(dat))
}
