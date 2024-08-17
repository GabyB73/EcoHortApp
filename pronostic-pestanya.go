package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Situar la imagen en el contenedor
func (app *Config) pronosticPestana() *fyne.Container {

	grafic := app.obtenirGrafic()
	contenidorGrafic := container.NewVBox(grafic)
	app.PronosticGraficContainer = contenidorGrafic
	return contenidorGrafic

}

// Obtener el gráfico y lo dimensionaremos
func (app *Config) obtenirGrafic() *canvas.Image {
	var img *canvas.Image
	url := "https://my.meteoblue.com/images/meteogram?temperature_units=C&windspeed_units=kmh&precipitation_units=mm&darkmode=false&iso2=es&lat=41.5168&lon=1.901&asl=111&tz=Europe%2FMadrid&dpi=72&apikey=jhMJTOUVRNvs25m4&lang=es&location_name=Abrera&sig=1040e23cda11792402d11328aa15eb36"
	err := app.descarregaArxiu(url, "pronostic.png")

	if err != nil {
		//No se puede obtener la imagen
		img = canvas.NewImageFromResource(resourceNodisponiblePng)
	} else {
		//Crear la imagen
		img = canvas.NewImageFromFile("pronostic.png")
	}

	img.SetMinSize(fyne.Size{
		Width:  770,
		Height: 480,
	})

	return img

}

// Descargar la imagen
func (app *Config) descarregaArxiu(url string, nomArxiu string) error {
	res, err := app.HTTPClient.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("Error al descargar l'arxiu") //aqui se usa la libreia errors
	}
	binari, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//recoger el flujo de bytes y decodificar la imagen
	img, _, err := image.Decode(bytes.NewReader(binari)) // aqui el guion bajo es pa ignorar el string
	if err != nil {
		return err
	}

	//crear el archivo a donde lo guardaremos en un futuro
	arxiu, err := os.Create(fmt.Sprintf("./%s", nomArxiu))
	if err != nil {
		return err
	}

	//condificación a PNG
	err = png.Encode(arxiu, img)
	if err != nil {
		return err
	}
	return nil

}
