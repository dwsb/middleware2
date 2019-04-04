package utils

import (
	"encoding/csv"
	"os"
)

/*
Etapas -
1 - Primeiro criar o arquivo csv com createFile dando um nome para ele
2 - Carregar o slice data com os dados capturados usando o addData
3 - Descarregar os dados do slice no csv criado usando saveData com o nome do arquivo criado
OBS: o arquivo receber toda carga de dados no final para ficarmos evitando IO escrevendo toda vez nele
quando tivermos uma captura de dado.
*/
var data = [][]string{}

//server para criar o arquivo csv que vai receber os dados capturados
func createFile(dir string) bool {
	f, err := os.Create(dir)
	if err != nil {
		return false
	}
	w := csv.NewWriter(f)
	w.Write([]string{"Protocol", "Time", "TypeMiddleware", "Function"})
	if err != nil {
		return false
	}
	w.Flush()
	f.Close()
	return true
}

//server para adicionar os dados do slice de dados
func addData(protocol string, time string, TypeMiddleware string, function string) {
	data = append(data, []string{protocol, time, TypeMiddleware, function})
}

//server pra descarregar o slice em um .csv que ja foi criado
func saveData(dir string) {
	f, err := os.OpenFile(dir, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	w := csv.NewWriter(f)
	w.WriteAll(data)
	if err != nil {
		return
	}
	w.Flush()
}

/*
func main() {

	createFile("data.csv")

	for i := 1000; i > 0; i-- {
		addData("TCP", "1ms", "RPC", "get")
	}

	saveData("data.csv")

}
*/
