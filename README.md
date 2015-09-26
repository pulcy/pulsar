# Pulcy Development Environment

## Requirements

* Node.js
* Npm
* Docker
* Git

## Environment setup

Clone the Pulcy development environment tools:
```
git clone ssh://git@arvika.pulcy.com/pulcy/pulcy.git
```

Install Pulcy development tools (inside development-environment repository): 

```
make
./scripts/pulcy install
```

Create an extra in /etc/hosts:
```
sudo echo "127.0.0.1	arvika-ssh" >> /etc/hosts
```

## Commonly used URL's

https://arvika.pulcy.com - GitLab
http://arvika-ssh:5000       - Subliminl private docker registry
http://arvika-ssh:16012      - Subliminl NPM proxy 
http://arvika-ssh:16073      - Subliminl private NPM repository
