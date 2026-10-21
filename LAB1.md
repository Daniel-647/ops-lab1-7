# 实验 1 记录：第 7 组

现场检查时将打开本文件，按其中的链接核对。证据以链接和实际终端输出给出。当前已补齐任务 0、初始 CI、主干制品追溯和任务 3 的规则验证；任务 4 之后的实验及分析仍待完成。main 受保护之后，修改本文件也通过 PR。

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
| 第一次 `make package` 得到的 BUILD_INFO（粘贴） | 见下方实际输出；两次的六行内容完全相同。 |
| 修改一行、不提交，再次打包：BUILD_INFO 是否变化？二进制文件的 sha256 是否变化？（贴出两次的 sha256） | `BUILD_INFO` 没有变化，二进制 SHA-256 发生变化，两次完整值见下表。 |
| 两次 `go version -m` 输出中的 vcs 行有何不同？ | `vcs`、`vcs.revision`、`vcs.time` 均相同，只有 `vcs.modified` 从 `false` 变为 `true`，实际输出见下方。 |
| 仅凭 BUILD_INFO，能否回答该包是什么、在哪台机器上、何时、用哪份代码构建？哪一项信息会出错？ | 能读出程序名、声明的版本号和构建号、HEAD 提交号及 Go 版本；没有记录构建机器和构建时间。第二次包的 `commit` 仍是原 HEAD，单凭它会漏掉未提交的注释修改，不能完整还原当时的源码状态。 |

### 实验条件和操作

2026-10-20 09:19:25（UTC+08:00）开始，在组内 Linux 主机的独立模板克隆中操作，来源为 [课程模板](https://github.com/zilinglius/operations-2026-lab1-sample)，提交号为 `7b63b41c13f061439997b0e6e4eaa2af922a7aa1`。这里的模板提交与小组仓库的提交历史不同，后面的 CI 追溯使用小组仓库。

两次均执行 `make package VERSION=1.0.0 BUILD_NO=1`，Go 版本为 `go1.26.9 linux/amd64`。第一次在干净工作区构建；第二次仅将 `conflict.go` 中的注释由“Weeks 表示一个教学班在哪些周上课。”改为“Weeks 表示教学班的上课周次（每周、单周或双周）。”，不提交，然后用相同参数构建。完成后已恢复注释，`git status --porcelain` 无输出。

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

我们原本只改了注释，但二进制的哈希还是变了。结合 VCS 输出看，这次构建记录的工作区状态从干净变成了有修改，所以不能把哈希不同直接解释成功能改变。`BUILD_INFO` 没有记录这个状态，第二次仍显示同一个提交号；如果只按这个提交检出代码，就会漏掉构建时存在的未提交修改。另一个要区分的是，`vcs.time` 是提交时间，不是本次构建时间。这里的 `1.0.0` 和构建号 `1` 都是手工传入，尚不表示已经发布了 GitHub Release。

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
| 未合入分支上的 tag 被 release gate 拒绝的运行 | <链接> |
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

仅凭构建号能否唯一确定是哪次运行？v1.0.0 的构建号在 `ci` 和 `release` 的运行列表中各对应哪次运行？不能。构建号在每个工作流中独立计数，v1.0.0 的构建号为 `1`，在 `ci` 中对应 [ci #1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39351112987)，在 `release` 中对应 [release #1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39375350868)。两次运行不同，需要用 `run_id` 定位构建该制品的具体运行。

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

| PR | 第一个提交的时间 | 合入时间 | 对应 Release 的发布时间 | 提交到发布的间隔 |
|---|---|---|---|---|
| #<编号> | <填写> | <填写> | <填写> | <填写> |
| #<编号> | <填写> | <填写> | <填写> | <填写> |
| #<编号> | <填写> | <填写> | <填写> | <填写> |

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

<填写>

**问题 2：PR 中的构建与 main 上的构建。** 在 PR 触发的运行日志中，`github.sha` 与 `pull_request.head.sha` 是否相同？它们能否在 main 的历史中找到？答案与合并方式有关，请写明本组使用的方式。合并之后，main 上这次改动的提交号是哪个？据此说明，为什么待发布的制品只能从 main 上的提交或 tag 构建。

<填写>

**问题 3：可重复构建。** 在本地同一个提交上按下表构建，记录哈希值。④、⑤ 需临时去掉 Makefile 中 tar 的 `--mtime=@0`，完成后用 `git checkout -- Makefile` 恢复，不要提交。

| 构建 | 条件 | 二进制 sha256 | tar.gz sha256 |
|---|---|---|---|
| ① | `make package`，`VERSION`、`BUILD_NO`、`RUN_ID`、`COMMIT` 均指定固定值 | <填写> | <填写> |
| ② | 与 ① 完全相同，再构建一次 | <填写> | <填写> |
| ③ | 仅将 `BUILD_NO` 换成另一个值 | <填写> | <填写> |
| ④ | 去掉 `--mtime=@0`，其余同 ① | <填写> | <填写> |
| ⑤ | 同 ④，间隔几秒后再构建一次 | <填写> | <填写> |
| Actions | 对 main 上的一次 push 运行执行 Re-run all jobs，比较两次尝试中 package 步骤打印的值 | <两次的值> | <两次的值> |

哪些相同、哪些不同，原因分别是什么？④ 的二进制文件与 ① 是否相同？如果不同，请结合任务 0 中看到的 vcs 信息分析。据此说明制品中应当包含什么、不应包含什么。另外，`release` 在同一个提交上再次编译，而没有直接发布 `ci` 产出的主干制品：为什么这样做？这与 L-05 第 11 页的“一次构建、多处部署”是否冲突，代价是什么？

<填写>

**问题 4：发布与回退。** 假设 v1.0.1 发布后发现缺陷，应当如何处理？为什么不应删除 v1.0.1 的 tag、再用同一版本号重新发布？任务 6 中将服务从 v1.0.1 切换回 v1.0.0 时，使用的是哪个文件，是否需要重新构建？请用两次 `/version` 的输出说明。如果改为从源码重新构建 v1.0.0，得到的文件与原附件是否相同（参考问题 3）？最后，如果有人删除 v1.0.0 的附件，换上一个新文件和对应的新 `.sha256`，`trace.sh` 能否发现？本组在第 4 节记录的值能否发现？可以查阅 GitHub 的 Immutable releases 设置的用途。

<填写>

## 7. 选做记录

| 选做 | 完成情况与证据 |
|---|---|
| O1 容器镜像 | <未做 / 链接、镜像摘要（digest）、运行容器的主机架构> |
| O2 缓存对反馈时间的影响 | <未做 / 预期、数据与结论> |
| O3 迁移到小组项目 | <未做 / 仓库链接> |

## 8. 成员贡献

每位成员一行，只填写 GitHub 用户名，无需填写姓名；行数应与小组人数一致。

| GitHub 用户名 | 发起并合入的 PR 编号 | 审批过的 PR 编号 |
|---|---|---|
| Daniel-647 | 待合入；计划负责任务 5 发布工作流 PR | 待填写 |
| QuQ-art | [#2](https://github.com/Daniel-647/ops-lab1-7/pull/2) | 待填写 |
| kanata-0527 | [#3](https://github.com/Daniel-647/ops-lab1-7/pull/3) | [#1](https://github.com/Daniel-647/ops-lab1-7/pull/1)、[#2](https://github.com/Daniel-647/ops-lab1-7/pull/2) |
| wenzhu4203 | [#5](https://github.com/Daniel-647/ops-lab1-7/pull/5) | 待填写 |
