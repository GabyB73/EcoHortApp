package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func (app *Config) getClimaText() (*canvas.Text, *canvas.Text, *canvas.Text, *canvas.Text) {

	var parte Diaria                                         // contenedor del struct Diaria
	var precipitacio, tempMax, tempMin, humitat *canvas.Text // contenedor de los valores de precipitación, temperatura máxima, temperatura mínima y humedad

	prediccio, err := parte.GetPrediccions() //ver si tiene que ir prediccions (en catala)
	if err != nil {

		//cuando falla la conexión con la API de AEMET
		gris := color.RGBA{R: 155, G: 155, B: 155, A: 255}
		precipitacio = canvas.NewText("Precipitación: No leido", gris)
		tempMax = canvas.NewText("Temperatura Máxima: No leido", gris)
		tempMin = canvas.NewText("Temperatura Mínima: No leido", gris)
		humitat = canvas.NewText("Humedad: No leido", gris)
	} else {
		colorTexto := color.RGBA{R: 0, G: 180, B: 0, A: 255} //definimos el color verde

		//condicion de la probabilidad de lluvia

		if prediccio.ProbPrecipitacio < 50 {
			colorTexto = color.RGBA{R: 180, G: 0, B: 0, A: 255} //definimos el color azul
		}

		//preparar los textos

		precipitacioText := fmt.Sprintf("Precipitació: %d%%", prediccio.ProbPrecipitacio)
		tempMaxText := fmt.Sprintf("Temperatura Máxima: %dºC", prediccio.TemperaturaMax)
		tempMinText := fmt.Sprintf("Temperatura Mínima: %dºC", prediccio.TemperaturaMin)
		humitatText := fmt.Sprintf("Humedad: %d%%", prediccio.HumitatRelativa)

		precipitacio = canvas.NewText(precipitacioText, colorTexto)
		tempMax = canvas.NewText(tempMaxText, nil)
		tempMin = canvas.NewText(tempMinText, nil)
		humitat = canvas.NewText(humitatText, colorTexto)

	}
	//alinear los textos
	precipitacio.Alignment = fyne.TextAlignLeading
	tempMax.Alignment = fyne.TextAlignCenter
	tempMin.Alignment = fyne.TextAlignCenter
	humitat.Alignment = fyne.TextAlignTrailing

	return precipitacio, tempMax, tempMin, humitat

}
