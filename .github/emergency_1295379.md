## Emergency_1295379

**This guide assumes you are running the chain using cosmosvisor as explained [here](https://docs.kyve.network/getting-started/chain-node.html).**

On block #1295379 an error occurred in the end_block_logic which caused the chain to halt.
To recover from this error an emergency fix is required.
To apply the emergency fix the following commands need to be executed.

Stop the current chain binary. If you are running using the system-daemon, do
```shell
sudo systemctl stop kyved
```

Move the patch-binary manually and prepare cosmovisor.
```shell
mkdir -p ~/.kyve/cosmovisor/upgrades/emergency_1295379/bin
cd ~/.kyve/cosmovisor/upgrades/emergency_1295379
echo '{"name":"emergency_1295381","info":""}' > upgrade-info.json
cd bin
wget https://github.com/KYVENetwork/chain/releases/download/v0.5.3/chain_linux_amd64.tar.gz
tar -xvzf chain_linux_amd64.tar.gz
```
Check that the sha256 sum is correct:
```
echo "1d93f530e438da9459b79c67a3ea7423aad7b0e814154eb310685500fdb8a758 chain_linux_amd64.tar.gz" | sha256sum -c
```

If there are issues with the disk-space, disable the backup creation of cosmovisor.
Add
```sh
# This line is optional
Environment="UNSAFE_SKIP_BACKUP=true"
```
to the other environment variables in `/etc/systemd/system/kyved.service` and reload the service:
```shell
sudo systemctl daemon-reload
```
Remember to remove this line once it's processed if you want to keep the backup option enabled.

Then start cosmovisor:
```shell
sudo systemctl start kyved
```
Watch the log with
```shell
sudo journalctl -u kyved -f 
```
and see if the upgrade passes successfully (i.e. the chain does not crash).

We will wait until `5th June 2022 - 12:00 UTC` until we start the validators again, to give everybody time to perform the upgrade.



