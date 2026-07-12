CREATE TABLE users (
                       id           VARCHAR(100) NOT NULL PRIMARY KEY,
                       role_id      INT NOT NULL,
                       username     VARCHAR(255) NOT NULL,
                       email        VARCHAR(255) NOT NULL,
                       password     VARCHAR(255) NOT NULL,
                       company_name VARCHAR(255) NOT NULL,
                       updated_at   BIGINT NOT NULL,
                       created_at   BIGINT NOT NULL
) ENGINE=InnoDB;

INSERT INTO users (id, role_id, username, email, password, company_name, updated_at, created_at) VALUES
                                                                                                     ('noc1', 1, 'admin_user', 'admin@greenet.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
                                                                                                     ('admin1', 1, 'admin_user', 'user@greenet.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
                                                                                                     ('teknisi1', 3, 'teknisi', 'mod@greenet.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
                                                                                                     ('driver1', 4, 'driver', 'mod@greenet.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
                                                                                                     ('superadmin', 99, 'superadmin', 'user@greenet.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
                                                                                                    ('customers', 2, 'admin_user', 'user@greenet.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
