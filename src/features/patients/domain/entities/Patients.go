package entities

type Patients struct {
	IDUsuario int          `json:"id_usuario"`
	Nombre    string       `json:"nombre"`
	Apellido  string       `json:"apellido"`
	Edad      int          `json:"edad"`
	Genero    string       `json:"genero"`
	NumeroContacto string  `json:"numero_contacto"`
}