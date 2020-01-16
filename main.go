package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

type day struct {
	Dia   int    `json:"dia"`
	Clima string `json:"clima"`
}

type rainregister struct {
	dia  int
	area float64
}

type planet struct {
	distance   int
	degrees    int
	angularvel int
}

var ferengi = planet{
	distance:   500,
	degrees:    0,
	angularvel: 1,
}

var betasoide = planet{
	distance:   2000,
	degrees:    0,
	angularvel: -3,
}

var vulcano = planet{
	distance:   1000,
	degrees:    0,
	angularvel: 5,
}

var (
	inputdia int
	//contadores
	contSequia     int
	contLluvia     int
	contPicoLluvia int
	contOptimo     int
	contNormal     int
	registerdays   []day
	rainsPeaks     []rainregister
)

func main() {

	var dia day
	var rains rainregister

	// var registerdays []day
	// var rainsPeaks []rainregister

	for i := 0; i < 3650; i++ {
		inputdia = i
		dia.Dia = i + 1

		posF, posV, posB := ConsultarDia(inputdia)

		if Sequia(posF, posV, posB) == true {
			contSequia++
			dia.Clima = "Sequia"
		} else {

			if DiaOptimo(posF, posV, posB) == true {
				contOptimo++
				dia.Clima = "Óptimo"

			} else {

				rainBool, peakValue := Lluvia(posF, posV, posB)

				if rainBool == true {
					rains.dia = i + 1
					rains.area = peakValue
					contLluvia++
					dia.Clima = "Lluvia"
					rainsPeaks = append(rainsPeaks, rains)

				} else {
					contNormal++
					dia.Clima = "Normal"
				}
			}
		}
		registerdays = append(registerdays, dia)
	}
	//Seteo el Área Máxima recorriendo todo el registro de días
	var maxArea rainregister
	for i := 0; i < len(rainsPeaks); i++ {
		aux := rainsPeaks[i].area
		if maxArea.area < aux {
			maxArea.area = aux
			maxArea.dia = rainsPeaks[i].dia
		}
	}
	//Guardo el array con los días y el area que tuvieron pico de lluvia
	var maxPeaks []rainregister
	for i := 0; i < len(rainsPeaks); i++ {
		if maxArea.area == rainsPeaks[i].area {
			maxPeaks = append(maxPeaks, rainsPeaks[i])
			contPicoLluvia++
		}
	}

	fmt.Println(registerdays)
	fmt.Println("Cantidad de días de Sequía: ", contSequia)
	fmt.Println("Cantidad de días Óptimos: ", contOptimo)
	fmt.Println("Cantidad de días de Lluvia: ", contLluvia)
	fmt.Println("Cantidad de días con Pico de Lluvia: ", contPicoLluvia)
	fmt.Println("El área del Triángulo formado por el máximo pico de Lluvia fue: ", maxArea.area)
	fmt.Println("los picos de lluvia fueron los días: ")
	for i := 0; i < len(maxPeaks); i++ {
		fmt.Print(maxPeaks[i].dia, ", ")
		if i == len(maxPeaks)-1 {
			fmt.Print("\n")
		}
	}
	fmt.Println("Cantidad de días Normales: ", contNormal)
	serverInit()

}

