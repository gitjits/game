package main

import (
	"fmt"
	"image/color"
	"io/ioutil"

	memfs "github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	memory "github.com/go-git/go-git/v5/storage/memory"
	bson "gopkg.in/mgo.v2/bson"
)

const (
	gridFileName = "gridfile"
)

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
	grid   TileGrid
	parent *GridTree
	next   *GridTree
	branch *GridTree
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
	worktree.Commit("commit 2", &git.CommitOptions{})
	worktree.Commit("commit 3", &git.CommitOptions{})
	err = createTestData(g.repo)
	if err != nil {
		fmt.Printf("Error creating test data: %v\n", err)
		return false
	}
	fmt.Printf("Initialized git repo and state file \"%s\"\n", gridFileName)

	return true
}

func gitCurrentGrid(g *Game) (TileGrid, error) {
	var grid TileGrid

	file, err := g.backingFS.Open(gridFileName)
	if err != nil {
		return grid, err
	}
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return grid, err
    }

	err = bson.Unmarshal(data, &grid)

	return grid, err
}

func iterCommits(g *Game) GridTree {
	var output GridTree

	commits, err := g.repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		fmt.Print("Error getting log iterator: ", err, "\n")
		return output
	}

	worktree, err := g.repo.Worktree()
	if err != nil {
		fmt.Print("Error getting current worktree!", err, "\n")
		return output
	}

	commit, err := commits.Next()
	if err != nil {
		fmt.Print("Error first commit: ", err, "\n")
		return output
	}

	// We need to know the initial HEAD to make sure we reset state before returning
	initialHash := commit.Hash

    commits.ForEach(func(commit *object.Commit) error {
		fmt.Print(commit)
		// Revert to this specific commit
		worktree.Checkout(&git.CheckoutOptions{Hash: commit.Hash})

		// Read grid data at this commit and update tree
		var newNode GridTree
		newNode.next = &output
		output.parent = &newNode
		grid, err := gitCurrentGrid(g)
        if err != nil {
            fmt.Println("Error on commit2", err)
            return err
        }
		newNode.grid = grid

		// Update root node
		output = newNode
        return nil
	})

	// Revert back to the original HEAD
	worktree.Checkout(&git.CheckoutOptions{Hash: initialHash})
	return output
}

func createTestData(repo *git.Repository) error {
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Helper function to create a commit
	createCommit := func(message string, c color.RGBA) (plumbing.Hash, error) {
        grid := createGrid(4, 4, 4, 4, 4, 4, c)
		file, err := w.Filesystem.Create(gridFileName)
		if err != nil {
			return plumbing.ZeroHash, err
		}
        bytes, err := bson.Marshal(grid)
        if err != nil {
            fmt.Printf("Serialization error: %v\n", err)
            return plumbing.ZeroHash, err
        }

        file.Write(bytes)
        file.Close()

		_, err = w.Add(gridFileName)
		if err != nil {
			return plumbing.ZeroHash, err
		}

		hash, err := w.Commit(message, &git.CommitOptions{})
		return hash, err
	}

	// Create initial commit on main
    _, err = createCommit("Initial commit on main", color.RGBA{R: 255, B: 255, G: 255})
	fmt.Print("Created a commit\n")
	if err != nil {
		return err
	}

	// Create feature1 branch from initial commit
	err = CreateBranch(repo, "feature1")
	if err != nil {
		return err
	}
	if err := CheckoutBranch(repo, "feature1"); err != nil {
		return err
	}

	// Add commits to feature1
	_, err = createCommit("First commit on feature1", color.RGBA{R: 255, B: 0, G: 0})
	if err != nil {
		return err
	}
	_, err = createCommit("Second commit on feature1", color.RGBA{R: 0, B: 0, G: 255})
	if err != nil {
		return err
	}

	// Back to main
	if err := CheckoutBranch(repo, "master"); err != nil {
		return err
	}

	// Add more commits to main
	_, err = createCommit("Second commit on main", color.RGBA{R: 0, B: 255, G: 255})
	if err != nil {
		return err
	}

	// Create feature2 from current main
	err = CreateBranch(repo, "feature2")
	if err != nil {
		return err
	}
	if err := CheckoutBranch(repo, "feature2"); err != nil {
		return err
	}

	// Add commit to feature2
	_, err = createCommit("First commit on feature2", color.RGBA{R: 255, B: 0, G: 255})
	if err != nil {
		return err
	}

	// Back to main for final commit
	if err := CheckoutBranch(repo, "master"); err != nil {
		return err
	}
	_, err = createCommit("Third commit on main", color.RGBA{R: 255, B: 255, G: 0})
	if err != nil {
		return err
	}

	return nil
}
