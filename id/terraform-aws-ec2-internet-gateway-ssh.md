# Terraform - Otomatisasi setup AWS EC2 dan juga akses Internet Gateway dan SSH

Pada post ini, kita akan belajar tentang implementasi terraform untuk mengotomatisasi setup AWS EC2 dengan akses internet dan SSH.

---

### 1. Kebutuhan

#### 1.1. Terraform CLI

Pastikan Terraform CLI tool tersedia. Jika belum, maka install terlebih dahulu.

#### 1.2. User IAM AWS

Siapkan satu buah user IAM baru dengan *programmatic access key* aktif dan akses penuh ke EC2 management. Kita akan gunakan `access_key` dan `secret_key`-nya pada tutorial ini. Bisa ikuti guide berikut untuk cara buat user IAM baru: [Membuat User IAM](aws-membuat-user-iam.md).

#### 1.3. Command `ssh-keygen` dan `ssh`

Pastikan CLI tools `ssh-keygen` dan `ssh` tersedia.

---

### 2. Persiapan

Buat folder baru (dengan nama bebas), isinya satu buah file bernama `infrastructure.tf`. Kita akan gunakan file ini untuk pendefinisian kode infrastruktur. Semua kode setup resource akan dituliskan dalam bahasa HCL dalam file tersebut, meliputi:

- Upload key pair (untuk keperluan SSH akses dari lokal ke EC2 instance)
- Pembuatan EC2 instance
- Pembuatan dan asosiasi security group ke VPC (dimana EC2 instance akan di-setup)
- Pembuatan public subnet
- Pembuatan internet gateway

Ok, mari kita mulai tutorialnya. Pertma siapkan folder dan file yang sudah disinggung di atas.

```bash
mkdir terraform-automate-aws-ec2-instance
cd terraform-automate-aws-ec2-instance
touch infrastructure.tf
```

Selanjutnya, buat public-key cryptography menggunakan CLI tool `ssh-keygen`. Dengan ini akan di-generate sebuah file public key `id_rsa.pub` dan private key `id_rsa`. Nantinya kita akan upload key public key-nya ke AWS dan menggunakan private key-nya untuk mengakses EC2 instance via `ssh`.

```bash
cd terraform-automate-aws-ec2-instance
ssh-keygen -t rsa -f ./id_rsa
```

