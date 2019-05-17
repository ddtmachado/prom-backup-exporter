package restic

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ddtmachado/prom-backup-exporter/collector"
)

var test1Json = []byte(`{
		"time": "2018-09-12T09:17:07.286761313-03:00",
    "tags": [ "test1" ],
    "short_id": "3fb55586"
  }`)

var test2Json = []byte(`{
	"time": "2018-09-12T09:18:32.229110942-03:00",
  "tags":[ "test2" ],
  "short_id": "f8da39e9"
}`)

func fakeExecCommand(command string, args ...string) *exec.Cmd {

	testHelper := getTestHelper(args...)
	cs := []string{"-test.run=" + testHelper, "--", command}
	cs = append(cs, args...)
	log.Println(cs)
	cmd := exec.Command(os.Args[0], cs...)
	return cmd
}

func getTestHelper(args ...string) string {
	helperProcess := "TestHelperProcess"

	cmdArgs := strings.Join(args, ", ")
	if strings.Contains(cmdArgs, "--tag") {
		helperProcess = "TestHelperTagProcess"
	}
	log.Println(helperProcess)
	return helperProcess
}

func fakeExecCommandError(command string, args ...string) *exec.Cmd {
	cmd := fakeExecCommand(command, args...)
	return cmd
}

func TestLatestSnapshots(t *testing.T) {

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	repo := &ResticRepository{
		Alias:    "test2",
		Password: "myPassword",
		Path:     "myPath",
	}

	snapshot, err := repo.LatestSnapshot()
	if err != nil {
		t.Fatalf("Unexpected error - '%s'", err.Error())
	}

	testResticSnapshot := &resticSnapshot{}
	err = json.Unmarshal(test2Json, &testResticSnapshot)
	if err != nil {
		t.Fatal(err)
	}
	compareResticSnapshots(t, testResticSnapshot, snapshot)
}

func TestHelperProcess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	os.Exit(0)
}

func TestHelperTagProcess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	for idx, arg := range os.Args[1:] {
		if arg == "--tag" {
			templateName := strings.Split(os.Args[idx+2], "=")[1]
			switch templateName {
			case "test1":
				fmt.Fprintf(os.Stdout, "[ %s ]", test1Json)
			case "test2":
				fmt.Fprintf(os.Stdout, "[ %s ]", test2Json)
			default:
				fmt.Fprintf(os.Stderr, "%s", "no json file template found")
				os.Exit(1)
			}
		}
	}
	os.Exit(0)
}

func compareResticSnapshots(t *testing.T, expected *resticSnapshot, returned *collector.BackupSnapshot) {
	if expected.Id != returned.Name {
		t.Errorf("Expected Id %s but got %s", expected.Id, returned.Name)
	}
	if strings.Compare(expected.Time, returned.DateString) == 0 {
		t.Errorf("Expected Time %s but got %s", expected.creationDateString(), returned.DateString)
	}
}
