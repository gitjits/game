package main

import (
	"fmt"

	memfs "github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	memory "github.com/go-git/go-git/v5/storage/memory"
	bson "gopkg.in/mgo.v2/bson"
)

const (
	gridFileName = "gridfile"
)

func gitSetup(g *Game) bool {
	// Setup git repo with in-memory filesystem
	repo, e := git.Init(memory.NewStorage(), g.backingFS)
	if e != nil {
		fmt.Printf("Error %v\n", e)
		return false
	}
	g.repo = repo

	// Create in-memory filesystem
	g.backingFS = memfs.New()

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
		return false
	}
	worktree.Add(gridFileName)
	worktree.Commit("Initial commit", &git.CommitOptions{})
	fmt.Printf("Initialized git repo and state file \"%s\"\n", gridFileName)
	return true
}
