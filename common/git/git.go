package git

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fanghongbo/ops-agent/utils"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"io/ioutil"
	"os"
	"path"
)

type NewGitClient struct {
	Url           string
	Path          string
	Username      string
	Password      string
	Repo          *git.Repository
	Buffer        *bytes.Buffer
	RepoType      int64  // 0 表示是公有仓库; 1 表示是私有仓库
	DefaultBranch string // 默认分支
	RemoteName    string // 远程名称, 默认是 origin
}

func (u *NewGitClient) auth() transport.AuthMethod {
	// 公有仓库不需要认证
	if u.RepoType == 0 {
		return nil
	} else {
		// 私有仓库返回认证信息
		return &http.BasicAuth{
			Username: u.Username,
			Password: u.Password,
		}
	}
}

func (u *NewGitClient) newDefaultBranch() string {
	// 默认master分支
	if u.DefaultBranch == "" {
		return "master"
	} else {
		return u.DefaultBranch
	}
}

func (u *NewGitClient) isGitDir() bool {
	var (
		gitDir     string
		objectsDir string
		refsDir    string
	)

	gitDir = path.Join(u.Path, ".git")
	if !utils.IsDir(gitDir) {
		return false
	}

	objectsDir = path.Join(gitDir, "objects")
	if !utils.IsDir(objectsDir) {
		return false
	}

	refsDir = path.Join(gitDir, "refs")
	if !utils.IsDir(refsDir) {
		return false
	}

	return true
}

func (u *NewGitClient) Update() error {
	var (
		files []os.FileInfo
		err   error
	)

	u.Buffer = new(bytes.Buffer)

	if utils.IsDir(u.Path) {
		if u.isGitDir() {
			return u.pull()
		} else {
			files, err = ioutil.ReadDir(u.Path)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				return u.clone()
			} else {
				return errors.New(fmt.Sprintf("%s 目录已经存在, 且是一个非空目录", u.Path))
			}
		}
	} else {
		return u.clone()
	}
}

func (u *NewGitClient) clone() error {
	var (
		err error
	)

	u.Repo, err = git.PlainClone(u.Path, false, &git.CloneOptions{
		RemoteName:    u.newRemoteName(),
		Auth:          u.auth(),
		URL:           u.Url,
		Progress:      u.Buffer,
		ReferenceName: plumbing.NewBranchReferenceName(u.newDefaultBranch()),
	})

	if err != nil {
		return err
	}
	return nil
}

func (u *NewGitClient) newRemoteName() string {
	if u.RemoteName == "" {
		return "origin"
	} else {
		return u.RemoteName
	}
}

func (u *NewGitClient) pull() error {
	var (
		err error
		w   *git.Worktree
	)

	if u.Repo, err = git.PlainOpen(u.Path); err != nil {
		return err
	}

	if w, err = u.Repo.Worktree(); err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName:    u.newRemoteName(),
		Auth:          u.auth(),
		Progress:      u.Buffer,
		ReferenceName: plumbing.NewBranchReferenceName(u.newDefaultBranch()),
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *NewGitClient) Head() (string, error) {
	var (
		repo *git.Repository
		ref  *plumbing.Reference
		err  error
	)

	if !u.isGitDir() {
		return "", errors.New("plugin dir is not exist")
	}

	repo, err = git.PlainOpen(u.Path)
	if err != nil {
		return "", err
	}

	ref, err = repo.Head()
	if err != nil {
		return "", err
	}

	return ref.Hash().String(), nil
}
