package engine

import (
	"fmt"
	"html/template"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Rakhulsr/foodcourt/internal/model"
	"github.com/gin-gonic/gin"
)

func SetupViewEngine(r *gin.Engine) {
	r.Static("/static", "./public")

	r.StaticFS("/uploads", gin.Dir("./public/uploads", false))

	r.SetFuncMap(template.FuncMap{

		"formatRupiah": func(amount int) string {
			s := fmt.Sprintf("%d", amount)
			if len(s) < 4 {
				return "Rp " + s
			}
			var res []byte
			count := 0
			for i := len(s) - 1; i >= 0; i-- {
				if count == 3 {
					res = append([]byte{'.'}, res...)
					count = 0
				}
				res = append([]byte{s[i]}, res...)
				count++
			}
			return "Rp " + string(res)
		},

		"formatDate": func(t time.Time) string {
			return t.Format("02 Jan 2006, 15:04")
		},

		"add": func(a, b int) int {
			return a + b
		},

		"statusColor": func(status string) string {
			switch strings.ToLower(status) {
			case "paid", "ready", "completed", "active":
				return "bg-green-100 text-green-800 border border-green-200"
			case "pending":
				return "bg-yellow-100 text-yellow-800 border border-yellow-200"
			case "confirmed", "preparing":
				return "bg-blue-100 text-blue-800 border border-blue-200"
			case "cancelled", "expired", "inactive":
				return "bg-red-100 text-red-800 border border-red-200"
			default:
				return "bg-gray-100 text-gray-800"
			}
		},

		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},

		"groupByBooth": func(order model.Order) []map[string]interface{} {
			grouped := make(map[uint]map[string]interface{})

			for _, item := range order.Items {
				boothID := item.BoothID

				if _, exists := grouped[boothID]; !exists {

					isNotified := false
					for _, log := range order.Logs {

						if log.BoothID != nil && *log.BoothID == boothID {
							isNotified = true
							break
						}
					}

					grouped[boothID] = map[string]interface{}{
						"Booth":      item.Booth,
						"Items":      []model.OrderItem{},
						"IsNotified": isNotified,
					}
				}

				currentItems := grouped[boothID]["Items"].([]model.OrderItem)
				grouped[boothID]["Items"] = append(currentItems, item)
			}

			var result []map[string]interface{}
			for _, v := range grouped {
				result = append(result, v)
			}
			return result
		},

		"buildGroupWALink": func(order model.Order, boothData interface{}, itemsData interface{}) template.URL {
			booth := boothData.(model.Booth)
			items := itemsData.([]model.OrderItem)

			re := regexp.MustCompile(`[^0-9]`)
			cleanNumber := re.ReplaceAllString(booth.WhatsApp, "")
			if strings.HasPrefix(cleanNumber, "0") {
				cleanNumber = "62" + cleanNumber[1:]
			}

			paymentStatus := "BELUM LUNAS ❌"
			if order.PaymentStatus == "paid" {
				paymentStatus = "SUDAH LUNAS ✅"
			}

			msgBuilder := strings.Builder{}
			fmt.Fprintf(&msgBuilder, "*PESANAN BARU!* 🔔\n")
			fmt.Fprintf(&msgBuilder, "Kepada: *%s*\n\n", booth.Name)
			fmt.Fprintf(&msgBuilder, "No. Order: *%s*\n", order.OrderCode)
			fmt.Fprintf(&msgBuilder, "Meja: *%s*\n", order.TableNumber)
			fmt.Fprintf(&msgBuilder, "Pemesan: *%s*\n", order.CustomerName)
			fmt.Fprintf(&msgBuilder, "Status: *%s*\n", paymentStatus)
			fmt.Fprintf(&msgBuilder, "--------------------------------\n")
			fmt.Fprintf(&msgBuilder, "🍽️ *DAFTAR MENU:*\n")

			for _, item := range items {
				notes := ""
				if item.Notes != "" {
					notes = fmt.Sprintf(" _(Catatan: %s)_", item.Notes)
				}
				fmt.Fprintf(&msgBuilder, "▪️ %dx *%s*%s\n", item.Quantity, item.Menu.Name, notes)
			}

			fmt.Fprintf(&msgBuilder, "--------------------------------\n\n")
			fmt.Fprintf(&msgBuilder, "Mohon diproses. Terima kasih!")

			fullLink := fmt.Sprintf("https://wa.me/%s?text=%s", cleanNumber, url.QueryEscape(msgBuilder.String()))
			return template.URL(fullLink)
		},
	})

	r.LoadHTMLGlob("views/**/*")
}
