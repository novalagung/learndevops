# AWS | Membuat User IAM

Pada tutorial ini, kita akan belajar cara membuat user IAM.

---

### 1. Definisi

IAM atau Identity and Access Management adalah salah satu fasilitas yang disediakan oleh AWS untuk mengontrol akses dari user individu atau group terhadap *resource* AWS. 

Kita bisa membuat banyak user IAM, dan juga memberikan ijin akses terhadap *resource* AWS ke user-user tersebut.

User IAM ini mirip seperti akun AWS, perbedaannya hanya pada user IAM ijin terhadap *resource* yang ada adalah dikontrol. Jadi user IAM tidak bisa secara leluasa mengakses *resource*, kecuali untuk *resources* yang diijinkan untuk diakses oleh user tersebut.

---

### 2. Membuat user IAM baru

Pertama login ke AWS console terlebih dahulu menggunakan akun AWS, lalu buka menu [AWS IAM page](https://console.aws.amazon.com/iam/home?region=ap-southeast-1#/home), klik menu **manage users**.

![AWS | Membuat user IAM baru](https://i.imgur.com/yx8dVAR.png)

Akan muncul halaman yang menampilkan semua user IAM. Selanjutnya, klik **Add user**, lalu isi namanya.

Jika user yang akan dibuat akan dipergunakan untuk manajemen resource menggunakan aplikasi *3rd party* ataupun AWS SDK, maka pastikan opsi **Programmatic Access** tercentang.

![AWS | Membuat user IAM baru - programmatic access](https://i.imgur.com/2V7shR9.png)

Kemudian klik *next* untuk memunculkan halaman user group. Pada laman ini buat group baru dengan ijin akses disesuaikan dengan kebutuhan. Sebagai contoh pada gambar berikut, group baru dibuat dengan nama `user` dengan ijin diberikan adalah akses penuh terhadap semua resource EC2.

![AWS | Membuat user IAM baru - iam user group](https://i.imgur.com/l46C9OQ.png)

Lalu klik tombol *next* beberapa kali hingga proses selesai.

![AWS | Membuat user IAM baru - access key and secret key](https://i.imgur.com/oqAAWZv.png)

Akan muncul informasi **access key ID** dan **secret access key** milik user yang baru saja dibuat. *Copy* informasi tersebut dan simpan, karena informasi ini tidak akan muncul lagi.

Ok, hanya itu saja. Keys bisa dipergunakan.
