package xsql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseQueries(t *testing.T) {
	testCases := []parseQueriesTestCase{
		{
			"single DO statement",
			query1,
			expected1,
		},
		{
			"DO statement + CREATE TABLE statement",
			query2,
			expected2,
		},
		{
			"$$ statement + CREATE TABLE statement",
			query3,
			expected3,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%s (%d)", tc.title, i), func(t *testing.T) {
			queries, err := ParseQueries(tc.query)
			require.Nil(t, err)
			require.Equal(t, tc.expected, queries)
		})
	}
}

type parseQueriesTestCase struct {
	title    string
	query    string
	expected []string
}

var query1 = `DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM "pg_type" WHERE typname = 'storage_type') THEN
		CREATE TYPE "storage_type" AS ENUM (
			'S3',
			'YT'
		);
    END IF;
END$$;
`

var expected1 = []string{`DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM "pg_type" WHERE typname = 'storage_type') THEN
		CREATE TYPE "storage_type" AS ENUM (
			'S3',
			'YT'
		);
    END IF;
END$$`}

var query2 = `DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM "pg_type" WHERE typname = 'storage_type') THEN
		CREATE TYPE "storage_type" AS ENUM (
			'S3',
			'YT'
		);
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS "resource" (
    "id" TEXT PRIMARY KEY,
	"byte_size" BIGINT NOT NULL,

	"file_name" TEXT,
	"checksum_md5" TEXT UNIQUE,
	"context" JSONB,

	"storage_type" storage_type NOT NULL,
	"s3_storage_info" JSONB,
	"yt_storage_info" JSONB,

	"creation_ts" BIGINT NOT NULL,
    "has_ownership" BOOLEAN NOT NULL,
	"draft_creation_ts" BIGINT
    
    -- CONSTRAINT checksum_is_required_for_non_external_resources CHECK ((NOT "has_ownership") OR ("checksum_md5" IS NOT NULL))
    CONSTRAINT storage_type_makes_corresponding_column_required CHECK (
		("storage_type" = 'YT') AND ("yt_storage_info" IS NOT NULL) OR
		("storage_type" = 'S3') AND ("s3_storage_info" IS NOT NULL)
	)
);`

var expected2 = []string{`DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM "pg_type" WHERE typname = 'storage_type') THEN
		CREATE TYPE "storage_type" AS ENUM (
			'S3',
			'YT'
		);
    END IF;
END$$`,
	`CREATE TABLE IF NOT EXISTS "resource" (
    "id" TEXT PRIMARY KEY,
	"byte_size" BIGINT NOT NULL,

	"file_name" TEXT,
	"checksum_md5" TEXT UNIQUE,
	"context" JSONB,

	"storage_type" storage_type NOT NULL,
	"s3_storage_info" JSONB,
	"yt_storage_info" JSONB,

	"creation_ts" BIGINT NOT NULL,
    "has_ownership" BOOLEAN NOT NULL,
	"draft_creation_ts" BIGINT
    
    -- CONSTRAINT checksum_is_required_for_non_external_resources CHECK ((NOT "has_ownership") OR ("checksum_md5" IS NOT NULL))
    CONSTRAINT storage_type_makes_corresponding_column_required CHECK (
		("storage_type" = 'YT') AND ("yt_storage_info" IS NOT NULL) OR
		("storage_type" = 'S3') AND ("s3_storage_info" IS NOT NULL)
	)
)`}

var query3 = `CREATE INDEX IF NOT EXISTS "resource_blob__blob_id" ON "resource_blob" USING HASH (
	"blob_id"
);

CREATE OR REPLACE FUNCTION build_blob_id(storage_type storage_type, s3_storage_info JSONB, yt_storage_info JSONB)
RETURNS text AS $$
BEGIN
    IF storage_type = 'S3' THEN
        RETURN storage_type || '_' || (s3_storage_info->>'endpoint') || '_' || (s3_storage_info->>'region') || '_' || (s3_storage_info->>'bucket') || '_' || (s3_storage_info->>'key');
    ELSE
		RETURN storage_type || '_' || (yt_storage_info->>'proxy') || '_' || (yt_storage_info->>'path');
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;`

var expected3 = []string{
	`CREATE INDEX IF NOT EXISTS "resource_blob__blob_id" ON "resource_blob" USING HASH (
	"blob_id"
)`,
	`CREATE OR REPLACE FUNCTION build_blob_id(storage_type storage_type, s3_storage_info JSONB, yt_storage_info JSONB)
RETURNS text AS $$
BEGIN
    IF storage_type = 'S3' THEN
        RETURN storage_type || '_' || (s3_storage_info->>'endpoint') || '_' || (s3_storage_info->>'region') || '_' || (s3_storage_info->>'bucket') || '_' || (s3_storage_info->>'key');
    ELSE
		RETURN storage_type || '_' || (yt_storage_info->>'proxy') || '_' || (yt_storage_info->>'path');
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql`}
