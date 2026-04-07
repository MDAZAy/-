from aiogram import F, Router
from aiogram.types import Message

from app.client import BackendClient, BackendError
from app.handlers.common import ensure_backend_user

router = Router()


@router.message(F.text == "Получить VPN ключ")
async def vpn_handler(message: Message, api: BackendClient) -> None:
    user = await ensure_backend_user(api, message.from_user)

    try:
        key = await api.issue_vpn_key(user["id"])
    except BackendError as exc:
        await message.answer(f"Не удалось выдать VPN-ключ: {exc}")
        return

    await message.answer(
        "\n".join(
            [
                "VPN ключ готов:",
                f"Провайдер: {key['provider']}",
                f"Название: {key['key_name']}",
                f"Ссылка: {key['access_url']}",
                "",
                "Если у пользователя уже был активный ключ, backend вернёт существующий.",
            ]
        )
    )

