"""PIN-based sharing utility functions."""

import hashlib
import hmac
import os
import secrets
from typing import Final

from .constants import PIN_PBKDF2_ITERATIONS, PIN_SALT_BYTES

_PIN_FILE_ID_MAX: Final[int] = 999999


def generate_salt() -> str:
    """Generate random salt. Returns 64-char hex string (32 bytes)."""
    return os.urandom(PIN_SALT_BYTES).hex()


def hash_pin(pin: str, salt: str) -> str:
    """Hash PIN with salt using PBKDF2-HMAC-SHA256. Returns 64-char hex hash."""
    return hashlib.pbkdf2_hmac(
        "sha256",
        pin.encode("utf-8"),
        bytes.fromhex(salt),
        PIN_PBKDF2_ITERATIONS,
    ).hex()


def verify_pin_hash(pin: str, salt: str, expected_hash: str) -> bool:
    """Verify PIN against stored hash using constant-time comparison."""
    actual_hash = hash_pin(pin, salt)
    return hmac.compare_digest(actual_hash, expected_hash)


def generate_pin_file_id() -> str:
    """Generate random 6-digit numeric file ID. Returns zero-padded string."""
    random_int = int.from_bytes(os.urandom(4), "big") % (_PIN_FILE_ID_MAX + 1)
    return f"{random_int:06d}"


def generate_short_file_id() -> str:
    """Generate random 8-char URL-safe file ID using base64url encoding."""
    return secrets.token_urlsafe(6)
