# Super Lulu

Super Lulu 是一个使用 Go 和 Ebitengine 开发的糖果风横版平台小游戏。

你将操作 Lulu 收集糖果、获取强化道具、投掷冰淇淋、通过关卡，并在终点前击败 Niuniu。

## 功能

- 横版平台跳跃玩法
- 多个糖果主题关卡
- 包含收集物、橙子、冰淇淋强化、敌人和 Boss 战
- 图片素材已内嵌到程序中，构建后的 Windows `.exe` 可以单独运行
- 已配置 GitHub Actions，可通过版本 tag 自动构建并发布 Release

## 操作

| 按键 | 操作 |
| --- | --- |
| `A` / `Left` | 向左移动 |
| `D` / `Right` | 向右移动 |
| `Space` / `W` / `Up` | 跳跃 |
| `J` | 拥有冰淇淋能力时投掷冰淇淋 |
| `Enter` | 开始 / 确认 / 继续 |
| `Up` / `Down` | 选择关卡 |
| `P` / `Esc` | 暂停 / 继续 |
| `R` | 重新开始当前关卡 |
| `M` | 返回菜单 |
| `F3` | 显示 / 隐藏调地图坐标层 |
| `1` - `6` | 调试用关卡快捷键 |

## 环境要求

- Go `1.26.2`
- 主要目标平台为 Windows

## 本地运行

```powershell
go run .
```

## 地图配置

关卡地图现在放在 `levels/levels.json` 中，并通过 `go:embed` 内嵌进程序。

常用字段：

- `Width`：关卡总宽度
- `Spawn`：玩家出生点
- `Platforms`：地面和平台矩形
- `Collect`：普通收集物
- `PowerUps`：冰淇淋能力
- `Oranges`：变大道具
- `Enemies`：敌人初始位置、巡逻范围和速度
- `Boss`：Boss 配置，没有 Boss 时可设为 `null`
- `Goal`：终点旗帜位置
- `Theme`：天空、地面和装饰色

坐标以左上角为原点，`X` 向右增加，`Y` 向下增加。手动编辑 JSON 时保留缩进；代码里导出关卡请使用 `marshalLevelsJSON`，它会统一用 `json.MarshalIndent` 生成可读格式。

调地图时可以在游戏里按 `F3` 打开坐标层。坐标层会显示玩家世界坐标、玩家碰撞框、相机位置，以及鼠标指向位置换算后的世界坐标；把鼠标放到想摆平台或道具的位置，就能直接参考 `Mouse world` 写入 JSON。

## 构建

构建 Windows GUI 可执行文件：

```powershell
go build -ldflags="-H windowsgui" -o dist\super-lulu.exe .
```

`assets/` 中的图片会通过 `go:embed` 在编译时内嵌进程序，因此运行时不需要额外携带素材文件夹。

## 发布

项目已包含 GitHub Actions workflow：`.github/workflows/release.yml`。

创建并推送版本 tag 后，会自动构建并发布 Release：

```powershell
git tag v1.0.0
git push origin v1.0.0
```

workflow 会上传一个单独的可执行文件，文件名类似：

```text
super-lulu-v1.0.0-windows-amd64.exe
```

也可以在 GitHub Actions 页面手动运行该 workflow，并填写要发布的版本 tag。
