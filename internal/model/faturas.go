// internal/model/faturas.go
// Define a entidade Fatura e seu mapeamento para a tabela financeiro.faturas.

package model

import "time"

// Fatura representa um registro financeiro de pagamento ou recebimento.
// Contém informações como valor, vencimento, categoria e situação.
type Fatura struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;column:id"`
	Vencimento   time.Time `gorm:"column:vencimento;not null"`
	Valor        float64   `gorm:"column:valor;not null"`
	NParcelas    string    `gorm:"column:n_parcelas"`
	Parcela      string    `gorm:"column:parcela"`
	Destinatario string    `gorm:"column:destinatario;type:varchar(100);not null"`
	Categoria    string    `gorm:"column:categoria;type:varchar(50)"`
	Situacao     string    `gorm:"column:situacao;type:varchar(20);default:'Pendente'"`
	TipoTransf   string    `gorm:"column:tipo_transf;type:varchar(50)"`
	NotaFiscal   string    `gorm:"column:nota_fiscal;type:varchar(50)"`
	Boleto       string    `gorm:"column:boleto;type:varchar(50)"`
	Empresa      string    `gorm:"column:empresa;type:varchar(100)"`
}

// TableName define o nome completo da tabela no banco de dados.
func (Fatura) TableName() string {
	return "financeiro.faturas"
}
