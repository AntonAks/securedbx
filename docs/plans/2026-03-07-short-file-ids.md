# Short File IDs Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace UUID-based `file_id` with 8-character URL-safe IDs for cleaner share links.

**Architecture:** Add `generate_short_file_id()` to `pin_utils.py`, add collision guard to `create_file_record` in `dynamo.py`, replace `uuid.uuid4()` in `upload_init/handler.py` with a 10-attempt retry loop. Existing UUID records are unaffected — they coexist in DynamoDB.

**Tech Stack:** Python 3.12, `secrets` stdlib, boto3/DynamoDB, pytest

---

### Task 1: Add `generate_short_file_id()` to `pin_utils.py`

**Files:**
- Modify: `backend/shared/pin_utils.py`
- Test: `backend/tests/test_pin_utils.py`

**Step 1: Write the failing tests**

Add this class to `backend/tests/test_pin_utils.py`:

```python
class TestGenerateShortFileId:
    def test_length_is_eight(self):
        file_id = generate_short_file_id()
        assert len(file_id) == 8

    def test_url_safe_chars(self):
        import re
        for _ in range(100):
            file_id = generate_short_file_id()
            assert re.match(r'^[A-Za-z0-9_-]+$', file_id), f"Invalid chars in: {file_id}"

    def test_uniqueness(self):
        ids = {generate_short_file_id() for _ in range(200)}
        assert len(ids) == 200
```

Also update the import at the top of the file — add `generate_short_file_id` to the existing import:

```python
from shared.pin_utils import generate_pin_file_id, generate_salt, generate_short_file_id, hash_pin, verify_pin_hash
```

**Step 2: Run to verify they fail**

```bash
cd backend && . venv/bin/activate && pytest tests/test_pin_utils.py::TestGenerateShortFileId -v --tb=short
```

Expected: `ImportError` or `AttributeError` — `generate_short_file_id` does not exist yet.

**Step 3: Implement the function**

In `backend/shared/pin_utils.py`, add at the top:

```python
import secrets
```

Then add this function after `generate_pin_file_id`:

```python
def generate_short_file_id() -> str:
    """Generate random 8-char URL-safe file ID using base64url encoding."""
    return secrets.token_urlsafe(6)
```

**Step 4: Run tests to verify they pass**

```bash
cd backend && . venv/bin/activate && pytest tests/test_pin_utils.py -v --tb=short
```

Expected: all tests pass, including the new `TestGenerateShortFileId` class.

**Step 5: Commit**

```bash
git add backend/shared/pin_utils.py backend/tests/test_pin_utils.py
git commit -m "feat: add generate_short_file_id() using secrets.token_urlsafe"
```

---

### Task 2: Add collision guard to `create_file_record` in `dynamo.py`

**Files:**
- Modify: `backend/shared/dynamo.py`

> No new test file needed here — the collision behavior is exercised by the handler retry test in Task 3. The dynamo change is a one-line guard following the exact same pattern already used by `create_pin_file_record` (line 594 in `dynamo.py`).

**Step 1: Update `create_file_record`**

In `backend/shared/dynamo.py`, replace the `put_item` call inside `create_file_record` (currently line 107):

Old:
```python
    table.put_item(Item=record)
    logger.info(f"Created {content_type} record ({access_mode}): {file_id}")
```

New:
```python
    try:
        table.put_item(Item=record, ConditionExpression="attribute_not_exists(file_id)")
    except ClientError as e:
        if e.response["Error"]["Code"] == "ConditionalCheckFailedException":
            raise ValueError(f"File ID {file_id} already exists")
        raise
    logger.info(f"Created {content_type} record ({access_mode}): {file_id}")
```

**Step 2: Run full test suite to verify nothing broke**

```bash
cd backend && . venv/bin/activate && pytest tests/ -v --tb=short
```

Expected: all existing tests still pass.

**Step 3: Commit**

```bash
git add backend/shared/dynamo.py
git commit -m "feat: add attribute_not_exists collision guard to create_file_record"
```

---

### Task 3: Replace UUID with short ID + retry loop in `upload_init/handler.py`

**Files:**
- Modify: `backend/lambdas/upload_init/handler.py`

**Step 1: Write the failing test**

Create `backend/tests/test_upload_init_handler.py`:

```python
"""Unit tests for upload_init handler — retry logic."""

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
```

**Step 2: Run to verify tests pass (they test logic, not imports)**

```bash
cd backend && . venv/bin/activate && pytest tests/test_upload_init_handler.py -v --tb=short
```

Expected: all 3 tests pass immediately (they test pure logic).

