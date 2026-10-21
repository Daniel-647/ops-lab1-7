# 分析题 3：实际实验数据

实验日期：2026-10-21（UTC+08:00）。仅整理数据和观察，四道分析题正文由组长撰写。

## 云主机五组构建

隔离实验目录：云主机 `/root/lab1-repro-20261021`，由已有小组仓库克隆；原仓库和正式 Release 附件未改动。五组使用同一提交 `63f498dfa4ee4f844ba648d176078fd52d8083c7`（v1.0.1），Go 1.26.9、Linux/amd64、GNU tar 1.35。

固定参数：`VERSION=repro BUILD_NO=1 RUN_ID=repro COMMIT=63f498dfa4ee4f844ba648d176078fd52d8083c7`。③仅将 BUILD_NO 改为2；④⑤仅临时删除 Makefile 的 `--mtime=@0`，恢复 BUILD_NO=1。⑤在④结束后等待3秒。实验后 Makefile 已恢复，`git status --porcelain` 无输出，`git diff --exit-code` 为0。

| 构建 | 二进制 SHA-256 | tar.gz SHA-256 |
|---|---|---|
| ① 固定全部参数 | `e1e5a82da6459da0107574de6de18af541ba8ff34778d94782325eb8e5251f81` | `a0ee06d98079a445da548b930d4ddc9e7aaab24bf6ac3d593bcec8f14db22a67` |
| ② 与①相同 | `e1e5a82da6459da0107574de6de18af541ba8ff34778d94782325eb8e5251f81` | `a0ee06d98079a445da548b930d4ddc9e7aaab24bf6ac3d593bcec8f14db22a67` |
| ③ 仅 BUILD_NO=2 | `e36a14065a781ba99c42de1ec5b7a1bc9d7b314f11900f5b2cf7519b9f1e09d4` | `586647f62d8673d2accaa8c7947f5ec97c6fff29acc6ee07fec1e610495d7fe2` |
| ④ 去掉 --mtime=@0，其余同① | `30fdd1624a55afbc0d43943a411d72a1e2681734912e47e1ea11f379c9e5207b` | `4be6432b809828571f32b63ad951055805d007b908b694d1451baa0bbc5d2fd2` |
| ⑤ 同④，等待3秒后再构建 | `30fdd1624a55afbc0d43943a411d72a1e2681734912e47e1ea11f379c9e5207b` | `3668cc915b03d5dc8889a19dbdd806602bf11961e909572392ef29fecd85db42` |

原始输出见 [reproducible-builds-20261021.txt](reproducible-builds-20261021.txt)。

观察：①②两项均相同；③两项均与①不同；④⑤二进制彼此相同、均与①不同，tar.gz彼此不同。①②③的 `vcs.modified=false`；④⑤为 `true`，模块版本也从 `v1.0.1` 变为 `v1.0.1+dirty`。所以①与④的二进制差异不能直接归因于tar时间；改动Makefile引入了Go记录的工作区状态变化。④⑤的BUILD_INFO内容相同，但文件时间不同；归档未再统一这些时间。上述是本次给定工具与参数下的观察，不能推出任意环境或任意Go版本都一定相同。

## GitHub Actions 重跑

本地五组使用v1.0.1提交；此CI重跑使用另一个main提交，两组实验分别内部比较，不互相比哈希。

- 运行：[ci #20 / run 39582697661](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39582697661)。
- 提交：`e05eae7164fc82040063ed96bbc80a5e16387c6d`。
- 两次参数均为 VERSION=0.0.0-main、BUILD_NO=20、RUN_ID=39582697661、COMMIT为上述提交；重跑的attempt从1变2。
- [attempt 1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39582697661/attempts/1) 与 [attempt 2](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39582697661/attempts/2) 均成功。

| 对象 | attempt 1 | attempt 2 |
|---|---|---|
| 二进制 | `3ea09987196f1da031ed371ab6358c31b7f5aacce275dea00a34e67b6d17aca0` | `3ea09987196f1da031ed371ab6358c31b7f5aacce275dea00a34e67b6d17aca0` |
| tar.gz | `76cdc9265e23280ee40fc2f8b43efd49e840950e02217d8e81624f4d3aabce37` | `76cdc9265e23280ee40fc2f8b43efd49e840950e02217d8e81624f4d3aabce37` |

日志中的上传artifact ZIP摘要不同；它是上传步骤生成的外层容器摘要，应与表中打包步骤的二进制、tar.gz哈希区分。原始日志摘录见 [ci-rerun-39582697661.txt](ci-rerun-39582697661.txt)，包含外层摘要供核查。

## 后续写作仍需覆盖

以上补齐题目3的实验数据；最终正文还需解释应记录的追溯信息、应排除的易变信息，以及为何release重新构建，与L-05第11页“一次构建、多处部署”的关系和代价。此文件不代替最终LAB1.md。
