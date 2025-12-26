"""Lambda function: Initialize file upload."""

import hashlib
import logging
import os
import time
import uuid
from typing import Any

from shared.constants import TTL_TO_SECONDS, UPLOAD_URL_EXPIRY_SECONDS
from shared.dynamo import create_file_record
from shared.exceptions import ValidationError
from shared.request_helpers import get_source_ip, parse_json_body
from shared.response import error_response, success_response
from shared.s3 import generate_upload_url
from shared.security import require_cloudfront_and_recaptcha
from shared.validation import validate_file_size, validate_ttl

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

# Environment variables
BUCKET_NAME = os.environ.get("BUCKET_NAME")
TABLE_NAME = os.environ.get("TABLE_NAME")


@require_cloudfront_and_recaptcha
def handler(event: dict[str, Any], context: Any) -> dict[str, Any]:
    """
    Initialize file or text secret upload.

    Security verification (CloudFront origin + reCAPTCHA) is handled by decorator.

    Expected request body (file):
    {
        "content_type": "file",
        "file_size": 1024,
        "ttl": "1h",
        "recaptcha_token": "token"
    }

    Expected request body (text):
    {
        "content_type": "text",
        "encrypted_text": "base64-encrypted-text",
        "ttl": "1h",
        "recaptcha_token": "token"
    }

    Returns:
    {
        "file_id": "uuid",
        "upload_url": "presigned-s3-url",  // Only for files
        "expires_at": 1234567890
    }
    """
    try:
        # Parse request body
        body = parse_json_body(event)
        content_type = body.get("content_type", "file")
        ttl = body.get("ttl")

        # Validate TTL
        validate_ttl(ttl)

        # Generate unique ID
        file_id = str(uuid.uuid4())

        # Calculate expiration timestamp
        ttl_seconds = TTL_TO_SECONDS[ttl]
        expires_at = int(time.time()) + ttl_seconds

        # Hash IP address (privacy)
        source_ip = get_source_ip(event)
        ip_hash = hashlib.sha256(source_ip.encode()).hexdigest()

        # Handle based on content type
        if content_type == "text":
            # Text secret
            encrypted_text = body.get("encrypted_text")
            if not encrypted_text:
                raise ValidationError("encrypted_text is required for text secrets")

            # Validate text length (base64 encoded)
            if len(encrypted_text) > 10000:  # ~7KB base64 = ~5KB encrypted ~= 1000+ chars
                raise ValidationError("Text secret too large (max 1000 characters)")

            # Create DynamoDB record with encrypted text
            create_file_record(
                table_name=TABLE_NAME,
                file_id=file_id,
                file_size=len(encrypted_text),
                expires_at=expires_at,
                ip_hash=ip_hash,
                content_type="text",
                encrypted_text=encrypted_text,
            )

            logger.info(f"Text secret created: file_id={file_id}, size={len(encrypted_text)}, ttl={ttl}")

            return success_response({
                "file_id": file_id,
                "expires_at": expires_at,
            })

        else:
            # File upload (default)
            file_size = body.get("file_size")
            validate_file_size(file_size)

            s3_key = f"files/{file_id}"

            # Create DynamoDB record
            create_file_record(
                table_name=TABLE_NAME,
                file_id=file_id,
                file_size=file_size,
                expires_at=expires_at,
                ip_hash=ip_hash,
                content_type="file",
                s3_key=s3_key,
            )

            # Generate presigned upload URL
            upload_url = generate_upload_url(
                bucket_name=BUCKET_NAME,
                s3_key=s3_key,
                expires_in=UPLOAD_URL_EXPIRY_SECONDS,
            )

            logger.info(f"File upload initialized: file_id={file_id}, size={file_size}, ttl={ttl}")

            return success_response({
                "file_id": file_id,
                "upload_url": upload_url,
                "expires_at": expires_at,
            })

    except ValidationError as e:
        logger.warning(f"Validation error: {e}")
        return error_response(str(e), 400)

    except Exception as e:
        logger.exception("Unexpected error in upload_init")
        return error_response("Internal server error", 500)
