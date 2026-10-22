# 实验 1 记录：第 7 组

本文件记录本组的构建基线、CI 与门禁演练、两个版本的发布追溯和运行回退，以及实验数据与四项分析。现场检查可按各节链接和实际终端输出核对。main 受保护，文档更新也通过 PR 合入。

## 1. 基本信息

| 项 | 内容 |
|---|---|
| 小组 | 第 7 组，共 4 人 |
| 仓库地址 | <https://github.com/Daniel-647/ops-lab1-7> |
| 代码载体 | 统一示例 `enroll-conflict`，Go |
| runner 与语言版本 | `ubuntu-24.04`，Go `1.26.9`；组内 Linux 主机为 Ubuntu 24.04.4 LTS、x86-64，Go 同为 `1.26.9` |
| 合并方式 | squash；仓库设置与规则集均只允许 Squash |

## 2. 任务 0：本地基线

| 问题 | 观察结果 |
|---|---|
| 第一次 `make package` 得到的 BUILD_INFO（粘贴） | 见下方实际输出 |
| 修改一行、不提交，再次打包：BUILD_INFO 是否变化？二进制文件的 sha256 是否变化？（贴出两次的 sha256） | `BUILD_INFO` 没有变化，二进制 SHA-256 发生变化，两次完整值见下表。 |
| 两次 `go version -m` 输出中的 vcs 行有何不同？ | `vcs`、`vcs.revision`、`vcs.time` 均相同，只有 `vcs.modified` 从 `false` 变为 `true`，实际输出见下方。 |
| 仅凭 BUILD_INFO，能否回答该包是什么、在哪台机器上、何时、用哪份代码构建？哪一项信息会出错？ | 该包是什么：能，程序为 `enroll-conflict`，版本为 `1.0.0`，构建号为 `1`。在哪台机器上构建：不能，没有记录构建机器。何时构建：不能，没有记录构建时间。用哪份代码构建：只能知道构建时的 HEAD 提交号，不能知道是否包含未提交的修改。会出错的项：对构建所用代码的判断；第二次构建包含未提交的修改，但 `commit` 没变，若将它当成完整的源码标识，就会误认为两次用了同一份代码。 |

两次 `BUILD_INFO` 的实际内容：

```text
name=enroll-conflict
version=1.0.0
build_no=1
run_id=local
commit=7b63b41c13f061439997b0e6e4eaa2af922a7aa1
go=go1.26.9
```

| 构建 | 二进制 SHA-256 |
|---|---|
| 第一次，工作区干净 | `f23447d406c42af7dd35e23c2d1b2f15484cfb709cc9b4bb1262643694b8528e` |
| 第二次，注释已改但未提交 | `5049293df8a6eb201d310b3f6cef0d405c26b5ece7f07efa71bfd71204905c90` |

第一次 `go version -m` 的 VCS 行：

```text
build vcs=git
build vcs.revision=7b63b41c13f061439997b0e6e4eaa2af922a7aa1
build vcs.time=2026-10-12T16:26:07Z
build vcs.modified=false
```

第二次：

```text
build vcs=git
build vcs.revision=7b63b41c13f061439997b0e6e4eaa2af922a7aa1
build vcs.time=2026-10-12T16:26:07Z
build vcs.modified=true
```
完整终端输出见 [任务 0 原始记录](evidence/task0-baseline.txt)。

## 3. 证据

