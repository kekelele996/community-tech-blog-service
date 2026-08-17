# Bug 复现说明

## Bug 是什么

重复关注、重复点赞没有返回“已关注/已点赞”，取消已存在的关注也会报错，并且关注和点赞后站内通知会被写入两份。

## 如何触发

```bash
cd backend
make follow
```

## 错误信息

```
--- FAIL: TestFollowLikeErrorAndNotificationComposite
    duplicate follow code = 50000, want 40905
    duplicate like code = 50000, want 40905
    alice unread notifications = 2, want 1
    bob unread notifications = 2, want 1
```