![Terraform - Otomatisasi setup AWS EC2 dengan akses Internet Gateway dan SSH - generate key pair](https://i.imgur.com/ZB16oJB.png)

---

### 3. Kode Infrastruktur

Sekarang, mari kita mulai penulisan kode infrastruktur. Silakan buka file `infrastructure.tf` menggunakan editor apa saja bebas.

#### 3.1. Set AWS sebagai provider

Definisikan blok kode provider, disini kita akan gunakan [AWS sebagai cloud provider](https://www.terraform.io/docs/providers/aws/index.html). Dalam blok kode provider, tulis informasi akses AWS seperti `region`, `access_key`, dan `secret_key`. Untuk keys nilainya kita isi menggunakan keys dari user IAM yang sudah dibuat (jadi silakan sesuaikan value-nya).

Ok, berikut adalah blok kode provider, silakan tulis pada file `infrastructure.tf`.

```bash
provider "aws" {
    region = "ap-southeast-1"
    access_key = "AKIAWLTS5CSXP7E3YLWG"
    secret_key = "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"
}
```

#### 3.2. Generate key pair baru, lalu upload ke AWS

Definisikan blok kode [resource `aws_key_pair`](https://www.terraform.io/docs/providers/aws/r/key_pair.html), namai blok tersebut dengan `my_instance_key_pair`, lalu tambahkan file public key yang sebelumnya sudah di-generate `id_rsa.pub` ke dalam blok kode ini.

```bash
resource "aws_key_pair" "my_instance_key_pair" {
    key_name = "terraform_learning_key_1"
    public_key = file("id_rsa.pub")
}
```

#### 3.3. Buat EC2 instance baru

Definisikan blok kode lagi, yaitu [resource `aws_instance`](https://www.terraform.io/docs/providers/aws/r/instance.html), namai dengan `my_instance`, lalu tulis spesifikasi VPC, instance type, key pair, security group, subnet, dan public IP dalam blok kode resource ini.

Each part of the code below is self-explanatory.

```bash
# buat instance EC2 baru
resource "aws_instance" "my_instance" {

    # ami => Amazon Linux 2 AMI (HVM), SSD Volume Type (ami-0f02b24005e4aec36).
    ami = "ami-0f02b24005e4aec36"

    # instance type => t2.micro.
    instance_type = "t2.micro"

    # key pair: terraform_learning_key_1.
    key_name = aws_key_pair.my_instance_key_pair.key_name

    # vpc security groups: my_vpc_security_group.
    vpc_security_group_ids = [aws_security_group.my_vpc_security_group.id]

    # public subnet: my_public_subnet.
    # subnet ini nantinya digunakan untuk keperluan internet gateway.
    subnet_id = aws_subnet.my_public_subnet.id

    # asosiasikan satu public IP kke instance ini.
    associate_public_ip_address = true
}
```

Property `key_name` kita isi dengan nilai diambil dari resource key pair `my_instance_key_pair` yang sebelumnya sudah dibuat. Statement `aws_key_pair.my_instance_key_pair.key_name` akan mengembalikan informasi `key_name` dari blok kode key pair tersebut (pada contoh ini, nilainya adalah `terraform_learning_key_1`).

Untuk property `vpc_security_group_ids` dan `subnet_id`, value-nya juga mengambil dari blok kode resource lainnya (mirip seperti `key_name`), hanya saja untuk kedua property ini kita belum definisikan blok kode resource-nya.

O iya, property `vpc_security_group_ids` hanya menerima value berupa valid array of string, itulah kenapa nilainya dibungkus `[]`. Meskipun hanya satu element array, tetap harus dituliskan dalam bentuk array.

#### 3.4. Alokasikan satu resource VPC dengan satu security group di dalamnya

Alokasikan satu buah [blok kode resource VPC](https://www.terraform.io/docs/providers/aws/r/vpc.html), dan definisikan juga [blok kode resource security](https://www.terraform.io/docs/providers/aws/r/security_group.html)group-nya.

```bash
# alokasikan VPC bernama my_vpc.
resource "aws_vpc" "my_vpc" {
    cidr_block = "10.0.0.0/16"
    enable_dns_hostnames = true
}

# buat satu buah security group baru untuk my_vpc.
resource "aws_security_group" "my_vpc_security_group" {

    # set blok kode security group ini untuk VPC my_vpc.
    vpc_id = aws_vpc.my_vpc.id

    # definisikan rule inbound baru: mengijinkan akses TCP/SSH dari manapun.
    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    # definisikan rule inbound baru: mengijinkan akses TCP/HTTML pada port 80 dari manapun.
    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    # definisikan rule outbound baru: mengijinkan jenis akses apapun dari manapun.
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}
```

Security group di atas dibuat untuk `my_vpc`, ditandai dengan statement `vpc_id = aws_vpc.my_vpc.id` dalam blok kode tersebut. VPC ini akan memiliki 3 buah rule inbound/outbound:

- Ijinkan akses SSH dari manapun. Nantinya kita perlu untuk konek ke EC2 instance yang dibuat, untuk mengecek apakah proses setup berjalan sesuai harapan atau tidak.
- Ijinkan akses ke port `80`. Ini mungkin penting, agar operasi-operasi standar linux seperti update, install dependensi, dan lainnya bisa dilakukan.
- Ijinkan semua jenis outgoing akses dari manapun. Dengan ini nantinya kita bisa melakukan operasi seperto download, remote, dan lainnya dari instance tersebut.

> Btw, ingress itu istilah lain untuk inbound, dan egress untuk outbound

#### 3.5. Alokasikan public subnet baru untuk `my_vpc`

Kita telah mempersiapkan `my_vpc` dengan alokasi blok CIDR `10.0.0.0/16`. Sekarang kita perlu melakukan [subnetting](https://www.terraform.io/docs/providers/aws/r/subnet.html) dari blok CIDR tersebut untuk keperluan akses publik, tentukan saja misal `10.0.0.0/24`.

```bash
# buat subnet baru untuk my_vpc.
resource "aws_subnet" "my_public_subnet" {
    vpc_id = aws_vpc.my_vpc.id
    cidr_block = "10.0.0.0/24"
}
```

Kalau kita kembali ke bagian pendefinisian blok resource `my_instance`, bisa dilihat bahwa resource subnet yang baru kita buat dipergunakan pada instance `my_instance`.

#### 3.6. Buat internet gateway dan asosiasi route table

Sekarang buat [internet gateway](https://www.terraform.io/docs/providers/aws/r/internet_gateway.html) baru untuk `my_vpc`. Lalu tempelkan ke [route table](https://www.terraform.io/docs/providers/aws/r/route_table.html) baru yang juga akan kita buat, untuk keperluan akses publik.

```bash
# but internet gateway, tempelkan ke `my_vpc`.
resource "aws_internet_gateway" "my_internet_gateway" {
    vpc_id = aws_vpc.my_vpc.id
}

# buat route table baru, route internet gateway ke `my_vpc`.
resource "aws_route_table" "my_public_route_table" {
    vpc_id = aws_vpc.my_vpc.id
    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = aws_internet_gateway.my_internet_gateway.id
    }
}
```

[Asosiasikan route table publik](https://www.terraform.io/docs/providers/aws/r/route_table_association.html) di atas ke `my_public_subnet`, agar kita bisa dapat internet akses pada EC2 instance `my_instance`.

```bash
# buat blok kode resource asosisasi route table untuk menghubungkan
# `my_public_subnet` dengan `my_public_route_table`.
resource "aws_route_table_association" "my_public_route_table_association" {
  subnet_id = aws_subnet.my_public_subnet.id
  route_table_id = aws_route_table.my_public_route_table.id
}
```

Sebenarnya bagian penulisan kode infrastruktur sudah cukup, tapi ada satu hal lagi yang perlu kita lakukan sebelum masuk ke bagian penerapan infra, yaitu menampilkan DNS dan IP publik dari instance yang sudah dibuat. Informasi tersebut nantinya kita pakai untuk testing, remote ssh ke instance tersebut.

Cara menampilkan output bisa dengan menggunakan [blok `output`](https://www.terraform.io/docs/configuration/outputs.html). Silakan tambahkan blok kode berikut.

```bash
output "public-dns" {
    value = aws_instance.my_instance.*.public_dns[0]
}
output "public-ip" {
    value = aws_instance.my_instance.public_ip
}
```

Ok, kode infra sudah siap, mari kita mulai proses *terraforming*.

---

### 4. Jalankan Terraform

#### 4.1. Terraform init

Pertama, jalankan command `terraform init`. Command ini akan menjalankan proses setup dan inisialisasi untuk kode infrastruktur yang telah dibuat. Pada contoh ini, command tersebut juga akan men-download beberapa dependencies seperti AWS provider yang digunakan.

```bash
cd terraform-automate-aws-ec2-instance
terraform init
```

![Terraform - Otomatisasi setup AWS EC2 dengan akses Internet Gateway dan SSH - terraform init](https://i.imgur.com/6PnpyNc.png)

#### 4.2. Terraform plan

Selanjutnya, jalankan command `terraform plan`, untuk melihat plan/perencanaan dari infrastruktur kita. Bagian ini sebenarnya opsional, hanya saja, mungkin penting untuk dilakukan agar kita bisa tau seperti apa nantinya output infrastruktur yang dihasilkan dari kode infra yang telah dibuat.

#### 4.3. Terraform apply

Terakhir, jalankan command `terraform apply` untuk mengeksekusi proses *terraforming*.

```bash
cd terraform-automate-aws-ec2-instance
terraform apply -auto-approve
```

O iya, flag `-auto-approve` pada command `terraform apply` adalah opsional. Dengan flag tersebut nantinya semua prompt konfirmasi di-skip.

![Terraform - Otomatisasi setup AWS EC2 dengan akses Internet Gateway dan SSH - terraform apply](https://i.imgur.com/rK1LX8c.png)

Pada infra file yang telah kita buat, dua blok output didefinisikan untuk memunculkan informasi DNS dan IP public. Kedua output tersebut nantinya muncul setelah proses terraforming selesai.

---

### 5. Testing

Sekarang, mari kita test instance yang sudah dibuat. Gunakan command `ssh` untuk remote konek ke instance tersebut. Bisa menggunakan DNS atau IP public, bebas.

```bash
ssh -i id_rsa ec2-user@ec2-18-140-245-218.ap-southeast-1.compute.amazonaws.com
```

![Terraform - Otomatisasi setup AWS EC2 dengan akses Internet Gateway dan SSH - ssh to ec2 instance](https://i.imgur.com/uL1TulT.png)

Bisa dilihat dari gambar di atas, kita bisa konek ke EC2 instance lewat SSH, dan instance tersebut adalah terhubung dengan internet.
