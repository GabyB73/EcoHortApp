package repository

import (
	"database/sql"
	"errors"
	"time"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

// Esta función devolverá el struct poblado con la conexión a la DB
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

//Desarrollar las funciones de la interfaz
/*Indicar en el receptor que haremos servir un puntero sobre el método NewSQLiteRepository
para utilizar la conexión establecida por esta acción*/
func (repo *SQLiteRepository) Migrate() error {
	peticio := `
		create table if not exist registres (
			id integer primary key autoincrement,
			data_registre integer not null,
			precipitacio integer not null,
			temp_maxima integer not null,
			temp_minima integer not null,
			humitat integer not null
		)
		`
	_, err := repo.Conn.Exec(peticio)
	return err

}

func (repo *SQLiteRepository) InsertRegistre(registres Registres) (*Registres, error) {
	//Preparar la instrucción para añadir un registro a la tabla registres
	sentencia := "insert into registres (data_registre, precipitacio, temp_maxima, temp_minima, humitat) values (?, ?, ?, ?, ?)"

	res, err := repo.Conn.Exec(sentencia, registres.Data.Unix, registres.Precipitacio, registres.TempMaxima, registres.TempMinima, registres.Humitat)
	if err != nil {
		return nil, err
	}

	//Añadir la llamada a la función LastInsertId() de la respuesta para obtener la id que se ha generado con la inserción
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	registres.ID = id

	return &registres, nil //Preparar el retornocon un objeto o nil en caso de error
}

func (repo *SQLiteRepository) LlegirTotsRegistres() ([]Registres, error) {
	//Formular la consulta para obtener todos los registros de la tabla registres
	sentencia := "select id, data_registre, precipitacio, temp_maxima, temp_minima, humitat from registres order by data_registre"
	//Ejecutar la consulta
	files, err := repo.Conn.Query(sentencia)
	if err != nil {
		return nil, err

	}
	defer files.Close() //Cerrar la conexión de la BD para optimizar

	var tots []Registres
	//Ejecutamos una estructura for para consultar uno de los resultados y lo incluiremos en el slice
	for files.Next() {
		var fila Registres
		var unixTime int64

		err := files.Scan(
			&fila.ID,
			&unixTime,
			&fila.Precipitacio,
			&fila.TempMaxima,
			&fila.TempMinima,
			&fila.Humitat,
		)
		//Si se produce algún error, lo devolveremos
		if err != nil {
			return nil, err
		}
		fila.Data = time.Unix(unixTime, 0)
		tots = append(tots, fila) //Aplicamos la inclusión del objeto dentro del slice
	}

	return tots, nil //Retornamos el slice o nil en caso de error

}

// Función para obtener datos por ID
func (repo *SQLiteRepository) LlegirRegistrePerID(id int64) (*Registres, error) {
	sentencia := "select id, data_registre, precipitacio, temp_maxima, temp_minima, humitat from registres where id = ?"

	resposta := repo.Conn.QueryRow(sentencia, id)

	var fila Registres
	var unixTime int64
	//Preparamos un struct de tipo Holdings con los datos obtenidos
	err := resposta.Scan(
		&fila.ID,
		&unixTime,
		&fila.Precipitacio,
		&fila.TempMaxima,
		&fila.TempMinima,
		&fila.Humitat,
	)

	if err != nil {
		return nil, err
	}

	fila.Data = time.Unix(unixTime, 0)

	return &fila, nil

}

// Realizar la actualización del registro falta hacer esta función
func (repo *SQLiteRepository) ActualitzarRegistre(id int64, actualitzar Registres) error {
	//Validar si la id es 0 y así controlar posibles errores
	if id == 0 {
		return errors.New("La ID a actualizar es incorrecta")
	}
	//Preparamos la petición para actualizar los datos
	sentencia := "Update registres set data_registre = ?, precipiacio = ?, temp_maxima = ?, temp_minima = ?, humitat = ?"
	res, err := repo.Conn.Exec(sentencia, actualitzar.Data.Unix(), actualitzar.Precipitacio, actualitzar.TempMaxima, actualitzar.TempMinima, actualitzar.Humitat, id)
	//Controlar posibles errores
	if err != nil {
		return err
	}
	//Comprobar el número de registros afectados durante la actualización
	afectats, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if afectats == 0 {
		return updateError
	}
	return nil
}

// Función para borrar un registro
func (repo *SQLiteRepository) BorrarRegistre(id int64) error {
	res, err := repo.Conn.Exec("delete from registres where id = ?", id)
	if err != nil {
		return err
	}

	//Comprobar el número de registros afectados durante la eliminación
	afectats, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if afectats == 0 {
		return deleteError

	}
	return nil
}
