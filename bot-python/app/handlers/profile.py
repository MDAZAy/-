from aiogram import F, Router
from aiogram.types import Message

from app.client import BackendClient, BackendError
from app.handlers.common import ensure_backend_user

router = Router()


@router.message(F.text == "Мой профиль")
async def profile_handler(message: Message, api: BackendClient) -> None:
    user = await ensure_backend_user(api, message.from_user)
    lines = [
        "Профиль:",
        f"ID: {user['id']}",
        f"Telegram ID: {user['telegram_id']}",
        f"Username: @{user['username']}" if user["username"] else "Username: не задан",
        f"Имя: {user['full_name']}",
    ]

    try:
        subscription = await api.get_active_subscription(user["id"])
        lines.extend(
            [
                "",
                "Активная подписка:",
                f"План ID: {subscription['plan_id']}",
                f"Статус: {subscription['status']}",
                f"До: {subscription['end_at']}",
            ]
        )
    except BackendError:
        lines.extend(["", "Активной подписки сейчас нет."])

    await message.answer("\n".join(lines))

