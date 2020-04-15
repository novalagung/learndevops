# Kubernetes - Deploy Aplikasi ke Kluster Minikube menggunakan Deployment controller, Service, dan Horizontal Pod Autoscaler

Pada tutorial kali ini, kita akan belajar cara men-*deploy* aplikasi *containerized* ke kluster kubernetes (minikube), meng-*enable* kapabilitas autoscaling, dan menyiapkan *service* untuk membuat aplikasi tersebut menjadi dapat diakses dari luar kluster.

Aplikasi yang akan di-deploy adalah aplikasi hello world sederhana dikembangkan menggunakan bahasa Go. Aplikasinya sudah jadi dan siap dipakai, sudah *dockerized* juga, dan image nya tersedia di [Docker Hub](https://hub.docker.com/repository/docker/novalagung/hello-world).

Jika teman-teman ingin men-deploy aplikasi lain, silakan. Siapkan saja image nya dan push ke Docker Hub. Bisa mengikuti panduan berikut untuk caranya [Push Image ke hub.docker.com](/docker-push-image-to-hub.html).

---

### 1. Kebutuhan

#### 1.1. Docker engine

Pastikan Docker engine adalah *running*. Jika di lokal belum ter-*install* Docker engine, maka *install* terlebih dahulu.

#### 1.2. Minikube

Pastikan Minikube sudah *running*. Jalankan command `minikube start` pada PowerShell yang dibuka dengan admin privilege. Jika belum ter-*install*, maka silakan *install* terlebih dahulu.

#### 1.3. Kubernetes CLI tool

Pastikan command `kubectl` sudah di-*install*. Jika belum ter-*install*, maka silakan *install* terlebih dahulu.

#### 1.4. `hey` tool (HTTP load generator)

*Install* tool ini di lokal. Panduan instalasi ada di https://github.com/rakyll/hey.

`hey` merupakan *benchmark tool*, mirip seperti Apache Benchmark. Nantinya kita akan menggunakan tool kecil ini untuk melakukan *stress test* pada aplikasi kita, untuk memastikan kapabilitas auto scaling berjalan sesuai harapan.

---

### 2. Persiapan

#### 2.1. Khusus untuk pengguna Windows, gunakan PowerShell yang dibuka dengan admin privilege

Tidak dianjurkan menggunakan Command Prompt atau CMD. Gunakan PowerShell yang dibuka dengan admin privilege.

#### 2.2. Buat file konfigurasi objek Kubernetes (dengan ekstensi `.yaml`)

Nantinya kita akan buat tiga buah objek Kubernetes, yaitu: deployment, horizontal pod auto scaler, dan service.

Untuk mempermudah proses pembelajaran, ketiga konfigurasi objek yang disebutkan di atas akan di tulis dalam satu file konfigurasi saja.

Ok, langsung saja, buat file bernama `k8s.yaml` (atau bisa gunakan nama lain, bebas). Buka file tersebut lewat editor favorit teman-teman.

---

### 3. Pendefinisian Objek

#### 3.1. Objek Deployment

Deployment adalah sebuah objek *controller* yang digunakan untuk manajemen pod dan replica sets. Pada bagian ini kita akan buat objek tersebut.

Pada file `k8s.yaml`, tulis konfigurasi di bawah ini. Silakan dibaca dan dipahami juga komentar yang ada pada bagian-bagian konfigurasi untuk tau kegunaan dari masing-masing bagian.

```yaml
---
# ada banyak API yang tersedia pada Kubernetes
# (untuk cek, bisa gunakan command `kubectl api-versions`).
# pada blok kode deployment berikut, kita pilih API `apps/v1`.
apiVersion: apps/v1

# alokasikan blok kode YAML ini untuk deployment object.
kind: Deployment

# namai objek dengan `my-app-deployment`.
metadata:
  name: my-app-deployment

# blok spec berikut berisi spesifikasi behaviour dari objek deployment ini.
spec:

  # selector.matchLabels digunakan untuk menentukan pod mana yg akan di-manage oleh deployment.
  # deployment akan me-manage semua pod yang sesuai dengan selektor ini.
  selector:
    matchLabels:
      app: my-app

  # blok kode template berikut berisi informasi bagaiaman pod akan dibuat.
  template:

    # tambahkan label pada pod dengan value adalah `my-app`.
    metadata:
      labels:
        app: my-app

    # blok spec berikut berisi spesifikasi behaviour dari objek pod `my-app`.
    spec:

      # list dari kontainer dalam pod `my-app`.
      containers:

          # alokasikan kontainer baru dengan nama `hello-world`.
        - name: hello-world

          # image dari container ini di-pull dari Docker Hub repo `novalagung/hello-world`.
          # jika image tersebut belum ada di lokal, maka akan otomatis di pull.
          image: novalagung/hello-world

          # beberapa environment variables yang dibutuhkan container ini.
          # lebih jelasnya bisa cek ke
          # https://hub.docker.com/repository/docker/novalagung/hello-world.
          env:
            - name: PORT
              value: "8081"
            - name: INSTANCE_ID
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name

          # pod ini hanya akan memiliki satu kontainer saja, yaitu `hello-world` di atas,
          # dan yang akan terjadi dalam kontainer ini adalah: sebuah web server di-start
          # pada port `8081`.
          # port tersebut di-expose agar bisa diakses dari luar pod dalam kluster.
          ports:
            - containerPort: 8081

          # spec container `hello-world`.
          resources:
            limits:
              cpu: 250m
              memory: 32Mi
```

Secara garis besar, blok kode di atas melakukan proses berikut (secara sekuensial):

- Membuat objek deployment bernama `my-app-deployment`.
- Mendefinisikan spesifikasi pod (dalam objek deployment) yang berisi satu buah kontainer.
- Kontainer tersebut bernama `hello-world`, image-nya di-pull dari Docker Hub.
- Environment informasi PORT dan Instance ID di-definisikan untuk keperluan proses build container. Port tersebut digunakan oleh web server dalam kontainer.
- Web server akan *listen to* port `8081` dan port ini di-expose ke luar. Artinya kita bisa mengakses web server tersebut dari luar pod (tapi masih dari dalam kluster).

Ok, sekarang mari kita terapkan file konfigurasi di atas menggunakan command berikut:

```bash
# terapkan konfigurasi
kubectl apply -f k8s.yaml

# munculkan semua objek deployment
kubectl get deployments

# munculkan semua objek pod
kubectl get pods
```

<p style="text-align: center;">
  <img src="https://i.imgur.com/VXlFDch.png" alt="Kubernetes - Deploy Aplikasi ke Kluster Minikube menggunakan Deployment controller, Service, dan Horizontal Autoscaler - run objek deployment">
</p>

#### 3.1. Test salah satu pod

Pada gambar di atas bisa dilihat, deployment berjalan sesuai harapan. Ada dua buah pod yang *running*.

> Kenapa ada dua pod? kenapa tidak satu, tiga, atau lainnya. Hal ini karena kita tidak mendefinisikan `spec.replicas` objek deployment. Nilai default replikai adalah dua. Jika kita mendefinisikan replika dengan nilai empat (misalnya), maka akan langsung ter-create 4 buah pod secara default.

Ok, mari kita coba lakukan testing. Kita akan coba *remotelly* konek ke salah satu pod, lalu mengecek apakah web server sudah *listen to* port `8081`.

```bash
# munculkan semua objek pod
kubectl get pods

# konek ke salah satu pod
kubectl exec -it <pod-name> -- /bin/sh

# cek apakah aplikasi/webserver yang listen ke port 8081
netstat -tulpn | grep :8081
```

<p style="text-align: center;">
  <img src="https://i.imgur.com/vdZaLf2.png" alt="Kubernetes - Deploy Aplikasi ke Kluster Minikube menggunakan Deployment controller, Service, dan Horizontal Autoscaler - konek ke pod">
</p>

Bisa dilihat pada image di atas, port tersebut digunakan oleh aplikasi/webserver.

#### 3.2. Terapkan perubahan pada objek deployment

Selain objek deployment, ada juga jenis objek kontroler lainnya yang tersedia di k8s.

Salah satu kelebihan dari objek deployment dibanding objek kontroler lainnya adalah, setiap perubahan konfigurasi pod yang diterapkan, maka perubahan tersebut akan diterapkan secara *seamless*.

Ok, sekarang mari kita buktikan bahwa statement di atas adalah valid, yaitu dengan cara mengubah beberapa konfigurasi berikut lalu kemudian menerapkannya.

- Perubahan pada `containers.env.value`, arakhan nilai environment `PORT` ke `8080`. Sebelumnya mengarah ke `8081`.
- Perubahan pada `containers.ports.containerPort`, arahkan nilainya ke `8080`. Sebelumnya mengarah ke `8081`.

Di bawah ini adalah tampilan dari konfigurasi deployment setelah perubahan di atas diterapkan.

```bash
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-deployment
spec:
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: hello-world
          image: novalagung/hello-world
          env:
            - name: PORT
              value: "8080" # <--- ubah dari 8081 ke 8080
            - name: INSTANCE_ID
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
          ports:
            - containerPort: 8080 # <--- ubah dari 8081 ke 8080
          resources:
            limits:
              cpu: 250m
              memory: 32Mi
```

Ok, sekarang terapkan konfigurasi yang baru di atas.

```bash
# terapkan konfigurasi
kubectl apply -f k8s.yaml

# munculkan semua pod
kubectl get pods

# konek ke salah satu pod
kubectl exec -it <pod-name> -- /bin/sh

# cek apakah ada yang listen ke port 8080
netstat -tulpn | grep :8080
```

<p style="text-align: center;">
  <img src="https://i.imgur.com/DZWCTSk.png" alt="Kubernetes - Deploy Aplikasi ke Kluster Minikube menggunakan Deployment controller, Service, dan Horizontal Autoscaler - terapkan perubahan objek deployment">
</p>

Bisa dilihat, perubahan yang kita lakukan pada pod diaplikasikan secara alusss oleh kontroler deployment. Dan web server sekarang listen ke port 8080.

> Tips! Gunakan command berikut untuk menampilkan error log pada pod. Ini cukup berguna ketika pod container dalam pod tidak berjalan seperti rencana. Dari sini kita bisa tau error log yang menjelaskan error.

> `kubectl get pods`<br />`kubectl describe pod <pod-name>`<br />`kubectl logs <pod-name>`

#### 3.2. Objek Service

Service digunakan untuk menjembatani akses antar pod, baik dalam cluster maupun dari luar cluster.

Pada bagian ini kita akan buat service baru, untuk meng-enable akses dari luar kluster ke dalam kluster, ke pod tujuan.

Ok, sekarang tambahakn konfigurasi berikut ke file `k8s.yaml`.

```bash
---
# pilih API `v1` untuk blok konfigurasi service.
apiVersion: v1

# alokasikan blok YAML berikut untuk service.
kind: Service

# namai dengan `my-service`.
metadata:
  name: my-service

# blok spec berikut berisi spesifikasi behaviour dari objek service ini.
spec:

  # tipe dari service yang dipilih adalah LoadBalancer.
  #
  # LoadBalancer service adalah pilihan standar untuk meng-expose sebuah service ke public.
  # Akan disiapkan sebuah network load balancer dengan satu buah public IP, dan nantinya
  # semua request yang masuk ke IP tersebut akan diproses oleh service.
  #
  # pada provider cloud, nantinya akan digenerate satu buah public IP.
  # tapi pada lokal (minikube), service akan bisa diakses lewat IP minikube.
  type: LoadBalancer

  # route trafik pada service ke pod dengan label berikut.
  selector:
    app: my-app

  # list port yang di expose di service ini.
  ports:

      # expose service port ke luar kluster (dalam konteks ini ke luar minikube).
      # jadi untuk mengakses service, gunakan public IP + nodePort berikut.
      #
      # cara untuk menampilkan exposed URL (public IP + nodePort):
      # `minikube service my-service --url`.
      #   => http://<cluster-public-ip>:<nodePort>
    - nodePort: 32199

      # request yang masuk dari luar ke nodePort, akan diarahkan ke port 80 milik service ini.
      #
      # cara untuk menampilkan service URL (dalam kluster):
      # `kubectl describe service my-service | findstr "IP"`.
      #   => http://<service-ip>:<port>
      port: 80

      # setelah itu, dari port 80 milik service akan di arahkan ke pod yang tersedia,
      # menggunakan algoritma round-robin (load balancer),
      # ke pod IP dengan port targetPort berikut.
      #
      #  => http://<pod-ip>:<targetPort>
      targetPort: 8080
```

`LoadBalancer` dipilih sebagai tipe service ini. Load balancer nantinya akan menerima request dari luar kluster (pada `<publicIP>:<nodePort>`), untuk diarahkan ke port `80` milik service dalam kluster, lalu di-teruskan ke pod (ke `<pod>:<targetPort>`) menggunakan algoritma umum load balancer (round-robin).

> Sebenarnya dalam contoh ini kita tidak harus menggunakan tipe `LoadBalancer`, tipe `NodePort` juga bisa dipergunakan.

Satu point penting yang perlu dicatat disini, karena kluster kita adalah Minikube, maka public IP disini adalah public IP dari Minikube. Untuk menampilkan IP Minikube, bisa gunakan command berikut:

```bash
# tampilkan minikube IP
minikube ip
```

Ok, sekarang mari kita terapkan konfigurasi `k8s.yaml` yang sudah ditambahkan service didalamnya.

```bash
# terapkan konfigurasi
kubectl apply -f k8s.yaml

# munculkan semua objek service
kubectl get services

# munculkan semua pod
kubectl get pods

# test service menggunakan `curl`
curl <minikubeIP>:<nodePort>
curl <minikubeIP>:32199
```

<p style="text-align: center;">
  <img src="https://i.imgur.com/IoEpMFH.jpg" alt="Kubernetes - Deploy Aplikasi ke Kluster Minikube menggunakan Deployment controller, Service, dan Horizontal Autoscaler - membuat objek service">
</p>

Bisa dilihat dari gambar di atas, kita telah men-*dispatch* beberapa HTTP request ke service yang sudah kita buat. Hasil response dari `curl` berbeda satu sama lain, hal ini karena pod yang menerima dan merespon request kita adalah berbeda di tiap request (efek dari load balancer).

> Tips! Cara yang lebih praktis untuk mendapatkan URL service, gunakan command berikut:

> `minikube service <service-name> --url`<br />`minikube service my-service --url`

#### 3.3. Objek Horizontal Pod Auto Scaler (HPA)

Pada bagian ini, kita akan menambahkan kapabilitas autoscaling pada pod. Jadi, semisal nantinya muncul lonjakan jumlah pengakses aplikasi, maka jumlah pod akan di scale secara otomatis untuk mengakomodir kebutuhan tersebut.

Salah satu cara untuk menerapkan autoscaling pada pod adalah dengan menambahkan objek HPA. Objek ini akan secara cerdas me-*manage* proses *scaling* pod, kapan harus disiapkan banyak replikasi pod, kapan harus sedikit. Kita bisa menambahkan kriteria scaling berdasarkan misalnya utilisasi CPU, atau lainnya.

Ok, sekarang tambahkan konfigurasi berikut ke file `k8s.yaml`.

```yaml
---
# pilih api `autoscaling/v2beta2` untuk HPA.
apiVersion: autoscaling/v2beta2

# alokasikan blok YAML berikut untuk HPA (HorizontalPodAutoscaler).
kind: HorizontalPodAutoscaler

# namai dengan `my-auto-scaler`.
metadata:
  name: my-auto-scaler

# blok spec berikut berisi spesifikasi behaviour dari objek autoscaler ini.
spec:

  # minimum replikasi pod.
  minReplicas: 3

  # maksimum replikasi pod.
  maxReplicas: 10

  # objek deployment yang ingin di-scale.
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app-deployment

  # blok metrics berikut berisi kriteria yang dijadikan acuan untuk scaling.
  metrics:

      # kriteria scaling: ketika utilisasi CPU 50%.
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 50
```

Penulis rasa penjelasan pada komentar tiap bagian di konfigurasi di atas cukup jelas. Intinya, proses scaling akan diterapkan ke objek deployment `my-app-deployment`, dengan kriteria utilisasi CPU 50%, dengan minimum pod 3 dan maksimum 10.

Ok, sekarang mari kita terapkan file konfigurasi yang sudah ditambahkan HPA di dalamnya.

```bash
# terapkan file konfigurasi
kubectl apply -f k8s.yaml

# munculkan semua HPA
kubectl get hpa

# munculkan detail HPA
kubectl describe hpa <hpa-name>
```

<p style="text-align: center;">
  <img src="https://i.imgur.com/R63y8dL.png" alt="Kubernetes - Deploy Aplikasi ke Kluster Minikube menggunakan Deployment controller, Service, dan Horizontal Autoscaler - objek horizontal pod auto scaler">
</p>


Sebelumnya kita punya dua pod *running*. Setelah HPA diterapkan, default pod yang running menjadi 3, hal ini karena `spec.minReplicas` kita isi nilainya 3.

#### 3.3.1. Stress test pada HPA

Mari kita test objek HPA yang sudah dibuat. Kita akan simulasikan trafik tinggi menggunakan `hey` tool. Sejumlah 50 concurrent request akan di-*dispatch* ke URL service selama 5 menit.

Jalankan command berikut pada window PowerShell baru.

```bash
# munculkan URL service
minikube service my-service --url

# simulasikan stress test
hey -c 50 -z 5m <service-URL>
```

Sekarang fokus ke PowerShell satunya, cek perilaku pod setelah beberapa menit.

```bash
# munculkan semua hpa dan pod
kubectl get hpa
kubectl get pods
```

<p style="text-align: center;">
  <img src="https://i.imgur.com/0lHYlxc.png" alt="Kubernetes - Deploy Aplikasi ke Kluster Minikube menggunakan Deployment controller, Service, dan Horizontal Autoscaler - objek HPA">
</p>

Setelah satu menit berlalu, tiba-tiba ada 6 buah pod yang running. Ini terjadi karena utilisasi CPU cukup tinggi, memenuhi kriteria scaling yang sudah didefinisikan. Bisa disimpulkan berarti proses autoscaling berjalan sesuai harapan.

HPA tidak hanya secara magis bisa me-manage jumlah pod ketika tinggi traffik, tapi juga ketika trafik rendah. Coba hentikan stress test, lalu tunggu beberapa menit kemudian cek lagi HPA dan pod, maka angka akan berubah seperti awal.