| 项 | 链接或文本 |
|---|---|
| 规则集里启用的规则 | [main-gate](https://github.com/Daniel-647/ops-lab1-7/rules)：Active，目标为默认分支 main，Bypass list 为空；启用 Restrict deletions、Block force pushes、Require a pull request before merging（1 人审批，仅 Squash）、Require status checks to pass（`ci`，来源 GitHub Actions）。 |
| 初始 main CI | [ci #1，attempt 1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39351112987)，push 事件，成功；这是设置门禁之前允许的初始 CI 提交。 |
| 任务 2 主干制品追溯 | SHA-256 校验通过，`trace.sh` 退出码为 0；[实际输出](evidence/ci-1-trace.txt)，详情见本节下方。 |
| 组长直接推送 main 被拒的终端输出 | 见下方实际输出及 [完整记录](evidence/main-gate-direct-push.txt)。 |
| 第一轮：失败的 PR 运行 | [PR #2](https://github.com/Daniel-647/ops-lab1-7/pull/2)；[失败运行 ci #4](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39359794537)，在 Run tests 步骤失败。  |
| 第一轮：修复后通过的 PR 运行 |  [修复提交 cd0f736](https://github.com/Daniel-647/ops-lab1-7/commit/cd0f73673752585c27cab631d88ce8762a0c7f4d)；[通过运行 ci #5](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39361867498)。 |
| 第一轮：合入 main 后的运行（含制品） | [main 上的成功运行 ci #6](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362282086)。 |
| 第二轮：找到的缺陷，以及现有测试为何无法发现它 | 将 `conflict.go` 的时段合法性判断由 `end < start` 改为 `end <= start`，使合法的单节课（如 `start=3, end=3`）被判为非法，接口返回 HTTP 400。原有合法样例未覆盖 `start == end` 的情况，因此旧测试仍然通过。 |
| 第二轮：仅含缺陷时仍然通过的运行 | [ci #7](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362790734) |
| 第二轮：补测试后失败的运行 | [ci #8](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362976846) |
| 第二轮：撤销缺陷后通过的运行 | [ci #9](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39363082268) |
| 未合入分支上的 tag 被 release gate 拒绝的运行 | [release #2：v0.0.1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39375781855)。tag 指向未合入 main 的提交 `a5945912371c6e5b0e354aacd417782d43d5b8e1`，在 Release gate 步骤以退出码 1 被拒绝，报错为 `Release rejected: tag commit is not in main history.`；后续打包和发布步骤跳过，未产生 Release，测试 tag 和临时分支已删除。 |
| Release v1.0.0 | [v1.0.0](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.0) |
| Release v1.0.1 | [v1.0.1](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.1) |

### 管理员直接推送验证

验证时 GitHub 认证身份为组长 `Daniel-647`。在独立临时克隆中创建一个空提交，执行 `git push origin HEAD:refs/heads/main`，尝试直接更新 main。实际返回：

```text
remote: error: GH013: Repository rule violations found for refs/heads/main.
remote: Review all repository rules at https://github.com/Daniel-647/ops-lab1-7/rules?ref=refs%2Fheads%2Fmain
remote:
remote: - Changes must be made through a pull request.
remote:
remote: - Required status check "ci" is expected.
remote:
To https://github.com/Daniel-647/ops-lab1-7.git
 ! [remote rejected] HEAD -> main (push declined due to repository rule violations)
error: failed to push some refs to 'https://github.com/Daniel-647/ops-lab1-7.git'
```

`git push` 退出码为 1，拒绝原因是仓库规则。验证前后远程 main 均为 `8a7ecb744b441462437e3c02e63a09803bf04eb1`，没有增加该空提交。这说明组长账号也受门禁约束。

### 初始 CI 制品检查

从初始 CI 的 Artifacts 下载 ZIP，解压后将其中的 `.tar.gz` 和 `.sha256` 上传到组内 Linux 主机。运行追溯前已拉取小组仓库最新 main；压缩包在 Mac 和云主机上算出的 SHA-256 均为：

```text
45cb31a4450f57b1f10ff44481c0c9199f55523f75a0d5ce917ec22d93ed71a1
```

`trace.sh` 读出的 `BUILD_INFO`：

```text
name=enroll-conflict
version=0.0.0-main
build_no=1
run_id=39351112987
commit=8a7ecb744b441462437e3c02e63a09803bf04eb1
go=go1.26.9
```

脚本确认校验和正确，提交位于 `origin/main` 历史上，退出码为 0；输出中的 [提交链接](https://github.com/Daniel-647/ops-lab1-7/commit/8a7ecb744b441462437e3c02e63a09803bf04eb1) 和 [构建链接](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39351112987) 与 GitHub 实际记录一致。这个提交是初始 CI 的直接提交，没有对应 PR；这是任务 2 的主干制品检查，任务 6 的两个 Release 仍需另行追溯。

## 4. 追溯与运行

对 v1.0.0 和 v1.0.1 各运行一次 `scripts/trace.sh`，粘贴输出，并填写下表。

| 项 | v1.0.0 | v1.0.1 |
|---|---|---|
| 版本号 | `1.0.0` | `1.0.1` |
| 构建号 | `1` | `3` |
| 提交号（SHA） | `f394c6b09e233af442c675eaa95f6702bb78e433` | `63f498dfa4ee4f844ba648d176078fd52d8083c7` |
| 引入该提交的 PR | [#4](https://github.com/Daniel-647/ops-lab1-7/pull/4) | [#5](https://github.com/Daniel-647/ops-lab1-7/pull/5) |
| 构建它的 Actions 运行（`trace.sh` 给出的链接） | [release #1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39375350868) | [release #3](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39433897261) |
| 制品 sha256（发布当时，取自 Release 里的 .sha256） | `40b6bbb4dc1b4875b84c7dfd3d2cec88fee137d4f64ac1999eb224bceb4048af` | `1bc77fb3b9095cb85d2ac5489a87c88b02953e241add2cb6318ff1032d10b827` |

仅凭构建号能否唯一确定是哪次运行？v1.0.0 的构建号在 `ci` 和 `release` 的运行列表中各对应哪次运行？
不能。构建号在每个工作流中独立计数，v1.0.0 的构建号为 `1`，在 `ci` 中对应 [ci #1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39351112987)，在 `release` 中对应 [release #1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39375350868)。两次运行不同，需要用 `run_id` 定位构建该制品的具体运行。

发布 v1.0.1 之后重新下载 v1.0.0 的附件，计算得到的 sha256：`40b6bbb4dc1b4875b84c7dfd3d2cec88fee137d4f64ac1999eb224bceb4048af`。与上表中发布当时的值是否一致？一致。

`trace.sh` 输出（v1.0.0）：

```text
== 1. 校验和
enroll-conflict_1.0.0_linux_amd64.tar.gz: OK
== 2. 制品里的 BUILD_INFO
name=enroll-conflict
version=1.0.0
build_no=1
run_id=39375350868
commit=f394c6b09e233af442c675eaa95f6702bb78e433
go=go1.26.9
== 3. 提交 f394c6b09e233af442c675eaa95f6702bb78e433
作者: Daniel-647 <li-yixi@sjtu.edu.cn>
时间: 2026-10-20T13:29:14+08:00
标题: ci: add gated release workflow (#4)
主干: 该提交在 origin/main 的历史上
指向它的 tag: v1.0.0
版本号: 1.0.0    构建号: 1
提交页面（含关联的 PR）: https://github.com/Daniel-647/ops-lab1-7/commit/f394c6b09e233af442c675eaa95f6702bb78e433
构建它的运行: https://github.com/Daniel-647/ops-lab1-7/actions/runs/39375350868
v1.0.0 trace退出码=0
```

`trace.sh` 输出（v1.0.1）：

```text
== 1. 校验和
enroll-conflict_1.0.1_linux_amd64.tar.gz: OK
== 2. 制品里的 BUILD_INFO
name=enroll-conflict
version=1.0.1
build_no=3
run_id=39433897261
commit=63f498dfa4ee4f844ba648d176078fd52d8083c7
go=go1.26.9
== 3. 提交 63f498dfa4ee4f844ba648d176078fd52d8083c7
作者: wenzhu4203 <wenzhu4203@gmail.com>
时间: 2026-10-20T21:47:03+08:00
标题: docs: add running version query instructions (#5)
主干: 该提交在 origin/main 的历史上
指向它的 tag: v1.0.1
版本号: 1.0.1    构建号: 3
提交页面（含关联的 PR）: https://github.com/Daniel-647/ops-lab1-7/commit/63f498dfa4ee4f844ba648d176078fd52d8083c7
构建它的运行: https://github.com/Daniel-647/ops-lab1-7/actions/runs/39433897261
v1.0.1 trace退出码=0
```

任务 6 第 4 步，两次 `/version` 的输出：

```text
运行 v1.0.1 时：
{"version":"1.0.1","commit":"63f498dfa4ee4f844ba648d176078fd52d8083c7","build_no":"3","go":"go1.26.9"}
换回 v1.0.0 后：
{"version":"1.0.0","commit":"f394c6b09e233af442c675eaa95f6702bb78e433","build_no":"1","go":"go1.26.9"}
```

## 5. 数据

时间均从 GitHub 页面读取，并注明时区。

所有时间均为 2026-10-20，时区为 UTC+08:00（北京时间）。
CI 结束时间统一采用 GitHub 记录的 ci 任务完成时间。

第一轮 gate 演练的时间线：

| 事件 | 时间 | 链接 |
|---|---|---|
| 发起 PR | 2026-10-20 10:35:57 | [PR #2](https://github.com/Daniel-647/ops-lab1-7/pull/2) |
| 第一次 ci 失败（运行结束） | 2026-10-20 10:36:16 | [ci #4](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39359794537) |
| 推送修复提交 | 2026-10-20 10:59:10 | [修复提交 cd0f736](https://github.com/Daniel-647/ops-lab1-7/commit/cd0f73673752585c27cab631d88ce8762a0c7f4d) |
| ci 通过（运行结束） | 2026-10-20 10:59:32 | [ci #5](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39361867498) |
| 合入 main | 2026-10-20 11:03:41 | [PR #2](https://github.com/Daniel-647/ops-lab1-7/pull/2) |

提交到发布的间隔，至少记录三个合入的 PR：

首个提交时间统一取作者时间（AuthorDate），对应 Release 取首次包含该 PR 改动的版本；间隔为发布时间减去首个提交时间。

| PR | 第一个提交的时间 | 合入时间 | 对应 Release 的发布时间 | 提交到发布的间隔 |
|---|---|---|---|---|
| [#2](https://github.com/Daniel-647/ops-lab1-7/pull/2) | 2026-10-20 10:31:32 | 2026-10-20 11:03:41 | [v1.0.0](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.0)，2026-10-20 13:31:42 | 3 小时 0 分 10 秒 |
| [#3](https://github.com/Daniel-647/ops-lab1-7/pull/3) | 2026-10-20 10:16:42 | 2026-10-20 11:54:16 | [v1.0.0](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.0)，2026-10-20 13:31:42 | 3 小时 15 分 0 秒 |
| [#4](https://github.com/Daniel-647/ops-lab1-7/pull/4) | 2026-10-20 13:20:14 | 2026-10-20 13:29:14 | [v1.0.0](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.0)，2026-10-20 13:31:42 | 11 分 28 秒 |
| [#5](https://github.com/Daniel-647/ops-lab1-7/pull/5) | 2026-10-20 21:29:00 | 2026-10-20 21:47:03 | [v1.0.1](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.1)，2026-10-20 21:50:45 | 21 分 45 秒 |

`ci` 工作流的运行时长（从 Actions 页面读取），至少 5 次，全部列出，不要只给出平均值：

| 运行 | 事件（PR / push） | 时长 | 备注 |
|---|---|---|---|
| [ci #1 / attempt 1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39351112987) | push → main | 总耗时 46 秒；`ci` job 42 秒 | 初始 CI，通过；总时段为 2026-10-20 09:01:47～09:02:33（UTC+08:00），job 为 09:01:50～09:02:32。|
| [ci #4](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39359794537) | PR #2 | ci 任务 12 秒 | 第一轮修复前，失败；ci 任务起止为 2026-10-20 10:36:04～10:36:16（UTC+08:00），失败步骤为 Run tests。 |
| [ci #5](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39361867498) | PR #2（推送修复后） | ci 任务 15 秒 | 第一轮修复后，通过；ci 任务起止为 2026-10-20 10:59:17～10:59:32（UTC+08:00）。 |
| [ci #6](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362282086) | push → main | ci 任务 14 秒 | PR #2 合入 main 后，通过；ci 任务起止为 2026-10-20 11:03:45～11:03:59（UTC+08:00）。 |
| [ci #7](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362790734) | PR #3 | ci 任务 10 秒 | 第二轮仅含缺陷，旧测试漏检，通过；ci 任务起止为 2026-10-20 11:09:23～11:09:33（UTC+08:00）。 |
| [ci #8](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362976846) | PR #3（补充测试后） | ci 任务 9 秒 | 第二轮补充测试后，失败；ci 任务起止为 2026-10-20 11:11:28～11:11:37（UTC+08:00），失败步骤为 Run tests。 |
| [ci #9](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39363082268) | PR #3（撤销缺陷后） | ci 任务 13 秒 | 第二轮撤销缺陷后，通过；ci 任务起止为 2026-10-20 11:12:41～11:12:54（UTC+08:00）。 |
| [ci #12](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39366725871) | push → main | ci 任务 17 秒 | PR #3 合入 main 后，通过；ci 任务起止为 2026-10-20 11:54:21～11:54:38（UTC+08:00）。 |

## 6. 分析问题

每题先列出所用数据（可引用上面的表），再用一段话（4～8 句）写出结论，并说明结论的成立条件。

**问题 1：gate 的收益与成本。** 两轮演练中，gate 拦住了什么，没有拦住什么？用第 5 节的时间线计算：第一轮因失败多花了多少时间，其中等待 `ci` 的时间和人工修改的时间各占多少？一次 `ci` 大约需要多久，在“提交到发布的间隔”中占多大比例？第二轮补充测试之后，还有哪类问题是测试本身无法发现的？

**所用数据：**第 3 节两轮演练记录、第 5 节时间线和 8 次 CI 时长；修复推送时间及计算依据见 [原始数据](evidence/analysis-data-20261021.json) 和 [时间拆分](evidence/analysis-data-20261021.md)。

第一轮把单周和双周错误地判为冲突，已有测试让 ci #4 失败，门禁阻止了合并，修复后 ci #5 才通过。第二轮把合法单节课判为非法，ci #7 却通过了，说明门禁只能拦住检查能够发现的问题；补上单节课测试后 ci #8 失败，撤销缺陷后 ci #9 通过。从首次失败的 10:36:16 到修复通过的 10:59:32，共用了 23 分 16 秒，这是本次观测到的修复反馈周期，不能直接当作与“没有失败”相比精确多花的时间。其中，失败后到修复推送用了 22 分 54 秒，占 98.42%，可以近似表示人工处理区间，但也包含沟通、空闲和本地验证。修复推送后等待 CI 反馈用了 22 秒，占 1.58%，其中推送到创建运行 5 秒、创建到 job 开始 2 秒、job 执行 15 秒。第 5 节 8 次 ci job 用时 9～42 秒，平均 16.5 秒；以一次 15 秒的 job 对比 PR #2 提交到发布的 10810 秒，占约 0.139%，这只是单次检查的占比。补充测试仍不能保证需求理解正确，也不能发现未覆盖的并发、性能或运行环境问题，还需要评审和针对性的验证。

**问题 2：PR 中的构建与 main 上的构建。** 在 PR 触发的运行日志中，`github.sha` 与 `pull_request.head.sha` 是否相同？它们能否在 main 的历史中找到？答案与合并方式有关，请写明本组使用的方式。合并之后，main 上这次改动的提交号是哪个？据此说明，为什么待发布的制品只能从 main 上的提交或 tag 构建。

**所用数据：**PR #2 的 [ci #5](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39361867498)、PR #3 的 [ci #9](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39363082268)，以及对应 main 运行 [ci #6](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362282086)、[ci #12](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39366725871)；[日志及祖先检查](evidence/pr-main-sha-20261021.txt)。

| 提交标识 | PR #2 | PR #3 |
|---|---|---|
| PR 运行的 `github.sha` | `248fa45ffa40c1667e6e9e03a6651dea5bb3cbeb` | `0b6d625d3be205e3b7ab7332c304f15f281ce7a4` |
| `pull_request.head.sha` | `cd0f73673752585c27cab631d88ce8762a0c7f4d` | `422dfa06fa3e5a67a6991a43a2dd4fb0edc0b3bd` |
| Squash 合入 main 后的 SHA | `392507cab237a19aa6f829ac918a59aaaf2ff3e0` | `3d94e8aa5601d5af5f3b68b2719b0263aab43408` |

本组两次 PR 运行中的 `github.sha` 都与 `pull_request.head.sha` 不同。前者对应 GitHub 把 PR 分支与当时目标分支临时合并后用于检查的提交，后者是 PR 分支本身的最新提交。我们使用 Squash 合并，GitHub 会把 PR 的改动压成一个新的 main 提交，所以上表每个 PR 的三个 SHA 都不同。对本组这四个 PR 阶段的 SHA 执行祖先检查，结果均为不在 main 历史中，而两个 Squash 后的 SHA 均在 main 历史中；这些提交对象本身都存在，并非查不到。PR 阶段通过检查，只能说明当时被检查的组合通过，不能把它的 SHA 当成最终 main 的提交号。因此，待发布的制品要从已合入 main 的提交或指向这类提交的 tag 构建，并在该提交上再次检查，才能让制品对应正式接受的源码。这一结论中的具体 SHA 关系取决于本组的 Squash 方式，不能直接套用到其他合并方式。

**问题 3：可重复构建。** 在本地同一个提交上按下表构建，记录哈希值。④、⑤ 需临时去掉 Makefile 中 tar 的 `--mtime=@0`，完成后用 `git checkout -- Makefile` 恢复，不要提交。

**实验条件与数据：**五组构建在组内 Linux 主机的独立克隆中完成，提交为 `63f498dfa4ee4f844ba648d176078fd52d8083c7`，Go 1.26.9、GNU tar 1.35；①②④⑤均固定 `VERSION=repro BUILD_NO=1 RUN_ID=repro COMMIT=上述提交`，③仅改为 `BUILD_NO=2`，⑤在④结束后等待 3 秒，完成后已恢复 Makefile。CI 重跑则使用提交 `e05eae7164fc82040063ed96bbc80a5e16387c6d`，两次均为 `VERSION=0.0.0-main BUILD_NO=20 RUN_ID=39582697661`；两组实验各自在相同输入下比较，不互相比哈希。[本地构建原始输出](evidence/reproducible-builds-20261021.txt)、[CI 重跑日志摘录](evidence/ci-rerun-39582697661.txt)。

| 构建 | 条件 | 二进制 sha256 | tar.gz sha256 |
|---|---|---|---|
| ① | `make package`，`VERSION`、`BUILD_NO`、`RUN_ID`、`COMMIT` 均指定固定值 | `e1e5a82da6459da0107574de6de18af541ba8ff34778d94782325eb8e5251f81` | `a0ee06d98079a445da548b930d4ddc9e7aaab24bf6ac3d593bcec8f14db22a67` |
| ② | 与 ① 完全相同，再构建一次 | `e1e5a82da6459da0107574de6de18af541ba8ff34778d94782325eb8e5251f81` | `a0ee06d98079a445da548b930d4ddc9e7aaab24bf6ac3d593bcec8f14db22a67` |
| ③ | 仅将 `BUILD_NO` 换成另一个值 | `e36a14065a781ba99c42de1ec5b7a1bc9d7b314f11900f5b2cf7519b9f1e09d4` | `586647f62d8673d2accaa8c7947f5ec97c6fff29acc6ee07fec1e610495d7fe2` |
| ④ | 去掉 `--mtime=@0`，其余同 ① | `30fdd1624a55afbc0d43943a411d72a1e2681734912e47e1ea11f379c9e5207b` | `4be6432b809828571f32b63ad951055805d007b908b694d1451baa0bbc5d2fd2` |
| ⑤ | 同 ④，间隔几秒后再构建一次 | `30fdd1624a55afbc0d43943a411d72a1e2681734912e47e1ea11f379c9e5207b` | `3668cc915b03d5dc8889a19dbdd806602bf11961e909572392ef29fecd85db42` |
| Actions | [ci #20 attempt 1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39582697661/attempts/1) 与 [attempt 2](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39582697661/attempts/2) | attempt 1：`3ea09987196f1da031ed371ab6358c31b7f5aacce275dea00a34e67b6d17aca0`<br>attempt 2：`3ea09987196f1da031ed371ab6358c31b7f5aacce275dea00a34e67b6d17aca0` | attempt 1：`76cdc9265e23280ee40fc2f8b43efd49e840950e02217d8e81624f4d3aabce37`<br>attempt 2：`76cdc9265e23280ee40fc2f8b43efd49e840950e02217d8e81624f4d3aabce37` |

哪些相同、哪些不同，原因分别是什么？④ 的二进制文件与 ① 是否相同？如果不同，请结合任务 0 中看到的 vcs 信息分析。据此说明制品中应当包含什么、不应包含什么。另外，`release` 在同一个提交上再次编译，而没有直接发布 `ci` 产出的主干制品：为什么这样做？这与 L-05 第 11 页的“一次构建、多处部署”是否冲突，代价是什么？

①②的二进制和压缩包哈希完全相同，CI 两次 attempt 的两项哈希也相同，说明在本次工具、源码状态和参数保持一致的条件下可以重复得到同样的输出。③只改变 BUILD_NO，但构建号既被写入二进制，也写入 BUILD_INFO，所以二进制和压缩包都变了。④的二进制也与①不同：虽然只改了 Makefile 的打包参数，工作区已经不干净，Go 记录的 `vcs.modified` 从 false 变成 true，模块版本带上 `+dirty`，这与任务 0 只改注释仍改变二进制的现象一致。④⑤的二进制彼此相同，但去掉 `--mtime=@0` 后归档保留了文件的修改时间，两次构建的时间不同，因此压缩包不同。正式制品应从干净工作区构建，保留版本、提交、构建号、运行 ID 和工具版本等追溯信息并附校验和，而不应把每次变化的当前时间、机器绝对路径等无关信息直接写入制品。我们的 main CI 使用 `0.0.0-main`，release 使用 tag 中的正式版本号，而且两个工作流各有构建号和运行 ID，因此当前实现会按发布参数重新编译，不能把 CI 包原样改名就当成 Release。对照 L-05 第 11 页，这没有严格做到“CI 制品直接作为发布制品”，多付出了一次构建、检查的时间，也需要承担两次构建环境变化造成差异的风险。发布包确定后再部署到不同环境或用于回退时，应使用同一个已校验的包，只改变环境配置，不再逐个环境重新构建。

**问题 4：发布与回退。** 假设 v1.0.1 发布后发现缺陷，应当如何处理？为什么不应删除 v1.0.1 的 tag、再用同一版本号重新发布？任务 6 中将服务从 v1.0.1 切换回 v1.0.0 时，使用的是哪个文件，是否需要重新构建？请用两次 `/version` 的输出说明。如果改为从源码重新构建 v1.0.0，得到的文件与原附件是否相同（参考问题 3）？最后，如果有人删除 v1.0.0 的附件，换上一个新文件和对应的新 `.sha256`，`trace.sh` 能否发现？本组在第 4 节记录的值能否发现？可以查阅 GitHub 的 Immutable releases 设置的用途。

**所用数据：**第 4 节两个 Release 的原始哈希、追溯结果和两次 `/version` 输出；问题 3 的重复构建结果。

如果 v1.0.1 发布后发现缺陷，应先确认旧版适用于当前环境，回退恢复服务，再通过 PR 修复并发布新版本，例如 v1.0.2。不能删除 v1.0.1 的 tag 再复用这个版本号，否则同一版本名会对应不同代码或文件，已经下载旧文件的人与后来下载的人拿到的内容可能不同，问题也难以复现。我们这次回退使用原先发布的 `enroll-conflict_1.0.0_linux_amd64.tar.gz` 中的可执行文件，不需要重新构建。切换前 `/version` 返回 `1.0.1 / 63f498d… / build_no=3`，切换后返回 `1.0.0 / f394c6b… / build_no=1`，说明本次运行已经切回旧版。如果从源码重建 v1.0.0，只有源码状态、工具和构建参数等输入都匹配时才可能得到原文件，仅检出同一个提交并不能保证哈希相同。如果有人同时替换附件和对应的 .sha256，并保留能通过主干检查的 BUILD_INFO，trace.sh 可能仍然通过，因为它只与当前提供的校验文件比较，不会核对发布时保存的旧值，也不会验证二进制确实由声明的源码构建。在第 4 节原始哈希记录可信且未被一起修改的前提下，将新下载文件与该记录比较可以发现文件被替换。[GitHub Immutable releases](https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases) 可在发布后锁定附件和关联 tag，防止直接修改、删除附件或移动 tag，从源头减少这种替换的机会。

## 7. 选做记录

| 选做 | 完成情况与证据 |
|---|---|
| O1 容器镜像 | 未做 |
| O2 缓存对反馈时间的影响 | 未做 |
| O3 迁移到小组项目 | 未做 |

## 8. 成员贡献

每位成员一行，只填写 GitHub 用户名，无需填写姓名；行数应与小组人数一致。

| GitHub 用户名 | 发起并合入的 PR 编号 | 审批过的 PR 编号 |
|---|---|---|
| Daniel-647 | [#1](https://github.com/Daniel-647/ops-lab1-7/pull/1)、[#4](https://github.com/Daniel-647/ops-lab1-7/pull/4) | [#5](https://github.com/Daniel-647/ops-lab1-7/pull/5#pullrequestreview-5585886164)、[#6](https://github.com/Daniel-647/ops-lab1-7/pull/6#pullrequestreview-5597215159) |
| QuQ-art | [#2](https://github.com/Daniel-647/ops-lab1-7/pull/2)、[#7](https://github.com/Daniel-647/ops-lab1-7/pull/7) | [#3](https://github.com/Daniel-647/ops-lab1-7/pull/3#pullrequestreview-5580572256)、[#4](https://github.com/Daniel-647/ops-lab1-7/pull/4#pullrequestreview-5581082529)、[#8](https://github.com/Daniel-647/ops-lab1-7/pull/8#pullrequestreview-5599977107) |
| kanata-0527 | [#3](https://github.com/Daniel-647/ops-lab1-7/pull/3)、[#8](https://github.com/Daniel-647/ops-lab1-7/pull/8) | [#1](https://github.com/Daniel-647/ops-lab1-7/pull/1#pullrequestreview-5579769549)、[#2](https://github.com/Daniel-647/ops-lab1-7/pull/2#pullrequestreview-5580286961)、[#7](https://github.com/Daniel-647/ops-lab1-7/pull/7#pullrequestreview-5598123857) |
| wenzhu4203 | [#5](https://github.com/Daniel-647/ops-lab1-7/pull/5)、[#6](https://github.com/Daniel-647/ops-lab1-7/pull/6) | [#9](https://github.com/Daniel-647/ops-lab1-7/pull/9) |
