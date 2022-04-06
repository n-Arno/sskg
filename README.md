StrongSwan Config Generator
===========================

Very quick and dirty way of defining multiple sites in a yaml file and generate config per site.

Configuration based on https://wiki.strongswan.org/projects/strongswan/wiki/UsableExamples#Site-To-Site-Scenario

Demo usage:
```
$ ./sskg
Usage: ./sskg <local site> <distant site>

$ cat sskg.yaml
sites:
- name: siteA
  public_ip: 1.2.3.4
  private_subnet: 192.168.0.0/24
- name: siteB
  public_ip: 4.3.2.1
  private_subnet: 192.168.1.0/24
- name: siteC
  public_ip: 9.8.7.6
  private_subnet: 192.168.2.0/24

$ ./sskg siteA siteB
# /etc/ipsec.conf - strongSwan IPsec configuration file

# basic configuration

config setup
              strictcrlpolicy=no
              uniqueids = yes
              charondebug = "all"

# Site to Site

conn siteA-to-siteB
              authby=secret
              left=%defaultroute
              leftid=1.2.3.4
              leftsubnet=192.168.0.0/24
              leftauth=psk
              right=4.3.2.1
              rightid=4.3.2.1
              rightsubnet=192.168.1.0/24
              rightauth=psk
              keyexchange=ikev2
              keyingtries=%forever
              fragmentation=yes
              ike=aes192gcm16-aes128gcm16-prfsha256-ecp256-ecp521,aes192-sha256-modp3072
              esp=aes192gcm16-aes128gcm16-ecp256-modp3072,aes192-sha256-ecp256-modp3072
              dpdaction=restart
              auto=route

# /etc/ipsec.secrets - This file holds shared secrets or RSA private keys for authentication.

1.2.3.4 : PSK "glMFZ5CmpIsUfWlyOJrcgxXVwli8-imbl"

# /etc/ipsec.conf - strongSwan IPsec configuration file

# basic configuration

config setup
              strictcrlpolicy=no
              uniqueids = yes
              charondebug = "all"

# Site to Site

conn siteB-to-siteA
              authby=secret
              left=%defaultroute
              leftid=4.3.2.1
              leftsubnet=192.168.1.0/24
              leftauth=psk
              right=1.2.3.4
              rightid=1.2.3.4
              rightsubnet=192.168.0.0/24
              rightauth=psk
              keyexchange=ikev2
              keyingtries=%forever
              fragmentation=yes
              ike=aes192gcm16-aes128gcm16-prfsha256-ecp256-ecp521,aes192-sha256-modp3072
              esp=aes192gcm16-aes128gcm16-ecp256-modp3072,aes192-sha256-ecp256-modp3072
              dpdaction=restart
              auto=route

# /etc/ipsec.secrets - This file holds shared secrets or RSA private keys for authentication.

4.3.2.1 : PSK "glMFZ5CmpIsUfWlyOJrcgxXVwli8-imbl"
```

