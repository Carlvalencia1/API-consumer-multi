package ports 

type IPatients interface {
	FindID(id_usuario int) (error)
}