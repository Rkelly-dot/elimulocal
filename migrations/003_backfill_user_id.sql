-- Migration 003: Backfill user_id from username match
-- Links legacy resources uploaded by a registered username to their user record.
UPDATE resources
SET user_id = (
    SELECT id FROM users WHERE username = resources.uploaded_by
)
WHERE user_id = 0
  AND EXISTS (
      SELECT 1 FROM users WHERE username = resources.uploaded_by
  );
