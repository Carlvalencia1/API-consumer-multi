package adapters

import (
    "apiconsumer/src/core"
    "apiconsumer/src/features/cases/domain/entities"
    "database/sql"
    "errors"
)

type MYSQL struct {
    conn *sql.DB
}

func NewMysql() *MYSQL {
    conn := core.NewMysql()
    return &MYSQL{
        conn: conn,
    }
}

// Cambiado de FindById a FindID para cumplir con la interfaz ICase
func (mysql *MYSQL) FindID(idExpediente int) (*entities.MedicalCase, error) {
    query, err := mysql.conn.Prepare("SELECT id_expediente, id_usuario, temperatura, peso, estatura, ritmo_cardiaco, fecha_registro FROM cases WHERE id_expediente = ?")
    if err != nil {
        return nil, err
    }
    defer query.Close()

    var medicalCase entities.MedicalCase

    err = query.QueryRow(idExpediente).Scan(
        &medicalCase.IDExpediente,
        &medicalCase.IDUsuario,
        &medicalCase.Temperatura,
        &medicalCase.Peso,
        &medicalCase.Estatura,
        &medicalCase.RitmoCardiaco,
        &medicalCase.FechaRegistro,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("no medical case found with the given ID")
        }
        return nil, err
    }

    return &medicalCase, nil
}