ALTER TABLE "release_pilots"
    DROP CONSTRAINT "pilot_uniq_combination",
    ADD CONSTRAINT "pilot_uniq_combination" UNIQUE ("flag_id", "external_id");

ALTER TABLE "release_pilots"
    RENAME COLUMN "flag_id" TO "feature_flag_id";

ALTER TABLE "release_pilots"
    RENAME TO "pilots";
