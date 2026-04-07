BACKEND_DIR=backend-go
BOT_DIR=bot-python

.PHONY: backend-build backend-run backend-stop backend-smoke bot-run

backend-build:
	powershell -ExecutionPolicy Bypass -File scripts/run-backend.ps1 -BuildOnly

backend-run:
	powershell -ExecutionPolicy Bypass -File scripts/run-backend.ps1

backend-stop:
	powershell -ExecutionPolicy Bypass -File scripts/stop-backend.ps1

backend-smoke:
	powershell -ExecutionPolicy Bypass -File scripts/smoke-backend.ps1

bot-run:
	powershell -ExecutionPolicy Bypass -File scripts/run-bot.ps1

