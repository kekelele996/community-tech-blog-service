# Bug 复现说明

## Bug 是什么

上下文取消后，按 ID 查文章、查全部关注、分页查通知、按 ID 查用户的查询仍然会继续执行，没有把取消信号传给数据库查询。

## 如何触发

```bash
cd backend
make cancel
```

## 错误信息

```
--- FAIL: TestContextCancellationPropagation
    article FindByID with canceled ctx err = find article by id: record not found, want context.Canceled
    AllFollowingIDs with canceled ctx err = <nil>, want context.Canceled
    notification ListByUser with canceled ctx err = <nil>, want context.Canceled
    user FindByID with canceled ctx err = find user by id: record not found, want context.Canceled
```
