# k8s-client-demo
research k8s client

## Controller

有时候 Controller 也被叫做 Operator。这两个术语的混用有时让人感到迷惑。Controller 是一个通用的术语，凡是遵循 “Watch K8s 资源并根据资源变化进行调谐” 模式的控制程序都可以叫做 Controller。而 Operator 是一种专用的 Controller，用于在 Kubernetes 中管理一些复杂的，有状态的应用程序。例如在 Kubernetes 中管理 MySQL 数据库的 MySQL Operator。

## k8s Code Generator

使用方式：https://github.com/kubernetes/code-generator/blob/master/generate-groups.sh

读取apis的时候，完整路径应该是 `<output-package>/<apis-package>`

## Reference

- [Kubernetes Controller 机制详解 1](https://www.zhaohuabing.com/post/2023-03-09-how-to-create-a-k8s-controller/)
- [Dynamic Client Sample](https://zhuanlan.zhihu.com/p/165970638)
- [Official Client Go Sample](https://github.com/kubernetes/client-go/blob/master/examples/README.md)
- [Official CRD Sample](https://github.com/kubernetes/apiextensions-apiserver/tree/master/examples/client-go)
- [Official CRD Controller](https://github.com/kubernetes/sample-controller)