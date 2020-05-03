ALTER TABLE "pilots"
    RENAME TO "release_pilots";

ALTER TABLE "release_pilots"
    RENAME COLUMN "feature_flag_id" TO "flag_id";

ALTER TABLE "release_pilots"
    DROP CONSTRAINT "pilot_uniq_combination",
    ADD CONSTRAINT "pilot_uniq_combination" UNIQUE ("flag_id", "env_id", "external_id");
