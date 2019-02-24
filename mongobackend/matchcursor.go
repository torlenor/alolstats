package mongobackend

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// MatchCursor is the cursor implementation for Match queries
type MatchCursor struct {
	cur mongo.Cursor
	ctx context.Context
}

// Next gets the next result from the cursor.
// Returns true if there were no errors and there is a next result.
func (c *MatchCursor) Next() bool {
	return c.cur.Next(c.ctx)
}

// Decode decodes the current result
func (c *MatchCursor) Decode(i interface{}) error {
	return c.cur.Decode(i)
}

// Close the cursor
func (c *MatchCursor) Close() error {
	return c.cur.Close(c.ctx)
}
