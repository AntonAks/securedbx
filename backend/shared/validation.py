"""Input validation utilities."""

import re
from typing import Any

from .constants import ALLOWED_TTL_VALUES, MAX_FILE_SIZE_BYTES, MAX_FILE_SIZE_MB
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


def validate_ttl(ttl: str) -> None:
    """
    Validate TTL value.

    Args:
        ttl: Time to live value (1h, 12h, 24h)

    Raises:
        ValidationError: If TTL is invalid
    """
    if ttl not in ALLOWED_TTL_VALUES:
        raise ValidationError(f"TTL must be one of {ALLOWED_TTL_VALUES}")
