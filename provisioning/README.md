# provisioning

ansible provisioning for umsatz

## Setup

Assuming you've got bower installed, just run

```
bower install
```

This will install all ansible 3rd plugin roles.

## Getting Started (Raspberry PI)

The fastest way to get up and running is by installing the prebuild releases.
This installs all dependencies, configures your raspberry pi & makes umsatz
available via HTTP.

``` bash
echo "ip.of.rasp.berry > raspberry"
ansible-playbook -i raspberry -u pi release-install.yml
```

## Getting Started (Vagrant)

Make sure to change group_vars/all `os_arch` to amd64

``` bash
vagrant up
ansible-playbook -i vagrant -u vagrant release-install.yml --private-key=~/.vagrant.d/insecure_private_key
```

## backup/ restore

``` bash
ansible-playbook -i hosts backup.yml -e "@backups/config.json"
ansible-playbook -i hosts restore.yml -e "@restore.json" -e "archive=backup.tar"
```