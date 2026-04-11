package script

import (
	"github.com/ecumeurs/upsiloncli/internal/dto"
)

type queueItem struct {
	pos  dto.Position
	path []dto.Position
}

// FindPath calculates the shortest path from start to end on the board.
// It avoids obstacles and cells occupied by living entities.
func FindPath(board *dto.BoardState, start, end dto.Position) []dto.Position {
	if board == nil {
		return nil
	}

	grid := board.Grid
	width := len(grid.Cells[0])  // Use actual slice length if width is not set
	height := len(grid.Cells)

	// BFS for shortest path
	queue := []queueItem{{pos: start, path: []dto.Position{}}}
	visited := make(map[dto.Position]bool)
	visited[start] = true

	// Pre-calculate blocked cells
	blocked := make(map[dto.Position]bool)
	for y, row := range grid.Cells {
		for x, cell := range row {
			if cell.Obstacle {
				blocked[dto.Position{X: x, Y: y}] = true
			}
		}
	}
	for _, entity := range board.Entities {
		if entity.HP > 0 {
			blocked[entity.Position] = true
		}
	}
	
	// Ensure the starting point isn't blocked by the acting unit itself
	delete(blocked, start)
	
	// Ensure the end point isn't blocked (though it usually is if it's an enemy, 
	// but the script should ask for a path to an empty tile).
	// If the user specifically asks for a path to a blocked tile, we return nil.

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.pos == end {
			return current.path
		}

		// Neighbors: Up, Down, Left, Right
		neighbors := []dto.Position{
			{X: current.pos.X, Y: current.pos.Y - 1},
			{X: current.pos.X, Y: current.pos.Y + 1},
			{X: current.pos.X - 1, Y: current.pos.Y},
			{X: current.pos.X + 1, Y: current.pos.Y},
		}

		for _, next := range neighbors {
			if next.X >= 0 && next.X < width && next.Y >= 0 && next.Y < height {
				if !visited[next] && !blocked[next] {
					visited[next] = true
					newPath := append([]dto.Position{}, current.path...)
					newPath = append(newPath, next)
					queue = append(queue, queueItem{pos: next, path: newPath})
				}
			}
		}
	}

	return nil
}
