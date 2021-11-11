# ic-github-auth
This is a git authorization server for Internet Computer.


<!-- Account-id -->
<!-- 4f535f57f5f27585a27f722a809d090c399575784105cc1fbef86067ca6075bf -->

<!-- principal -->
<!-- ybrmj-jknrv-qf4f7-vk3jf-cevyz-ipheh-d5ucz-hqas7-3ilbn-4kyek-yae -->


dfx ledger --network ic create-canister <principal-identifier> --amount <icp-tokens>
dfx ledger --network ic create-canister ybrmj-jknrv-qf4f7-vk3jf-cevyz-ipheh-d5ucz-hqas7-3ilbn-4kyek-yae --amount 0.25

Transfer sent at BlockHeight: 1268939
Canister created with id: "ba3ps-giaaa-aaaag-qaafq-cai"

dfx identity --network ic deploy-wallet <canister-identifer>
dfx identity --network ic deploy-wallet ba3ps-giaaa-aaaag-qaafq-cai

Creating a wallet canister on the ic network.
The wallet canister on the "ic" network for user "default" is "ba3ps-giaaa-aaaag-qaafq-cai"

ba3ps-giaaa-aaaag-qaafq-cai

➜  testdeploy git:(master) dfx wallet --network ic balance
8484064427224 cycles.

https://<WALLET-CANISTER-ID>.raw.ic0.app
https://ba3ps-giaaa-aaaag-qaafq-cai.raw.ic0.app