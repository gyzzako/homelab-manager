package git

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"gopkg.in/yaml.v3"

	"homelab-manager/internal/hosts/providers"
)

type GitConfig struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

const HOST_DATA_FILE_NAME = "host_entries.yml"

func PushToGit(entries []providers.HostEntry, gitCfg GitConfig) error {
	dir := os.TempDir()
	localRepoPath := filepath.Join(dir, "homelab-git")

	_ = os.RemoveAll(localRepoPath)

	repo, err := cloneRepo(localRepoPath, gitCfg)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err := writeDataToFile(localRepoPath, entries); err != nil {
		return err
	}

	if err := doChangesExist(w); err != nil {
		return err
	}

	_, _ = w.Add(HOST_DATA_FILE_NAME)
	msg := fmt.Sprintf("update host entries %s", time.Now().Format(time.RFC3339))
	_, err = w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Homelab manager",
			Email: "",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "homelab-manager",
			Password: gitCfg.Token,
		},
	})
}

func cloneRepo(localRepoPath string, gitCfg GitConfig) (*git.Repository, error) {
	repo, err := git.PlainClone(localRepoPath, false, &git.CloneOptions{
		URL: gitCfg.URL,
		Auth: &http.BasicAuth{
			Username: "homelab-manager",
			Password: gitCfg.Token,
		},
	})

	if err != nil {
		if err == transport.ErrEmptyRemoteRepository {
			return handleEmptyRepo(localRepoPath, gitCfg)
		}

		return nil, fmt.Errorf("clone failed: %w", err)
	}

	return repo, nil
}

func handleEmptyRepo(localRepoPath string, gitCfg GitConfig) (*git.Repository, error) {
	repo, err := git.PlainInit(localRepoPath, false)
	if err != nil {
		return repo, err
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{gitCfg.URL},
	})
	if err != nil {
		return repo, err
	}

	return repo, nil
}

func doChangesExist(w *git.Worktree) error {
	status, err := w.Status()
	if err != nil {
		return err
	}

	if status.IsClean() {
		return fmt.Errorf("nothing to commit")
	}

	return nil
}

func writeDataToFile(localRepoPath string, entries []providers.HostEntry) error {
	ymlPath := filepath.Join(localRepoPath, HOST_DATA_FILE_NAME)
	file, err := os.Create(ymlPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := yaml.NewEncoder(file).Encode(entries); err != nil {
		return err
	}

	return nil
}
