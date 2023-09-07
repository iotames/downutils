## 依赖项

- Go(v1.12+)语言工具(略)
- C编译器(用于连接系统图形驱动)
- 系统图形驱动程序

`Linux-Debian`

gcc, 图形库头文件(graphics library header files)

```shell
    sudo apt-get install gcc libgl1-mesa-dev xorg-dev
```

`Windows`

下载并安装TDM-GCC(gcc编译器)

> https://github.com/jmeubank/tdm-gcc/releases/download/v10.3.0-tdm64-2/tdm64-gcc-10.3.0-2.exe

其他操作系统请参看官方文档


## 生成应用图标

打包后的应用任务栏和窗口左上角图标都是空的。必须打包图片静态资源到应用程序中。

测试发现 `powershell` 命令生成的 `bundled.go` 文件为 `UTF-16` 编码，无法使用。
打开 git-cli 使用 `Bash` 命令模式生成 `bundled.go` 则正常。

    fyne bundle ./resource/images/logo.png >> bundled.go
    app := app.New()
    // app.SetIcon(theme.FyneLogo())
    app.SetIcon(StaticResource)

> https://blog.csdn.net/raoxiaoya/article/details/121626549


## 构建应用

1. 下载安装fyne命令工具

```shell
    go get -d fyne.io/fyne/v2/cmd/fyne
```
or
```
    go install fyne.io/fyne/v2/cmd/fyne@latest
```
2. 使用fyne打包命令

```
    fyne package -os windows -icon resource/images/logo.png

```
默认读取 `FyneApp.toml` 配置文件打包。命令行参数覆盖配置文件参数。
- `-os` 参数值: `darwin` `linux` `windows`
- `-icon` 应用程序图标路径


## 中文字体支持

> https://github.com/lusingander/fyne-font-example

由于思源黑体 /思源宋体都是 OTF 格式，把字体做成 TTF 格式了，有人已经做好了转换，并把字库放在了 GitHub 上面：

https://github.com/be5invis/source-han-sans-ttf/releases

https://github.com/junmer/source-han-serif-ttf

https://github.com/Pal3love/Source-Han-TrueType