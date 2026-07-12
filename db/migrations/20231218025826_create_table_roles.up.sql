
CREATE TABLE roles
(
    id   int  NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(100) not null,
    PRIMARY KEY (id)
) ENGINE = InnoDB;

insert into roles(id, name, description) VALUE (2, 'customers', 'customers');
insert into roles(id, name, description) VALUE (1, 'admin', 'Admin');
insert into roles(id, name, description) VALUE (99, 'superadmin', 'Admin');
insert into roles(id, name, description) VALUE (50, 'nonaktif', 'user not aktif');
insert into roles(id, name, description) VALUE (3, 'technician', 'readuser aktif');
insert into roles(id, name, description) VALUE (4, 'driver', 'readuser aktif');
insert into roles(id, name, description) VALUE (5, 'karyawan', 'readuser aktif');
insert into roles(id, name, description)
    VALUE (8, 'manager', 'manage team and approve tasks');

insert into roles(id, name, description)
    VALUE (9, 'supervisor', 'monitor operational activities');

insert into roles(id, name, description)
    VALUE (10, 'support', 'handle customer support and tickets');
insert into roles(id, name, description)
    VALUE (16, 'finance', 'financial and billing access');

insert into roles(id, name, description)
    VALUE (17, 'hr', 'human resource management');

