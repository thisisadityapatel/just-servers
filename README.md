## just-servers

Low-level network servers (TCP|UDP) in Go, managed on AWS EC2 using Terraform for provisioning, deployment and termination.

#### Challenges
[Protohacker Challenges](https://protohackers.com)

- [0: Echo Smoke Test](https://github.com/thisisadityapatel/just-servers/tree/main/servers/echo)
- [1: Prime Time](https://github.com/thisisadityapatel/just-servers/tree/main/servers/primetime)
- [2: Means to an End](https://github.com/thisisadityapatel/just-servers/tree/main/servers/means_to_an_end)
- [3: Budget Chat](https://github.com/thisisadityapatel/just-servers/tree/main/servers/budget_chat)

#### Set-up

Locally
```shell
chmod +x deploy-locally.sh
./deploy-locally.sh
```

AWS EC2 (t2-micro)
```shell
terraform init
terraform apply
```
