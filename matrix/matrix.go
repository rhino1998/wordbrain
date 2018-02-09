package matrix

//Position describes a position in a 2d matrix with the bottom left as 0,0
type Position struct {
	X, Y int
}

//Add returns the sum of two positions
func (p Position) Add(op Position) Position {
	return Position{X: p.X + op.X, Y: p.Y + op.Y}
}

//InMatrix returns true if the position is valid for the given matrix
func (p Position) InMatrix(m Matrix) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < len(m) && p.Y < len(m[p.X])
}

//Matrix describes a matrix of strings
type Matrix [][]Space

//New creates a new matrix of dimension x,y
func New(x, y int) {
	m := make(Matrix, x)
	for i, c := range m {
		m[i] = make([]Space, len(c))
	}
}

func (m Matrix) copy() Matrix {
	nm := make(Matrix, len(m))
	for i, c := range m {
		nm[i] = append(make([]Space, 0, len(c)), c...)
	}
	return nm
}

//Remove removes positions from a copy of m and returns it
func (m Matrix) Remove(ps ...Position) Matrix {
	nm := m.copy()
	for _, p := range ps {
		if p.InMatrix(m) {
			nm[p.X][p.Y] = EmptySpace
		}
	}

	for i := 0; i < len(nm[0]); i++ {
		for y := range nm {
			for x, space := range nm[y] {
				if space != EmptySpace {
					continue
				}
				if y+1 < len(nm[x]) {
					nm[y][x] = nm[y+1][x]
					nm[y+1][x] = EmptySpace
				} else {
					nm[y][x] = EmptySpace
				}
			}
		}
	}
	return nm
}

//ValidMoves lists the valid moves from a position
func (m Matrix) ValidMoves(p Position) []Position {
	moves := make([]Position, 0)

	addIfInMatrix := func(p Position) {
		if p.InMatrix(m) {
			moves = append(moves, p)
		}
	}
	addIfInMatrix(p.Add(Position{X: 1, Y: 0}))
	addIfInMatrix(p.Add(Position{X: 1, Y: 1}))
	addIfInMatrix(p.Add(Position{X: 0, Y: 1}))
	addIfInMatrix(p.Add(Position{X: -1, Y: 1}))
	addIfInMatrix(p.Add(Position{X: -1, Y: 0}))
	addIfInMatrix(p.Add(Position{X: -1, Y: -1}))
	addIfInMatrix(p.Add(Position{X: 0, Y: -1}))
	addIfInMatrix(p.Add(Position{X: 1, Y: -1}))

	return moves
}

func (m Matrix) String() string {
	out := ""
	for x := range m {
		out += "| "
		for _, char := range m[len(m)-x-1] {
			out += string(char) + " "
		}
		out += "|\n"
	}
	return out
}
