Name:           go_web_server
Version:        0.1
Release:        1.all
Summary:        GO WEB Server

License:        Apache License 2.0
URL:            https://github.com/arloor/go_web_server
#Source0:

buildroot:      %_topdir/BUILDROOT
BuildRequires:  git golang
#Requires:

%description
Rust Http Proxy which is based on hyper and Tls-Listener.

%prep
if [ -d /tmp/go_web_server ]; then
        cd /tmp/go_web_server;
          git pull --ff-only || {
            echo "git pull 失败，重新clone"
            cd /tmp
            rm -rf /tmp/go_web_server
            git clone %{URL} /tmp/go_web_server
          }
else
        git clone %{URL} /tmp/go_web_server
fi

%build
cd /tmp/go_web_server
go mod tidy
go build go_web_server/cmd/go_web_server

%install
cd /tmp/go_web_server
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/lib/systemd/system
mkdir -p %{buildroot}/etc/go_web_server
mkdir -p %{buildroot}/var/go_web_server
install  -m755 go_web_server %{buildroot}/usr/bin/go_web_server
install  -m755 rpm/go_web_server.service %{buildroot}/lib/systemd/system/go_web_server.service
install  -m755 rpm/env %{buildroot}/etc/go_web_server/env
install  -m755 favicon.ico %{buildroot}/var/go_web_server/favicon.ico

%check

%pre


%post
[ ! -d /usr/share/go_web_server ]&&{
  mkdir -p /usr/share/go_web_server
}
[ ! -f /usr/share/go_web_server/privkey.pem -o ! -f /usr/share/go_web_server/cert.pem ]&&{
  if [ -f /usr/share/go_web_server/cert.pem ]; then
    rm -f /usr/share/go_web_server/cert.pem
  fi
  if [ -f /usr/share/go_web_server/privkey.pem ]; then
      rm -f /usr/share/go_web_server/privkey.pem
  fi
  openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout /usr/share/go_web_server/privkey.pem -out /usr/share/go_web_server/cert.pem -days 3650 -subj "/C=cn/ST=hl/L=sd/O=op/OU=as/CN=example.com"
}
systemctl daemon-reload

%files
/usr/bin/go_web_server
%config(noreplace) /var/go_web_server/favicon.ico
%config /lib/systemd/system/go_web_server.service
%config(noreplace) /etc/go_web_server/env



%changelog
* Sun May 07 2023 root
- init