from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    bot_token: str
    backend_base_url: str = "http://localhost:8080"
    request_timeout: float = 10.0
    support_url: str = "https://t.me/your_support"

    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8")

