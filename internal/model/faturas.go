// Package model
package model

import "time"

type Faturas struct {
	FaturaID 	uint		`csv:"faturaid"` 	
	Vencimento	time.Time 	`csv:"vencimento"`	
	Valor		float64 	`csv:"valor"`
	NParcelas	string		`csv:"nParcelas"`	
	Parcela		string		`csv:"parcela"`		
	Destinatario 	string		`csv:"destinatario"`		
	Categoria	string		`csv:"categoria"`	
	Situacao	string		`csv:"situacao"`	
	TipoTransf	string		`csv:"tipoTransf"`	
	NotaFiscal	string		`csv:"notaFiscal"`	
	Boleto		string		`csv:"boleto"`		
	Empresa 	string		`csv:"empresa"`		
}
