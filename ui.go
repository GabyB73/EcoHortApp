package main

import (
	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func (app *Config) makeUI() {

	//conexión a la AEMET
	precipitacio, tempMax, tempMin, humitat := app.getClimaText()

	//Insertar los datos en el contenedor

	contenidorClimaDades := container.NewGridWithColumns(4,
		precipitacio,
		tempMax,
		tempMin,
		humitat,
	)

	//Incluir el contenedor a la ventana
	app.ClimaDadesContainer = contenidorClimaDades

	//Obtener la barra de herramientas
	toolBar := app.getToolBar(app.MainWindow)

	pronosticContenidor := app.pronosticPestana()

	//Las pestañas
	pestanes := container.NewAppTabs(
		container.NewTabItemWithIcon("Pronóstic",
			theme.HomeIcon(), pronosticContenidor(),
			container.NewTabItemWithIcon("Diari Meteorologic"),
			theme.InfoIcon(), canvas.NewText("jhkjhkhjkjkj", nil)),
	)

	pestanes.SetTabLocation(container.SetTabLocationTop)

	contenidorFinal := container.NewVBox(contenidorClimaDades, toolBar, pestanes)

	//Incluir el contenedor a la ventana
	app.MainWindow.SetContent(contenidorFinal)

	go func() {
		for range time.Tick(time.Second * 30) {
			app.actualitzarClimaDadesContent()
		}

	}()

}

func (app *Config) actualitzarClimaDadesContent() {
	app.InfoLog.Println("Cargando la info del clima")
	precipitacio, tempMax, tempMin, humitat := getClimaText()
	app.ClimaDadesContainer.Objects = []fyne.CanvasObject{precipitacio, tempMax, tempMin, humitat}
	app.ClimaDadesContainer.Refresh()

	grafic := app.obtenirGrafic()
	app.PronosticGraficContainer.Objects = []fyne.CanvasObject{grafic}
	app.PronosticGraficContainer.Refresh()

}
