# Terraform | AWS EC2 + Load Balancer + Auto Scaling

In this post, we are going to learn about the usage of terraform to automate the setup of AWS EC2 instance on an auto scaling environment with application load balancer applied.

Since we will be using auto-scaling feature, then the app within the instance needs to be deployed in automated manner.

The application is a simple go app, currently hosted on Github in a private repo. We will clone the app using Github token, we will talk about it in details in some part of this tutorial.

---

### 1. Prerequisites

#### 1.1. Terraform CLI

Ensure terraform CLI is available. If not, then follow guide on [Terraform Installation](terraform-cli-installation.md).

#### 1.2. Individual AWS IAM user

Prepare a new individual IAM user with programmatic access key enabled and has access to EC2 management. We will use the `access_key` and `secret_key` on this tutorial. If you haven't create the IAM user, then follow guide on [Create Individual IAM User](aws-create-individual-iam-user.md).

#### 1.3. `ssh-keygen` and `ssh` commands

Ensure both `ssh-keygen` and `ssh` command are available.

---

### 2. Initialization

Create a new folder contains a file named `infrastructure.tf`. We will use the file as the infrastructure code. Every setup will be written in HCL language inside the file, including: 

- Uploading key pair (for ssh access to the instance)
- Creating EC2 instance
- ...

Ok, let's back to the tutorial. Now create the infrastructure file.

```bash
mkdir terraform-automate-aws-ec2-instance
cd terraform-automate-aws-ec2-instance
touch infrastructure.tf
```

Next, create new key pair using `ssh-keygen` command below. This will generate the `id_rsa.pub` public key, and `id_rsa` private key. Later we will upload the public key into aws and use the private key to perform `ssh` access into the newly created EC2 instance.

```bash
cd terraform-automate-aws-ec2-instance
ssh-keygen -t rsa -f ./id_rsa
```

