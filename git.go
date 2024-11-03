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
	node := GridTree{
		grid:       grid,
		prev:       &g.gridTree,
		next:       nil,
		generation: g.gridTree.generation,
	}
	if branch {
		node.generation++
	}
	g.gridTree.next = &node
	g.gridTree = node
	// fmt.Println("prev node %p, replacing with %p. old's next is %p\n", node.prev, &node, node.prev.next)

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
	hash := gitCommitGrid(g, createGrid(4, 4, 5, 5, 4, 4, color.RGBA{R: 255, B: 255, G: 255, A: 1}), false)

	fmt.Println("Created a commit on master", hash)

	// Add commits to feature1
	hash = gitCommitGrid(g, createGrid(4, 4, 4, 4, 4, 4, color.RGBA{R: 255, B: 0, G: 0, A: 1}), true)
	fmt.Println("Created a commit on feature1", hash)

	hash = gitCommitGrid(g, createGrid(4, 4, 4, 4, 4, 4, color.RGBA{R: 0, B: 0, G: 255, A: 1}), false)
	fmt.Println("Created a commit on feature1", hash)

	// Add commit to feature2
	hash = gitCommitGrid(g, createGrid(4, 4, 4, 4, 4, 4, color.RGBA{R: 255, B: 0, G: 255, A: 1}), true)
	fmt.Println("Created a commit on feature2", hash)

	return nil
}
