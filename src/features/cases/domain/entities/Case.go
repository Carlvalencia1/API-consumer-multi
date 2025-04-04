package entities


import  "time"

type MedicalCase struct {
	IDExpediente  int       `json:"id_expediente"`
	IDUsuario     int       `json:"id_usuario"`
	Temperatura   float64   `json:"temperatura"`
	Peso          float64   `json:"peso"`
	Estatura      float64   `json:"estatura"`
	RitmoCardiaco int       `json:"ritmo_cardiaco"`
	FechaRegistro time.Time `json:"fecha_registro"`
}