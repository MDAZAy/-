from aiogram import F, Router
from aiogram.types import Message

from app.client import BackendClient, BackendError
from app.handlers.common import ensure_backend_user

router = Router()


@router.message(F.text == "Моя подписка")
async def subscription_handler(message: Message, api: BackendClient) -> None:
    user = await ensure_backend_user(api, message.from_user)

    try:
        subscription = await api.get_active_subscription(user["id"])
    except BackendError:
        await message.answer("Активной подписки пока нет.")
        return

    await message.answer(
        "\n".join(
            [
                "Моя подписка:",
                f"ID: {subscription['id']}",
                f"План ID: {subscription['plan_id']}",
                f"Статус: {subscription['status']}",
                f"Старт: {subscription['start_at']}",
                f"Окончание: {subscription['end_at']}",
            ]
        )
    )

