# Bug 复现说明

## Bug 是什么

对不存在的文章下架、删除不存在的关注、给不存在的通知标记已读、从收藏夹移除不存在的文章，这些本应报“不存在”的操作全部静默返回成功。

## 如何触发

```bash
cd backend
make missing
```

## 错误信息

```
--- FAIL: TestRepositoryMissingRecordErrorsComposite
    update missing article status err = <nil>, want ErrNotFound
    delete missing follow err = <nil>, want ErrNotFound
    mark missing notification read err = <nil>, want ErrNotFound
    remove missing collection article err = <nil>, want ErrNotFound
```
