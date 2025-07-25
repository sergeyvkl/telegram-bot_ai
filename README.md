# Telegram Bot with AI

Telegram Bot для общения с OpenRouter AI на чистом Go<br>
Найди в Telegram [@KsvAIBot](https://t.me/KsvAIBot)<br>

![Изображение](/images/bot_ai.png)<br><br>
![Изображение](/images/1.PNG)
![Изображение](/images/2.PNG)<br>
![Изображение](/images/3.PNG)
![Изображение](/images/4.PNG)

Инструкция по настройке и запуску<br>
Создайте файл config.json в той же директории, что и программа:<br>

json<br>
{<br>
  "telegram_token": "ВАШ_TELEGRAM_BOT_TOKEN",<br>
  "openrouter_key": "ВАШ_OPENROUTER_API_KEY",<br>
  "openrouter_model": "openai/gpt-3.5-turbo"<br>
}<br>

Замените ВАШ_TELEGRAM_BOT_TOKEN на токен вашего бота, полученный от @BotFather в Telegram.<br>

Замените ВАШ_OPENROUTER_API_KEY на ваш API-ключ от OpenRouter (можно получить на openrouter.ai).<br>

Вы можете изменить модель по умолчанию в конфигурации (например, на "anthropic/claude-2" или другую поддерживаемую модель).<br>

Как это работает<br>
    Бот получает обновления от Telegram через long polling.<br>
    Когда пользователь отправляет сообщение, бот формирует запрос к OpenRouter API.<br>
    OpenRouter обрабатывает запрос и возвращает ответ ИИ.<br>
    Бот отправляет ответ пользователю в Telegram.<br>

Особенности реализации<br>
    Чистый Go без сторонних библиотек<br>
    Обработка ошибок API<br>
    Конфигурация через JSON-файл<br>
    Поддержка разных моделей ИИ через OpenRouter<br>

