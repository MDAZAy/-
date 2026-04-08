from aiogram import F, Router
from aiogram.types import CallbackQuery, InlineKeyboardButton, InlineKeyboardMarkup, Message

from app.client import BackendClient
from app.handlers.common import ensure_backend_user

router = Router()


@router.message(F.text == "Купить тариф")
async def buy_tariff_handler(message: Message, api: BackendClient) -> None:
    plans = await api.get_plans()
    if not plans:
        await message.answer("Активные тарифы недоступны.")
        return

    keyboard = InlineKeyboardMarkup(
        inline_keyboard=[
            [InlineKeyboardButton(text=f"{plan['name']} — {plan['price']} RUB", callback_data=f"buy_plan:{plan['id']}")]
            for plan in plans
        ]
    )
    await message.answer("Выберите тариф для покупки:", reply_markup=keyboard)


@router.callback_query(F.data.startswith("buy_plan:"))
async def buy_tariff_callback(callback: CallbackQuery, api: BackendClient) -> None:
    plan_id = int(callback.data.split(":")[1])
    user = await ensure_backend_user(api, callback.from_user)
    payment = await api.create_payment(user_id=user["id"], plan_id=plan_id)

    keyboard = InlineKeyboardMarkup(
        inline_keyboard=[
            [InlineKeyboardButton(text="Открыть оплату", url=payment["payment_url"])],
        ]
    )
    await callback.message.answer(
        "Платёж создан. Откройте страницу оплаты и завершите оплату. После подтверждения подписка активируется автоматически.",
        reply_markup=keyboard,
    )
    await callback.answer()
