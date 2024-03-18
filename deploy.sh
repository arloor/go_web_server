#! /bin/bash
hosts="hk.arloor.dev sg.arloor.dev di.arloor.dev us.arloor.dev gg.arloor.dev"
# echo "" > ~/.ssh/known_hosts
# for i in ${hosts}; do
#     ssh-keyscan -H ${i} >> ~/.ssh/known_hosts
# done
for i in ${hosts}; do
    ssh -o StrictHostKeyChecking=no root@${i} '
            hostname;
            systemctl restart proxygo;
            podman image prune -f 2>/dev/null
            podman images --digests |grep arloor/go_web_server|awk "{print \$4\" \"\$3}";
            '
done
