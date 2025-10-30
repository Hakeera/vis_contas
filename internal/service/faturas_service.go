// Package service
package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"vis_contas/config"
	"vis_contas/internal/model"
)

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
