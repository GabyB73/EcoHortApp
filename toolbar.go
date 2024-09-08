package main

import (
	"ecohortapp/repository"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) getToolBar(win fyne.Window) *widget.Toolbar { //recordar que el * es un puntero
	toolBar := widget.NewToolbar(
		widget.NewToolbarSpacer(), //creamos un espaciador que empujará los elementos a la derecha
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addRegistresDialog()
		}), //Creamos una nueva acción indicando qué ícono y qué función se ejecutará al hacer clic
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			app.actualitzarClimaDadesContent()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			//Realizar la llamada a la función que contiene el diálogo para modificar las preferencias de la aplicación
			w := app.mostrarPreferencies()
			//Definir el tamaño de esta ventana
			w.Resize(fyne.NewSize(300, 200))
			//Mostrar la ventana
			w.Show()
		}),
	)
	return toolBar

}

func (app *Config) mostrarPreferencies() dialog.Dialog {
	//win := app.App.NewWindow("Ajustes")

	//Definir las variables donde guardaremos el resultado del método de entrada
	dadaMunicipi := widget.NewEntry()
	dadaApiKey := widget.NewEntry()
	dadaMunicipi.Text = municipi
	dadaApiKey.Text = apiKey
	esStringValidador := func(s string) error {
		_, err := fmt.Print(reflect.TypeOf(s))
		if err != nil {
			return err
		}
		return nil
	}
	dadaMunicipi.Validator = esStringValidador
	dadaApiKey.Validator = esStringValidador

	fmt.Print(dadaMunicipi.Text)
	//Crear el diálogo creando un formulario
	addForm := dialog.NewForm(
		"Configurar ajustes",
		"Guardar",
		"Cancelar",
		//Añadir las etiquetas en forma de item para el formulario
		[]*widget.FormItem{
			{Text: "Código Municipio", Widget: dadaMunicipi},
			{Text: "API Key", Widget: dadaApiKey},
		},
		//A continuación rea lizamos la validación de los datos
		func(valid bool) {
			if valid {
				municipi = dadaMunicipi.Text
				apiKey = dadaApiKey.Text
				fmt.Print(dadaMunicipi.SelectedText())
				//Desarrollar un filtrado, convertir los datos al formato de la bd
				app.App.Preferences().SetString("municipi", dadaMunicipi.Text)
				app.App.Preferences().SetString("apiKey", dadaApiKey.Text)

				//Invocar el parámetro del struc Config para permitir que refresque el widget de la tabla con el nuevo registro
				app.actualitzarClimaDadesContent()

			}
		},
		app.MainWindow)

	return addForm
}

// Función para añadir Resitres donde referenciamos el struct Config
func (app *Config) addRegistresDialog() dialog.Dialog {
	//Definir las variables donde guardaremos el resultado del método de entrada
	dataRegistreEntrada := widget.NewEntry()
	precipitacioEntrada := widget.NewEntry()
	tempMaximaEntrada := widget.NewEntry()
	tempMinimaEntrada := widget.NewEntry()
	humitatEntrada := widget.NewEntry()

	app.AfegirRegistresDataRegistreEntrada = dataRegistreEntrada
	app.AfegirRegistresPrecipitacioEntrada = precipitacioEntrada
	app.AfegirRegistresTempMaximaEntrada = tempMaximaEntrada
	app.AfegirRegistresTempMinimaEntrada = tempMinimaEntrada
	app.AfegirRegistresHumitatEntrada = humitatEntrada

	validacioData := func(s string) error {
		//Aplicamos el formato de la fecha para que sea correcto
		//Se debe refenciar el formato con el layout standard 2006-01-02
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return err //Si hay un error, devolverá el error
		}
		return nil //Si no hay error, devolverá nil
	}
	dataRegistreEntrada.Validator = validacioData

	esIntValidador := func(s string) error {
		_, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		return nil
	}

	precipitacioEntrada.Validator = esIntValidador
	tempMaximaEntrada.Validator = esIntValidador
	tempMinimaEntrada.Validator = esIntValidador
	humitatEntrada.Validator = esIntValidador

	//Definir un placeholder y facilitar la usabilidad en el campo data
	dataRegistreEntrada.PlaceHolder = "yyyy-mm-dd"

	//Crear el diálogo creando un formulario
	addForm := dialog.NewForm(
		"Añadir Registres",
		"Añadir",
		"Cancelar",
		//Añadir las etiquetas en forma de item para el formulario
		[]*widget.FormItem{
			{Text: "Data Registre", Widget: dataRegistreEntrada},
			{Text: "Probabilitat de Precipitació", Widget: precipitacioEntrada},
			{Text: "Temperatura Màxima", Widget: tempMaximaEntrada},
			{Text: "Temperatura Mínima", Widget: tempMinimaEntrada},
			{Text: "Humedad", Widget: humitatEntrada},
		},
		//A continuación realizamos la validación de los datos
		func(valid bool) {
			if valid {
				//Desarrollar un filtrado, convertir los datos al formato de la bd
				//Se debe referenciar el formato de la fecha con el layout standard 2006-01-02
				dataRegistre, _ := time.Parse("2006-01-02", dataRegistreEntrada.Text)
				precipitacio, _ := strconv.Atoi(precipitacioEntrada.Text)
				tempMaxima, _ := strconv.Atoi(tempMaximaEntrada.Text)
				tempMinima, _ := strconv.Atoi(tempMinimaEntrada.Text)
				humitat, _ := strconv.Atoi(humitatEntrada.Text)

				//Invocar el método de la bd para insertar registros y que poblaremos con los datos formateados
				_, err := app.DB.InsertRegistre(repository.Registres{
					Data:         dataRegistre,
					Precipitacio: precipitacio,
					TempMaxima:   tempMaxima,
					TempMinima:   tempMinima,
					Humitat:      humitat,
				})
				//Capturar el posible error en el log de errores
				if err != nil {
					app.ErrorLog.Println("Error al insertar el registro: ", err)
				}
				//Invocar el parámetro del struc Config para permitir que refresque el widget de la tabla con el nuevo registro
				app.actualitzarRegistresTable()

			}
		},
		app.MainWindow)

	//Establecemos el tamaño de la ventana y mostramos el diálogo
	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()

	return addForm

}
