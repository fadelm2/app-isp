ALTER TABLE users ADD COLUMN region_id VARCHAR(50) NULL;

CREATE TABLE IF NOT EXISTS internet_packages (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    speed_mbps INT NOT NULL,
    price DOUBLE NOT NULL,
    installation_fee DOUBLE NOT NULL,
    tax_rate DOUBLE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS registrations (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    nik VARCHAR(50) NOT NULL,
    birth_place VARCHAR(100) NOT NULL,
    birth_date VARCHAR(50) NOT NULL,
    gender VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    installation_address TEXT NOT NULL,
    billing_address TEXT NOT NULL,
    package_id VARCHAR(100) NOT NULL,
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,
    notes TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    ktp_path VARCHAR(555) NOT NULL,
    selfie_path VARCHAR(555),
    house_path VARCHAR(555),
    installation_path VARCHAR(555),
    supporting_doc_path VARCHAR(555),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    FOREIGN KEY (package_id) REFERENCES internet_packages(id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS routers (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    username VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'offline',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS customers (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    registration_id VARCHAR(100),
    user_id VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    package_id VARCHAR(100) NOT NULL,
    router_id VARCHAR(100),
    ppp_username VARCHAR(100),
    ppp_password VARCHAR(100),
    radius_username VARCHAR(100),
    radius_password VARCHAR(100),
    due_date_day INT NOT NULL DEFAULT 10,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (package_id) REFERENCES internet_packages(id),
    FOREIGN KEY (router_id) REFERENCES routers(id) ON DELETE SET NULL,
    FOREIGN KEY (registration_id) REFERENCES registrations(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    customer_id VARCHAR(100) NOT NULL,
    due_date BIGINT NOT NULL,
    period_month INT NOT NULL,
    period_year INT NOT NULL,
    amount DOUBLE NOT NULL,
    tax_amount DOUBLE NOT NULL,
    installation_fee DOUBLE NOT NULL,
    total_amount DOUBLE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    snap_token VARCHAR(255),
    paid_at BIGINT,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    invoice_id VARCHAR(100) NOT NULL,
    transaction_id VARCHAR(255) NOT NULL,
    payment_type VARCHAR(100) NOT NULL,
    paid_amount DOUBLE NOT NULL,
    status VARCHAR(100) NOT NULL,
    paid_at BIGINT NOT NULL,
    raw_response TEXT,
    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS customer_histories (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    customer_id VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    notes TEXT,
    created_by VARCHAR(100) NOT NULL,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB;

-- FreeRADIUS tables simulation/integration
CREATE TABLE IF NOT EXISTS radcheck (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL DEFAULT '',
    attribute VARCHAR(64) NOT NULL DEFAULT '',
    op VARCHAR(2) NOT NULL DEFAULT '==',
    value VARCHAR(253) NOT NULL DEFAULT '',
    KEY username (username(32))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS radreply (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL DEFAULT '',
    attribute VARCHAR(64) NOT NULL DEFAULT '',
    op VARCHAR(2) NOT NULL DEFAULT '=',
    value VARCHAR(253) NOT NULL DEFAULT '',
    KEY username (username(32))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS radacct (
    radacctid BIGINT AUTO_INCREMENT PRIMARY KEY,
    acctsessionid VARCHAR(64) NOT NULL DEFAULT '',
    acctuniqueid VARCHAR(32) NOT NULL DEFAULT '',
    username VARCHAR(64) NOT NULL DEFAULT '',
    groupname VARCHAR(64) NOT NULL DEFAULT '',
    realm VARCHAR(64) DEFAULT '',
    nasipaddress VARCHAR(15) NOT NULL DEFAULT '',
    nasportid VARCHAR(15) DEFAULT NULL,
    nasporttype VARCHAR(32) DEFAULT NULL,
    acctstarttime DATETIME DEFAULT NULL,
    acctupdatetime DATETIME DEFAULT NULL,
    acctstoptime DATETIME DEFAULT NULL,
    acctinterval INT DEFAULT NULL,
    acctsessiontime INT UNSIGNED DEFAULT NULL,
    acctauthentic VARCHAR(32) DEFAULT NULL,
    connectinfo_start VARCHAR(50) DEFAULT NULL,
    connectinfo_stop VARCHAR(50) DEFAULT NULL,
    acctinputoctets BIGINT DEFAULT NULL,
    acctoutputoctets BIGINT DEFAULT NULL,
    calledstationid VARCHAR(50) NOT NULL DEFAULT '',
    callingstationid VARCHAR(50) NOT NULL DEFAULT '',
    acctterminatecause VARCHAR(32) NOT NULL DEFAULT '',
    framedipaddress VARCHAR(15) NOT NULL DEFAULT '',
    acctstartdelay INT DEFAULT NULL,
    acctstopdelay INT DEFAULT NULL,
    xascendsessionsvrkey VARCHAR(10) DEFAULT NULL,
    KEY username (username),
    KEY acctsessionid (acctsessionid),
    KEY acctuniqueid (acctuniqueid),
    KEY acctstarttime (acctstarttime),
    KEY acctstoptime (acctstoptime),
    KEY nasipaddress (nasipaddress)
) ENGINE=InnoDB;
