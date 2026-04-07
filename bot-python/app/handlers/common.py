from aiogram.types import User

from app.client import BackendClient


async def ensure_backend_user(api: BackendClient, tg_user: User) -> dict:
    full_name = " ".join(part for part in [tg_user.first_name, tg_user.last_name] if part).strip()
    return await api.ensure_user(
        telegram_id=tg_user.id,
        username=tg_user.username,
        full_name=full_name or tg_user.username or str(tg_user.id),
    )

