package usecase

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	_ "modernc.org/sqlite"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppUsecase struct {
	Client *whatsmeow.Client
}

func NewWhattsAppUsecase() *WhatsAppUsecase {
	dbLog := waLog.Stdout("Database", "ERROR", true)

	connectionString := "file:wa_session.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous=NORMAL"
	container, err := sqlstore.New(context.Background(), "sqlite", connectionString, dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	go connectWA(client)

	return &WhatsAppUsecase{Client: client}
}

func connectWA(client *whatsmeow.Client) {
	if client.Store.ID == nil {

		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {

				fmt.Println("\n\n===================================================")
				fmt.Println("SCAN QR CODE INI UNTUK LOGIN WA ADMIN:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("===================================================\n\n")
			} else {
				fmt.Println("Login Event:", evt.Event)
			}
		}
	} else {

		err := client.Connect()
		if err != nil {
			panic(err)
		}
		fmt.Println("✅ WhatsApp Service Terhubung!")
	}
}

func (s *WhatsAppUsecase) SendMessage(ctx context.Context, phone string, message string) error {

	re := regexp.MustCompile(`[^0-9]`)
	phone = re.ReplaceAllString(phone, "")

	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}

	fmt.Printf("DEBUG: Sending to Clean JID: %s@s.whatsapp.net\n", phone)

	if !s.Client.IsConnected() {
		return fmt.Errorf("whatsapp belum terhubung")
	}

	jid, _ := types.ParseJID(phone + "@s.whatsapp.net")

	var err error
	for i := 0; i < 3; i++ {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		_, err = s.Client.SendMessage(ctx, jid, &waE2E.Message{
			Conversation: &message,
		})
		cancel()

		if err == nil {
			return nil
		}

		fmt.Printf("⚠️ Gagal kirim ke %s (Percobaan %d/3): %v. Mencoba lagi...\n", phone, i+1, err)
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("gagal mengirim pesan setelah 3x percobaan: %v", err)
}
