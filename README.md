# ISP Management System

## Description

ISP Management System adalah aplikasi OSS/BSS untuk mengelola operasional Internet Service Provider (ISP). Sistem ini mendukung manajemen pelanggan, integrasi MikroTik & FreeRADIUS, pembayaran otomatis menggunakan Midtrans, serta proses suspend dan reaktivasi layanan internet secara otomatis berdasarkan status pembayaran.

Project ini dibangun menggunakan **Golang Clean Architecture** sehingga mudah dikembangkan, diuji, dan dipelihara.

---

# Features

## Authentication

* Login
* Logout
* Refresh Token
* JWT Authentication
* Role Based Access Control (RBAC)

---

## User Management

* User Management
* Role Management
* Permission Management
* Audit Log

---

## Customer Registration

Pelanggan dapat melakukan pendaftaran secara online.

### Registration Form

* Nama Lengkap
* NIK
* Tempat Lahir
* Tanggal Lahir
* Jenis Kelamin
* Email
* Nomor HP
* Alamat Instalasi
* Alamat Penagihan
* Paket Internet
* Google Maps Location
* Catatan

### Upload Document

* Foto KTP
* Selfie dengan KTP (Opsional)
* Foto Rumah (Opsional)
* Foto Lokasi Instalasi
* Dokumen Pendukung
* PDF

Supported Format

* JPG
* JPEG
* PNG
* PDF

---

## Registration Workflow

```text
Customer

↓

Isi Form Registrasi

↓

Upload Dokumen

↓

Submit

↓

Admin Review

↓

Survey

↓

Approve

↓

Generate Customer

↓

Generate Subscription

↓

Generate Radius User

↓

Installasi

↓

Internet Aktif
```

---

## Customer Management

* Create Customer
* Update Customer
* Suspend Customer
* Unsuspend Customer
* Terminate Customer
* Customer History
* Customer Notes

---

## Internet Package

* Package Management
* Internet Speed
* Price
* Installation Fee
* Tax
* Active / Non Active Package

---

## MikroTik Integration

* Router Management
* Router Monitoring
* PPP Secret Management
* Queue Management
* Disconnect Active Session
* Router API Integration

---

## FreeRadius Integration

* Radius User
* Radius Authentication
* Radius Accounting
* Active Session
* Session History

---

## Billing

* Generate Monthly Invoice
* Invoice History
* Invoice Status
* Due Date
* Payment History

Invoice Status

* Pending
* Paid
* Owed
* Expired
* Cancelled

---

## Midtrans Payment

Customer dapat membayar melalui Midtrans.

Flow

```text
Generate Invoice

↓

Customer Bayar

↓

Midtrans

↓

Webhook

↓

Verifikasi Signature

↓

Invoice PAID

↓

Customer ACTIVE

↓

Enable Radius User

↓

Enable PPP Secret

↓

Internet Aktif
```

---

## Automatic Suspension

Scheduler akan melakukan pengecekan invoice setiap hari.

Jika invoice melewati jatuh tempo:

```text
Invoice Owed

↓

Suspend Customer

↓

Disable Radius User

↓

Disable PPP Secret

↓

Disconnect Active Session

↓

Internet Mati
```

---

## Automatic Reactivation

Setelah pembayaran diterima:

```text
Webhook Midtrans

↓

Invoice Paid

↓

Enable Radius User

↓

Enable PPP Secret

↓

Customer Active

↓

Internet Hidup
```

---

## Dashboard

* Total Customer
* Active Customer
* Suspended Customer
* Owed Customer
* Today's Payment
* Monthly Revenue
* Router Status
* Online User
* Offline User

---

## Notification

* Email
* WhatsApp
* Telegram

---

# Clean Architecture

![Clean Architecture](architecture.png)

### Request Flow

