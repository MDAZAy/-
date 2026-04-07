from aiogram import F, Router
from aiogram.types import Message

from app.client import BackendClient

router = Router()


@router.message(F.text == "Тарифы")
async def plans_handler(message: Message, api: BackendClient) -> None:
    plans = await api.get_plans()
    if not plans:
        await message.answer("Активных тарифов пока нет.")
        return

    lines = ["Доступные тарифы:"]
    for plan in plans:
        lines.append(
            f"\n#{plan['id']} {plan['name']}\n"
            f"Цена: {plan['price']} RUB\n"
            f"Срок: {plan['duration_days']} дней\n"
            f"{plan['description']}"
        )

    await message.answer("\n".join(lines))

