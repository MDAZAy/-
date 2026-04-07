from __future__ import annotations

from typing import Any

import httpx


class BackendError(RuntimeError):
    pass


class BackendClient:
    def __init__(self, base_url: str, timeout: float) -> None:
        self._client = httpx.AsyncClient(base_url=base_url.rstrip("/"), timeout=timeout)

    async def close(self) -> None:
        await self._client.aclose()

    async def health(self) -> dict[str, Any]:
        return await self._request("GET", "/health")

    async def ensure_user(self, telegram_id: int, username: str | None, full_name: str) -> dict[str, Any]:
        return await self._request(
            "POST",
            "/api/v1/users/ensure",
            json={
                "telegram_id": telegram_id,
                "username": username or "",
                "full_name": full_name,
            },
        )

    async def get_plans(self) -> list[dict[str, Any]]:
        return await self._request("GET", "/api/v1/plans")

    async def create_payment(self, user_id: int, plan_id: int) -> dict[str, Any]:
        return await self._request(
            "POST",
            "/api/v1/payments/create",
            json={"user_id": user_id, "plan_id": plan_id},
        )

    async def get_active_subscription(self, user_id: int) -> dict[str, Any]:
        return await self._request("GET", f"/api/v1/subscriptions/active/{user_id}")

    async def issue_vpn_key(self, user_id: int) -> dict[str, Any]:
        return await self._request("POST", "/api/v1/vpn/issue", json={"user_id": user_id})

    async def _request(self, method: str, path: str, **kwargs: Any) -> Any:
        response = await self._client.request(method, path, **kwargs)
        if response.is_success:
            return response.json()

        message = f"backend error: {response.status_code}"
        try:
            payload = response.json()
            message = payload.get("error", message)
        except ValueError:
            pass
        raise BackendError(message)

