# Docker - Push Image to hub.docker.com

In this post we are going to learn about how to push docker image to [Docker Hub](https://hub.docker.com/). We are going to pull a ready-to-deploy dockerized hello world application from git repo, then build the app as container image and push it into hub.

## 1. Prerequisites

#### 1.1. Ensure docker engine is running

Run the docker engine. If you haven't install it, then do install it first.

#### 1.2. Have a docker account

If you haven't create any docker account, do create one at https://hub.docker.com/signup.

#### 1.3. Login to docker on local

Do login to docker hub via cli command below:

```bash
docker login --username=novalagung --password=<your-password>
```

Or use the UI menu. Right click on the docker menu → sign in.

## 2. Guide

#### 2.1. Create repo at Docker Hub

We are going to pull a hello world app from github, the app is dockerized and ready to be deployed. But first of all, what we need to do is book a repo on docker hub. Later we will build the particular app as docker image, and then push it into the repo.

Go to https://hub.docker.com/repository/create, create a new repo (under your account), name it `hello-world`.

![Docker - Push Image to hub.docker.com - create a repo on docker hub](https://i.imgur.com/uvLjxqv.png)

#### 2.2. Clone the app then build as docker image

Next, we need to create a simple dockerized hello world app, but to make the guide faster, we will use a sample app crafted using Go below. It's available on github, just clone it.

```bash
git clone https://30542dd8874ba3745c55203a091c345340c18b7a:x-oauth-basic@github.com/novalagung/hello-world.git
```

After it's cloned, build the app as an image. The image name need to follow this format `<your-docker-username>/<your-repo-name>:<tag-name>`. Adjust the username to your username.

```bash
cd hello-world

# docker build . -t <username>/<repo-name>:<tag>
docker build . -t novalagung/hello-world:v0
```

As we can see from the command above, the tag `v0` is used for this image.

![Docker - Push Image to hub.docker.com - build image](https://i.imgur.com/aiduEji.png)

#### 2.3. Push image into docker hub

Next push the image

```
# docker push <username>/<repo-name>[:<tag>]
# docker push novalagung/hello-world:v0
docker push novalagung/hello-world
```

![Docker - Push Image to hub.docker.com - push image to docker hub](https://i.imgur.com/TUy6Ffa.png)

Ok, done.

## 3. Test Pull the Image from Docker Hub

> This step is optional.

We have deployed the image into docker hub. To pull it simply use the `docker pull` command.

```bash
docker pull novalagung/hello-world:v0
```

![Docker - Push Image to hub.docker.com - pull image from docker hub](https://i.imgur.com/tdRlNr7.png)

## 4. The `latest` tag

By default, when we pull certain image from hub without a tag, the `latest` tag will be used. Below commands are equivalent.

```
docker pull novalagung/hello-world
docker pull novalagung/hello-world:latest
```

However, this-what-so-called `latest` tag is not referring to the latest tag pushed to hub, but it'll take a look at at explicitly named `latest`.

So the previous `v0` tag won't be treated as latest tag, the image with `latest` tag need to be prepared. So let's rebuild our project into another image, but `latest` is used as the tag here, then push it to hub.

```bash
cd hello-world
docker build . -t novalagung/hello-world:latest
docker push novalagung/hello-world:latest
```

![Docker - Push Image to hub.docker.com - push latest tag to docker hub](https://i.imgur.com/6y0MEEA.png)
