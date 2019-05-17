package elasticsearch

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var jsonOk = []byte(`{  
	"snapshots":[  
	   {  
		  "snapshot":"backup-es-latest",
		  "repository":"my_backup",         
		  "stats":{             
			 "start_time_in_millis":1536243158017,             
			 "total_size_in_bytes":1371           
		  }                      
	   }
	]
 }`)

var jsonUnknowRepostitory = []byte(`{
	"error": {
	  "root_cause": [
		{
		  "type": "repository_missing_exception",
		  "reason": "[unknowRepostitory] missing"
		}
	  ],
	  "type": "repository_missing_exception",
	  "reason": "[unknowRepostitory] missing"
	},
	"status": 404
  }`)

func setupTest(t *testing.T, json []byte, status int) (*ElasticSearchRepo, func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(json)
	}))

	esRepo := OpenRepository("testRepo", ts.URL, "my_backup")
	// Test teardown - return a closure for use by 'defer'
	return esRepo, func() {
		ts.Close()
	}
}

func TestLatestSnapshot(t *testing.T) {
	t.Log("Testing a success case ")
	repo, teardown := setupTest(t, jsonOk, http.StatusOK)
	defer teardown()

	elasticSearchSnapshot, esError := repo.LatestSnapshot()
	if esError != nil {
		t.Fatalf("Unexpected error: %s", esError.Error())
	}

	if elasticSearchSnapshot.Name != snapshotName {
		t.Errorf("Name - Expected %s but got %s", snapshotName, elasticSearchSnapshot.Name)
	}

	if elasticSearchSnapshot.Size != float64(1371) {
		t.Errorf("Size - Expected %f but got %f", float64(1371), elasticSearchSnapshot.Size)
	}

	if elasticSearchSnapshot.DateString != "Thu Sep  6 14:12:38 UTC 2018" {
		t.Errorf("TimeInMillis - Expected %s but got %s", "Thu Sep  6 14:12:38 UTC 2018", elasticSearchSnapshot.DateString)
	}
}

func TestLatestSnapshotUnknowURL(t *testing.T) {
	t.Log("Testing an unknow URL")
	repo := OpenRepository("testRepo", "", "my_backup")
	_, err := repo.LatestSnapshot()
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestLatestSnapshotUnknowRespository(t *testing.T) {
	t.Log("Testing an unknow repository")

	repo, teardown := setupTest(t, jsonUnknowRepostitory, http.StatusNotFound)
	defer teardown()

	_, esError := repo.LatestSnapshot()
	if esError == nil {
		t.Errorf("Expected an error")
	}
}
