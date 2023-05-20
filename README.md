```shell
## 打包
yum install -y rpm-build
yum install -y rpmdevtools
rm -rf ~/rpmbuild
rpmdev-setuptree
```

```shell
## 打包
wget https://raw.githubusercontent.com/arloor/go_web_server/master/rpm/go_web_server.spec -O /var/go_web_server.spec
rpmbuild -bb /var/go_web_server.spec

## 安装
version=0.1
release=3.all
echo RPM信息
rpm -qpi ~/rpmbuild/RPMS/x86_64/go_web_server-${version}-${release}.x86_64.rpm
echo 配置文件
rpm -qpc ~/rpmbuild/RPMS/x86_64/go_web_server-${version}-${release}.x86_64.rpm
echo 所有文件
rpm -qpl ~/rpmbuild/RPMS/x86_64/go_web_server-${version}-${release}.x86_64.rpm
systemctl stop go_web_server
# rpm -ivh在安装新版本时会报错文件冲突，原因是他没有进行更新或降级的能力，而yum install可以处理可执行文件的更新或降级
yum install -y ~/rpmbuild/RPMS/x86_64/go_web_server-${version}-${release}.x86_64.rpm
## 启动
systemctl daemon-reload
systemctl start go_web_server
systemctl status go_web_server --no-page
```