// Package service
package service

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"vis_contas/config"
	"vis_contas/internal/model"
)

const layout = "2/1/2006"

// CSVtoSQL converte csv para o banco SQL
func CSVtoSQL(csvPath string) error {

	// Ler o arquivo CSV
	faturas, err := ReadFaturasCSV(csvPath)
	if err != nil {
		return fmt.Errorf("erro ao ler CSV: %w", err)
	}

	// Obter conexão com o banco
	db := config.GetDB()

	// Inserir dados em batch para melhor performance
	batchSize := 100
	for i := 0; i < len(faturas); i += batchSize {
		end := i + batchSize
		if end > len(faturas) {
			end = len(faturas)
		}

		batch := faturas[i:end]
		
		// Usar CreateInBatches para inserção eficiente
		if err := db.CreateInBatches(batch, batchSize).Error; err != nil {
			return fmt.Errorf("erro ao inserir batch %d-%d: %w", i, end-1, err)
		}

		log.Printf("Inseridas %d faturas (batch %d-%d)", len(batch), i, end-1)
	}

	log.Printf("Sucesso! Total de %d faturas inseridas no banco de dados", len(faturas))
	return nil
}

// ReadFaturasCSV obtém os dados do csv retorna em memória
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
			Empresa:     	row[11],
		}
		faturas = append(faturas, f)
	}

	return faturas, nil
}

// FilterFaturas filtra faturas no banco de dados de forma dinâmica
func FilterFaturas(categoria, situacao, dataStr, empresa string) ([]model.Faturas, error) {
	db := config.GetDB()
	
	// Começar com query base
	query := db.Model(&model.Faturas{})
	
	// Construir filtros dinamicamente
	conditions := []string{}
	args := []any{}
	
	// Filtro de Categoria
	if categoria != "" {
		conditions = append(conditions, "LOWER(categoria) = LOWER(?)")
		args = append(args, categoria)
	}
	
	// Filtro de Empresa
	if empresa != "" {
		conditions = append(conditions, "LOWER(empresa) = LOWER(?)")
		args = append(args, empresa)
	}
	
	// Filtro de Situação
	if situacao != "" {
		conditions = append(conditions, "LOWER(situacao) = LOWER(?)")
		args = append(args, situacao)
	}
	
	// Filtro de Data
	if dataStr != "" {
		filtroData, err := time.Parse("2/1/2006", dataStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao parsear data '%s': %w", dataStr, err)
		}
		
		// Filtrar por data específica (ignorando hora)
		startOfDay := time.Date(filtroData.Year(), filtroData.Month(), filtroData.Day(), 0, 0, 0, 0, filtroData.Location())
		endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)
		
		conditions = append(conditions, "vencimento >= ? AND vencimento <= ?")
		args = append(args, startOfDay, endOfDay)
	}
	
	// Aplicar todos os filtros
	if len(conditions) > 0 {
		whereClause := strings.Join(conditions, " AND ")
		query = query.Where(whereClause, args...)
	}
	
	// Executar query
	var filtradas []model.Faturas
	result := query.Find(&filtradas)
	
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao filtrar faturas: %w", result.Error)
	}
	
	return filtradas, nil
}

// AtualizarSituacao alterna a situação da fatura entre Pago e Pendente.
func AtualizarSituacao(situacao, idStr string) []model.Faturas {

	faturas, err := ReadFaturasCSV("data/csv/faturas.csv")
	if err != nil {
		fmt.Println("Erro ao Atualizar:", err)
		return nil
	}
	// converter o id da string para inteiro
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		fmt.Println("ID não compatível")	
	}

	// percorrer a lista e alterar a fatura correspondente
	for i, f := range faturas {
		if f.FaturaID == uint(id) {
			if situacao == "Pago" {
				faturas[i].Situacao = "Pendente"
			} else {
				faturas[i].Situacao = "Pago"
			}
			break
		}
	}
	return faturas
}
