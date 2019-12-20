# aws-iam-authenticator-pwmgr

fetch credentials for aws-iam-authenticator from password management

current implimentations:

- lpass ([LastPass](#lastpass))

## Lastpass

Prerequisites:

- [lpass](https://github.com/lastpass/lastpass-cli/blob/master/README.md)
- [aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html)

Create an entry in lastpass with username the AWS_ACCESS_KEY_ID and password the AWS_SECRET_ACCESS_KEY

Then configure your kubeconfig as normal but use `aws-iam-authenticator-lpass` as the command

```yaml
...
kind: Config
preferences: {}
users:
- name: aws
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: aws-iam-authenticator-lpass # <--- note the different command here
      env:
      - name: "AWS_PROFILE"
        value: "REPLACE_ME_WITH_YOUR_LASTPASS_ENTRY_INCLUDING_FOLDER"
      args:
        - "token"
        - "-i"
        - "REPLACE_ME_WITH_YOUR_CLUSTER_ID"
        - "-r"
        - "REPLACE_ME_WITH_YOUR_ROLE_ARN"
```

See https://github.com/kubernetes-sigs/aws-iam-authenticator for more details
