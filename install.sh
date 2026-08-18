cd
pkill myray
rm -rf myray
mkdir myray
cd myray
mkdir html

export MYRAY_IPV4=$(curl -s4 ip.sb)
export MYRAY_IPV6=$(curl -s6 ip.sb)
export MYRAY_RAND=$(openssl rand 16 | basenc --base64url | tr -d =)

curl -Lo tmp.tar.gz https://github.com/caddyserver/caddy/releases/download/v2.11.4/caddy_2.11.4_linux_amd64.tar.gz
tar xf tmp.tar.gz caddy
rm -f tmp.tar.gz
mv caddy mycaddy

printf '
{
        auto_https disable_redirects
        ocsp_stapling off
        servers {
                protocols h1
        }
        default_sni %s
}

%s %s {
        tls {
                protocols tls1.3 tls1.3
                issuer acme {
                        profile shortlived
                }
        }
        encode
        file_server {
                root html
        }
        header -Server
        handle_errors {
                header -Server
        }
        handle_path /%s/* {
                reverse_proxy unix/@myray
        }
}
' $MYRAY_IPV4 $MYRAY_IPV4 $MYRAY_IPV6 $MYRAY_RAND > Caddyfile

echo '<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>' > html/index.html

sudo setcap cap_net_bind_service=+ep ./mycaddy

curl -o myray https://raw.githubusercontent.com/bocduck/myray/refs/heads/main/test
chmod +x myray

echo '
pkill mycaddy
pkill myray
cd ~/myray
nohup ./mycaddy run > log_mycaddy 2>&1 &
nohup ./myray -net unix -addr @myray >/dev/null 2>&1 &
' > run
chmod +x run

echo "vless://00000000-0000-0000-0000-000000000000@${MYRAY_IPV4}:443?security=tls&type=httpupgrade&path=%2F${MYRAY_RAND}%2F&sni=${MYRAY_IPV4}&udp=0"

echo "vless://00000000-0000-0000-0000-000000000000@[${MYRAY_IPV6}]:443?security=tls&type=httpupgrade&path=%2F${MYRAY_RAND}%2F&sni=${MYRAY_IPV6}&udp=0"

cd
myray/run
