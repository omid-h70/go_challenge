package models

type TableENUMS int

const (
	E_TARDE_TABLE TableENUMS = iota
	E_INSTROMENT_TABLE
)

type TradeReports struct {
	Name string
	Date string
}

type DBInterface interface {
	ConnectToDB() *DBInterface
	CreateTables() *DBInterface
	InsertDBSeeds() *DBInterface
	InsertDummySeeds(count int) *DBInterface
	GetTradeReport() []TradeReports
	GetDBLog(TableENUMS) string
	CloseDB() error
}
