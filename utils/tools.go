package utils

// Função para construir a string de data com base nos parâmetros year e month
func ConstruirData(year, month string) string {
	var dataStr string

	if year != "" && month != "" {
		// Ambos ano e mês especificados: formato YYYY-MM
		dataStr = year + "-" + month
	} else if year != "" {
		// Apenas ano especificado: formato YYYY
		dataStr = year
	} else if month != "" {
		// Apenas mês especificado: formato MM
		dataStr = month
	}

	// Se ambos estiverem vazios, dataStr permanece vazio
	return dataStr
}
