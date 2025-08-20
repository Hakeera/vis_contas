// Package utils
package utils

import (
	"fmt"
	"html/template"
	"strings"
)

// TemplateFunctions define as funções de template personalizadas
var TemplateFunctions = template.FuncMap{
    // Função lower - converte string para minúscula
    "lower": strings.ToLower,
    
    // Função upper - converte string para maiúscula
    "upper": strings.ToUpper,
    
    // Função formatMoney - formata valores monetários
    "formatMoney": func(value float64) string {
        return fmt.Sprintf("%.2f", value)
    },
    
    // Função replace - substitui strings
    "replace": func(old, new, s string) string {
        return strings.ReplaceAll(s, old, new)
    },
    
    // Função contains - verifica se contém substring
    "contains": strings.Contains,
}

