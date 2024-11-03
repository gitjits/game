package main

import (
	"fmt"
	"image/color"
	"io"

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

func ListBranches(r *git.Repository) ([]plumbing.Hash, error) {
	branches := []plumbing.Hash{}
	refs, err := r.References()
	if err != nil {
		return nil, err
	}

	refs.ForEach(func(ref *plumbing.Reference) error {
        if !ref.Name().IsBranch() {
            return nil
        }
        branches = append(branches, ref.Hash())
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

func CheckoutBranch(g *Game, branchName string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	opts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Force:  true,
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

func gitSetup(g *Game) error {
	// Create in-memory filesystem
	g.backingFS = memfs.New()

	// Setup git repo with in-memory filesystem
	repo, err := git.Init(memory.NewStorage(), g.backingFS)
	if err != nil {
		return err
	}
	g.repo = repo

	// Save the grid to an in-memory file
	file, err := g.backingFS.Create(gridFileName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return err
	}
	bytes, err := bson.Marshal(g.selectedGrid)
	if err != nil {
		fmt.Printf("Serialization error: %v\n", err)
		return err
	}

	file.Write(bytes)
	file.Close()

	w, err := g.repo.Worktree()
	if err != nil {
		fmt.Print("Error getting current worktree!", err, "\n")
		return err
	}
    _, err = w.Add(gridFileName)
    if err != nil {
        return err
    }

	err = createTestData(g)
	if err != nil {
		fmt.Printf("Error creating test data: %v\n", err)
		return err
	}
	fmt.Printf("Initialized git repo and state file \"%s\"\n", gridFileName)

	return nil
}

func gitCurrentGrid(g *Game) (TileGrid, error) {
	var grid TileGrid

	file, err := g.backingFS.Open(gridFileName)
	if err != nil {
		return grid, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return grid, err
	}

	err = bson.Unmarshal(data, &grid)

	return grid, err
}

func buildCommitTree(g *Game) (*GridTree, error) {
	var head *GridTree // Keep track of the head of our list

	worktree, err := g.repo.Worktree()
	worktree.Checkout(&git.CheckoutOptions{})
	if err != nil {
		fmt.Print("Error getting current worktree!", err, "\n")
		return &GridTree{}, err
	}

	branches, err := ListBranches(g.repo)
	if err != nil {
		return head, err
	}
	for _, branch := range branches {
		commits, err := g.repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime, From: branch})
		if err != nil {
			fmt.Print("Error getting log iterator: ", err, "\n")
			return &GridTree{}, err
		}

		commits.ForEach(func(commit *object.Commit) error {
			worktree.Checkout(&git.CheckoutOptions{Hash: commit.Hash})

			newNode := &GridTree{}

			grid, err := gitCurrentGrid(g)
			if err != nil {
				fmt.Println("Error on commit2", err)
				return err
			}
			newNode.grid = grid

			// Link the new node properly
			if head != nil {
				newNode.next = head
				head.parent = newNode
			}
			head = newNode

			return nil
		})
	}

	return head, err
}

func createTestData(g *Game) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Helper function to create a commit
	createCommit := func(message string, grid TileGrid) (plumbing.Hash, error) {
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
	_, err = createCommit("Initial commit on main", createGrid(4, 4, 5, 5, 4, 4, color.RGBA{R: 255, B: 255, G: 255, A: 1}))
	fmt.Print("Created a commit\n")
	if err != nil {
		return err
	}

	// Create feature1 branch from initial commit
	err = CreateBranch(g.repo, "feature1")
	if err != nil {
		return err
	}
	if err := CheckoutBranch(g, "feature1"); err != nil {
		return err
	}

	// Add commits to feature1
	_, err = createCommit("first com mit on feature1", createGrid(4, 4, 4, 4, 4, 4, color.RGBA{R: 255, B: 0, G: 0, A: 1}))
	if err != nil {
		return err
	}
	_, err = createCommit("Second commit on feature1", createGrid(4, 4, 4, 4, 4, 4, color.RGBA{R: 0, B: 0, G: 255, A: 1}))
	if err != nil {
		return err
	}

	// Back to main
	if err := CheckoutBranch(g, "master"); err != nil {
		return err
	}

	// Add more commits to main
	_, err = createCommit("Second commit on main", createGrid(4, 4, 4, 4, 4, 4, color.RGBA{R: 0, B: 255, G: 255, A: 1}))
	if err != nil {
		return err
	}

	// Create feature2 from current main
	err = CreateBranch(g.repo, "feature2")
	if err != nil {
		return err
	}
	if err := CheckoutBranch(g, "feature2"); err != nil {
		return err
	}

	// Add commit to feature2
	_, err = createCommit("First commit on feature2", createGrid(4, 4, 4, 4, 4, 4, color.RGBA{R: 255, B: 0, G: 255, A: 1}))
	if err != nil {
		return err
	}

	// Back to main for final commit
	if err := CheckoutBranch(g, "master"); err != nil {
		return err
	}
	_, err = createCommit("Third commit on main", createGrid(4, 4, 4, 4, 4, 4, color.RGBA{R: 255, B: 255, G: 0, A: 1}))
	if err != nil {
		return err
	}

	return nil
}