func clima(w http.ResponseWriter, r *http.Request) {

	// name := r.PostFormValue("vulcanDay")
	// name := r.URL.Query()

	values := r.URL.Query().Get("dia")
	params, _ := strconv.Atoi(values)
	params = params - 1

	if params >= 0 && params < 3650 {
		w.Header().Add("Content-Type", "application/json")
		jeison, _ := json.Marshal(registerdays[params])
		w.Write(jeison)
	} else {
		fmt.Fprintln(w, "Ingrese un dia correcto.")
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	// tmpl := template.Must(template.ParseFiles("index.html"))
	// var data string
	// tmpl.Execute(w, data)
	http.ServeFile(w, r, "index.html")
}

func serverInit() {
	http.HandleFunc("/clima", clima)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

//ConsultarDia Recibe el día solicitado y Devuelve las posiciones de los planetas
func ConsultarDia(consultedDay int) (posFerengi, posVulcano, posBetasoide int) {
	var posF, posV, posB int

	var auxF int
	auxF = consultedDay

	posF = -auxF
	auxF = posF % 360 //modulo de los grados del circulo
	posF = auxF
	posF = auxF + 360
	if consultedDay == 0 || auxF == 0 {
		posF = posF - 360
	}
	// fmt.Println("posicion del planeta F: ", posF)

	var auxV int
	auxV = consultedDay

	posV = auxV * 5
	auxV = posV % 360 //modulo de los grados del circulo
	posV = auxV
	// fmt.Println("posicion del planeta V: ", posV)

	var auxB int
	auxB = consultedDay

	posB = auxB * -3
	auxB = posB % 360 //modulo de los grados del circulo
	posB = auxB + 360
	if consultedDay == 0 || auxB == 0 {
		posB = posB - 360
	}

	// fmt.Println("posicion del planeta B: ", posB)
	return posF, posV, posB
}

//ReverseAng Devuelve el ángulo inverso
func ReverseAng(grade int) int {
	if grade < 180 {
		return grade + 180
	}
	if grade > 180 {
		return grade - 180
	}
	return 0
}

//Sequia Devuelve un boolean, confirmando si los planetas estan alineados con el origen (Sol)
func Sequia(posF, posV, posB int) bool {

	if posF == posV && posF == posB {
		return true
	}
	if posF == posV && posF == ReverseAng(posB) {
		return true
	}
	if posF == posB && posF == ReverseAng(posV) {
		return true
	}
	if posV == posB && posV == ReverseAng(posF) {
		return true
	}

	return false
}

//DiaOptimo Recibe los grados de los planetas, y con ellos calcula las coordenadas (X,Y) , hace 2 pendientes y si tienen la misma inclinación implica que están alineados
func DiaOptimo(fer, vul, bet int) bool {

	coordinateesFerengiX, coordinateesFerengiY := cartesiano(ferengi, fer)
	coordinateesVulcanoX, coordinateesVulcanoY := cartesiano(vulcano, vul)
	coordinateesBetasoideX, coordinateesBetasoideY := cartesiano(betasoide, bet)

	// (y2 - y1) / (x2 - x1) == (y3 - y1) / (x3 - x1)

	m1 := (coordinateesVulcanoY - coordinateesFerengiY) / (coordinateesVulcanoX - coordinateesFerengiX)
	m2 := (coordinateesBetasoideY - coordinateesFerengiY) / (coordinateesBetasoideX - coordinateesFerengiX)

	m1 = math.Round(m1*10) / 10
	m2 = math.Round(m2*10) / 10
	// m1 := math.Round((coordinateesVulcanoY-coordinateesFerengiY)/(coordinateesVulcanoX-coordinateesFerengiX)*10) / 10
	// m2 := math.Round((coordinateesBetasoideY-coordinateesFerengiY)/(coordinateesBetasoideX-coordinateesFerengiX)*10) / 10

	// math.Abs((coordinateesVulcanoY-coordinateesFerengiY)/(coordinateesVulcanoX-coordinateesFerengiX)) == math.Abs((coordinateesBetasoideY-coordinateesFerengiY)/(coordinateesBetasoideX-coordinateesFerengiX))

	// if math.Abs((coordinateesVulcanoY-coordinateesFerengiY)/(coordinateesVulcanoX-coordinateesFerengiX)) == math.Abs((coordinateesBetasoideY-coordinateesFerengiY)/(coordinateesBetasoideX-coordinateesFerengiX)) {
	// 	return true
	// }
	num1 := math.Abs(m1)
	num2 := math.Abs(m2)
	if num1 == num2 {
		return true
	}
	return false
}

//Lluvia Recibe los grados de los planetas, y con ellos calcula el área del triangulo formado por los 3 puntos y las áreas interiores
func Lluvia(fer, vul, bet int) (bool, float64) {

	coordFerengiX, coordFerengiY := cartesiano(ferengi, fer)
	coordVulcanoX, coordVulcanoY := cartesiano(vulcano, vul)
	coordBetasoideX, coordBetasoideY := cartesiano(betasoide, bet)
	// Se calcula el área total de los 3 puntos obtenidos
	totalArea := TriangleArea(coordFerengiX, coordFerengiY, coordVulcanoX, coordVulcanoY, coordBetasoideX, coordBetasoideY)
	// Se reemplaza en cada planeta por el origen de coordenadas para obtener los triangulos interiores
	planet12Area := TriangleArea(coordFerengiX, coordFerengiY, coordVulcanoX, coordVulcanoY, 0.0, 0.0)
	planet13Area := TriangleArea(coordFerengiX, coordFerengiY, 0.0, 0.0, coordBetasoideX, coordBetasoideY)
	planet23Area := TriangleArea(0.0, 0.0, coordVulcanoX, coordVulcanoY, coordBetasoideX, coordBetasoideY)
	sumaTotal := planet12Area + planet13Area + planet23Area

	if sumaTotal <= totalArea {
		return true, sumaTotal
	}
	return false, sumaTotal
}

//TriangleArea Calcula el área de un triangulo recibiendo 3 puntos
func TriangleArea(x1, y1, x2, y2, x3, y3 float64) float64 {
	totalArea := (((x1*y2 + y1*x3 + x2*y3) - (x3*y2 + x1*y3 + y1*x2)) / 2)
	totalArea = math.Abs(totalArea) //Las áreas nunca pueden ser negativas
	return totalArea
}

//cartesiano Obtiene las coordenadas x,y usando coordenadas polares, recibe el planeta y el grado donde se encuentra
func cartesiano(p planet, grad int) (float64, float64) {

	rad := (float64(grad) * math.Pi) / 180                   //convierte el ángulo a radianes
	xcord := math.Trunc(float64(p.distance) * math.Cos(rad)) // X = (distancia del planeta al origen) * coseno del ángulo
	ycord := math.Trunc(float64(p.distance) * math.Sin(rad)) // Y = (distancia del planeta al origen) * seno del ángulo
	// xcord := float64(p.distance) * math.Cos(rad) // X = (distancia del planeta al origen) * coseno del ángulo
	// ycord := float64(p.distance) * math.Sin(rad) // Y = (distancia del planeta al origen) * seno del ángulo
	return xcord, ycord
}

/*
	1. ¿Cuántos períodos de sequía habrá? Habra uno por cada vez que esten alineados los planetas
	3. ¿Cuántos períodos de condiciones óptimas de presión y temperatura habrá? Habra uno por cada vez que esten
	alineados los planetas sin estar alineados con el sol
*/
/*
	2. ¿Cuántos períodos de lluvia habrá y qué día será el pico máximo de lluvia?

*/
