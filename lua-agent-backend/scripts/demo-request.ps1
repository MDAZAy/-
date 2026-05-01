$body = @{
  prompt = "Напиши Lua-скрипт, который суммирует массив чисел и печатает результат"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://127.0.0.1:8080/generate" -Method Post -ContentType "application/json" -Body $body
