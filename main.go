package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"io/fs"
	"log"
	"os"
)

// Main calls the correct function as per the CLI arguments passed.
func main() {
	args := os.Args
	count := len(args)
	if count > 3 {
		log.Fatal("Error: too many arguments")
	} else if count == 1 {
		sync()
	} else if count == 2 && args[1] == "config" {
		config()
	} else if count == 3 && args[1] == "push" {
		push(args[2])
	} else {
		log.Fatal("Unsupported command: ", args)
	}
}

// Sync uploads changes if local repo is up-to-date with remote.
func sync() {
	conf = conf.load()
	if !pull() {
		push(conf.Dir)
	}
}

// Push calls commit, then uploads the specified dir.
func push(dir string) {
	repo := open(dir)
	commit(repo)
	err := repo.Push(&git.PushOptions{})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		log.Fatal("repo.Push failed: ", err)
	}
}

// Pull tries to get the latest commit from remote
// Returns false if local repo is up-to-date.
func pull() (update bool) {
	repo := open(conf.Dir)
	wtree, err := repo.Worktree()
	if err != nil {
		log.Fatal("repo.Worktree failed: ", err)
	}
	err = wtree.Pull(&git.PullOptions{})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return false
		} else if errors.Is(err, git.ErrUnstagedChanges) {
			log.Fatal("Local repo conflicts with remote")
		} else {
			log.Fatal("wtree.Pull failed: ", err)
		}
	}
	return true
}

// Commit adds the local dir to a git commit for pushing.
func commit(repo *git.Repository) {
	wtree, err := repo.Worktree()
	if err != nil {
		log.Fatal("repo.Worktree failed: ", err)
	}
	_, err = wtree.Add(".")
	if err != nil {
		log.Fatal("wtree.Add failed: ", err)
	}
	_, err = wtree.Commit(
		"Uploaded with syncnotes", &git.CommitOptions{})
	if err != nil {
		log.Fatal("wtree.Commit failed: ", err)
	}
}

// Open loads and returns a git repo from the notes directory.
func open(dir string) *git.Repository {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			log.Fatal("Local repo not found, please run 'syncnotes config'")
		} else {
			log.Fatal("git.PlainOpen failed: ", err)
		}
	}
	return repo
}

// Clone runs the git clone command inside the notes directory.
func clone() {
	options := &git.CloneOptions{
		URL:   conf.URL,
		Depth: 0,
	}
	err := options.Validate()
	if err != nil {
		log.Fatal("options.Validate failed: ", err)
	}
	_, err = git.PlainClone(conf.Dir, false, options)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			fmt.Println()
			fmt.Println("Warning: repo already exists inside specified dir")
			fmt.Println("If you wish to overwrite, delete .git manually")
		} else {
			log.Fatal("git.PlainClone failed: ", err)
		}
	}
}

// Config gets user configuration options, saves them by calling load(),
// then clones the notes repository into the specified directory.
func config() {
	fmt.Print("Local notes directory: ")
	_, err := fmt.Scanln(&conf.Dir)
	if err != nil {
		log.Fatal("fmt.Scanln failed: ", err)
	}
	fmt.Print("Github username: ")
	_, err = fmt.Scanln(&conf.Username)
	if err != nil {
		log.Fatal("fmt.Scanln failed: ", err)
	}
	fmt.Print("Github repo name: ")
	_, err = fmt.Scanln(&conf.Repo)
	if err != nil {
		log.Fatal("fmt.Scanln failed: ", err)
	}
	fmt.Print("Private repository? [y/n]: ")
	var ssh, URL string
	_, err = fmt.Scanln(&ssh)
	if err != nil {
		log.Fatal("fmt.Scanln failed: ", err)
	}
	if ssh == "y" {
		conf.Method = "ssh"
		URL = "git@github.com:"
	} else {
		conf.Method = "https"
		URL = "https://github.com/"
	}
	conf.URL = fmt.Sprint(URL, conf.Username, "/", conf.Repo, ".git")
	conf.write()
	clone()
}

// Conf seems to be necessary only for clone, basically obsolete.
type Conf struct {
	Dir      string
	Username string
	Repo     string
	Method   string
	URL      string
}

var conf = Conf{}

// Load opens the JSON configuration file from the program folder,
// then unmarshals it and returns it as a Conf struct.
func (c Conf) load() Conf {
	confPath, _ := os.Executable()
	confPath = fmt.Sprint(confPath, ".conf")
	data, err := os.ReadFile(confPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Fatal("Config file not found, please run 'syncnotes config'")
		} else {
			log.Fatal("os.ReadFile failed: ", err)
		}
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Fatal("json.Unmarshal failed: ", err)
	}
	return c
}

// Write marshals the Conf struct to JSON and writes it to a
// 'syncnotes.conf' file in the program folder.
func (c Conf) write() {
	confPath, _ := os.Executable()
	confPath = fmt.Sprint(confPath, ".conf")
	data, err := json.Marshal(c)
	if err != nil {
		log.Fatal("json.Marshal failed: ", err)
	}
	err = os.WriteFile(confPath, data, 0666)
	if err != nil {
		log.Fatal("os.WriteFile failed: ", err)
	}
}
