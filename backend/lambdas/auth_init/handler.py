"""Lambda function: Generate PoW challenge for CLI auth."""

import json
import logging
import os
import secrets
import time
from typing import Any

import boto3

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

CHALLENGE_TTL = 300  # 5 minutes


def handler(event: dict[str, Any], context: Any) -> dict[str, Any]:
    """
    Generate a PoW challenge and store it in DynamoDB (TTL 5 min).

    No CloudFront/auth check — this endpoint is intentionally open
    (rate-limited at API Gateway level).

    Response: { challenge: "hex32chars", difficulty: 4 }
    """
    table_name = os.environ.get("AUTH_TABLE_NAME")
    difficulty = int(os.environ.get("POW_DIFFICULTY", "4"))

    challenge = secrets.token_hex(16)  # 32 hex chars
    expires_at = int(time.time()) + CHALLENGE_TTL

    if table_name:
        dynamodb = boto3.resource("dynamodb")
        table = dynamodb.Table(table_name)
        table.put_item(
            Item={
                "pk": f"challenge#{challenge}",
                "challenge": challenge,
                "expires_at": expires_at,
            }
        )
    else:
        logger.warning("AUTH_TABLE_NAME not set — challenge not persisted (dev mode)")

    logger.info(f"PoW challenge issued: difficulty={difficulty}")

    return {
        "statusCode": 200,
        "headers": {"Content-Type": "application/json"},
        "body": json.dumps({"challenge": challenge, "difficulty": difficulty}),
    }
