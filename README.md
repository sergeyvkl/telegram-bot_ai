Инструкция по настройке и запуску
Создайте файл config.json в той же директории, что и программа:

json
{
  "telegram_token": "ВАШ_TELEGRAM_BOT_TOKEN",
  "openrouter_key": "ВАШ_OPENROUTER_API_KEY",
  "openrouter_model": "openai/gpt-3.5-turbo"
}

Замените ВАШ_TELEGRAM_BOT_TOKEN на токен вашего бота, полученный от @BotFather в Telegram.

Замените ВАШ_OPENROUTER_API_KEY на ваш API-ключ от OpenRouter (можно получить на openrouter.ai).

Вы можете изменить модель по умолчанию в конфигурации (например, на "anthropic/claude-2" или другую поддерживаемую модель).

Как это работает
    Бот получает обновления от Telegram через long polling.
    Когда пользователь отправляет сообщение, бот формирует запрос к OpenRouter API.
    OpenRouter обрабатывает запрос и возвращает ответ ИИ.
    Бот отправляет ответ пользователю в Telegram.

Особенности реализации
    Чистый Go без сторонних библиотек
    Обработка ошибок API
    Конфигурация через JSON-файл
    Поддержка разных моделей ИИ через OpenRouter

