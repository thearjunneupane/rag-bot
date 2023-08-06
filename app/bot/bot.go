package bot

import (
	"errors"
	"fmt"
	"strings"

	"database/sql"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"github.com/thearjnep/rag-bot/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	dgBot *discordgo.Session
	db    *sql.DB
)

func Initialize() {
	var err error
	db, err = sql.Open("postgres", "postgresql://username:password@localhost:port/databasename?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	dgBot, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dgBot.AddHandler(messageCreate)

	dgBot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = dgBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running...!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		caser := cases.Title(language.English)
		arg := caser.String(strings.TrimSpace(strings.TrimPrefix(m.Content, "rag")))

		if strings.HasPrefix(strings.ToLower(arg), "-add") {
			// Split the message content using the ":" delimiter
			parts := strings.SplitN(strings.ToLower(arg), ":", 2)
			if len(parts) != 2 {
				s.ChannelMessageSend(m.ChannelID, "Invalid command format. Use `-add:newSlang`.")
				return
			}

			newSlang := strings.TrimSpace(parts[1])
			err := addSlang(newSlang)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Failed to add the slang: "+err.Error())
			} else {
				s.ChannelMessageSend(m.ChannelID, "Slang added: "+newSlang)
			}
			return
		}

		switch arg {
		case "Help":
			s.ChannelMessageSend(m.ChannelID, "If you know ashmin then write `rag ashmin`")
			break
		default:
			if containsAllCharacters(arg, "arjun") {
				s.ChannelMessageSend(m.ChannelID, "Arjun is the god, you got it म द फ क !")
				return
			}
			randomRag, err := genRandomRag()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Failed to fetch a random slang: "+err.Error())
				return
			}
			s.ChannelMessageSend(m.ChannelID, arg+" is a "+randomRag+".")
			break
		}
	}
}

func genRandomRag() (string, error) {
	var word string
	err := db.QueryRow("SELECT word FROM slangs ORDER BY RANDOM() LIMIT 1").Scan(&word)
	if err != nil {
		return "", err
	}
	return word, nil
}

func containsAllCharacters(input, target string) bool {
	lowerInput := strings.ToLower(input)
	for _, char := range target {
		if !strings.ContainsRune(lowerInput, char) {
			return false
		}
	}
	return true
}

func addSlang(newSlang string) error {
	// Check if the slang already exists in the database
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM slangs WHERE word = $1", newSlang).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("slang already exists")
	}

	// If the slang doesn't exist, insert it into the database
	_, err = db.Exec("INSERT INTO slangs (word) VALUES ($1)", newSlang)
	return err

}
