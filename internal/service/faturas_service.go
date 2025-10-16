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
func ReadFaturasCSV(path string) ([]model.Fatura, error) {
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

	var faturas []model.Fatura

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

		f := model.Fatura{
			ID:           uint(id),
			Vencimento:   venc,
			Valor:        valor,
			NParcelas:    row[3],
			Parcela:      row[4],
			Destinatario: row[5],
			Categoria:    row[6],
			Situacao:     row[7],
			TipoTransf:   row[8],
			NotaFiscal:   row[9],
			Boleto:       row[10],
			Empresa:      row[11],
		}
		faturas = append(faturas, f)
	}

	return faturas, nil
}

// FilterFaturas filtra faturas no banco de dados de forma dinâmica
func FilterFaturas(categoria, situacao, dataStr, empresa string) ([]model.Fatura, error) {
	db := config.GetDB()

	// Começar com query base
	query := db.Model(&model.Fatura{})

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

	// Filtro de Data - Agora suporta diferentes formatos
	if dataStr != "" {
		if strings.Contains(dataStr, "-") {
			// Formato YYYY-MM (ano e mês específicos)
			parts := strings.Split(dataStr, "-")
			if len(parts) == 2 {
				year, errYear := strconv.Atoi(parts[0])
				month, errMonth := strconv.Atoi(parts[1])

				if errYear != nil || errMonth != nil {
					return nil, fmt.Errorf("erro ao parsear data '%s': formato inválido", dataStr)
				}

				// Primeiro dia do mês às 00:00:00
				startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
				// Último dia do mês às 23:59:59
				endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

				conditions = append(conditions, "vencimento >= ? AND vencimento <= ?")
				args = append(args, startOfMonth, endOfMonth)
			}
		} else if len(dataStr) == 4 {
			// Formato YYYY (ano completo)
			year, err := strconv.Atoi(dataStr)
			if err != nil {
				return nil, fmt.Errorf("erro ao parsear ano '%s': %w", dataStr, err)
			}

			// Primeiro dia do ano às 00:00:00
			startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
			// Último dia do ano às 23:59:59
			endOfYear := startOfYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

			conditions = append(conditions, "vencimento >= ? AND vencimento <= ?")
			args = append(args, startOfYear, endOfYear)
		} else if len(dataStr) == 2 {
			// Formato MM (mês em todos os anos)
			month, err := strconv.Atoi(dataStr)
			if err != nil || month < 1 || month > 12 {
				return nil, fmt.Errorf("erro ao parsear mês '%s': mês inválido", dataStr)
			}

			// Usar EXTRACT para filtrar por mês
			conditions = append(conditions, "EXTRACT(MONTH FROM vencimento) = ?")
			args = append(args, month)
		}
	}

	// Aplicar todos os filtros
	if len(conditions) > 0 {
		whereClause := strings.Join(conditions, " AND ")
		query = query.Where(whereClause, args...)
	}

	// Executar query
	var filtradas []model.Fatura
	result := query.Find(&filtradas)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao filtrar faturas: %w", result.Error)
	}

	return filtradas, nil
}

// AtualizarSituacao alterna a situação da fatura entre Pago e Pendente.
func AtualizarSituacao(idStr string) error {
	db := config.GetDB()

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return fmt.Errorf("ID inválido: %v", idStr)
	}

	var fatura model.Fatura
	if err := db.First(&fatura, id).Error; err != nil {
		return fmt.Errorf("fatura não encontrada: %w", err)
	}

	// Alternar valor
	if fatura.Situacao == "Pago" {
		fatura.Situacao = "Pendente"
	} else {
		fatura.Situacao = "Pago"
	}

	if err := db.Save(&fatura).Error; err != nil {
		return fmt.Errorf("erro ao atualizar situação: %w", err)
	}

	var faturas []model.Fatura
	if err := db.Order("vencimento ASC").Find(&faturas).Error; err != nil {
		return fmt.Errorf("erro ao buscar faturas: %w", err)
	}

	return nil
}
