package main

import (
	"fmt"
	"image/color"
	"math/rand"
)

type GridTree struct {
	grid TileGrid
	prev *GridTree
	next *GridTree

	generation int
	commitHash string
}

func generateRandomHash() string {
    b := make([]byte, 20)
    rand.Read(b)
    return fmt.Sprintf("%x", b)
}

func gitCommitGrid(g *Game, grid TileGrid, branch bool, cls bool) string {
	old := g.gridTree
	if old.prev != nil {
		old.prev.next = &old
	}
	node := GridTree{
		grid:       grid,
		prev:       &old,
		next:       nil,
		generation: old.generation,
        commitHash: generateRandomHash(),
	}
    if branch && node.generation == 4 {
        g.logger.AddMessage("[!] ", "Maximum allowed branches", true)
        return ""
    }
	if branch {
		node.generation++
        g.logger.AddMessage("you$ ", fmt.Sprintf("git checkout -b branch%d", node.generation), false)
	}
	g.gridTree = node
	old.next = &g.gridTree
	g.autoScroll = true
	g.selected = &grid
	g.selected.IsSelectedGrid = true

    if !cls {
        g.logger.AddMessage("you$ ", "git commit -m 'move a piece'", false)
        g.logger.AddMessage("", fmt.Sprintf("[main %s] move a piece", node.commitHash[0:8]), true)
        g.logger.AddMessage("", "1 files changed, 1 insertions(+), 0 deletions(-)", true)
    }
	return  node.commitHash
}

func mergeCurrentBranch(g *Game) {
    fmt.Println("merge", g.gridTree.generation)
    if g.gridTree.generation == 0 {
        g.logger.AddMessage("[!] ", "You cannot merge", false)
        return
    }
    if g.gridTree.generation - 1 > 0 {
        g.logger.AddMessage("you$ ", fmt.Sprintf("git checkout branch%d", g.gridTree.generation - 1), false)
    } else {
        g.logger.AddMessage("you$ ", "git checkout main", false)
    }
    g.logger.AddMessage("you$ ", "git merge " + g.gridTree.commitHash, false)
    g.logger.AddMessage("", "Updating " + g.gridTree.prev.commitHash, true)
    g.logger.AddMessage("", "Fast-forward", true)
    g.logger.AddMessage("", " board.bson | 1 +", true)
    g.logger.AddMessage("", " 1 file changed, 1 insertions(+)", true)
	if g.gridTree.generation == 0 {
		// There's nothing to merge up into if we're first generation.
		return
	}
	node := &g.gridTree
	genOG := node.generation
	for node.generation == genOG {
		node.generation--
		node = node.prev
	}
}

func nukeCurrentBranch(g *Game) {
	if g.gridTree.grid.SizeX == 0 || g.gridTree.prev == nil || g.gridTree.prev.grid.SizeX == 0 || g.gridTree.prev.prev == nil || g.gridTree.prev.prev.grid.SizeX == 0 {
        g.logger.AddMessage("[!] ", "Can't revert", false)
		return
	} else {
        g.logger.AddMessage("[!] ", "Reverting", false)
    }

    //node := g.gridTree.prev
	//g.gridTree = *node
    prevNode := g.gridTree.prev
    prevNode.next = nil  // Break the circular reference
    g.gridTree = *prevNode
}

func gitSetup(g *Game) {
	g.gridTree = GridTree{}
}

func gitCurrentGrid(g *Game) TileGrid {
	return g.gridTree.grid
}

func commitTestData(g *Game) error {
	// Create initial commit on main
	hash := gitCommitGrid(g, createGrid(4, 4, 9, 9, 4, 4, color.RGBA{R: 255, B: 255, G: 255, A: 1}), false, false)

	fmt.Println("Created a commit on master", hash)

	// Add commits to feature1
	hash = gitCommitGrid(g, createGrid(4, 4, 9, 9, 4, 4, color.RGBA{R: 255, B: 0, G: 0, A: 1}), true, false)
	fmt.Println("Created a commit on feature1", hash)

	hash = gitCommitGrid(g, createGrid(4, 4, 9, 9, 4, 4, color.RGBA{R: 0, B: 0, G: 255, A: 1}), false, false)
	fmt.Println("Created a commit on feature1", hash)

	// Add commit to feature2
	hash = gitCommitGrid(g, createGrid(4, 4, 9, 9, 4, 4, color.RGBA{R: 0, B: 255, G: 0, A: 1}), true, false)
	fmt.Println("Created a commit on feature2", hash)

	return nil
}