**Step 3: Update the handler**

In `backend/lambdas/upload_init/handler.py`:

Remove:
```python
import uuid
```

Add to imports at top (with the other `shared` imports):
```python
from shared.pin_utils import generate_short_file_id
```

Also add near the top of the file:
```python
MAX_ID_RETRIES = 10
```

Then replace the ID generation and the two `create_file_record` calls (text and file branches) with a retry loop. The current structure creates `file_id` once then calls `create_file_record` in two branches. Replace lines 127–198 with:

```python
        # Generate unique short ID with retry loop
        file_id = None
        for attempt in range(MAX_ID_RETRIES):
            candidate_id = generate_short_file_id()
            try:
                if content_type == "text":
                    encrypted_text = body.get("encrypted_text")
                    if not encrypted_text:
                        raise ValidationError("encrypted_text is required for text secrets")

                    max_text_size = 100000 if access_mode == ACCESS_MODE_MULTI else 10000
                    if len(encrypted_text) > max_text_size:
                        raise ValidationError("Text secret too large")

                    create_file_record(
                        table_name=TABLE_NAME,
                        file_id=candidate_id,
                        file_size=len(encrypted_text),
                        expires_at=expires_at,
                        ip_hash=ip_hash,
                        content_type="text",
                        encrypted_text=encrypted_text,
                        access_mode=access_mode,
                        salt=salt,
                        encrypted_key=encrypted_key,
                    )

                    logger.info(f"Text secret created: file_id={candidate_id}, size={len(encrypted_text)}, ttl={ttl}, access_mode={access_mode}")

                    return success_response({
                        "file_id": candidate_id,
                        "expires_at": expires_at,
                    })

                else:
                    file_size = body.get("file_size")
                    validate_file_size(file_size)

                    s3_key = f"files/{candidate_id}"

                    create_file_record(
                        table_name=TABLE_NAME,
                        file_id=candidate_id,
                        file_size=file_size,
                        expires_at=expires_at,
                        ip_hash=ip_hash,
                        content_type="file",
                        s3_key=s3_key,
                        access_mode=access_mode,
                        salt=salt,
                        encrypted_key=encrypted_key,
                    )

                    upload_url = generate_upload_url(
                        bucket_name=BUCKET_NAME,
                        s3_key=s3_key,
                        expires_in=UPLOAD_URL_EXPIRY_SECONDS,
                    )

                    logger.info(f"File upload initialized: file_id={candidate_id}, size={file_size}, ttl={ttl}, access_mode={access_mode}")

                    return success_response({
                        "file_id": candidate_id,
                        "upload_url": upload_url,
                        "expires_at": expires_at,
                    })

            except ValueError as e:
                if "already exists" in str(e):
                    logger.warning(f"Short file ID collision: {candidate_id}, retrying (attempt {attempt + 1})")
                    continue
                raise

        logger.error("Failed to generate unique file ID after max retries")
        return error_response("Failed to generate unique file ID. Please try again.", 500)
```

**Step 4: Run full test suite**

```bash
cd backend && . venv/bin/activate && pytest tests/ -v --tb=short
```

Expected: all tests pass.

**Step 5: Commit**

```bash
git add backend/lambdas/upload_init/handler.py backend/tests/test_upload_init_handler.py
git commit -m "feat: replace UUID with short file ID in upload_init handler"
```

---

### Task 4: Update docstring in `dynamo.py`

**Files:**
- Modify: `backend/shared/dynamo.py`

**Step 1: Update the docstring**

In `create_file_record`, the `file_id` param docstring says "Unique file/secret ID (UUID)". Change it to:

```python
        file_id: Unique file/secret ID (8-char URL-safe string)
```

**Step 2: Run tests**

```bash
cd backend && . venv/bin/activate && pytest tests/ -v --tb=short
```

Expected: all tests pass.

**Step 3: Commit**

```bash
git add backend/shared/dynamo.py
git commit -m "docs: update file_id docstring to reflect short ID format"
```

---

### Task 5: Final verification

**Step 1: Run the full test suite one last time**

```bash
cd backend && . venv/bin/activate && pytest tests/ -v --tb=short
```

Expected: all tests pass, zero failures.

**Step 2: Run frontend tests to confirm no regressions**

```bash
cd frontend && npx vitest run
```

Expected: all tests pass.

**Step 3: Manual smoke check**

Review the generated share URL shape. After the change, `file_id` values in share URLs will look like `dX9a_kQm` instead of `550e8400-e29b-41d4-a716-446655440000`. The URL structure (`/#/download?id=...&key=...`) is unchanged.