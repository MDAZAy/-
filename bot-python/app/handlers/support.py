from aiogram import F, Router
from aiogram.types import Message

from app.config import Settings

router = Router()


@router.message(F.text == "Поддержка")
async def support_handler(message: Message, settings: Settings) -> None:
    await message.answer(f"Поддержка: {settings.support_url}")

