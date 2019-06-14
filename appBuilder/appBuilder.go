package appBuilder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jetilling/projectBuilder/configVars"
)

type AppDetails struct {
	Name              string `json:"name"`
	UniqueID          string `json:"uniqueID"`
	frontEndFramework string `json:"frontEndFramework"`
	// Models   []Models `json:"models"`
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

type GeneralResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}

func Build(w http.ResponseWriter, r *http.Request) {

	// INITIALIZE CONFIG VARIABLES
	configVars.InitConfigVars()

	data := AppDetails{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	fmt.Println("DATA: ", data)
	rootPath := configVars.Config.ROOT_PATH
	projectFolderString := fmt.Sprintf("%s_%s", data.Name, data.UniqueID)
	projectFilePath := fmt.Sprintf("%s/%s/%s", rootPath, projectFolderString, data.Name)

	setupNewProject(data, projectFolderString, projectFilePath)

	switch data.frontEndFramework {
	case "react":
		addReact(projectFilePath, projectFolderString, data.Name)
	case "vue":
	case "jquery":
	case "reasonml":
		fmt.Println("framework not ready")
	}

	pushToGithub(projectFolderString, data.Name)

	fmt.Println("Initial Project Built")

	// readFile(fmt.Sprintf("/Users/jetilling/projects/%s/%s/app/User.php", projectFolderString, data.Name))

	json.NewEncoder(w).Encode(GeneralResponse{Success: true, ErrorMessage: ""})

}

func setupNewProject(data AppDetails, projectFolderString, projectFilePath string) {

	createProjectDirectory(data)
	downloadLaravel(projectFolderString)
	renameLaravelApp(projectFolderString, data.Name)
	copyEnvironmentFile(projectFilePath, projectFolderString, data.Name)
	copyDockerFiles(projectFilePath, projectFolderString, data.Name)

	removeTemplateFiles(projectFolderString, data.Name)
}
