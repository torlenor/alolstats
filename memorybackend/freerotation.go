package memorybackend

import (
	"time"

	"github.com/torlenor/alolstats/riotclient"
)

// GetFreeRotation gets the stored free champions rotation
func (b *Backend) GetFreeRotation() (*riotclient.FreeRotation, error) {
	b.log.Debugln("Getting Free Rotation List from storage")

	b.mutex.Lock()
	defer b.mutex.Unlock()

	freeRotation := b.freeRotation

	return &freeRotation, nil
}

// GetFreeRotationTimeStamp gets the timestamp of the stored free champions rotation
func (b *Backend) GetFreeRotationTimeStamp() time.Time {
	b.log.Debugln("Getting Free Rotation List TimeStamp from storage")

	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.freeRotation.Timestamp
}

// StoreFreeRotation stores a new free champions rotation list
func (b *Backend) StoreFreeRotation(freeRotation *riotclient.FreeRotation) error {
	b.log.Debugln("Storing new Free Rotation List in storage")

	b.mutex.Lock()
	b.freeRotation = *freeRotation
	b.mutex.Unlock()

	return nil
}
