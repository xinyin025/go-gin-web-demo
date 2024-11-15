CREATE TABLE "public"."users" (
      "id" int8 NOT NULL DEFAULT nextval('seq_user'::regclass),
      "created_at" timestamp(6),
      "updated_at" timestamp(6),
      "deleted_at" timestamp(6),
      "username" varchar(255),
      "password" varchar(255),
      CONSTRAINT "users_pkey" PRIMARY KEY ("id")
)
;