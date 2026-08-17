# Bug 复现说明

## Bug 是什么

通知的触发者不存在时，组装通知列表会解引用空用户指针，接口直接 panic。

## 如何触发

```bash
cd backend
make notice
```

## 错误信息

```
panic: runtime error: invalid memory address or nil pointer dereference
goroutine 11 [running]:
github.com/techblog/community/internal/service.(*UserService).BuildProfile(...)
github.com/techblog/community/internal/service.(*NotificationService).buildDTO(...)
github.com/techblog/community/internal/service.(*NotificationService).List(...)
```
