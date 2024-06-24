CREATE TABLE "booked_appointments"(
                                      "id" SERIAL PRIMARY KEY NOT NULL,
                                      "department_id" UUID NOT NULL,
                                      "doctor_id" UUID NULL,
                                      "patient_id" UUID NOT NULL,
                                      "doctor_service_id" UUID NOT NULL,
                                      "appointment_date" DATE NOT NULL,
                                      "appointment_time" TIME(0) WITHOUT TIME ZONE NOT NULL,
                                      "duration" BIGINT NOT NULL,
                                      "key" VARCHAR(20) NOT NULL,
                                      "expires_at" TIMESTAMP(0) WITHOUT TIME ZONE,
                                      "patient_problem" TEXT NOT NULL
                                      "status" VARCHAR(255) NOT NULL CHECK ("status" IN ('attended', 'cancelled', 'no_show', 'waiting')),
                                      "payment_type" VARCHAR(255) NOT NULL CHECK ("payment_type" IN ('cash', 'card', 'insurance')),
                                      "payment_amount" DOUBLE PRECISION NOT NULL,
                                      "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE ,
                                      "deleted_at" TIMESTAMP(0) WITHOUT TIME ZONE
);