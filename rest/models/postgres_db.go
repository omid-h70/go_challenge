package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"sync"
)

var lock = &sync.Mutex{}

type PostgresDBStruct struct {
	conn *pgx.Conn
}

var postgresHandle PostgresDBStruct
var DbString string = "postgres://root:postgres@localhost:5432"
var DBName string = "go_challenge"

func (*PostgresDBStruct) GetInstance() *PostgresDBStruct {
	if postgresHandle.conn == nil {
		lock.Lock()
		defer lock.Unlock()
		postgresHandle.conn, _ = pgx.Connect(context.Background(), DbString)
		if postgresHandle.conn != nil {
			_, err := postgresHandle.conn.Exec(context.Background(), "DROP DATABASE "+DBName)
			if err != nil {
				log.Println(err)
			}
			_, err = postgresHandle.conn.Exec(context.Background(), "CREATE DATABASE "+DBName)
			if err != nil {
				log.Println(err)
			} else {
				postgresHandle.conn.Close(context.Background())
				NewDBString := DbString + "/" + DBName

				postgresHandle.conn, err = pgx.Connect(context.Background(), NewDBString)
				if err != nil {
					log.Println(err)
				}
			}

		}
	} else {
		fmt.Println("Single instance already created. Connection Is Open")
	}
	return &postgresHandle
}

func (*PostgresDBStruct) CreateTables() *PostgresDBStruct {
	err := createInstrumentTable(postgresHandle.conn)
	if err != nil {
		log.Println(err)
	}
	err = createTradeTable(postgresHandle.conn)
	if err != nil {
		log.Println(err)
	}
	return &postgresHandle
}

func createInstrumentTable(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS Instrument(
		instrument_id INT,
		name VARCHAR(255)
	)`)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func createTradeTable(conn *pgx.Conn) error {
	//_, err := conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS Trade(
	//	trade_id INT PRIMARY KEY,
	//	instrument_id INT NOT NULL,
	//	date_en date,
	//	open INT,
	//	high INT,
	//	low INT,
	//	close INT,
	//	CONSTRAINT Trade_FK FOREIGN KEY(instrument_id) REFERENCES Instrument(instrument_id)
	//)`)

	_, err := conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS Trade(
		trade_id INT,
		instrument_id INT,
		date_en date,
		open INT,
		high INT,
		low INT,
		close INT
	)`)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func (*PostgresDBStruct) InsertDBSeeds() *PostgresDBStruct {
	err := insertInstrumentSeeds(postgresHandle.conn)
	if err != nil {
		log.Println(err)
	}

	err = insertTradeSeeds(postgresHandle.conn)
	if err != nil {
		log.Println(err)
	}
	return &postgresHandle
}

func (*PostgresDBStruct) InsertDummySeeds(count int) *PostgresDBStruct {
	var dbScheme string = `INSERT INTO Trade(trade_id, instrument_id, date_en, open, high, low, close)
	SELECT random()*10, 
    random()*10, 
	NOW() - '1 day'::INTERVAL * (RANDOM()::int * 100),
    random()*10, 
    random()*10, 
    random()*10, 
    random()*10 
	FROM generate_series(1,%d) id`

	dbScheme = fmt.Sprintf(dbScheme, count)

	_, err := postgresHandle.conn.Exec(context.Background(), dbScheme)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "QueryRow2 failed: %v\n", err)
		if err != nil {
			log.Println(err)
		}
	}
	return &postgresHandle
}

func insertTradeSeeds(conn *pgx.Conn) error {
	var dbScheme string = `
	INSERT INTO Trade VALUES
	(1, 1, '2020-01-01', 1001, 2001, 301, 401),
	(1, 1, '2020-01-02', 1002, 2002, 302, 402),
	(1, 1, '2020-01-03', 1003, 2003, 303, 403),
	(1, 2, '2020-01-01', 1004, 2004, 304, 404),
	(1, 2, '2020-01-03', 1005, 2005, 305, 405),
	(1, 5, '2020-01-01', 1006, 2006, 306, 406),
	(1, 1, '2021-01-01', 1007, 2007, 307, 407)
	`
	_, err := conn.Exec(context.Background(), dbScheme)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "QueryRow2 failed: %v\n", err)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func insertInstrumentSeeds(conn *pgx.Conn) error {
	var dbScheme string = "INSERT INTO Instrument values(1,'AAPL'), (2,'GOOGL')"
	_, err := conn.Exec(context.Background(), dbScheme)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "QueryRow1 failed: %v\n", err)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func (*PostgresDBStruct) GetDBLog(eTableEnum TableENUMS) string {
	var val string
	var scheme string
	var jsonData string = "["

	if eTableEnum == E_TARDE_TABLE {
		scheme = "select row_to_json(Trade) FROM Trade LIMIT 100"
	} else if eTableEnum == E_INSTROMENT_TABLE {
		scheme = "select row_to_json(Instrument) FROM Instrument LIMIT 100"
	}

	if len(scheme) > 0 {
		rows, err := postgresHandle.conn.Query(context.Background(), scheme)
		defer rows.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		} else {
			for rows.Next() {
				rows.Scan(&val)
				jsonData += val
				jsonData += ","
				fmt.Println(val)
			}
		}
	}
	/*Fix json Data*/
	jsonDatalen := len(jsonData)
	if jsonDatalen > 0 && jsonData[jsonDatalen-1] == ',' {
		jsonData = jsonData[:jsonDatalen-1]
	}

	jsonData += "]"
	return jsonData
}

func (*PostgresDBStruct) GetTradeReport(tradeName string) []TradeReports {
	var stReport []TradeReports
	var scheme string = `
	SELECT to_char(Trade.date_en, 'YYYY/MM/DD'), Instrument.name
	FROM Instrument
	RIGHT JOIN Trade On Instrument.instrument_id = Trade.instrument_id
    WHERE Instrument.name = '%s'
	ORDER BY Trade.date_en DESC LIMIT 1
	`
	scheme = fmt.Sprintf(scheme, tradeName)

	var stTmp TradeReports
	err := postgresHandle.conn.QueryRow(context.Background(), scheme).Scan(&stTmp.Date, &stTmp.Name)
	stReport = append(stReport, stTmp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}

	//rows, err := postgresHandle.conn.Query(context.Background(), scheme)
	//defer rows.Close()

	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	//} else {
	//	for i := 0; rows.Next(); i++ {
	//		var stTmp TradeReports
	//		rows.Scan(&stTmp.Date)
	//		rows.Scan(&stTmp.Name)
	//		stReport = append(stReport, stTmp)
	//	}
	//}
	return stReport
}

func (*PostgresDBStruct) CloseDB() error {
	var err error
	if postgresHandle.conn != nil {
		err = postgresHandle.conn.Close(context.Background())
	}
	return err
}
