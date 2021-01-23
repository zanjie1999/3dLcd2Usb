# 使用3D打印机面板的LCD2USB（USB2LCD）

这是一款价格极其低廉的大屏usb lcd监控面板，制作所使用的成本不会超过15元人民币  
如果希望更低的成本和焊接，可以看一下 [lcd1602的版本](https://github.com/zanjie1999/usb2lcd) 只需要不到10元人民币

## 如何使用
-----
由于使用的是国产的单片机，坑很多，因此请严格按照本md来操作  
首先在你喜欢的编辑器比如VsCode安装好PlatfromIO的插件  
先在插件的Boards页面搜索lgt8f328p把支持库安装好
再到 [官网](http://www.lgtic.com/downloads/) 底部找到 LGT8FX8D/P系列Arduino硬件支持包 进行下载  
并且将压缩包中的 hardware\LGT\avr 文件夹里面的内容替换到你的系统用户文件夹下的.platformio/packages/framework-lgt8fx  
然后就可以直接编译上传了

## 上位机
-----
控制方式非常简单  
\r会回到左上角第一个字符  
\n会换行  
按下方按键串口会回传1  
编码器按钮会回传0
顺时针旋转回传+
逆时针旋转回传-
  
在monitor目录下有一个使用go写的上位机
