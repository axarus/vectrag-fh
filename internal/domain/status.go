package domain

import "fmt"

type Status string

const (
	StatusDraft   Status = "draft"
	StatusPublish Status = "publish"
	StatusDelete  Status = "delete"
)

func ValidateStatus(status Status) error {
	if status != StatusDraft && status != StatusPublish && status != StatusDelete {
		return fmt.Errorf("must be either 'draft', 'publish', or 'delete'")
	}
	return nil
}
