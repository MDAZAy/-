from aiogram.types import KeyboardButton, ReplyKeyboardMarkup


def main_menu() -> ReplyKeyboardMarkup:
    return ReplyKeyboardMarkup(
        keyboard=[
            [KeyboardButton(text="Мой профиль"), KeyboardButton(text="Тарифы")],
            [KeyboardButton(text="Купить тариф"), KeyboardButton(text="Моя подписка")],
            [KeyboardButton(text="Получить VPN ключ"), KeyboardButton(text="Поддержка")],
        ],
        resize_keyboard=True,
    )