1. Client mengirim request (HTTP, gRPC, Worker, Scheduler, Messaging)
2. Delivery melakukan parsing request menjadi DTO
3. Delivery memanggil Use Case
4. Use Case menjalankan business logic
5. Use Case membentuk Entity
6. Repository melakukan operasi database
7. Database menyimpan atau mengambil data
8. Use Case memanggil Gateway jika diperlukan
9. Gateway melakukan komunikasi ke sistem eksternal
10. Gateway mengembalikan hasil ke Use Case
11. Delivery mengembalikan Response ke Client

---

# Project Structure

```text
cmd/

├── api/
├── worker/
├── scheduler/

internal/

├── auth/
├── user/
├── role/
├── permission/
├── customer/
├── registration/
├── package/
├── invoice/
├── payment/
├── radius/
├── mikrotik/
├── router/
├── dashboard/
├── notification/
├── webhook/
├── scheduler/
├── middleware/
├── helper/
├── validator/
├── logger/

db/

├── migrations/

configs/

storage/

├── uploads/
│   ├── ktp/
│   ├── selfie/
│   ├── installation/
│   ├── documents/
│   └── pdf/

docs/

api/
```

---

# Tech Stack

## Backend

* Golang 1.25+
* Go Fiber
* PostgreSQL
* GORM
* Redis
* JWT
* Viper
* Logrus
* Validator

## Database

* PostgreSQL

## Payment Gateway

* Midtrans

## Authentication

* JWT

## Radius

* FreeRadius

## Router

* MikroTik RouterOS API

## Object Storage

* Local Storage
* MinIO (Optional)
* Amazon S3 (Optional)

---

# Framework & Library

* Fiber
* GORM
* Viper
* Golang Migrate
* Validator
* Logrus
* JWT
* Redis Client
* Midtrans SDK
* MikroTik API Client
* FreeRadius

---

# Configuration

Semua konfigurasi berada pada file:

```text
configs/config.json
```

Contoh konfigurasi:

```json
{
  "app": {
    "name": "ISP Management",
    "port": 8080
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "postgres",
    "password": "password",
    "database": "isp_management"
  },
  "midtrans": {
    "server_key": "",
    "client_key": "",
    "is_production": false
  },
  "mikrotik": {
    "host": "",
    "username": "",
    "password": ""
  }
}
```

---

# API Documentation

Seluruh dokumentasi API berada pada folder:

```text
api/
```

---

# Database Migration

Migration berada pada folder:

```text
db/migrations
```

## Create Migration

```bash
migrate create -ext sql -dir db/migrations create_table_customer
```

## Run Migration

```bash
migrate \
-database "postgres://postgres:password@localhost:5432/isp_management?sslmode=disable" \
-path db/migrations up
```

## Rollback

```bash
migrate \
-database "postgres://postgres:password@localhost:5432/isp_management?sslmode=disable" \
-path db/migrations down
```

---

# Scheduler

Scheduler berjalan otomatis untuk:

* Generate Monthly Invoice
* Suspend Owed Customer
* Reactivate Paid Customer
* Disconnect Active Session
* Reminder Payment
* Backup Router
* Radius Synchronization

---

# Storage

```text
storage/

uploads/

├── ktp/
├── selfie/
├── house/
├── installation/
├── documents/
└── pdf/
```

---

# Run Application

## Install Dependency

```bash
go mod tidy
```

## Run API

```bash
go run cmd/api/main.go
```

## Run Worker

```bash
go run cmd/worker/main.go
```

## Run Scheduler

```bash
go run cmd/scheduler/main.go
```

---

# Unit Test

```bash
go test -v ./...
```

---

# Future Roadmap

* Customer Portal
* Mobile Apps
* Multi Branch ISP
* Multi Tenant
* CRM
* Ticketing System
* WhatsApp Gateway
* Telegram Bot
* ONU Management
* OLT Management
* PPPoE Monitoring
* SNMP Monitoring
* Radius Dashboard
* Mikrotik Monitoring
* AI Customer Support
* AI Revenue Analytics
* AI Network Monitoring

---

# License

MIT License.
