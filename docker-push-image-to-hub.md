# Docker - Push Image to hub.docker.com

In this post, we are going to learn about how to push a Docker image to [Docker Hub](https://hub.docker.com/).

---

## 1. Prerequisites

### 1.1. Docker engine

Ensure the Docker engine is running. If you haven't installed it, then install it first.

### 1.2. Docker Hub account

Prepare a Docker Hub account. If you don't have it, then follow a guide on [Create Docker Hub Account](docker-hub-create-account.html).

### 1.3. Login to Docker Hub on the local machine

Do log in to Docker Hub via CLI command below:

```bash
docker login --username=novalagung --password=<your-password>
```

Or use the UI menu. It is available by doing a right-click on the docker menu → sign in.

---

## 2. Guide

### 2.1. Create repo at Docker Hub

First of all, we need to book a repo on Docker Hub. Later we will push the image to that particular repo.

Go to https://hub.docker.com/repository/create, create a new repo (under your account), name it `hello-world` (or anything).

![Docker - Push Image to hub.docker.com - create a repo on Docker Hub](https://i.imgur.com/uvLjxqv.png)

### 2.2. Clone the example app then build as Docker image

Next, we need to create a simple dockerized hello world app. But to make the thing faster, we will use a ready-to-deploy-dockerized hello world app crafted using Go. It's available on Github (via Github token), just run the command below.

```bash
git clone https://30542dd8874ba3745c55203a091c345340c18b7a:x-oauth-basic@github.com/novalagung/hello-world.git
```

After the cloning process is finished, build the app as Docker image with a name in this format `<your-docker-username>/<your-repo-name>:<tag-name>`. Adjust the value of `<your-docker-username>` to use your actual Docker Hub username.

```bash
cd hello-world

# docker build . -t <username>/<repo-name>:<tag>
docker build . -t novalagung/hello-world:v0
```

As we can see from the command above, the tag `v0` is used on this image.

![Docker - Push Image to hub.docker.com - build image](https://i.imgur.com/aiduEji.png)

### 2.3. Push image into Docker Hub

Next, use `docker push` command below to push the image that we just built.

```
# docker push <username>/<repo-name>[:<tag>]
docker push novalagung/hello-world
```

![Docker - Push Image to hub.docker.com - push image to Docker Hub](https://i.imgur.com/TUy6Ffa.png)

Ok, done.

---

## 3. Test - Pull the Image from Docker Hub

> This step is optional.

We have pushed the image into Docker Hub. To pull it, use the `docker pull` command.

```bash
# docker pull <username>/<repo-name>[:<tag>]
docker pull novalagung/hello-world:v0
```

![Docker - Push Image to hub.docker.com - pull image from Docker Hub](https://i.imgur.com/tdRlNr7.png)

---

## 4. The `latest` tag

By default, when we pull a certain image from the Hub without a tag specified, then the `latest` tag of the particular image will be pulled.

Take a look at two commands below, they are equivalent.

```
docker pull novalagung/hello-world
docker pull novalagung/hello-world:latest
```

The funny thing about this what-so-called `latest` tag is, it is actually not referring to the latest tag pushed to the Hub, it'll look for a tag with explicit name `latest`.

The previous `v0` tag won't be treated as the latest tag. To have the latest tag, we shall rebuild our project into another image then push it to the Hub, but this time during the build we will do it using `latest` as the tag.

```bash
cd hello-world
docker build . -t novalagung/hello-world:latest
docker push novalagung/hello-world:latest
```

![Docker - Push Image to hub.docker.com - push latest tag to Docker Hub](https://i.imgur.com/6y0MEEA.png)
