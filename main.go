package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Конфигурация бота
type Config struct {
	TelegramToken   string `json:"telegram_token"`
	OpenRouterKey   string `json:"openrouter_key"`
	OpenRouterModel string `json:"openrouter_model"`
}

// Структуры для Telegram API
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type SendMessageRequest struct {
	ChatID           int    `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageID int    `json:"reply_to_message_id,omitempty"`
}

// Структуры для OpenRouter API
type OpenRouterRequest struct {
	Model    string               `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
}

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

var (
	config       Config
	lastUpdateID = 0
)

func main() {
	// Загрузка конфигурации
	loadConfig()

	// Проверка обязательных параметров
	if config.TelegramToken == "" || config.OpenRouterKey == "" {
		log.Fatal("Требуются telegram_token и openrouter_key в config.json")
	}

	if config.OpenRouterModel == "" {
		config.OpenRouterModel = "gryphe/mythomax-l2-13b" // Модель по умолчанию
	}

	// URL для Telegram API
	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/", config.TelegramToken)

	log.Println("Бот запущен. Ожидание сообщений...")

	// Основной цикл бота
	for {
		updates, err := getUpdates(telegramAPI, lastUpdateID)
		if err != nil {
			log.Println("Ошибка получения обновлений:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			lastUpdateID = update.UpdateID + 1

			// Игнорируем не текстовые сообщения
			if update.Message.Text == "" {
				continue
			}

			log.Printf("Новое сообщение от @%s: %s", update.Message.From.Username, update.Message.Text)

			// Отправляем действие "печатает"
			sendChatAction(telegramAPI, update.Message.Chat.ID, "typing")

			// Получаем ответ от OpenRouter
			aiResponse, err := getOpenRouterResponse(update.Message.Text)
			if err != nil {
				log.Println("Ошибка получения ответа от ИИ:", err)
				sendTextMessage(telegramAPI, update.Message.Chat.ID, "⚠️ Произошла ошибка при обработке вашего запроса. Пожалуйста, попробуйте позже.")
				continue
			}

			// Отправляем ответ пользователю
			err = sendTextMessage(telegramAPI, update.Message.Chat.ID, aiResponse)
			if err != nil {
				log.Println("Ошибка отправки сообщения:", err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func loadConfig() {
	// Попробуем прочитать из config.json
	file, err := os.Open("config.json")
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			log.Println("Ошибка чтения config.json:", err)
		}
	}

	// Переопределяем переменными окружения (если есть)
	if token := os.Getenv("TELEGRAM_TOKEN"); token != "" {
		config.TelegramToken = token
	}
	if key := os.Getenv("OPENROUTER_KEY"); key != "" {
		config.OpenRouterKey = key
	}
	if model := os.Getenv("OPENROUTER_MODEL"); model != "" {
		config.OpenRouterModel = model
	}
}

func getUpdates(apiURL string, offset int) ([]Update, error) {
	url := fmt.Sprintf("%sgetUpdates?timeout=30&offset=%d", apiURL, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP статус %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if !response.OK {
		return nil, fmt.Errorf("ответ Telegram не OK")
	}

	return response.Result, nil
}

func sendTextMessage(apiURL string, chatID int, text string) error {
	message := SendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%ssendMessage", apiURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка отправки сообщения: %s", string(body))
	}

	return nil
}

func sendChatAction(apiURL string, chatID int, action string) error {
	url := fmt.Sprintf("%ssendChatAction?chat_id=%d&action=%s", apiURL, chatID, action)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка отправки действия чата")
	}

	return nil
}

func getOpenRouterResponse(prompt string) (string, error) {
	requestData := OpenRouterRequest{
		Model: config.OpenRouterModel,
		Messages: []OpenRouterMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("ошибка формирования запроса: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.OpenRouterKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/my-telegram-bot") // Обязательно для OpenRouter
	req.Header.Set("X-Title", "Telegram AI Bot")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error.Message != "" {
			return "", fmt.Errorf("ошибка API: %s", errorResp.Error.Message)
		}
		return "", fmt.Errorf("ошибка API (статус %d): %s", resp.StatusCode, string(body))
	}

	var response OpenRouterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("ошибка разбора ответа: %v", err)
	}

	if len(response.Choices) == 0 {
		if response.Error.Message != "" {
			return "", fmt.Errorf("ошибка в ответе: %s", response.Error.Message)
		}
		return "", fmt.Errorf("пустой ответ от API")
	}

	return response.Choices[0].Message.Content, nil
}