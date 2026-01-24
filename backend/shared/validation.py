"""Input validation utilities."""

import re
from typing import Any

from .constants import (
    ALLOWED_TTL_VALUES,
    MAX_CUSTOM_TTL_MINUTES,
    MAX_FILE_SIZE_BYTES,
    MAX_FILE_SIZE_MB,
    MIN_CUSTOM_TTL_MINUTES,
)
from .exceptions import ValidationError

# UUID pattern for file ID validation
UUID_PATTERN = re.compile(r"^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")


def validate_file_id(file_id: str) -> None:
    """
    Validate file ID format (UUID v4).

    Args:
        file_id: File ID to validate

    Raises:
        ValidationError: If file ID is invalid
    """
    if not file_id:
        raise ValidationError("File ID is required")

    if not UUID_PATTERN.match(file_id.lower()):
        raise ValidationError("Invalid file ID format")


def validate_file_size(file_size: int) -> None:
    """
    Validate file size.

    Args:
        file_size: File size in bytes

    Raises:
        ValidationError: If file size is invalid
    """
    if not isinstance(file_size, int):
        raise ValidationError("File size must be an integer")

    if file_size <= 0:
        raise ValidationError("File size must be positive")

    if file_size > MAX_FILE_SIZE_BYTES:
        raise ValidationError(f"File size exceeds maximum limit ({MAX_FILE_SIZE_MB} MB)")


def validate_ttl(ttl: Any) -> None:
    """
    Validate TTL value.

    Accepts either:
    - Preset strings: "1h", "12h", "24h"
    - Numeric minutes: 5 to 10080 (5 min to 7 days)

    Args:
        ttl: Time to live value

    Raises:
        ValidationError: If TTL is invalid
    """
    # Accept preset strings
    if isinstance(ttl, str) and ttl in ALLOWED_TTL_VALUES:
        return

    # Accept numeric minutes
    if isinstance(ttl, (int, float)):
        minutes = int(ttl)
        if minutes < MIN_CUSTOM_TTL_MINUTES:
            raise ValidationError(
                f"Custom TTL must be at least {MIN_CUSTOM_TTL_MINUTES} minutes"
            )
        if minutes > MAX_CUSTOM_TTL_MINUTES:
            raise ValidationError(
                f"Custom TTL cannot exceed {MAX_CUSTOM_TTL_MINUTES} minutes (7 days)"
            )
        return

    # Invalid format
    raise ValidationError(
        f"TTL must be one of {ALLOWED_TTL_VALUES} or a number of minutes ({MIN_CUSTOM_TTL_MINUTES}-{MAX_CUSTOM_TTL_MINUTES})"
    )
