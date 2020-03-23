# Docker - Push Image ke hub.docker.com

Pada post kali ini, kita akan belajar tentang cara *push* image Docker ke [Docker Hub](https://hub.docker.com/).

---

### 1. Kebutuhan

#### 1.1. Docker engine

Pastikan Docker engine adalah *running*. Jika di lokal belum ter-*install* Docker engine, maka *install* terlebih dahulu.

#### 1.2. Akun Docker Hub

Siapkan akun Docker Hub. Jika belum punya, maka buat terlebih dahulu. Bisa ikuti petunjuk [Membuat Akun Docker Hub](docker-hub-membuat-akun.html).

#### 1.3. Login ke Docker Hub di lokal

Untuk login, gunakan command CLI berikut:

```bash
docker login --username=novalagung --password=<your-password>
```

Atau bisa juga menggunakan menu login yang tersedia. Klik kanan pada ikon Docker, lalu pilih *sign in*.

---

### 2. Panduan

#### 2.1. Buat repo di Docker Hub

Sebelum memulai proses, pastinya kita perlu mem-*booking* sebuah repo di Docker Hub. Nantinya kita akan push image docker ke repo ini.

Buka https://hub.docker.com/repository/create, buat repo baru repo, namai dengan `hello-world` (atau lainnya, bebas).

![Docker - Push Image ke hub.docker.com - buat repo baru di Docker Hub](https://i.imgur.com/uvLjxqv.png)

#### 2.2. Clone aplikasi example, lalu build sebagai image Docker

Selanjutnya, kita perlu mempersiapkan sebuah aplikasi hello world yang sudah *dockerized*. Tapi untuk mempercepat tutorial, kita akan gunakan sebuah aplikasi hello world yang sudah siap pakai yang dikembangkan menggunakan bahasa Go berikut. Aplikasinya bisa di pull dari repo Github menggunakan token (karena *visibility* repo nya adalah *private*).

```bash
git clone https://30542dd8874ba3745c55203a091c345340c18b7a:x-oauth-basic@github.com/novalagung/hello-world.git
```

Setelah proses *cloning* selesai, lakukan proses *build* ke bentuk image Docker dengan nama mengikuti format `<your-docker-username>/<your-repo-name>:<tag-name>`. Silakan sesuaikan *value* dari `<your-docker-username>` untuk menggunakan user Docker Hub masing-masing.

```bash
cd hello-world

# docker build . -t <username>/<repo-name>:<tag>
docker build . -t novalagung/hello-world:v0
```

Bisa dilihat dari command di atas, tag `v0` digunakan pada image ini.

![Docker - Push Image ke hub.docker.com - build image](https://i.imgur.com/aiduEji.png)

#### 2.3. Push image ke Docker Hub

Selanjutnya, gunakan command `docker push` untuk *push* image (yang barusan kita `build`) ke Docker Hub.

```
# docker push <username>/<repo-name>[:<tag>]
docker push novalagung/hello-world
```

![Docker - Push Image ke hub.docker.com - push image ke Docker Hub](https://i.imgur.com/TUy6Ffa.png)

Ok, selesai.

---

### 3. Test - Pull image dari Docker Hub

> Bagian ini adalah opsional.

Sebelumnya kita telah belajar tentang cara *push* image ke Docker Hub. Untuk cara *pull*, gunakan command `docker pull`.

```bash
# docker pull <username>/<repo-name>[:<tag>]
docker pull novalagung/hello-world:v0
```

![Docker - Push Image ke hub.docker.com - pull image dari Docker Hub](https://i.imgur.com/tdRlNr7.png)

---

### 4. Tag `latest`

Secara *default*, ketika kita pull sebuah image dari Hub tanpa menuliskan spesifik tag-nya, maka tag `latest` adalah yang akan di-*pull*.

Silakan perhatikan dua command berikut. Keduanya adalah ekuivalen.

```
docker pull novalagung/hello-world
docker pull novalagung/hello-world:latest
```

Yang menarik dan penting untuk diketahui dari tag `latest` ini, bahwa tag ini tidak mengarah ke tag yang terakhir di push, melainkan mengarah ke tag yang namanya secara eksplisit adalah `latest`.

Tag sebelumnya yang sudah kita push, `v0`, tidak adakan dianggap sebagai tag `latest`. Untuk membuat tag latest, maka kita perlu *rebuild* ulang aplikasi yang sama sebagai tag `latest`, lalu push ke Docker Hub.

```bash
cd hello-world
docker build . -t novalagung/hello-world:latest
docker push novalagung/hello-world:latest
```

![Docker - Push Image ke hub.docker.com - push tag latest ke Docker Hub](https://i.imgur.com/6y0MEEA.png)
