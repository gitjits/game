package main

import (
	"fmt"
	"image/color"
)

type GridTree struct {
	grid TileGrid
	prev *GridTree
	next *GridTree

	generation int
	commitHash string
}

func gitCommitGrid(g *Game, grid TileGrid, branch bool) string {
	old := g.gridTree
	if old.prev != nil {
		old.prev.next = &old
	}
	node := GridTree{
		grid:       grid,
		prev:       &old,
		next:       nil,
		generation: old.generation,
	}
	if branch {
		node.generation++
	}
	g.gridTree = node
	old.next = &g.gridTree

	return "blah"
}

func mergeCurrentBranch(g *Game) {
	if g.gridTree.generation == 0 {
		// There's nothing to merge up into if we're first generation.
		return
	}
	node := g.gridTree
	genOG := node.generation
	for node.generation == genOG {
		node.generation--
		node = *node.prev
	}

}

func nukeCurrentBranch(g *Game) {
	if g.gridTree.generation == 0 {
		// There's nothing to nuke if we're first generation.
		return
	}

	// Loop back through the nodes until we hit the previous branch
	node := g.gridTree
	genOG := node.generation
	for node.generation == genOG {
		node = *node.prev
	}

	g.gridTree = node
}

func gitSetup(g *Game) {
	g.gridTree = GridTree{}
}

func gitCurrentGrid(g *Game) TileGrid {
	return g.gridTree.grid
}

func commitTestData(g *Game) error {
	// Create initial commit on main
	hash := gitCommitGrid(g, createGrid(4, 4, 9, 9, 4, 4, color.RGBA{R: 255, B: 255, G: 255, A: 1}), false)

	fmt.Println("Created a commit on master", hash)

	// Add commits to feature1
	hash = gitCommitGrid(g, createGrid(4, 4, 9, 9, 4, 4, color.RGBA{R: 255, B: 0, G: 0, A: 1}), true)
	fmt.Println("Created a commit on feature1", hash)

	hash = gitCommitGrid(g, createGrid(4, 4, 9, 9, 4, 4, color.RGBA{R: 0, B: 0, G: 255, A: 1}), false)
	fmt.Println("Created a commit on feature1", hash)

	// Add commit to feature2
	hash = gitCommitGrid(g, createGrid(4, 4, 9, 9, 4, 4, color.RGBA{R: 0, B: 255, G: 0, A: 1}), true)
	fmt.Println("Created a commit on feature2", hash)

	return nil
}
