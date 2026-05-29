## INPUT

Isi bagian ini sebelum mengirim prompt:

```
MODULE_NAME   : {isi nama Go module, contoh: my-awesome-api}
DOMAIN        : {isi nama domain, contoh: auth, user, product, order}
FEATURE       : {isi nama fitur, contoh: register, profile, create-product}
DESCRIPTION   : {isi deskripsi singkat fitur ini}
ENDPOINT      : {isi HTTP method dan path, contoh: POST /auth/register}
MIDDLEWARE    : {isi: none / protected / optional}

REQUEST FIELDS:
  - {nama_field}: {tipe} | {aturan validasi} | {keterangan}
  - contoh → email: string | required,email | Email pengguna
  - contoh → password: string | required,min=6 | Password minimal 6 karakter
  - contoh → name: string | required,max=100 | Nama lengkap

RESPONSE FIELDS:
  - {nama_field}: {tipe} | {keterangan}
  - contoh → token: string | JWT token
  - contoh → user: object | Data user yang dibuat

BUSINESS LOGIC:
  {jelaskan alur logic bisnis secara singkat}
  contoh:
  1. Cek apakah email sudah terdaftar
  2. Hash password
  3. Simpan user baru
  4. Assign role default
  5. Return data user

DB OPERATIONS (untuk repository):
  - {nama_method}: {operasi} | {keterangan}
  - contoh → FindByEmail(email): SELECT | Cari user berdasarkan email
  - contoh → Create(user): INSERT | Simpan user baru
  - contoh → FindDefaultRole(): SELECT | Ambil role dengan is_default=true

ADDITIONAL NOTES:
  {catatan tambahan jika ada, atau tulis "none"}
```

---

## CONTOH INPUT YANG SUDAH DIISI

Berikut contoh input untuk fitur register:

```
MODULE_NAME   : my-awesome-api
DOMAIN        : auth
FEATURE       : register
DESCRIPTION   : Fitur registrasi user baru dengan auto-assign role default
ENDPOINT      : POST /auth/register
MIDDLEWARE    : none

REQUEST FIELDS:
  - email: string | required,email | Email pengguna
  - password: string | required,min=6 | Password minimal 6 karakter
  - name: string | required,max=100 | Nama lengkap pengguna

RESPONSE FIELDS:
  - id: uint | ID user yang baru dibuat
  - email: string | Email user
  - name: string | Nama user

BUSINESS LOGIC:
  1. Cek apakah email sudah terdaftar, jika ada return error
  2. Hash password menggunakan util.HashPassword
  3. Simpan user baru ke database
  4. Ambil role default (is_default=true)
  5. Assign role default ke user via UserRole
  6. Return data user tanpa password

DB OPERATIONS:
  - FindByEmail(email): SELECT | Cek apakah email sudah ada
  - Create(user): INSERT | Simpan user baru
  - FindDefaultRole(): SELECT | Ambil role dengan is_default=true
  - CreateUserRole(userRole): INSERT | Assign role ke user

ADDITIONAL NOTES:
  - Jika tidak ada role default di DB, tetap buat user tanpa role
  - Password tidak boleh ada di response
```

---

## MULAI

Setelah mengisi bagian INPUT di atas, kirim prompt ini ke AI.
AI akan mulai dengan file **[1/9]** dan menunggu perintah `next` untuk melanjutkan.