package main

import (
	"ecohortapp/repository"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Realizamos una función que retornará un contenedir de Fyne con el contenido de la pestaña de registros
func (app *Config) registresPestanya() *fyne.Container {
	//Invocamos la función anterior para cargar la estructura de dtos con la interfaz de slice de slices
	app.Registres = app.getRegistresSlice()
	//También invocamos el método getRegistresTable() y lo asignamos al atributo RegistresTable del struct Config
	app.RegistresTable = app.getRegistresTable()
	//Creamos un contenedor con una caja vertical y donde situamos el widget que hemos generado de la tabla Registres
	registresContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		//Definimos un contenedor que permitirá realizar cuadrículas adaptativas y que indicaremos con dos parámetros: el nombre de las filas/columnas y el objetos que situaremos
		container.NewAdaptiveGrid(1, app.RegistresTable),
	)
	return registresContainer
}

// Realizar una función que retorna un contenedor de Fyne con el contenido de la pestaña de registros
func (app *Config) getRegistresTable() *widget.Table {
	//Definimos la estructura del widget para crear una nueva tabla con fyne
	t := widget.NewTable(
		func() (int, int) {
			return len(app.Registres), len(app.Registres[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Col == (len(app.Registres[0])-1) && i.Row != 0 {
				//Ultima cel.la - situa un botò
				w := widget.NewButtonWithIcon("Borrar", theme.DeleteIcon(), func() {
					//Presentem un dialeg de confirmació
					dialog.ShowConfirm("Borrar?", "", func(deleted bool) {
						if deleted {
							id, _ := strconv.Atoi(app.Registres[i.Row][0].(string)) //Transformem el identificador a decimal sencer
							err := app.DB.BorrarRegistre(int64(id))                 //Invoquem el metode per borrar a partir d'un id
							//Capturem possibles errors
							if err != nil {
								app.ErrorLog.Println(err)
							}
						}
						//Forcem el refresc de la taula
						app.actualitzarRegistresTable()
					}, app.MainWindow)
				})
				//Creem un widget d'alta importancia per mostrar un missatge destacat
				w.Importance = widget.HighImportance

				//Definim el contenidor a on situarem el objecte corresponent a el boto.
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					w,
				}
			} else {
				//situarem la informació rebuda en el slice, recordem que primer gestiona la fila i després la columna
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(app.Registres[i.Row][i.Col].(string)),
				}
			}
		})

	//Establim el ample de les diferents celdes
	colWidths := []float32{50, 100, 100, 100, 100, 100, 110}
	//Executem una estructura for per aplicar cada un de els amples amb el metode SetColumnWidth
	for i := 0; i < len(colWidths); i++ {
		t.SetColumnWidth(i, colWidths[i])
	}

	return t
}

// Realizar una función adicional que retorne el puntero en forma de tabla y donde se situan los datos
func (app *Config) getRegistresSlice() [][]interface{} {
	var slice [][]interface{}
	//Invocamos el método inferior registresActuals()
	registres, err := app.registresActuals()
	if err != nil {
		//Capturamos el posible error en el log de errores
		app.ErrorLog.Println("Error al leer los registros: ", err)
	}
	//Realizamos un append para incluir los registros obtenidos en forma de filas
	//y definimos las etiquetas de cada columna para la fila inicial
	slice = append(slice, []interface{}{"ID", "Data", "Precipitació", "Temp. Màxima", "Temp. Mínima", "Humedad", "Opcions"})

	//Ejecutamos un bucle for para recorrer los registros y añadirlos al slice
	for _, x := range registres {
		//Creamos una interfaz vacía para la fila actual
		var filaActual []interface{}

		//Vamos añadiendo a la fila actual cada uno de los valores que correspondan a cada columna definida en el inicio
		filaActual = append(filaActual, strconv.FormatInt(x.ID, 10))
		//Transformamos el valor numérico a String en base 10
		filaActual = append(filaActual, x.Data.Format("2006-01-02"))
		//Formateamos la fecha sl standard americano
		filaActual = append(filaActual, fmt.Sprintf("%d%%", x.Precipitacio))
		//Formateamos la salida a un valor decimal entero
		filaActual = append(filaActual, fmt.Sprintf("%dºC", x.TempMaxima))
		//Formateamos la salida a un valor decimal entero
		filaActual = append(filaActual, fmt.Sprintf("%dºC", x.TempMinima))
		//Formateamos la salida a un valor decimal entero
		filaActual = append(filaActual, fmt.Sprintf("%d%%", x.Humitat))
		//Formateamos la salida a un valor decimal entero
		filaActual = append(filaActual, widget.NewButton("Borrar", func() {})) //Definimos el botón para eliminar y que invocará una función que ya definiremos

		//Añadimos la fila actual al slice de filas
		slice = append(slice, filaActual)
	}

	return slice

}

// Realizar una función para obtener todos los registros con un slice de nuestro repositorio en la BD
func (app *Config) registresActuals() ([]repository.Registres, error) {
	registres, err := app.DB.LlegirTotsRegistres()
	if err != nil {
		//Capturamos el posible error en el log de errores
		app.ErrorLog.Println("Error al leer los registros: ", err)
		return nil, err
	}
	return registres, nil
}
