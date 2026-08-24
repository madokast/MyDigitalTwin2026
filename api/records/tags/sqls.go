package tags

const getRecordTagsSQL = `
SELECT tags FROM records WHERE id=$1;
`

const getAllTagCountSQL = `
SELECT
    tag,
    COUNT(*) AS cnt
FROM records
CROSS JOIN LATERAL unnest(tags) AS tag
GROUP BY tag
ORDER BY cnt DESC, tag;
`

const (
	attachResultTagAttached      = 1
	attachResultRecordNotFound   = 2
	attachResultTagAlreadyExists = 3
	detachResultTagDetached      = 4
	detachResultRecordNotFound   = 5
	detachResultTagNotExists     = 6
	renameResultTagRenamed       = 7
	renameResultRecordNotFound   = 8
	renameResultTagNotExists     = 9
	renameResultNewTagSameAsOld  = 10
	renameResultNewTagExists     = 11
)

const attachTagSQL = `
WITH target AS (
    SELECT tags
    FROM records
    WHERE id = $1
),
updated AS (
    UPDATE records
    SET tags = array_append(tags, $2)
    WHERE id = $1
      AND NOT ($2 = ANY(tags))
    RETURNING tags
)
SELECT
    CASE
        WHEN EXISTS (SELECT 1 FROM updated) THEN 1 -- tag attached
        WHEN NOT EXISTS (SELECT 1 FROM target) THEN 2 -- record not found
        ELSE 3 -- tag already exists
    END AS result,
    COALESCE(
        (SELECT tags FROM updated),
        (SELECT tags FROM target),
        ARRAY[]::TEXT[]
    ) AS tags;
`

const detachTagSQL = `
WITH target AS (
    SELECT tags
    FROM records
    WHERE id = $1
),
updated AS (
    UPDATE records
    SET tags = array_remove(tags, $2)
    WHERE id = $1
      AND $2 = ANY(tags)
    RETURNING tags
)
SELECT
    CASE
        WHEN EXISTS (SELECT 1 FROM updated) THEN 4 -- tag detached
        WHEN NOT EXISTS (SELECT 1 FROM target) THEN 5 -- record not found
        ELSE 6 -- tag not exists
    END AS result,
    COALESCE(
        (SELECT tags FROM updated),
        (SELECT tags FROM target),
        ARRAY[]::TEXT[]
    ) AS tags;
`

const renameTagSQL = `
WITH target AS (
    SELECT tags
    FROM records
    WHERE id = $1
),
updated AS (
    UPDATE records
    SET tags = array_replace(tags, $2, $3)
    WHERE id = $1
      AND $2 = ANY(tags)
      AND $2 <> $3
      AND NOT ($3 = ANY(tags))
    RETURNING tags
)
SELECT
    CASE
        WHEN EXISTS (SELECT 1 FROM updated) THEN 7 -- tag renamed
        WHEN NOT EXISTS (SELECT 1 FROM target) THEN 8 -- record not found
        WHEN NOT EXISTS (SELECT 1 FROM target WHERE $2 = ANY(target.tags)) THEN 9 -- tag not exists
        WHEN $2 = $3 THEN 10 -- new tag same as old
        ELSE 11 -- new tag already exists
    END AS result,
    COALESCE(
        (SELECT tags FROM updated),
        (SELECT tags FROM target),
        ARRAY[]::TEXT[]
    ) AS tags;
`
