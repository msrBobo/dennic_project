INSERT INTO patients (id, first_name, last_name, birth_date, gender, blood_group, phone_number, address, city, country, patient_problem, created_at, updated_at, deleted_at)
VALUES
('123e4567-e89b-12d3-a456-426614370001', 'John', 'Doe', '1990-05-15', 'male', 'A+', '1234567890', '123 Main St', 'Anytown', 'USA', 'Headache', '2024-05-09 08:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370002', 'Jane', 'Doe', '1985-08-20', 'female', 'B+', '9876543210', '456 Elm St', 'Othertown', 'UK', 'Fever', '2024-05-09 09:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370003', 'Michael', 'Smith', '1978-03-10', 'male', 'AB+', '5554443333', '789 Oak St', 'Smalltown', 'Canada', 'Cough', '2024-05-09 10:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370004', 'Emily', 'Johnson', '1995-12-25', 'female', 'O-', '2223334444', '101 Pine St', 'Villagetown', 'Australia', 'Sore throat', '2024-05-09 11:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370005', 'David', 'Brown', '1982-07-05', 'male', 'B-', '7778889999', '202 Maple St', 'Greenvillage', 'New Zealand', 'Back pain', '2024-05-09 12:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370006', 'Olivia', 'Wilson', '1992-09-18', 'female', 'A-', '6665554444', '303 Cedar St', 'Largetown', 'Germany', 'Stomach ache', '2024-05-09 13:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370007', 'James', 'Taylor', '1980-11-30', 'male', 'O+', '9990001111', '404 Oak St', 'Hometown', 'France', 'Allergy', '2024-05-09 14:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370008', 'Sophia', 'Lee', '1988-04-12', 'female', 'AB-', '1112223333', '505 Elm St', 'Bigtown', 'Italy', 'Fatigue', '2024-05-09 15:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370009', 'Daniel', 'Martinez', '1976-01-08', 'male', 'A-', '4445556666', '606 Pine St', 'Cityville', 'Spain', 'Nausea', '2024-05-09 16:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614370010', 'Isabella', 'Garcia', '1987-06-22', 'female', 'B+', '8889990000', '707 Maple St', 'Townsville', 'Japan', 'Fever', '2024-05-09 17:00:00', NULL, NULL);

----------------------------------------------------------------------------------------------------------------------------------------------

INSERT INTO doctor_availability (id, department_id, doctor_id, doctor_date, start_time, end_time, status, created_at, updated_at, deleted_at)
VALUES
(1,'123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614274001', '2024-05-09', '08:00:00', '12:00:00', 'available', '2024-05-09 08:00:00', NULL, NULL),
(2,'123e4567-e89b-12d3-a456-426614174002', '123e4567-e89b-12d3-a456-426614274002', '2024-05-09', '09:00:00', '13:00:00', 'available', '2024-05-09 09:00:00', NULL, NULL),
(3,'123e4567-e89b-12d3-a456-426614174003', '123e4567-e89b-12d3-a456-426614274003', '2024-05-09', '10:00:00', '14:00:00', 'available', '2024-05-09 10:00:00', NULL, NULL),
(4,'123e4567-e89b-12d3-a456-426614174004', '123e4567-e89b-12d3-a456-426614274004', '2024-05-09', '11:00:00', '15:00:00', 'available', '2024-05-09 11:00:00', NULL, NULL),
(5,'123e4567-e89b-12d3-a456-426614174005', '123e4567-e89b-12d3-a456-426614274005', '2024-05-09', '12:00:00', '16:00:00', 'available', '2024-05-09 12:00:00', NULL, NULL),
(6,'123e4567-e89b-12d3-a456-426614174006', '123e4567-e89b-12d3-a456-426614274006', '2024-05-09', '13:00:00', '17:00:00', 'available', '2024-05-09 13:00:00', NULL, NULL),
(7,'123e4567-e89b-12d3-a456-426614174007', '123e4567-e89b-12d3-a456-426614274007', '2024-05-09', '14:00:00', '18:00:00', 'available', '2024-05-09 14:00:00', NULL, NULL),
(8,'123e4567-e89b-12d3-a456-426614174008', '123e4567-e89b-12d3-a456-426614274008', '2024-05-09', '15:00:00', '19:00:00', 'available', '2024-05-09 15:00:00', NULL, NULL),
(9,'123e4567-e89b-12d3-a456-426614174009', '123e4567-e89b-12d3-a456-426614274009', '2024-05-09', '16:00:00', '20:00:00', 'available', '2024-05-09 16:00:00', NULL, NULL),
(10,'123e4567-e89b-12d3-a456-426614174010', '123e4567-e89b-12d3-a456-426614274010', '2024-05-09', '17:00:00', '21:00:00', 'available', '2024-05-09 17:00:00', NULL, NULL);

