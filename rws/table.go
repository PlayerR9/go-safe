package rws

import (
	"iter"
	"sync"

	"github.com/PlayerR9/go-safe/common"
)

// Table is a boundless table. This means that any operation done
// to out-of-bounds cells will not cause any error. The table is
// safe for concurrent use.
type Table[T any] struct {
	// table is the underlying table.
	table [][]T

	// width is the width of the table.
	width int

	// height is the height of the table.
	height int

	// mu is the mutex for the table.
	mu sync.RWMutex
}

// NewTable creates a new Table with a width and height.
//
// Parameters:
//   - width: The width of the table.
//   - height: The height of the table.
//
// Returns:
//   - *Table: The new Table.
//   - error: If the table could not be created.
//
// Errors:
//   - errors.BadParameterError: If width or height is negative.
func NewTable[T any](width, height int) (*Table[T], error) {
	if width < 0 {
		return nil, common.NewErrBadParam("width", "must not non-negative")
	} else if height < 0 {
		return nil, common.NewErrBadParam("height", "must be non-negative")
	}

	table := make([][]T, 0, height)

	for i := 0; i < height; i++ {
		table = append(table, make([]T, width, width))
	}

	return &Table[T]{
		table:  table,
		width:  width,
		height: height,
	}, nil
}

// Height returns the height of the table.
//
// Returns:
//   - int: The height of the table. 0 if the receiver is nil.
func (t *Table[T]) Height() int {
	if t == nil {
		return 0
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.height
}

// Width returns the width of the table.
//
// Returns:
//   - int: The width of the table. 0 if the receiver is nil.
func (t *Table[T]) Width() int {
	if t == nil {
		return 0
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.width
}

// CellAt returns the cell at the specified position.
//
// Parameters:
//   - x: The x position of the cell.
//   - y: The y position of the cell.
//
// Returns:
//   - T: The cell at the specified position. The zero value if the receiver is nil
//     or the position is out of bounds.
func (t *Table[T]) CellAt(x, y int) T {
	zero := *new(T)

	if t == nil || y < 0 || x < 0 {
		return zero
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if y >= t.height || x >= t.width {
		return zero
	}

	return t.table[y][x]
}

// ResizeWidth resizes the width of the table. The width is not
// resized if the receiver is nil or the new width is the same as the
// current width.
//
// Parameters:
//   - new_width: The new width of the table.
//
// Returns:
//   - error: If the table could not be resized.
//
// Errors:
//   - gers.BadParameterError: If new_width is negative.
func (t *Table[T]) ResizeWidth(new_width int) error {
	if t == nil {
		return common.ErrNilReceiver
	} else if new_width < 0 {
		return common.NewErrBadParam("new_width", "must be non-negative")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if new_width == t.width {
		return nil
	}

	if new_width < t.width {
		for i := 0; i < t.height; i++ {
			t.table[i] = t.table[i][:new_width:new_width]
		}
	} else {
		extension := make([]T, new_width-t.width)

		for i := 0; i < t.height; i++ {
			t.table[i] = append(t.table[i], extension...)
		}
	}

	t.width = new_width

	return nil
}

// ResizeHeight resizes the height of the table. The height is not
// resized if the receiver is nil or the new height is the same as the
// current height.
//
// Parameters:
//   - new_height: The new height of the table.
//
// Returns:
//   - error: If the table could not be resized.
//
// Errors:
//   - gers.BadParameterError: If new_height is negative.
func (t *Table[T]) ResizeHeight(new_height int) error {
	if t == nil {
		return common.ErrNilReceiver
	} else if new_height < 0 {
		return common.NewErrBadParam("new_height", "be non-negative")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if new_height == t.height {
		return nil
	}

	if new_height < t.height {
		t.table = t.table[:new_height:new_height]
	} else {
		for i := t.height; i < new_height; i++ {
			t.table = append(t.table, make([]T, t.width))
		}
	}

	t.height = new_height

	return nil
}

// SetCellAt sets the cell at the specified position. The cell is not
// set if the receiver is nil or the position is out of bounds.
//
// Parameters:
//   - cell: The cell to set.
//   - x: The x position of the cell.
//   - y: The y position of the cell.
func (t *Table[T]) SetCellAt(cell T, x, y int) {
	if t == nil || y < 0 || x < 0 {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if y >= t.height || x >= t.width {
		return
	}

	t.table[y][x] = cell
}

// Row returns an iterator over the rows in the table.
//
// Returns:
//   - iter.Seq2[int, []T]: An iterator over the rows in the table. Never returns nil.
func (t *Table[T]) Row() iter.Seq2[int, []T] {
	if t == nil {
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	return func(yield func(int, []T) bool) {
		for i := 0; i < t.height; i++ {
			if !yield(i, t.table[i]) {
				return
			}
		}
	}
}

// Free releases any resources associated with the Table.
func (t *Table[T]) Free() {
	if t == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.table) == 0 {
		t.width = 0
		t.height = 0

		return
	}

	for i := 0; i < t.height; i++ {
		clear(t.table[i])
		t.table[i] = nil
	}

	clear(t.table)
	t.table = nil

	t.width = 0
	t.height = 0
}
