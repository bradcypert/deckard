package migrations

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//Migration defines a structure for holding metadata and queries to be ran against a database
type Migration struct {
	Queries []Query
}

//Query defines a name/value pair for queries where the name is expected to be a filename and the value is the SQL query
type Query struct {
	Name  string
	Value string
}

/*FindInPath finds migration files in a specific directory. `dir` param is the directory you want to search.
 *isUp is a bool representing if you're searching for up or down queries.
 */
func FindInPath(dir string, isUp bool) Migration {
	queries := make([]Query, 0)

	var suffix string
	if isUp == true {
		suffix = ".up.sql"
	} else {
		suffix = ".down.sql"
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			contents, _ := ioutil.ReadFile(file.Name())
			queries = append(queries, Query{
				Name:  file.Name(),
				Value: string(contents),
			})
		}
	}
	return Migration{
		Queries: queries,
	}
}

//Create creates a new migration with the provided name at the given directory.
func Create(outputDir string, name string) {
	// Add in the / suffix if it wasn't added by the user
	if len(outputDir) > 0 && !strings.HasSuffix(outputDir, "/") {
		outputDir += "/"
	}

	filepath := outputDir + makeTimestamp(time.Now()) + "__" + name
	upError := createFile(filepath + ".up.sql")
	downError := createFile(filepath + ".down.sql")

	fmt.Printf("Created file %s\n", filepath+".up.sql")
	fmt.Printf("Created file %s\n", filepath+".down.sql")

	if upError != nil {
		log.Fatal(upError)
	}

	if downError != nil {
		log.Fatal(downError)
	}
}

func createFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()

	return err
}

func makeTimestamp(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()/int64(time.Millisecond), 10)
}
