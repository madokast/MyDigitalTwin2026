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
