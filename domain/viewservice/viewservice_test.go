package viewservice_test

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/rerost/bqv/domain/viewmanager"
	"github.com/rerost/bqv/domain/viewservice"
)

func PrepareDirForTest(dir string) error {
	projectRoot := os.Getenv("PROJECT_ROOT")
	cmd := exec.Command("cp", "-a", path.Join(projectRoot, "example/dataset"), dir)
	_, err := cmd.CombinedOutput()
	return err
}

func DiffDir(dir1, dir2 string) (string, error) {
	cmd := exec.Command("diff", dir1, dir2)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func TestViewServiceCopy(t *testing.T) {
	ctx := context.Background()
	dirTest1, err := ioutil.TempDir("", "test1")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirTest1)
	PrepareDirForTest(dirTest1)

	dirTest2, err := ioutil.TempDir("", "test1")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirTest2)

	fileManager1 := viewmanager.NewFileManager(dirTest1)
	fileManager2 := viewmanager.NewFileManager(dirTest2)
	service := viewservice.NewService()
	if err := service.Copy(ctx, fileManager1, fileManager2); err != nil {
		t.Error(err)
		return
	}

	if diff, err := service.Diff(ctx, fileManager1, fileManager2); err != nil {
		if err != nil {
			t.Error(err)
			return
		}
		if len(diff) != 0 {
			t.Error(diff)
		}
	}

	if diff, err := DiffDir(dirTest1, dirTest2); err != nil || diff != "" {
		if err != nil {
			t.Error(err)
			t.Error(diff)
			return
		}
	}
}
