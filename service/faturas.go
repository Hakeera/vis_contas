// Package service
package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"vis_contas/model"
)

const layout = "2/1/2006"

func ReadFaturasCSV(path string) ([]model.Faturas, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	records, err := reader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("erro ao ler csv: %w", err)
	}

	var faturas []model.Faturas

	for i, row := range records {
		if i == 0 {
			continue
		}

		// parse seguro dos campos
		id, _ := strconv.Atoi(row[0])
		venc, err := time.Parse("2/1/2006", row[1])
		if err != nil {
			fmt.Println("Erro ao parsear data:", row[1], err)
		}
		valor, _ := strconv.ParseFloat(row[2], 64)

		f := model.Faturas{
			FaturaID:   	uint(id),
			Vencimento: 	venc,
			Valor:      	valor,
			NParcelas:  	row[3],
			Parcela:    	row[4],
			Destinatario:   row[5],
			Categoria:  	row[6],
			Situacao:   	row[7],
			TipoTransf: 	row[8],
			NotaFiscal: 	row[9],
			Boleto:     	row[10],
		}
		faturas = append(faturas, f)
	}

	return faturas, nil
}


func FilterFaturas(faturas []model.Faturas, categoria, situacao, dataStr string) []model.Faturas {
    var filtradas []model.Faturas

    var filtroData time.Time
    var err error
    if dataStr != "" {
        filtroData, err = time.Parse("2/1/2006", dataStr)
        if err != nil {
            fmt.Println("Erro ao parsear data do filtro:", err)
            filtroData = time.Time{} // ignora se inválida
        }
    }

    for _, f := range faturas {
        if categoria != "" && !strings.EqualFold(f.Categoria, categoria) {
            continue
        }
        if situacao != "" && !strings.EqualFold(f.Situacao, situacao) {
            continue
        }
        if !filtroData.IsZero() {
            // comparando apenas data, sem hora
            y1, m1, d1 := f.Vencimento.Date()
            y2, m2, d2 := filtroData.Date()
            if y1 != y2 || m1 != m2 || d1 != d2 {
                continue
            }
        }
        filtradas = append(filtradas, f)
    }

    return filtradas
}
