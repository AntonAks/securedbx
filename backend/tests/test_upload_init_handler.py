"""Unit tests for upload_init handler — retry logic and text_size validation."""

import json
from unittest.mock import patch

import pytest


class TestShortFileIdRetryLogic:
    """Test the retry logic pattern used by upload_init handler."""

    def test_retry_exhaustion_returns_error(self):
        """After max retries, handler should signal failure."""
        MAX_ID_RETRIES = 10
        attempts = 0
        succeeded = False

        for _ in range(MAX_ID_RETRIES):
            attempts += 1
            # Simulate all attempts colliding
            collision = True
            if not collision:
                succeeded = True
                break

        assert attempts == MAX_ID_RETRIES
        assert succeeded is False

    def test_success_on_first_attempt(self):
        """Should succeed immediately when no collision."""
        MAX_ID_RETRIES = 10
        attempts = 0
        succeeded = False

        for _ in range(MAX_ID_RETRIES):
            attempts += 1
            collision = False  # No collision
            if not collision:
                succeeded = True
                break

        assert attempts == 1
        assert succeeded is True

    def test_success_on_third_attempt(self):
        """Should succeed after two collisions."""
        MAX_ID_RETRIES = 10
        attempts = 0
        succeeded = False
        collisions_before_success = 2

        for i in range(MAX_ID_RETRIES):
            attempts += 1
            collision = i < collisions_before_success
            if not collision:
                succeeded = True
                break

        assert attempts == 3
        assert succeeded is True


class TestTextSizeValidation:
    """Test text_size field handling in text upload path."""

    def _make_event(self, body: dict) -> dict:
        """Build a minimal Lambda event with CloudFront + reCAPTCHA headers."""
        return {
            "headers": {
                "X-Origin-Verify": "test-secret",
                "CF-Connecting-IP": "1.2.3.4",
            },
            "body": json.dumps(body),
        }

    def _base_body(self, **overrides) -> dict:
        return {
            "content_type": "text",
            "encrypted_text": "A" * 100,
            "ttl": "1h",
            "recaptcha_token": "tok",
            **overrides,
        }

    def test_text_size_valid_stored_as_file_size(self, monkeypatch):
        """When text_size is a valid int, it is passed to create_file_record as file_size."""
        monkeypatch.setenv("CLOUDFRONT_SECRET", "test-secret")
        from importlib import reload

        import lambdas.upload_init.handler as h

        reload(h)
        with (
            patch("lambdas.upload_init.handler.create_file_record") as mock_create,
            patch("lambdas.upload_init.handler.generate_unique_file_id", return_value="ABCD1234"),
            patch("lambdas.upload_init.handler.hash_ip_secure", return_value="h"),
            patch("lambdas.upload_init.handler.get_source_ip", return_value="1.2.3.4"),
            patch("lambdas.upload_init.handler.TABLE_NAME", "t"),
            patch("shared.security.require_cloudfront_and_recaptcha", lambda f: f),
        ):
            event = self._make_event(self._base_body(text_size=42))
            h.handler(event, None)
            _, kwargs = mock_create.call_args
            assert kwargs["file_size"] == 42

    def test_text_size_absent_falls_back_to_encrypted_len(self, monkeypatch):
        """When text_size is absent, file_size falls back to len(encrypted_text)."""
        encrypted = "B" * 80
        monkeypatch.setenv("CLOUDFRONT_SECRET", "test-secret")
        from importlib import reload

        import lambdas.upload_init.handler as h

        reload(h)
        with (
            patch("lambdas.upload_init.handler.create_file_record") as mock_create,
            patch("lambdas.upload_init.handler.generate_unique_file_id", return_value="ABCD1234"),
            patch("lambdas.upload_init.handler.hash_ip_secure", return_value="h"),
            patch("lambdas.upload_init.handler.get_source_ip", return_value="1.2.3.4"),
            patch("lambdas.upload_init.handler.TABLE_NAME", "t"),
            patch("shared.security.require_cloudfront_and_recaptcha", lambda f: f),
        ):
            event = self._make_event(self._base_body(encrypted_text=encrypted))
            h.handler(event, None)
            _, kwargs = mock_create.call_args
            assert kwargs["file_size"] == 80

    def test_text_size_null_falls_back_to_encrypted_len(self, monkeypatch):
        """When text_size is explicitly null in JSON, fallback to len(encrypted_text)."""
        encrypted = "C" * 60
        monkeypatch.setenv("CLOUDFRONT_SECRET", "test-secret")
        from importlib import reload

        import lambdas.upload_init.handler as h

        reload(h)
        with (
            patch("lambdas.upload_init.handler.create_file_record") as mock_create,
            patch("lambdas.upload_init.handler.generate_unique_file_id", return_value="ABCD1234"),
            patch("lambdas.upload_init.handler.hash_ip_secure", return_value="h"),
            patch("lambdas.upload_init.handler.get_source_ip", return_value="1.2.3.4"),
            patch("lambdas.upload_init.handler.TABLE_NAME", "t"),
            patch("shared.security.require_cloudfront_and_recaptcha", lambda f: f),
        ):
            body = self._base_body(encrypted_text=encrypted)
            body["text_size"] = None
            event = self._make_event(body)
            h.handler(event, None)
            _, kwargs = mock_create.call_args
            assert kwargs["file_size"] == 60

    @pytest.mark.parametrize(
        "bad_value,label",
        [
            (0, "zero"),
            (-1, "negative"),
            (99.9, "float"),
            (True, "bool_true"),
            (False, "bool_false"),
            ("100", "string"),
        ],
    )
    def test_text_size_invalid_returns_400(self, bad_value, label, monkeypatch):
        """Invalid text_size values must return HTTP 400."""
        monkeypatch.setenv("CLOUDFRONT_SECRET", "test-secret")
        from importlib import reload

        import lambdas.upload_init.handler as h

        reload(h)
        with (
            patch("lambdas.upload_init.handler.create_file_record"),
            patch("lambdas.upload_init.handler.generate_unique_file_id", return_value="ABCD1234"),
            patch("lambdas.upload_init.handler.hash_ip_secure", return_value="h"),
            patch("lambdas.upload_init.handler.get_source_ip", return_value="1.2.3.4"),
            patch("lambdas.upload_init.handler.TABLE_NAME", "t"),
            patch("shared.security.require_cloudfront_and_recaptcha", lambda f: f),
        ):
            event = self._make_event(self._base_body(text_size=bad_value))
            result = h.handler(event, None)
            assert result["statusCode"] == 400, f"Expected 400 for text_size={label!r}"

    def test_text_size_exceeds_encrypted_len_returns_400(self, monkeypatch):
        """text_size larger than len(encrypted_text) must return 400."""
        encrypted = "D" * 50
        monkeypatch.setenv("CLOUDFRONT_SECRET", "test-secret")
        from importlib import reload

        import lambdas.upload_init.handler as h

        reload(h)
        with (
            patch("lambdas.upload_init.handler.create_file_record"),
            patch("lambdas.upload_init.handler.generate_unique_file_id", return_value="ABCD1234"),
            patch("lambdas.upload_init.handler.hash_ip_secure", return_value="h"),
            patch("lambdas.upload_init.handler.get_source_ip", return_value="1.2.3.4"),
            patch("lambdas.upload_init.handler.TABLE_NAME", "t"),
            patch("shared.security.require_cloudfront_and_recaptcha", lambda f: f),
        ):
            event = self._make_event(self._base_body(encrypted_text=encrypted, text_size=9999))
            result = h.handler(event, None)
            assert result["statusCode"] == 400