----------------------------------------------------------------------------------------------------------------------------------------------

INSERT INTO archive (id, doctor_availability_id, start_time, patient_problem, end_time, status, payment_type, payment_amount, created_at, updated_at, deleted_at)
VALUES
(1, 1, '08:00:00', 'Headache', '09:00:00', 'attended', 'card', 50.00, '2024-05-09 08:00:00', NULL, NULL),
(2, 2, '09:00:00', 'Fever', '10:00:00', 'attended', 'cash', 40.00, '2024-05-09 09:00:00', NULL, NULL),
(3, 3, '10:00:00', 'Cough', '11:00:00', 'cancelled', 'insurance', 60.00, '2024-05-09 10:00:00', NULL, NULL),
(4, 4, '11:00:00', 'Sore throat', '12:00:00', 'no_show', 'card', 45.00, '2024-05-09 11:00:00', NULL, NULL),
(5, 5, '12:00:00', 'Back pain', '13:00:00', 'attended', 'cash', 55.00, '2024-05-09 12:00:00', NULL, NULL),
(6, 6, '13:00:00', 'Stomach ache', '14:00:00', 'attended', 'card', 65.00, '2024-05-09 13:00:00', NULL, NULL),
(7, 7, '14:00:00', 'Allergy', '15:00:00', 'cancelled', 'cash', 70.00, '2024-05-09 14:00:00', NULL, NULL),
(8, 8, '15:00:00', 'Fatigue', '16:00:00', 'no_show', 'insurance', 80.00, '2024-05-09 15:00:00', NULL, NULL),
(9, 9, '16:00:00', 'Nausea', '17:00:00', 'attended', 'card', 90.00, '2024-05-09 16:00:00', NULL, NULL),
(10,10, '17:00:00', 'Fever', '18:00:00', 'attended', 'cash', 100.00, '2024-05-09 17:00:00', NULL, NULL);

----------------------------------------------------------------------------------------------------------------------------------------------

INSERT INTO doctor_notes (appointment_id, doctor_id, patient_id, prescription, created_at, updated_at, deleted_at)
VALUES
(1, '123e4567-e89b-12d3-a456-426614274001', '123e4567-e89b-12d3-a456-426614370001', 'Prescription for patient 1', '2024-05-09 08:00:00', NULL, NULL),
(2, '123e4567-e89b-12d3-a456-426614274002', '123e4567-e89b-12d3-a456-426614370002', 'Prescription for patient 2', '2024-05-09 09:00:00', NULL, NULL),
(3, '123e4567-e89b-12d3-a456-426614274003', '123e4567-e89b-12d3-a456-426614370003', 'Prescription for patient 3', '2024-05-09 10:00:00', NULL, NULL),
(4, '123e4567-e89b-12d3-a456-426614274004', '123e4567-e89b-12d3-a456-426614370004', 'Prescription for patient 4', '2024-05-09 11:00:00', NULL, NULL),
(5, '123e4567-e89b-12d3-a456-426614274005', '123e4567-e89b-12d3-a456-426614370005', 'Prescription for patient 5', '2024-05-09 12:00:00', NULL, NULL),
(6, '123e4567-e89b-12d3-a456-426614274006', '123e4567-e89b-12d3-a456-426614370006', 'Prescription for patient 6', '2024-05-09 13:00:00', NULL, NULL),
(7, '123e4567-e89b-12d3-a456-426614274007', '123e4567-e89b-12d3-a456-426614370007', 'Prescription for patient 7', '2024-05-09 14:00:00', NULL, NULL),
(8, '123e4567-e89b-12d3-a456-426614274008', '123e4567-e89b-12d3-a456-426614370008', 'Prescription for patient 8', '2024-05-09 15:00:00', NULL, NULL),
(9, '123e4567-e89b-12d3-a456-426614274009', '123e4567-e89b-12d3-a456-426614370009', 'Prescription for patient 9', '2024-05-09 16:00:00', NULL, NULL),
(10,'123e4567-e89b-12d3-a456-426614274010',  '123e4567-e89b-12d3-a456-426614370010', 'Prescription for patient 10', '2024-05-09 17:00:00', NULL, NULL);

