package migrations

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bradcypert/deckard/lib/db"
)

/*
* Up or down?
 */
func FindInPath(dir string, isUp bool) db.Migration {
	queries := make([]db.Query, 0)

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
			queries = append(queries, db.Query{
				Name:  file.Name(),
				Value: string(contents),
			})
		}
	}
	return db.Migration{
		Queries: queries,
	}
}

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
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()

	return err
}

func makeTimestamp(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()/int64(time.Millisecond), 10)
}
