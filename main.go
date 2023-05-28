package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"test/bitcoin-rate-in-uah/conf"
	"test/bitcoin-rate-in-uah/services/mailgun"
	"test/bitcoin-rate-in-uah/utils"
)

type CoinGeckoResponse struct {
	Bitcoin struct {
		UAH float64 `json:"uah"`
	} `json:"bitcoin"`
}

var mutex sync.Mutex

func main() {
	http.HandleFunc("/rate", getCurrentRateHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/sendEmails", sendEmailsHandler)

	// Make emails.txt file or create new one
	filePath := "emails.txt"
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal("Creating file error:", err)
		}
		defer file.Close()
	}

	var port = conf.GetEnvConst("PORT")
	var serverInfo = fmt.Sprintf("The service is running. Available at %s", port)
	log.Println(serverInfo)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getCurrentRateHandler(w http.ResponseWriter, r *http.Request) {
	// getting current BTC rate
	rate, err := utils.GetBitcoinRate()
	if err != nil {
		http.Error(w, "Failed to get Bitcoin rate", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "The current exchange rate of Bitcoin in UAH: %.2f UAH", rate)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	// email validation
	if !utils.ValidateEmail(email) {
		http.Error(w, "Error: wrong email", http.StatusBadRequest)
		return
	}

	filePath := "emails.txt"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// checking if new email exists or not
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if email == line {
			var error = fmt.Sprintf("Email %s already subscribed", email)
			http.Error(w, error, http.StatusConflict)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	// writing new email to file
	_, err = fmt.Fprintln(file, email)
	if err != nil {
		log.Println("Error writing to file:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Email added: %s", email)
}

func sendEmailsHandler(w http.ResponseWriter, r *http.Request) {
	emails, err := utils.GetEmailList("emails.txt")
	if err != nil {
		log.Fatal("Reding file error:", err)
	}

	// getting current BTC rate
	rate, err := utils.GetBitcoinRate()
	if err != nil {
		log.Fatal("Failed getting rate:", err)
	}

	infoEmail := conf.GetEnvConst("INFO_EMAIL")
	message := fmt.Sprintf("The current exchange rate of Bitcoin in UAH: %.2f UAH", rate)

	mutex.Lock()
	defer mutex.Unlock()

	for _, email := range emails {
		// send email via mailgun
		_, err = mailgun.SendMail(
			email,
			infoEmail,
			"Rate of Bitcoin in UAH",
			message,
		)

		if err != nil {
			var error = fmt.Sprintf("Email sent error to %s, error: %s", email, err.Error())
			http.Error(w, error, http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Email sent: %s", email)
	}
}
