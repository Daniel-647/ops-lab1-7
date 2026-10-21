# 分析题 1、2 与报告第 5 节的数据

核查日期：2026-10-21。此文件整理实测数据、可复核计算和写作边界，分析题正文仍由组长撰写；不代替最终LAB1.md。全部展示时间均为UTC+08:00（北京时间），原始API时间见同目录 [analysis-data-20261021.json](analysis-data-20261021.json)。

## 三个以上 PR 的提交到发布间隔

统计口径：起点采用每个PR最终保留提交列表中首个提交的AuthorDate（作者时间），终点采用最早包含该PR合并提交的Release的published_at。Git祖先检查确认PR #2、#3、#4的合并提交均包含在v1.0.0对应提交中，PR #5对应v1.0.1。采用提交时间作为代理，不能视为首次写代码或开始工作的真实时间。

| PR | 首个提交的作者时间 | 合入时间 | 首个包含该改动的Release及发布时间 | 提交到发布间隔 |
|---|---|---|---|---|
| [#2](https://github.com/Daniel-647/ops-lab1-7/pull/2) | 2026-10-20 10:31:32 | 2026-10-20 11:03:41 | [v1.0.0](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.0)：2026-10-20 13:31:42 | 3小时0分10秒 |
| [#3](https://github.com/Daniel-647/ops-lab1-7/pull/3) | 2026-10-20 10:16:42 | 2026-10-20 11:54:16 | [v1.0.0](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.0)：2026-10-20 13:31:42 | 3小时15分0秒 |
| [#4](https://github.com/Daniel-647/ops-lab1-7/pull/4) | 2026-10-20 13:20:14 | 2026-10-20 13:29:14 | [v1.0.0](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.0)：2026-10-20 13:31:42 | 0小时11分28秒 |
| [#5](https://github.com/Daniel-647/ops-lab1-7/pull/5) | 2026-10-20 21:29:00 | 2026-10-20 21:47:03 | [v1.0.1](https://github.com/Daniel-647/ops-lab1-7/releases/tag/v1.0.1)：2026-10-20 21:50:45 | 0小时21分45秒 |

PR #3 的首个提交 `592373f01abd35ff07bc80884e9e0239c0c3acf7`，AuthorDate为2026-10-20 10:16:42，CommitDate（提交者时间）为11:08:18。上表统一采用AuthorDate，所以为3小时15分；如果统一改用CommitDate，则该行应为2小时23分24秒，其余三行不变。两种时间不能不注明就混用，也不能把其中任意一种当成推送时间。该数据集包含流程演练和发布等待，不能据此推断一般项目的生产交付效率。

## CI 时长：十二次实际运行

每行均为attempt 1。时长统一按ci job的completed_at减started_at计算；排队列为job started_at减workflow created_at，与运行执行时间分列。样本覆盖初始构建、两轮演练与发布准备，不是缓存开关对照，不用于替代O2实验。

| 运行 | 事件 | 结果 | job开始 | job结束 | job时长 | 创建到job开始 |
|---|---|---|---|---|---|---|
| [ci #1](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39351112987) | push → main | success | 2026-10-20 09:01:50 | 2026-10-20 09:02:32 | 42 秒 | 3 秒 |
| [ci #4](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39359794537) | PR | failure | 2026-10-20 10:36:04 | 2026-10-20 10:36:16 | 12 秒 | 2 秒 |
| [ci #5](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39361867498) | PR | success | 2026-10-20 10:59:17 | 2026-10-20 10:59:32 | 15 秒 | 2 秒 |
| [ci #6](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362282086) | push → main | success | 2026-10-20 11:03:45 | 2026-10-20 11:03:59 | 14 秒 | 2 秒 |
| [ci #7](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362790734) | PR | success | 2026-10-20 11:09:23 | 2026-10-20 11:09:33 | 10 秒 | 2 秒 |
| [ci #8](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362976846) | PR | failure | 2026-10-20 11:11:28 | 2026-10-20 11:11:37 | 9 秒 | 2 秒 |
| [ci #9](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39363082268) | PR | success | 2026-10-20 11:12:41 | 2026-10-20 11:12:54 | 13 秒 | 2 秒 |
| [ci #11](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39365760963) | PR | success | 2026-10-20 11:43:15 | 2026-10-20 11:43:34 | 19 秒 | 3 秒 |
| [ci #12](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39366725871) | push → main | success | 2026-10-20 11:54:21 | 2026-10-20 11:54:38 | 17 秒 | 2 秒 |
| [ci #13](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39374359390) | PR | success | 2026-10-20 13:20:24 | 2026-10-20 13:20:38 | 14 秒 | 2 秒 |
| [ci #14](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39375149163) | push → main | success | 2026-10-20 13:29:19 | 2026-10-20 13:29:35 | 16 秒 | 3 秒 |
| [ci #16](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39433405883) | push → main | success | 2026-10-20 21:47:09 | 2026-10-20 21:47:28 | 19 秒 | 3 秒 |

本组十二个样本范围9～42秒，中位数14.5秒，平均约16.67秒；全部样本已列出，不能只保留平均数。不同事件、是否打包、缓存和共享runner负载均可能影响时长，这些样本不构成受控性能比较。原报告初始CI所列46秒是workflow层面的时间，42秒才是对应job时间；引用时保持口径一致。

## 第一轮修复时间的可复核拆分

| 事件 | 北京时间 | 来源 |
|---|---|---|
| PR #2创建 | 2026-10-20 10:35:57 | PR createdAt |
| 首次ci开始 | 2026-10-20 10:36:04 | ci #4 job started_at |
| 首次ci失败结束 | 2026-10-20 10:36:16 | ci #4 job completed_at |
| 修复提交的作者/提交者时间 | 2026-10-20 10:59:05 | cd0f736提交元数据 |
| 修复推送 | 2026-10-20 10:59:10 | PushEvent，push_id=47082368430，head=cd0f73673752585c27cab631d88ce8762a0c7f4d |
| 修复ci创建 | 2026-10-20 10:59:15 | ci #5 created_at |
| 修复ci开始 | 2026-10-20 10:59:17 | ci #5 job started_at |
| 修复ci通过结束 | 2026-10-20 10:59:32 | ci #5 job completed_at |
| 合入main | 2026-10-20 11:03:41 | PR mergedAt |

计算：

- 首次失败结束→修复通过：23分16秒（1396秒）。这是观测到的修复反馈周期，不是严格测得的“如果没有失败能省多少时间”。
- 首次失败结束→修复推送：22分54秒（1374秒）。这是人工处理的代理区间，包含阅读、沟通、空闲、本地验证和推送，不能声称1374秒全部用于编辑代码。
- 修复推送→ci创建：5秒；ci创建→job开始：2秒；job执行：15秒。后面三项共22秒，加1374秒正好1396秒。
- 失败ci和修复ci的job执行时间分别为12秒与15秒，合计27秒。
- 若以一次15秒ci job与PR #2的10810秒提交到发布间隔相比，占比约0.139%。这是单次job占比，不能当成该PR全部CI运行成本的占比。

## PR SHA 与 main SHA

日志来自PR #2修复通过的 [ci #5](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39361867498)、PR #3恢复后的 [ci #9](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39363082268)，以及合入main后的 [ci #6](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39362282086)、[ci #12](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39366725871)。

| 对象 | PR #2 | PR #3 |
|---|---|---|
| PR运行ref | refs/pull/2/merge | refs/pull/3/merge |
| PR运行github.sha | `248fa45ffa40c1667e6e9e03a6651dea5bb3cbeb` | `0b6d625d3be205e3b7ab7332c304f15f281ce7a4` |
| pull_request.head.sha | `cd0f73673752585c27cab631d88ce8762a0c7f4d` | `422dfa06fa3e5a67a6991a43a2dd4fb0edc0b3bd` |
| Squash后main提交 | `392507cab237a19aa6f829ac918a59aaaf2ff3e0` | `3d94e8aa5601d5af5f3b68b2719b0263aab43408` |
| PR运行github.sha是main祖先吗 | 否，exit=1 | 否，exit=1 |
| PR head.sha是main祖先吗 | 否，exit=1 | 否，exit=1 |
| main提交是main祖先吗 | 是，exit=0 | 是，exit=0 |

本次祖先检查时origin/main为 `e05eae7164fc82040063ed96bbc80a5e16387c6d`。四个PR提交对象均已成功取得，exit=1表示不是祖先，而不是对象不存在；另两个main提交exit=0。临时merge提交有两个父提交，分别是当时目标分支提交和PR head；Squash生成的新main提交只有一个父提交。准确表述是“不在main历史中”，不是“GitHub上找不到这些提交”。这些结果针对本组Squash方式，不应泛化到所有合并方式。

原始日志摘录和Git命令结果见 [pr-main-sha-20261021.txt](pr-main-sha-20261021.txt)。

## 写作时还要结合的已有材料

- 问题3：见 [reproducible-builds-20261021.md](reproducible-builds-20261021.md) 和对应原始日志。
- 问题4：最新main中的LAB1.md第4节已包含两个版本的追溯、sha256与/version输出；本文件不重复组员内容。
- 发布门禁拒绝：[release #2](https://github.com/Daniel-647/ops-lab1-7/actions/runs/39375781855)，错误为“Release rejected: tag commit is not in main history.”；后续打包、发布步骤跳过。
- 第二轮证据由组员继续补充。最后统稿时再把本文件的相关数据放入LAB1.md，填完整贡献表并更新开头进度描述。
