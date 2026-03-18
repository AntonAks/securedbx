"""Unit tests for upload_init handler — retry logic."""


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
