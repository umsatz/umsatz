package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testManager struct {
	server *httptest.Server
}

var TestManager testManager

func TestMain(m *testing.M) {
	provider, _ := newBackupProvider("./dummy")
	restorer := &DummyBackupRestorer{}
	TestManager.server = httptest.NewServer(newBackupServer(provider, restorer))
	ret := m.Run()
	os.Exit(ret)
}

func TestListAvailableBackups(t *testing.T) {
	response, err := http.Get(TestManager.server.URL + "/")
	if err != nil {
		t.Fatal("list request failed: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tcode: %v", "200", response.StatusCode)
	}

	decoder := json.NewDecoder(response.Body)
	var backups []Backup

	if err := decoder.Decode(&backups); err != nil {
		t.Fatalf("Unable to decode json response: %#v", err)
	}

	if len(backups) != 1 {
		t.Fatalf("Did not respond correct dummy responses")
	}

	backup := backups[0]
	if backup.BaseName() != "backup-2014-09-08" {
		t.Fatalf("Wrong backup name %v", backup.BaseName())
	}
	if backup.CreatedAt != "2014-09-08" {
		t.Fatalf("Wrong creation date %v", backup.CreatedAt)
	}
}

func TestDeleteBackup(t *testing.T) {
	os.Create("./dummy/backup-2014-09-09.tar")

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%v/backup-2014-09-09", TestManager.server.URL), nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("list request failed: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tcode: %v, %v", "200", response.StatusCode, response.Body)
	}
}

type DummyBackupRestorer struct{}

func (d *DummyBackupRestorer) Restore(*Backup) error {
	return nil
}

func TestRestoreBackup(t *testing.T) {
	response, err := http.Post(fmt.Sprintf("%v/backup-2014-09-08/restore", TestManager.server.URL), "application/json", nil)
	if err != nil {
		t.Fatal("list request failed: %v", err)
	}

	if response.StatusCode != http.StatusAccepted {
		t.Fatalf("Response body did not contain expected %v:\n\tcode: %v, %v", "200", response.StatusCode, response.Body)
	}
}
