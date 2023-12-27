hosts="di.arloor.dev us.arloor.dev ti.arloor.dev hk.arloor.dev sg.arloor.dev gg.arloor.dev bwg.arloor.dev"
for i in $hosts;do
ssh root@$i '
    hostname
    systemctl restart proxy
'
done

ssh root@us.arloor.dev '
    hostname
    systemctl restart guest
'