package schema

type Active string

const (
	DataStatusActive Active = "ACTIVE"
	DataStatusInActive Active = "INACTIVE"
)

type Table struct {
	id string
	createdTimeStamp string
	updatedTimeStamp string
	dataActiveStatus Active
}