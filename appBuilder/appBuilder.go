package appBuilder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jetilling/projectBuilder/configVars"
)

type AppDetails struct {
	Name       string     `json:"name"`
	UniqueID   string     `json:"uniqueID"`
	DevDetails DevDetails `json:"devDetails"`
	// Models   []Models `json:"models"`
}

type DevDetails struct {
	UseReact bool `json:"useReact"`
	UseRedux bool `json:"useRedux"`
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

	rootPath := configVars.Config.ROOT_PATH
	projectFolderString := fmt.Sprintf("%s_%s", data.Name, data.UniqueID)
	projectFilePath := fmt.Sprintf("%s/%s/%s", rootPath, projectFolderString, data.Name)

	setupNewProject(data, projectFolderString, projectFilePath)

	if data.DevDetails.UseReact {
		addReact(projectFilePath, projectFolderString, data.Name)
	}

	// pushToGithub(projectFolderString, data.Name)

	fmt.Println("Initial Project Built")

	json.NewEncoder(w).Encode(GeneralResponse{Success: true, ErrorMessage: ""})

	// readFile(fmt.Sprintf("/Users/jetilling/projects/%s/%s/app/User.php", projectFolderString, data.Name))

}

func setupNewProject(data AppDetails, projectFolderString, projectFilePath string) {

	// inputEnvironmentFile := fmt.Sprintf("%s/.env.template", projectFilePath)
	// outputEnvironmentFile := fmt.Sprintf("%s/.env", projectFilePath)
	// inputDockerFile := fmt.Sprintf("%s/docker-compose_template.yml", projectFilePath)
	// outputDockerFile := fmt.Sprintf("%s/docker-compose.yml", projectFilePath)
	// inputTestDBFile := fmt.Sprintf("%s/create-testing-db_template.sql", projectFilePath)
	// outputTestDBFile := fmt.Sprintf("%s/create-testing-db.sql", projectFilePath)

	createProjectDirectory(data)
	downloadLaravel(projectFolderString)
	renameLaravelApp(projectFolderString, data.Name)
	copyEnvironmentFile(projectFilePath, projectFolderString, data.Name)
	copyDockerFiles(projectFilePath, projectFolderString, data.Name)

	// Here we need to find and replace {{project_name}} with the actual project name
	// We need to do this in the .env.example file and the docker-compose file

	// Environment
	// findAndReplace(inputEnvironmentFile, outputEnvironmentFile, "{{project_name}}", strings.ToLower(data.Name))

	// Docker Compose
	// findAndReplace(inputDockerFile, outputDockerFile, "{{project_name}}", strings.ToLower(data.Name))

	// Create Testing DB
	// findAndReplace(inputTestDBFile, outputTestDBFile, "{{project_name}}", strings.ToLower(data.Name))

	removeTemplateFiles(projectFolderString, data.Name)
}
