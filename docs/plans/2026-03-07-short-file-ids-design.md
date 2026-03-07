# Short File IDs — Design

**Date:** 2026-03-07
**Branch:** `Short-URLs-Shorter-file-IDs-for-cleaner-links`

## Problem

Regular uploads currently use a full UUID v4 as `file_id` (36 chars, e.g. `550e8400-e29b-41d4-a716-446655440000`). This makes share URLs long and hard to read/share verbally. PIN uploads already use 6-digit numeric IDs, so the infra supports short IDs.

## Decision

Replace UUID with `secrets.token_urlsafe(6)` — an 8-character URL-safe base64 string (charset: A-Z, a-z, 0-9, `-`, `_`). 64^8 ~ 281 trillion combinations.

This is Option A: the short ID becomes the primary key everywhere (DynamoDB partition key, S3 object key, share URL). No secondary short-code field, no GSI.

## Backward Compatibility

No migration needed. Existing UUID records remain in DynamoDB and S3 unchanged. Old share links continue to work — `get_file_record` is a simple key lookup by string, indifferent to format. New uploads get short IDs; old records keep UUIDs. They coexist in the same table.

## Collision Handling

`create_file_record` will use `ConditionExpression="attribute_not_exists(file_id)"` (same pattern already used by `create_pin_file_record`). On `ConditionalCheckFailedException`, raise `ValueError`. The upload handler retries up to 10 times before returning a 500 error. With 281 trillion combinations, collisions will not occur in practice.

## Files Changed

| File | Change |
|------|--------|
| `backend/shared/pin_utils.py` | Add `generate_short_file_id()` using `secrets.token_urlsafe(6)` |
| `backend/shared/dynamo.py` | Add `attribute_not_exists` condition + `ValueError` to `create_file_record` |
| `backend/lambdas/upload_init/handler.py` | Replace `uuid.uuid4()` with 10-attempt retry loop; remove `import uuid` |

Frontend: no changes. Tests: update assertions that check UUID format.

## Non-Goals

- No changes to PIN upload flow (already uses short numeric IDs)
- No URL redirects or alias table
- No migration of existing records
