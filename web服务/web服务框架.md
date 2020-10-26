# web服务框架

一般情况下，如果只是很少的api，或者功能需求非常简单，我们通常不建议使用web框架。但是如果api众多，功能复杂，那么选用合适的一个web框架，对项目的帮助是非常大的。

golang的生态中有好多个非常流行的web服务框架:

- [beego](https://github.com/astaxie/beego)
- [echo](https://github.com/labstack/echo)
- [iris](https://github.com/kataras/iris)
- [gin](https://github.com/gin-gonic/gin)

这些框架的核心功能都大同小异，会有自己定义的`router`，自己定义的各种web组件或者中间件。

一般我们选定一款web框架后会研究它的实现，方便自己在实现项目运用中排查问题。

