package mongobackend

import (
	"fmt"
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// GetFreeRotation gets the stored free champions rotation
func (b *Backend) GetFreeRotation() (riotclient.FreeRotation, error) {
	return riotclient.FreeRotation{}, fmt.Errorf("Not implemented")
}

// GetFreeRotationTimeStamp gets the timestamp of the stored free champions rotation
func (b *Backend) GetFreeRotationTimeStamp() time.Time {
	return time.Time{}
}

// StoreFreeRotation stores a new free champions rotation list
func (b *Backend) StoreFreeRotation(freeRotation riotclient.FreeRotation) error {
	return fmt.Errorf("Not implemented")
}
