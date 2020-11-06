# AWS - Create Individual IAM User

In this tutorial we will learn how to create an individual IAM user.

---

### 1. Definitions

IAM or Identity and Access Management is used to securely control individual and group access to your AWS resources. You can create and manage user identities (IAM users) and grant permissions for those IAM users to access your resources.

The IAM user is similar to a AWS account user, the only differences are IAM user's permission towards AWS resources are controlled (by the AWS account user).

---

### 2. Create new IAM User

To create IAM user, you (as the owner of AWS account) need to login to AWS console first. Then do open [AWS IAM page](https://console.aws.amazon.com/iam/home?region=ap-southeast-1#/home) and click the **manage users** menu.

![AWS - Create Individual IAM User](https://i.imgur.com/yx8dVAR.png)

You will directed to new page that show list of created IAM users. Next, click the **Add user**, then fill the name.

If the particular user will be used on 3rd party or AWS SDK, then do check check the **Programmatic Access** option.

![AWS - Create Individual IAM User - programmatic access](https://i.imgur.com/2V7shR9.png)

Then click next to open the user group page. In here do create new group with certain access checked. For example in the image below, a new group named `user` is created with full access to EC2 features.

![AWS - Create Individual IAM User - iam user group](https://i.imgur.com/l46C9OQ.png)

Do click next few times, then the user creation process will be done.

![AWS - Create Individual IAM User - access key and secret key](https://i.imgur.com/oqAAWZv.png)

Copy the **access key ID** and **secret access key**, save it into some notes because you won't be able to see the secret key again.

Ok, that's it. The keys are ready to use.
