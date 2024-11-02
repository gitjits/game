package main

import (
	"fmt"

	memfs "github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	memory "github.com/go-git/go-git/v5/storage/memory"
	bson "gopkg.in/mgo.v2/bson"
)

const (
	gridFileName = "gridfile"
)

func CreateBranch(r *git.Repository, branchName string) error {
	headRef, err := r.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), headRef.Hash())

	err = r.Storer.SetReference(ref)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

func ListBranches(r *git.Repository) ([]string, error) {
	branches := []string{}
	refs, err := r.References()
	if err != nil {
		return nil, err
	}

	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branches = append(branches, ref.Name().Short())
		}
		return nil
	})

	return branches, nil
}

func CheckoutBranch(r *git.Repository, branchName string) error {
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	opts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Force:  false,
	}

	if err := w.Checkout(opts); err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	return nil
}

type GridTree struct {
	grid TileGrid
	prev *GridTree
	next *GridTree
}

func gitSetup(g *Game) bool {
	// Create in-memory filesystem
	g.backingFS = memfs.New()

	// Setup git repo with in-memory filesystem
	repo, e := git.Init(memory.NewStorage(), g.backingFS)
	if e != nil {
		fmt.Printf("Error %v\n", e)
		return false
	}
	g.repo = repo

	// Save the grid to an in-memory file
	file, err := g.backingFS.Create(gridFileName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return false
	}
	bytes, err := bson.Marshal(g.grid)
	if err != nil {
		fmt.Printf("Serialization error: %v\n", err)
		return false
	}

	file.Write(bytes)
	file.Close()

	// Commit initial game state to current branch
	worktree, err := g.repo.Worktree()
	if err != nil {
		fmt.Printf("Error opening worktree: %v\n", err)
		return false
	}
	worktree.Add(gridFileName)
	worktree.Commit("Initial commit", &git.CommitOptions{})
	fmt.Printf("Initialized git repo and state file \"%s\"\n", gridFileName)

	return true
}

func iterCommits(repo *git.Repository) {
	commits, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		fmt.Print("Error getting log iterator: ", err, "\n")
		return
	}

	fmt.Print("Got commit log!\n")
	commit, err := commits.Next()
	for err == nil {
		fmt.Print(commit)
		commit, err = commits.Next()
	}
}
