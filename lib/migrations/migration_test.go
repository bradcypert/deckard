package migrations

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMakeTimestamp(t *testing.T) {
	timestamp := makeTimestamp(time.Date(2000, 4, 17, 4, 44, 44, 0, time.UTC))
	if timestamp != "955946684000" {
		t.Error("Expected to find:", 0, "instead got:", timestamp)
	}
}

func TestCreate(t *testing.T) {
	os.Mkdir("./tmp", 0777)
	defer os.RemoveAll("./tmp")
	Create("./tmp", "test_create")

	files, err := ioutil.ReadDir("./tmp")

	if err != nil {
		t.Error("Unable to create file")
	}

	if !strings.Contains(files[0].Name(), "test_create") {
		t.Error("Unable to find created file. Creation may have failed.")
	}

	if len(files) < 2 {
		t.Error("Did not create an up and a down migration.")
	}
}

func TestFindInPath(t *testing.T) {
	os.Mkdir("./tmp", 0777)
	defer os.RemoveAll("./tmp")
	os.Create("./tmp/123__foo.up.sql")
	os.Create("./tmp/123__foo.down.sql")

	migration := FindInPath("./tmp", true)

	if len(migration.Queries) != 1 {
		t.Error("Did not find exactly 1 up query")
	}

	migration = FindInPath("./tmp", false)

	if len(migration.Queries) != 1 {
		t.Error("Did not find exactly 1 down query")
	}
}
