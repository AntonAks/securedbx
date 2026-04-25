"""Unit tests for auth_init Lambda."""

import json
from unittest.mock import MagicMock, patch


def make_event(headers=None):
    return {
        "headers": headers or {"x-origin-verify": "test-secret"},
        "body": "{}",
        "requestContext": {"identity": {"sourceIp": "1.2.3.4"}},
    }


class TestAuthInit:
    def test_returns_challenge_and_difficulty(self):
        with patch.dict(
            "os.environ",
            {
                "AUTH_TABLE_NAME": "auth-table",
                "POW_DIFFICULTY": "4",
                "CLOUDFRONT_SECRET": "test-secret",
            },
        ):
            mock_table = MagicMock()
            with patch("boto3.resource") as mock_boto:
                mock_boto.return_value.Table.return_value = mock_table

                from lambdas.auth_init.handler import handler

                result = handler(make_event(), None)

        assert result["statusCode"] == 200
        body = json.loads(result["body"])
        assert "challenge" in body
        assert body["difficulty"] == 4
        assert len(body["challenge"]) == 32  # 16 bytes hex

    def test_challenge_stored_in_dynamo(self):
        with patch.dict(
            "os.environ",
            {
                "AUTH_TABLE_NAME": "auth-table",
                "POW_DIFFICULTY": "4",
                "CLOUDFRONT_SECRET": "test-secret",
            },
        ):
            mock_table = MagicMock()
            with patch("boto3.resource") as mock_boto:
                mock_boto.return_value.Table.return_value = mock_table

                from lambdas.auth_init.handler import handler

                handler(make_event(), None)

            mock_table.put_item.assert_called_once()
            call_args = mock_table.put_item.call_args[1]["Item"]
            assert "pk" in call_args
            assert "challenge" in call_args
            assert "expires_at" in call_args
