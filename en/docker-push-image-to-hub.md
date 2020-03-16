# Docker | Push Image to hub.docker.com

In this post we are going to learn about how to push a docker image to [Docker Hub](https://hub.docker.com/).

## 1. Prerequisites

#### 1.1. Ensure docker engine is running

Run the docker engine. If you haven't installed it, then do it.

#### 1.2. Have a docker account

If you haven't created any docker account, do create one at https://hub.docker.com/signup.

#### 1.3. Login to docker hub on the local machine

Do log in to docker hub via CLI command below:

```bash
docker login --username=novalagung --password=<your-password>
```

Or use the UI menu. It is available by doing a right click on the docker menu → sign in.

## 2. Guide

#### 2.1. Create repo at Docker Hub

First of all, we need to book a repo on Docker hub. Later we will push the image to that particular repo.

Go to https://hub.docker.com/repository/create, create a new repo (under your account), name it `hello-world` (or anything).

![Docker | Push Image to hub.docker.com | create a repo on docker hub](https://i.imgur.com/uvLjxqv.png)

#### 2.2. Clone the app then build as docker image

Next, we need to create a simple dockerized hello world app. But to make the thing faster, we will use a ready-to-deploy-dockerized hello world app crafted using Go. It's available on Github (via Github token), just run the command below.

```bash
git clone https://30542dd8874ba3745c55203a091c345340c18b7a:x-oauth-basic@github.com/novalagung/hello-world.git
```

After the cloning process is finished, build the app as Docker image with a name in this format `<your-docker-username>/<your-repo-name>:<tag-name>`. Adjust the value of `<your-docker-username>` to use your actual Docker hub username.

```bash
cd hello-world

# docker build . -t <username>/<repo-name>:<tag>
docker build . -t novalagung/hello-world:v0
```

As we can see from the command above, the tag `v0` is used on this image.

![Docker | Push Image to hub.docker.com | build image](https://i.imgur.com/aiduEji.png)

#### 2.3. Push image into docker hub

Next, use `docker push` command below to push the image that we just built.

```
# docker push <username>/<repo-name>[:<tag>]
docker push novalagung/hello-world
```

![Docker | Push Image to hub.docker.com | push image to docker hub](https://i.imgur.com/TUy6Ffa.png)

Ok, done.

## 3. Test - Pull the Image from Docker Hub

> This step is optional.

We have pushed the image into Docker hub. To pull it, use the `docker pull` command.

```bash
# docker pull <username>/<repo-name>[:<tag>]
docker pull novalagung/hello-world:v0
```

![Docker | Push Image to hub.docker.com | pull image from docker hub](https://i.imgur.com/tdRlNr7.png)

## 4. The `latest` tag

By default, when we pull a certain image from hub without a tag specified, then the `latest` tag of the particular image will be pulled.

Take a look at two commands below, they are equivalent.

```
docker pull novalagung/hello-world
docker pull novalagung/hello-world:latest
```

The funny thing about this what-so-called `latest` tag is, it is actually not referring to the latest tag pushed to the hub, it'll look for a tag with explicit name `latest`.

The previous `v0` tag won't be treated as the latest tag. To have a latest tag, we shall rebuild our project into another image then push it to hub, but this time during the build we will do it using `latest` as the tag.

```bash
cd hello-world
docker build . -t novalagung/hello-world:latest
docker push novalagung/hello-world:latest
```

![Docker | Push Image to hub.docker.com | push latest tag to docker hub](https://i.imgur.com/6y0MEEA.png)
