"""Lambda function: Verify PoW and issue CLI API key."""

import hashlib
import json
import logging
import os
import secrets
import time
from typing import Any

import boto3

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

KEY_TTL = 86400  # 24 hours


def _error(message: str, status: int = 400) -> dict:
    return {
        "statusCode": status,
        "headers": {"Content-Type": "application/json"},
        "body": json.dumps({"error": message}),
    }


def handler(event: dict[str, Any], context: Any) -> dict[str, Any]:
    """
    Verify PoW solution and issue a 24-hour CLI API key.

    Request: { challenge: "hex", nonce: 12345 }
    Response: { api_key: "sdbx_cli_...", expires_at: unix_ts }
    """
    table_name = os.environ.get("AUTH_TABLE_NAME")
    difficulty = int(os.environ.get("POW_DIFFICULTY", "4"))
    prefix = "0" * difficulty

    try:
        body = json.loads(event.get("body", "{}"))
        challenge = body.get("challenge", "")
        nonce = body.get("nonce")
        if not challenge or nonce is None:
            return _error("challenge and nonce are required")
    except (json.JSONDecodeError, KeyError):
        return _error("invalid request body")

    dynamodb = boto3.resource("dynamodb")
    table = dynamodb.Table(table_name)

    # Verify challenge exists and not expired
    response = table.get_item(Key={"pk": f"challenge#{challenge}"})
    item = response.get("Item")
    if not item or int(time.time()) > item.get("expires_at", 0):
        return _error("invalid or expired challenge")

    # Delete challenge (single use)
    table.delete_item(Key={"pk": f"challenge#{challenge}"})

    # Verify PoW solution
    h = hashlib.sha256(f"{challenge}{nonce}".encode()).hexdigest()
    if not h.startswith(prefix):
        return _error("invalid PoW solution")

    # Issue API key
    api_key = f"sdbx_cli_{secrets.token_hex(16)}"
    expires_at = int(time.time()) + KEY_TTL

    table.put_item(
        Item={
            "pk": f"apikey#{api_key}",
            "api_key": api_key,
            "expires_at": expires_at,
        }
    )

    logger.info(f"CLI API key issued, expires_at={expires_at}")

    return {
        "statusCode": 200,
        "headers": {"Content-Type": "application/json"},
        "body": json.dumps({"api_key": api_key, "expires_at": expires_at}),
    }
