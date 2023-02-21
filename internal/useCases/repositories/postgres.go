package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/MaxKut3/BudgetBot/config"
	"github.com/MaxKut3/BudgetBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx"
)

type Repository interface {
	Insert(msg *models.Message)
	GetTotalAmount(update tgbotapi.Update)
}

type repository struct {
	Conn *pgx.Conn
}

func NewRepository(cfg *config.TgBotConfig) *repository {

	conn, err := pgx.Connect(cfg.Connect)
	if err != nil {
		log.Panic(fmt.Errorf("connection failed: %e", err))
		return nil
	}

	return &repository{
		Conn: conn,
	}
}

func (r *repository) Insert(msg *models.Message) {

	date := time.Now()

	_, err := r.Conn.Exec("INSERT INTO BudgetBotData (chatId, category, amount, cur, amountRub, month, year) VALUES ($1::bigint, $2::varchar(50), $3::integer, $4::varchar(3), $5::integer, $6::integer, $7::integer )", msg.ChatID, msg.Category, msg.Sum, msg.Cur, msg.SumRub, int(date.Month()), date.Year())
	if err != nil {
		log.Println("Данные в базу не записались")
	}
}

func (r *repository) GetTotalAmount(update tgbotapi.Update) int {

	date := time.Now()

	rows, err := r.Conn.Query("SELECT sum(amountRub) as total FROM BudgetBotData WHERE chatId = $1 and month = $2  and  year = $3", update.Message.Chat.ID, int(date.Month()), date.Year())
	if err != nil {
		log.Println("Данные с базы не пришли")
	}

	defer rows.Close()

	var total int

	for rows.Next() {
		fmt.Println(rows)
		errScan := rows.Scan(&total)
		if errScan != nil {
			log.Printf("Строка запроса не распарсилась: %e", errScan)
		}
	}
	return total / 100
}
