package main

import (
	"database/sql"
	"ecohortapp/repository"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	_ "github.com/glebarez/go-sqlite" //importar el paquete de la BBDD
)

type Config struct {
	App                                fyne.App    //El atributo a donde alamcenaremos la configuración base de la app
	InfoLog                            *log.Logger //El atributo a donde almacenaremos el log de información
	ErrorLog                           *log.Logger //El atributo a donde almacenaremos el log de errores
	MainWindow                         fyne.Window
	ClimaDadesContainer                *fyne.Container
	HTTPClient                         http.Client
	PronosticGraficContainer           *fyne.Container
	DB                                 repository.Repository //Definir la referencia a la conexión de la BBDD SQLite
	Registres                          [][]interface{}       //Para almacenar el slice de slices en forma de interfaz donde están contenidos los datos obtenidos de la BD
	RegistresTable                     *widget.Table         //Para almacenar la referencia al puntero que corresponde al widget de la tabla
	AfegirRegistresDataRegistreEntrada *widget.Entry         //Añadir la referencia a la entrada del valor data registre para nuevos registros que guardemos en la bd
	AfegirRegistresPrecipitacioEntrada *widget.Entry         //Añadir la referencia a la entrada del valor precipitacio para nuevos registros que guardemos en la bd
	AfegirRegistresTempMaximaEntrada   *widget.Entry         //Añadir la referencia a la entrada del valor tempMaxima para nuevos registros que guardemos en la bd
	AfegirRegistresTempMinimaEntrada   *widget.Entry         //Añadir la referencia a la entrada del valor tempMinima para nuevos registros que guardemos en la bd
	AfegirRegistresHumitatEntrada      *widget.Entry         //Añadir la referencia a la entrada del valor humitat para nuevos registros que guardemos en la bd
	municipi                           string                //Añadimos la referencia a este valor de configuración para que se guarde en la base de datos
	apiKey                             string                //Añadimos la referencia a este valor de configuración para que se guarde en la base de datos
}

var myApp Config

func main() {

	//Crear una App con Fyne
	fyneApp := app.NewWithID("cat.cibernarium.ecohortapp")
	myApp.App = fyneApp

	//Desarrollar los logs
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)        //Log normal
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Lshortfile) //Log de errores

	//Conexión con la BBDD
	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err) //Recordar que Panic es equivalente a Print pero acompañado de una llamada a Panic
	}

	//Crear un repositorio de BBDD
	myApp.setupDB(sqlDB)

	//Definir la capacidad de que el usuario modifique el municipio y la apiKey
	municipi = fyneApp.Preferences().StringWithFallback("municipi", "08001")
	apiKey = fyneApp.Preferences().StringWithFallback("apiKey", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJvZGlnaW9jaW9AZ21haWwuY29tIiwianRpIjoiYjRlZTViMjctZDhhMS00YmIxLWFiZjgtYmFjYTViOTc5ZDhjIiwiaXNzIjoiQUVNRVQiLCJpYXQiOjE2NzU2MTY3OTIsInVzZXJJZCI6ImI0ZWU1YjI3LWQ4YTEtNGJiMS1hYmY4LWJhY2E1Yjk3OWQ4YyIsInJvbGUiOiIifQ.y-WKC8DkAJ4O__aNkvWS60AwmYl6dVHcBZKcowfmNKs")

	//Configurar el tamaño y características de la pantalla principal de fyne
	myApp.MainWindow = fyneApp.NewWindow("Eco Hort App")
	myApp.MainWindow.Resize(fyne.NewSize(800, 500)) //el tamaño de la nueva ventana
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster() //definiendo como pantalla princial, si cerramos esta ventana la aplicación se cerrará

	myApp.makeUI() //Crear una invocación a una función externa que creará la interfaz gráfica a partir del contenido

	//Mostrar y ejecutar la App
	myApp.MainWindow.ShowAndRun()

}

// Realizar la función para invocar la conexión a la BBDD
func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	if os.Getenv("DB_PATH") != "" {
		///Recuperar la variable de entorno
		path = os.Getenv("DB_PATH")
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		app.InfoLog.Println("La base de datos se guardará en: ", path)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil

}

// Crear un repositorio para la BBDD
func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqlDB)

	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Println("Error al migrar la base de datos", err)
		log.Panic(err)
	}
}
