# CI/CD - Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre

In this tutorial we are going to create an ebook instance using Github, then publish it into github pages in automated manner via Github Actions, as well as generate the PDF, EPUB, and Mobi file.

So, for every push happen to the particular repo, the Github Actions (CI/CD) will trigger certain process including compilation and web-ebook generation, and then the result will be pushed into `gh-pages` branch, so then the web version of the ebook will publicly accessible.

---

### 1. Prerequisites

#### 1.1. Gitbook CLI

Install gitbook CLI (if you haven't). Do follow the guide on https://github.com/GitbookIO/gitbook-cli.

#### 1.2. Github Account

Ensure you have a Github account.

#### 1.3. Git bash

Ensure you have git bash client installed in your local machine.

---

### 2. Guide

#### 2.1. Create a Github repo

First, create a new repo in your Github account, it can be private one or public, doesn't matter. I will name the repo `softwareengineering`

![Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre - create github repo](https://i.imgur.com/diIHwxE.png)

#### 2.2. Create a new gitbook project

Next, use `gitbook` command line to initialize a new project. Use any name as the project name. In below I'll use `softwareengineering` as the name.

After the project created, try to test it locally.

```bash
gitbook init softwareengineering
cd softwareengineering
gitbook serve
```

![Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre - gitbook init project](https://i.imgur.com/99Q5kvv.png)

As we can see, the web version of the book is running up.

#### 2.2. Prepare ssh github deploy key

We are going to use Github Action plugin [peaceiris/actions-gh-pages](https://github.com/peaceiris/actions-gh-pages) to make pushing into `gh-pages` branch easier. To make it happen, first generate ssh deploy key using command below (run it in your local machine).

```
ssh-keygen -t rsa -b 4096 -C "$(git config user.email)" -f gh-pages -N ""
# You will get 2 files:
#   gh-pages.pub (public key)
#   gh-pages     (private key)
```

The above command generate two files:

- `gh-pages.pub` file as the public key
- `gh-pages` file as the private key

Do upload these two files into repo's project keys and secret menu respectively. To do that, open the repo, click `Settings`, then do follow steps below:

![Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre - prepare github deploy key](https://i.imgur.com/t8RVwN7.png)

#### 2.3. Create Github workflow CI/CD file for generating the web version of ebook

We are going to make this project automatically deploy the web version of the ebook on every push, including the first push.

Create a workflow a new file called `deploy.yml` place it in `<yourproject>/.github/workflows`, then fill it with configuration below:

```yaml
# file softwareengineering/.github/workflow/deploy.yml

name: 'deploy website and ebooks'

on:
  push:
    branches:
      - master

jobs:
  job_deploy_website:
    name: 'deploy website'
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v1
    - uses: actions/setup-node@v1
      with:
        node-version: '10.x'
    - name: 'Installing gitbook cli'
      run: npm install -g gitbook-cli
    - name: 'Generating distributable files'
      run: |
        gitbook install
        gitbook build
    - uses: peaceiris/actions-gh-pages@v2.5.0
      env:
        ACTIONS_DEPLOY_KEY: ${{ secrets.ACTIONS_DEPLOY_KEY }}
        PUBLISH_BRANCH: gh-pages
        PUBLISH_DIR: ./_book
```

In summary the workflow above will do these things sequentially:

- Trigger this workflow on every push happen on `master` branch.
- Install `nodejs`.
- Install `gitbook` CLI.
- Build the project.
- use `peaceiris/actions-gh-pages` plugin to deploy the built result to `gh-pages` branch.

The previous github deploy key is used on the push-to-gh-pages process.

#### 2.4. Push project to Github repo

```bash
cd softwareengineering

# ignore certain directory
touch .gitignore
echo '_book' >> .gitignore

# init git repo
git init
git add .
git commit -m "init"
git remote add origin git@github.com:novalagung/softwareengineering.git

# push
git push origin master
```

Navigate to browser, open your github repo, click `Actions`, watch a workflow process that currently is running.

![Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre - github workflow](https://i.imgur.com/SZfwqZs.png)

After the workflow is complete, then try to open in the browser `https://<github-username>.github.io/<repo-name>`. In this example it is `https://novalagung.github.io/softwareengineering/`.

![Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre - open web version of the book](https://i.imgur.com/HzCygaX.png)

If you are still not sure about the URL, open `Settings` menu of your Github repo, then scroll down little bit until `Github Pages` section appear. The Github Pages URL will appear there.

![Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre - github pages url](https://i.imgur.com/eD5BmPv.jpg)

#### 2.5. Modify the workflow file to be able to enable generate the file version

Open the previous `deploy.yml` file, add new job called `job_deploy_ebooks` below.

```yaml
# file softwareengineering/.github/workflow/deploy.yml

name: 'deploy website and ebooks'

on:
  push:
    branches:
      - master

env:
  ebook_name: 'softwareengineeringtutorial'

jobs:
  job_deploy_website:
    # ...
  job_deploy_ebooks:
    name: 'deploy ebooks'
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v1
    - uses: actions/setup-node@v1
      with:
        node-version: '10.x'
    - name: 'Installing gitbook cli'
      run: npm install -g gitbook-cli
    - name: 'Installing calibre'
      run: |
        sudo -v
        wget -nv -O- https://download.calibre-ebook.com/linux-installer.sh | sudo sh /dev/stdin
    - name: 'Preparing for ebooks generations'
      run: |
        gitbook install
        mkdir _book
    - name: 'Generating ebook in pdf'
      run: gitbook pdf ./ ./_book/${{ env.ebook_name }}.pdf
    - name: 'Generating ebook in epub'
      run: gitbook epub ./ ./_book/${{ env.ebook_name }}.epub
    - name: 'Generating ebook in mobi'
      run: gitbook mobi ./ ./_book/${{ env.ebook_name }}.mobi
    - uses: peaceiris/actions-gh-pages@v2.5.0
      env:
        ACTIONS_DEPLOY_KEY: ${{ secrets.ACTIONS_DEPLOY_KEY }}
        PUBLISH_BRANCH: ebooks
        PUBLISH_DIR: ./_book
```

The `job_deploy_website` that we have created responsible for generating the web base version of the ebook. This newly created `job_deploy_ebooks` created for different purpose, to generate the file version (pdf, epub, mobi). The generated file will be pushed to branch named `ebooks`.

The ebook file generated handle by a library called `calibre`.

Ok, now let's update the repo with recent changes.

```bash
git add .
git commit -m "update"
git push origin master
```

![Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre - workflow to generate ebook files](https://i.imgur.com/iXd7bnr.png)

After the process complete, the ebooks will be available for download in these following URLs. Please adjust it to follow your github profile and repo name.

```bash
https://github.com/novalagung/softwareengineering/raw/ebooks/softwareengineeringtutorial.pdf
https://github.com/novalagung/softwareengineering/raw/ebooks/softwareengineeringtutorial.epub
https://github.com/novalagung/softwareengineering/raw/ebooks/softwareengineeringtutorial.mobi
```

FYI! Since the ebook files are accessible through github direct link, this mean the visibility of the repo need to be public (not private). If you want the repo to be in private but keep the files accessible, then do push the files into `gh-pages` branch.

#### 2.6. Add custom domain

Now we are going to add a custom domain to our Github Page. To do that, do navigate to your domain control panel, then add new CNAME record that point to your Github page domain `<github-username>.github.io`.

![Serverless Ebook using Gitbook, Github Pages, Github Actions, and Calibre - custom domain github pages](https://i.imgur.com/a1vF2Xk.png)

FYI, In this exaxmple we pick subdomain `softwareengineering.novalagung.com`.

Then in your gitbook project, add new file called `CNAME`, fill it with the subdomain URL. After that push the update to the upstream.

```bash
echo 'softwareengineering.novalagung.com' > CNAME

git add .
git commit -m "update"
git push origin master
```

