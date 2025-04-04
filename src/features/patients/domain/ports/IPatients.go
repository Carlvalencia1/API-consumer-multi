package ports 

type IPatients interface {
	FindID(id int) (error)
}