package tags

const getRecordTagsSQL = `
SELECT tags FROM records WHERE id=$1;
`
