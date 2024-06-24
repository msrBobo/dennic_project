package usecase

import (
	"fmt"
	"strings"
	"time"
)

type BaseUseCase struct{}

func (u *BaseUseCase) Error(msg string, err error) error {
	if len(strings.TrimSpace(msg)) != 0 {
		return fmt.Errorf("%v: %w", msg, err)
	}
	return err
}

func (u *BaseUseCase) BeforeRequest(guid *string, createdAt *time.Time, updatedAt *time.Time) {
	_ = guid

	if createdAt != nil {
		*createdAt = time.Now().UTC()
	}

	if updatedAt != nil {
		*updatedAt = time.Now().UTC()
	}
}
