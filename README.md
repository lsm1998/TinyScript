# TinyScript
Self made script language

使用Golang语言实现的简单脚本语言，支持基本的数据类型，支持函数，支持闭包，支持类，提供了内置的系统函数

## 依赖环境

- Golang 1.16

## 起源

一次偶然的机会看到了一篇文章，讲述了如何使用Go语言实现一个简单的脚本语言，于是我试图将其扩展，就有了这个项目。

源文地址：https://www.poetries.cn/crafting/welcome/the-lox-language.html

## 扩展点

- 增加了let关键字，用于声明变量
- 增加了一些注释，方便学习
- 增加了系统内置函数，通过import关键字引入
- 支持function关键字，和fn关键字作用一致（和JavaScript一致）
- 支持数组
