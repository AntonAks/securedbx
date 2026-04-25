"""Unit tests for auth_confirm Lambda."""

import hashlib
import json
import time
from unittest.mock import MagicMock, patch


def solve_pow(challenge: str, difficulty: int) -> int:
    prefix = "0" * difficulty
    nonce = 0
    while True:
        h = hashlib.sha256(f"{challenge}{nonce}".encode()).hexdigest()
        if h.startswith(prefix):
            return nonce
        nonce += 1


def make_event(body: dict) -> dict:
    return {"headers": {}, "body": json.dumps(body), "requestContext": {}}


class TestAuthConfirm:
    def _mock_table(self, challenge: str, expires_at: int | None = None):
        mock_table = MagicMock()
        mock_table.get_item.return_value = {
            "Item": {
                "pk": f"challenge#{challenge}",
                "challenge": challenge,
                "expires_at": expires_at or int(time.time()) + 300,
            }
        }
        return mock_table

    def test_valid_pow_returns_api_key(self):
        challenge = "abc123testchallenge"
        difficulty = 2
        nonce = solve_pow(challenge, difficulty)

        with patch.dict(
            "os.environ",
            {
                "AUTH_TABLE_NAME": "auth-table",
                "POW_DIFFICULTY": str(difficulty),
            },
        ):
            mock_table = self._mock_table(challenge)
            with patch("boto3.resource") as mock_boto:
                mock_boto.return_value.Table.return_value = mock_table

                from lambdas.auth_confirm.handler import handler

                result = handler(make_event({"challenge": challenge, "nonce": nonce}), None)

        assert result["statusCode"] == 200
        body = json.loads(result["body"])
        assert body["api_key"].startswith("sdbx_cli_")
        assert "expires_at" in body

    def test_invalid_pow_rejected(self):
        challenge = "abc123"
        with patch.dict(
            "os.environ",
            {
                "AUTH_TABLE_NAME": "auth-table",
                "POW_DIFFICULTY": "4",
            },
        ):
            mock_table = self._mock_table(challenge)
            with patch("boto3.resource") as mock_boto:
                mock_boto.return_value.Table.return_value = mock_table

                from lambdas.auth_confirm.handler import handler

                result = handler(make_event({"challenge": challenge, "nonce": 0}), None)

        assert result["statusCode"] == 400
        body = json.loads(result["body"])
        assert "invalid" in body.get("error", "").lower()

    def test_expired_challenge_rejected(self):
        challenge = "expired123"
        with patch.dict(
            "os.environ",
            {
                "AUTH_TABLE_NAME": "auth-table",
                "POW_DIFFICULTY": "2",
            },
        ):
            mock_table = MagicMock()
            mock_table.get_item.return_value = {
                "Item": {"pk": f"challenge#{challenge}", "challenge": challenge, "expires_at": 1}
            }
            with patch("boto3.resource") as mock_boto:
                mock_boto.return_value.Table.return_value = mock_table

                from lambdas.auth_confirm.handler import handler

                result = handler(make_event({"challenge": challenge, "nonce": 0}), None)

        assert result["statusCode"] == 400

    def test_missing_challenge_rejected(self):
        with patch.dict("os.environ", {"AUTH_TABLE_NAME": "auth-table", "POW_DIFFICULTY": "2"}):
            mock_table = MagicMock()
            mock_table.get_item.return_value = {"Item": None}
            with patch("boto3.resource") as mock_boto:
                mock_boto.return_value.Table.return_value = mock_table

                from lambdas.auth_confirm.handler import handler

                result = handler(make_event({"challenge": "notfound", "nonce": 0}), None)

        assert result["statusCode"] == 400