![Terraform | AWS EC2 + Load Balancer + Auto Scaling | generate key pair](https://i.imgur.com/ZB16oJB.png)

---

### 3. Infrastructure Code

Now we shall start writing the infrastructure config. Open `infrastructure.tf` in any editor.

#### 3.1. Define AWS provider

Define the provider block with [AWS as chosen cloud provider](https://www.terraform.io/docs/providers/aws/index.html). Also define these properties: `region`, `access_key`, and `secret_key`; with values derived from the created IAM user.

Write a block of code below into `infrastructure.tf`

```bash
provider "aws" {
    region = "ap-southeast-1"
    access_key = "AKIAWLTS5CSXP7E3YLWG"
    secret_key = "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"
}
```

#### 3.2. Generate new key pair then upload to AWS

Define new [`aws_key_pair` resource](https://www.terraform.io/docs/providers/aws/r/key_pair.html) block with local name: `my_instance_key_pair`. Put the previously generated `id_rsa.pub` public key inside the block to upload it to AWS.

```bash
resource "aws_key_pair" "my_instance_key_pair" {
    key_name = "terraform_learning_key_1"
    public_key = file("id_rsa.pub")
}
```

#### 3.3. Book a VPC, and enable internet gateway on it

Book a VPC, name it `my_vpc`. Then enable internet gateway on it. Each part of code below is self-explanatory.

```bash
# allocate a vpc named my_vpc.
resource "aws_vpc" "my_vpc" {
    cidr_block = "10.0.0.0/16"
    enable_dns_hostnames = true
}

# setup internet gateway for my_vpc.
resource "aws_internet_gateway" "my_vpc_igw" {
    vpc_id = aws_vpc.my_vpc.id
}

# attach the internet gateway my_vpc_igw into my_vpc.
resource "aws_route_table" "my_public_route_table" {
    vpc_id = aws_vpc.my_vpc.id
    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = aws_internet_gateway.my_vpc_igw.id
    }
}
```

#### 3.4. Allocate two different subnets on two different availability zones (within same region)

Application Load Balancer or ALB requires two subnets setup on two availability zones (within same region).

In this example, the region we used is `ap-southeast-1`, as defined in the provider block above (see 3.1). There are two zones available within this region, `ap-southeast-1a` and `ap-southeast-1b`. The ALB (not classic network load balancer) requires at least to be enabled on two different zones, so we will use those two.

```bash
# prepare a subnet for availability zone ap-southeast-1a.
resource "aws_subnet" "my_subnet_public_southeast_1a" {
    vpc_id = aws_vpc.my_vpc.id
    cidr_block = "10.0.0.0/24"
    availability_zone = "ap-southeast-1a"
}
# associate the internet gateway into newly created subnet for ap-southeast-1a
resource "aws_route_table_association" "my_public_route_association_for_southeast_1a" {
    subnet_id = aws_subnet.my_subnet_public_southeast_1a.id
    route_table_id = aws_route_table.my_public_route_table.id
}

# prepare a subnet for availability zone ap-southeast-1b
resource "aws_subnet" "my_subnet_public_southeast_1b" {
    vpc_id = aws_vpc.my_vpc.id
    cidr_block = "10.0.1.0/24"
    availability_zone = "ap-southeast-1b"
}
# associate the internet gateway into newly created subnet for ap-southeast-1b
resource "aws_route_table_association" "my_public_route_association_for_southeast_1b" {
    subnet_id = aws_subnet.my_subnet_public_southeast_1b.id
    route_table_id = aws_route_table.my_public_route_table.id
}
```

The internet gateway associated with two zones that we just created. In this example, it is required for the application hosted within instances on these zones to be able to connect to the internet.

#### 3.5. Define ALB resource block, listener, security group, and target group

The ALB will be created with two subnets attached (subnets from `ap-southeast-1a` and `ap-southeast-1b`).

```bash
# create an application load balancer.
# attach the previous availability zones' subnets into this load balancer.
resource "aws_lb" "my_alb" {
    name = "my-alb"
    internal = false # set lb for public access
    load_balancer_type = "application" # use application load balancer
    security_groups = [aws_security_group.my_alb_security_group.id]
    subnets = [ # attach the availability zones' subnets.
        aws_subnet.my_subnet_public_southeast_1a.id,
        aws_subnet.my_subnet_public_southeast_1b.id 
    ]
}
```

The security group for our load balancer has only two rules.

- Allow only incoming TCP/HTTP request on port `80`.
- Allow every kind of outgoing request.

```bash
# prepare a security group for our load balancer my_alb.
resource "aws_security_group" "my_alb_security_group" {
    vpc_id = aws_vpc.my_vpc.id
    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}
```

Next, we shall prepare the ALB listener. The load balancer will listen for every incoming request to port `80`, and then the particular request will be directed towards port `8080` on the instance.

Port `8080` is chosen here because the application (that will be deployed later) will listen to this port.

```bash
# create an alb listener for my_alb.
# forward rule: only accept incoming HTTP request on port 80,
# then it'll be forwarded to port target:8080.
resource "aws_lb_listener" "my_alb_listener" {  
    load_balancer_arn = aws_lb.my_alb.arn
    port = 80  
    protocol = "HTTP"
    default_action {    
        target_group_arn = aws_lb_target_group.my_alb_target_group.arn
        type = "forward"  
    }
}

# my_alb will forward the request to a particular app,
# that listen on 8080 within instances on my_vpc.
resource "aws_lb_target_group" "my_alb_target_group" {
    port = 8080
    protocol = "HTTP"
    vpc_id = aws_vpc.my_vpc.id
}
```

#### 3.6. Define launch config (and it's required dependencies) for auto-scaling

We are not going to simply create an instance then deploy the application into it. Instead, the instance creation and app deployment will be automated using AWS auto-scaling feature.

In the resource block below, we will set up the launch configuration for the auto-scaling. This launch config is the one that decides how the instance will be created.

- The *Amazon Linux 2 AMI t2.micro* is used here.
- The launched instance will have a public IP attached, this is better to be set to `false`, but in here we might need it for testing purpose.
- The previous allocated key pair will also be used on the instance, to make it accessible through SSH access. This part is also for testing purpose.

Other than that, there is one point left that is very important, the `user_data`. The user data is a block of bash script that will be executed during instance bootstrap. We will use this to automate the deployment of our application. The whole script is stored in a file named `deployment.sh`, we will prepare it later.

```bash
# setup launch configuration for the auto-scaling.
resource "aws_launch_configuration" "my_launch_configuration" {

    # Amazon Linux 2 AMI (HVM), SSD Volume Type (ami-0f02b24005e4aec36).
    image_id = "ami-0f02b24005e4aec36"

    instance_type = "t2.micro"
    key_name = aws_key_pair.my_instance_key_pair.key_name # terraform_learning_key_2
    security_groups = [aws_security_group.my_launch_config_security_group.id]

    # set to false on prod stage.
    # otherwise true, because ssh access might be needed to the instance.
    associate_public_ip_address = true
    lifecycle {
        # ensure the new instance is only created before the other one is destroyed.
        create_before_destroy = true
    }

    # execute bash scripts inside deployment.sh on instance's bootstrap.
    # what the bash scripts going to do in summary:
    # fetch a hello world app from Github repo, then deploy it in the instance.
    user_data = file("deployment.sh")
}
```

Below is the launch config security group. In this block, we define the security group specifically for the instances that will be created by the auto scale launch config. Three rules defined here:

- Allow incoming `TCP/SSH` access on port `22`.
- Allow `TCP/HTTP` access on port `8080`.
- Allow every kind of outgoing requests.

```bash
# security group for launch config my_launch_configuration.
resource "aws_security_group" "my_launch_config_security_group" {
    vpc_id = aws_vpc.my_vpc.id
    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }
    ingress {
        from_port = 8080
        to_port = 8080
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}
```

Ok, the auto scale launch config is ready, now we shall attach it into our ALB.

```bash
# create an autoscaling then attach it into my_alb_target_group.
resource "aws_autoscaling_attachment" "my_aws_autoscaling_attachment" {
    alb_target_group_arn = aws_lb_target_group.my_alb_target_group.arn
    autoscaling_group_name = aws_autoscaling_group.my_autoscaling_group.id
}
```

Next, we shall prepare the auto scaling group config. This resource is used to determine when or on what condition the scaling process run.

- As per below config, the auto-scaling will have minimum of 2 instances alive, and 5 max.
- The `ELB` health check is enabled.
- The previous two subnets on `ap-southeast-1a` and `ap-southeast-1b` are applied.

```bash
# define the autoscaling group.
# attach my_launch_configuration into this newly created autoscaling group below.
resource "aws_autoscaling_group" "my_autoscaling_group" {
    name = "my-autoscaling-group"
    desired_capacity = 2 # ideal number of instance alive
    min_size = 2 # min number of instance alive
    max_size = 5 # max number of instance alive
    health_check_type = "ELB"

    # allows deleting the autoscaling group without waiting
    # for all instances in the pool to terminate
    force_delete = true

    launch_configuration = aws_launch_configuration.my_launch_configuration.id
    vpc_zone_identifier = [
        aws_subnet.my_subnet_public_southeast_1a.id,
        aws_subnet.my_subnet_public_southeast_1b.id 
    ]
    timeouts {
        delete = "15m" # timeout duration for instances
    }
    lifecycle {
        # ensure the new instance is only created before the other one is destroyed.
        create_before_destroy = true
    }
}
```

#### 3.7. Print the ALB public DNS

Everything is pretty much done, except we need to print the ALB public DNS, so then we can do the testing.

```bash
# print load balancer's DNS, test it using curl.
#
# curl my-alb-625362998.ap-southeast-1.elb.amazonaws.com
output "alb-url" {
    value = aws_lb.my_alb.dns_name
}
```

---

### 4. App Deployment Script

We have done with the infrastructure code, next prepare the deployment script.

Create a file named `deployment.sh` in the same directory where infra code is placed. It will contain bash scripts for automating app deployment. This file will be used by auto-scaling launcher to automate app setup during instance bootstrap.

The application is written in Go, and the AMI *Amazon Linux 2 AMI t2.micro* that used here does not have any Go tools ready, that's why we need to set it up.

> **Deploying app** means that the app is ready (has been built into binary), so what we need is simply just run the binary.

> However to make our learning process better, in this example, we are going to fetch the app source code and perform the build and deploy processes within the instance.

Ok, here we go, the bash script.

```bash
#!/bin/bash

# install git
sudo yum -y install git

# download go, then install it
wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.14.linux-amd64.tar.gz

# clone the hello world app.
# The app is hosted on private repo,
# that's why the github token is used on cloning the repo
github_token=30542dd8874ba3745c55203a091c345340c18b7a
git clone https://$github_token:x-oauth-basic@github.com/novalagung/hello-world.git \
    && echo "cloned" \
    || echo "clone failed"

# export certain variables required by go
export GO111MODULE=on
export GOROOT=/usr/local/go
export GOCACHE=~/gocache
mkdir -p $GOCACHE
export GOPATH=~/goapp
mkdir -p $GOPATH

# create local vars specifically for the app
export PORT=8080
export INSTANCE_ID=`curl -s http://169.254.169.254/latest/meta-data/instance-id`

# build the app
cd hello-world
/usr/local/go/bin/go env
/usr/local/go/bin/go mod tidy
/usr/local/go/bin/go build -o binary

# run the app with nohup
nohup ./binary &
```

---

### 5. Run Terraform

#### 5.1. Terraform initialization

First, run the `terraform init` command. This command will do some setup/initialization, certain dependencies (like AWS provider that we used) will be downloaded.

```bash
cd terraform-automate-aws-ec2-instance
terraform init
```

#### 5.2. Terraform plan

Next, run `terraform plan`, to see the plan of our infrastructure. This step is optional, however, might be useful for us to see the outcome from the infra file.

#### 5.3. Terraform apply

Last, run the `terraform apply` command to execute the infrastructure plan.

```bash
cd terraform-automate-aws-ec2-instance
terraform apply -auto-approve
```

The `-auto-approve` flag is optional, it will skip the confirmation prompt during execution.

After the process is done, public DNS shall appear. Next, we shall test the instance. 

---

### 6. Test Instance

Use the `curl` command to make HTTP request to the ALB public DNS instance.

```bash
curl -X GET my-alb-613171058.ap-southeast-1.elb.amazonaws.com
```

![Terraform | AWS EC2 + Load Balancer + Auto Scaling | curl to load balancer](https://i.imgur.com/5jonEG2.png)

We can see from the image above, the HTTP response is different from one another across those multiple `curl` commands. The load balancer manages the traffic, sometimes we will get the instance A, sometimes B, etc.

In the AWS console, the instances that up and running are visible.

![Terraform | AWS EC2 + Load Balancer + Auto Scaling | aws console](https://i.imgur.com/iETYwfw.png)
