CREATE TABLE "schema_migration" (
"version" TEXT NOT NULL
);
CREATE UNIQUE INDEX "version_idx" ON "schema_migration" (version);
CREATE TABLE "users" (
"id" TEXT PRIMARY KEY,
"username" TEXT NOT NULL,
"email" TEXT NOT NULL,
"admin" NUMERIC NOT NULL,
"password_hash" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE "posts" (
"id" TEXT PRIMARY KEY,
"title" TEXT NOT NULL,
"content" text NOT NULL,
"author_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE "comments" (
"id" TEXT PRIMARY KEY,
"content" text NOT NULL,
"author_id" char(36) NOT NULL,
"post_id" char(36) NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
