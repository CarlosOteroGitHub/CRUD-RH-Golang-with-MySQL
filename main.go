package main

//Especificación de las dependencias de Golang.
import (
	"log"
	"net/http"
	"text/template"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

//Variable que ejecuta cada una de las plantillas de visualización y diseño de la carpeta "plantillas"
var plantilla = template.Must(template.ParseGlob("plantillas/*"))

//Función que ejecuta la conexión a la base de datos MySQL.
func conexionDB()(conexion *sql.DB){
	Driver:="mysql"
	Usuario:="root"
	Contrasena:=""
	Nombre:="empleados_angular_db"

	conexion, err := sql.Open(Driver, Usuario + ":" + Contrasena + "@tcp(127.0.0.1)/" + Nombre)
	if err != nil{
		panic(err.Error())
	}

	return conexion
}

//Función para ejecutar la plantilla de visualización y diseño llamada "inicio"
func Inicio(w http.ResponseWriter, r *http.Request){
	plantilla.ExecuteTemplate(w, "inicio", nil)
}

//Estructura para recibir datos de los formularios para la inserción y actualización de datos.
type Empleado struct {
	Id int
	Nombre string
	Nacimiento string
	Correo string
	Descripcion string
}

//Función para ejecutar la plantilla de visualización y diseño llamada "empleados" Tambien ejecuta una consulta SQL para listar todos los registros del modelo "Empleados"
func Empleados(w http.ResponseWriter, r *http.Request){
	conexion := conexionDB()
	mostrarDatos, err := conexion.Query("select id, nombre, nacimiento, correo, descripcion from empleados;")

	if err != nil{
		panic(err.Error())
	} else {
		empleado := Empleado{}
		arregloEmpleados := []Empleado{}

		for mostrarDatos.Next(){
			var id int
			var nombre, nacimiento, correo, descripcion string
			err = mostrarDatos.Scan(&id, &nombre, &nacimiento, &correo, &descripcion)

			if err != nil{
				panic(err.Error())
			} else {
				empleado.Id = id
				empleado.Nombre = nombre
				empleado.Nacimiento = nacimiento
				empleado.Correo = correo
				empleado.Descripcion = descripcion

				arregloEmpleados=append(arregloEmpleados, empleado)
			}
		}

		plantilla.ExecuteTemplate(w, "empleados", arregloEmpleados)
	}
}

//Función para ejecutar la plantilla de visualización y diseño llamada "crear"
func Crear(w http.ResponseWriter, r *http.Request){
	plantilla.ExecuteTemplate(w, "crear", nil)
}

//Función para ejecutar una consulta SQL para permita insertar un registro al modelo "Empleados"
func Insertar(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		nombre := r.FormValue("nombre")
		nacimiento := r.FormValue("nacimiento")
		correo := r.FormValue("correo")
		descripcion := r.FormValue("descripcion")

		conexion := conexionDB()
		insertarDatos, err := conexion.Prepare("insert into empleados(nombre, nacimiento, correo, descripcion) values(?, ?, ?, ?);")

		if err != nil{
			panic(err.Error())
		} else {
			insertarDatos.Exec(nombre, nacimiento, correo, descripcion)
			http.Redirect(w, r, "/empleados", 301)
		}
	}
}

//Función para ejecutar una consulta SQL para permita borrar un registro al modelo "Empleados"
func Borrar(w http.ResponseWriter, r *http.Request){
	idEmpleado := r.URL.Query().Get("id")

	conexion := conexionDB()
	borrarDatos, err := conexion.Prepare("delete from empleados where id = ?;")

	if err != nil{
		panic(err.Error())
	} else {
		borrarDatos.Exec(idEmpleado)
		http.Redirect(w, r, "/empleados", 301)
	}
}

//Función para ejecutar la plantilla de visualización y diseño llamada "editar"
func Editar(w http.ResponseWriter, r *http.Request){
	idEmpleado := r.URL.Query().Get("id")

	conexion := conexionDB()
	mostrarDatos, err := conexion.Query("select id, nombre, nacimiento, correo, descripcion from empleados where id = ?;", idEmpleado)

	if err != nil{
		panic(err.Error())
	} else {
		empleado := Empleado{}

		for mostrarDatos.Next(){
			var id int
			var nombre, nacimiento, correo, descripcion string
			err = mostrarDatos.Scan(&id, &nombre, &nacimiento, &correo, &descripcion)

			if err != nil{
				panic(err.Error())
			} else {
				empleado.Id = id
				empleado.Nombre = nombre
				empleado.Nacimiento = nacimiento
				empleado.Correo = correo
				empleado.Descripcion = descripcion
			}
		}

		plantilla.ExecuteTemplate(w, "editar", empleado)
	}
}

//Función para ejecutar una consulta SQL para permita actualizar un registro al modelo "Empleados"
func Actualizar(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		id_empleado := r.FormValue("id_empleado")
		nombre := r.FormValue("nombre")
		nacimiento := r.FormValue("nacimiento")
		correo := r.FormValue("correo")
		descripcion := r.FormValue("descripcion")

		conexion := conexionDB()
		actualizarDatos, err := conexion.Prepare("update empleados set nombre=?, nacimiento=?, correo=?, descripcion=? where id = ?;")

		if err != nil{
			panic(err.Error())
		} else {
			actualizarDatos.Exec(nombre, nacimiento, correo, descripcion, id_empleado)
			http.Redirect(w, r, "/empleados", 301)
		}
	}
}

//Función principal en donde se gestionan las rutas en el servidor, asi como tambien la especificación del puerto de conexión.
func main(){
	http.HandleFunc("/inicio", Inicio)
	http.HandleFunc("/empleados", Empleados)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor corriendo...")

	http.ListenAndServe(":8080", nil)
}