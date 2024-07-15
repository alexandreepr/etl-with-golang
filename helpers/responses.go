package helpers

type ReportResponse struct {
	TotalRows							int64 `json:"total_rows"`
	InvalidCPFCount                     int64 `json:"invalid_cpf_count"`
	InvalidLojaMaisFrequenteCNPJCount   int64 `json:"invalid_loja_mais_frequente_cnpj_count"`
	InvalidLojaUltimaCompraCNPJCount    int64 `json:"invalid_loja_ultima_compra_cnpj_count"`
}