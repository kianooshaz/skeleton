package iop

import "github.com/google/uuid"

type OrganizationID uuid.UUID

func (o OrganizationID) String() string {
	return uuid.UUID(o).String()
}
