import asyncio
import logging

from aiogram import Bot, Dispatcher

from app.client import BackendClient
from app.config import Settings
from app.handlers import payments, plans, profile, start, subscriptions, support, vpn


async def main() -> None:
    logging.basicConfig(level=logging.INFO)

    settings = Settings()
    bot = Bot(token=settings.bot_token)
    api = BackendClient(base_url=settings.backend_base_url, timeout=settings.request_timeout)

    dispatcher = Dispatcher()
    dispatcher["api"] = api
    dispatcher["settings"] = settings

    dispatcher.include_router(start.router)
    dispatcher.include_router(profile.router)
    dispatcher.include_router(plans.router)
    dispatcher.include_router(payments.router)
    dispatcher.include_router(subscriptions.router)
    dispatcher.include_router(vpn.router)
    dispatcher.include_router(support.router)

    try:
        await dispatcher.start_polling(bot)
    finally:
        await api.close()
        await bot.session.close()


if __name__ == "__main__":
    asyncio.run(main())