----------------------------------------------------------------------------------------------------------------------------------------------

INSERT INTO booked_appointments (department_id, doctor_id, patient_id, appointment_date, appointment_time, duration, key, expires_at, patient_status, created_at, updated_at, deleted_at)
VALUES
('123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614274001', '123e4567-e89b-12d3-a456-426614370001', '2024-05-10', '08:00:00', 30, 'key123', '2024-05-10 07:30:00', TRUE, '2024-05-09 08:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174002', '123e4567-e89b-12d3-a456-426614274002', '123e4567-e89b-12d3-a456-426614370002', '2024-05-10', '09:00:00', 30, 'key124', '2024-05-10 08:30:00', TRUE, '2024-05-09 09:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174003', '123e4567-e89b-12d3-a456-426614274003', '123e4567-e89b-12d3-a456-426614370003', '2024-05-10', '10:00:00', 30, 'key125', '2024-05-10 09:30:00', TRUE, '2024-05-09 10:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174004', '123e4567-e89b-12d3-a456-426614274004', '123e4567-e89b-12d3-a456-426614370004', '2024-05-10', '11:00:00', 30, 'key126', '2024-05-10 10:30:00', TRUE, '2024-05-09 11:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174005', '123e4567-e89b-12d3-a456-426614274005', '123e4567-e89b-12d3-a456-426614370005', '2024-05-10', '12:00:00', 30, 'key127', '2024-05-10 11:30:00', TRUE, '2024-05-09 12:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174006', '123e4567-e89b-12d3-a456-426614274006', '123e4567-e89b-12d3-a456-426614370006', '2024-05-10', '13:00:00', 30, 'key128', '2024-05-10 12:30:00', TRUE, '2024-05-09 13:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174007', '123e4567-e89b-12d3-a456-426614274007', '123e4567-e89b-12d3-a456-426614370007', '2024-05-10', '14:00:00', 30, 'key129', '2024-05-10 13:30:00', TRUE, '2024-05-09 14:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174008', '123e4567-e89b-12d3-a456-426614274008', '123e4567-e89b-12d3-a456-426614370008', '2024-05-10', '15:00:00', 30, 'key130', '2024-05-10 14:30:00', TRUE, '2024-05-09 15:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174009', '123e4567-e89b-12d3-a456-426614274009', '123e4567-e89b-12d3-a456-426614370009', '2024-05-10', '16:00:00', 30, 'key131', '2024-05-10 15:30:00', TRUE, '2024-05-09 16:00:00', NULL, NULL),
('123e4567-e89b-12d3-a456-426614174010', '123e4567-e89b-12d3-a456-426614274010', '123e4567-e89b-12d3-a456-426614370010', '2024-05-10', '17:00:00', 30, 'key132', '2024-05-10 16:30:00', TRUE, '2024-05-09 17:00:00', NULL, NULL);