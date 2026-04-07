from aiogram import Router
from aiogram.filters import Command
from aiogram.types import Message

from app.client import BackendClient
from app.handlers.common import ensure_backend_user
from app.keyboards.main import main_menu
from app.texts import HELP_TEXT, WELCOME_TEXT

router = Router()


@router.message(Command("start"))
async def start_handler(message: Message, api: BackendClient) -> None:
    await ensure_backend_user(api, message.from_user)
    await message.answer(WELCOME_TEXT, reply_markup=main_menu())


@router.message(Command("menu"))
async def menu_handler(message: Message) -> None:
    await message.answer("Главное меню:", reply_markup=main_menu())


@router.message(Command("help"))
async def help_handler(message: Message) -> None:
    await message.answer(HELP_TEXT, reply_markup=main_menu())

