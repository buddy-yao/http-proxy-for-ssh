# http-proxy-for-ssh
ssh proxy via http connect

# usage:
## command usage
`http-proxy-for-ssh <proxyhost> <proxyport> <dsthost> <dstport> [authfile]`

## ssh proxy usage
edit your ~/.ssh/config
```
Host <host-pattern>
    ProxyCommand http-proxy-for-ssh <proxyhost> <proxyport> %h %p [authfile]
```
