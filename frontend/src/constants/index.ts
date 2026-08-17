// 与后端 constants 对应的枚举/常量（枚举出现位置清单见 README）
export const ROLE = {
  USER: 1,
  ADMIN: 2,
} as const

export const USER_STATUS = {
  DISABLED: 0,
  ACTIVE: 1,
} as const

export const ARTICLE_STATUS = {
  DRAFT: 0,
  PUBLISHED: 1,
  OFFLINE: 2,
} as const

export const VISIBILITY = {
  PRIVATE: 0,
  PUBLIC: 1,
} as const

export const NOTIFICATION_TYPE = {
  FOLLOW: 'follow',
  LIKE: 'like',
  COMMENT: 'comment',
  SYSTEM: 'system',
} as const

export const SORT_TYPE = {
  LATEST: 'latest',
  HOTTEST: 'hottest',
} as const

export const ARTICLE_STATUS_TEXT: Record<number, string> = {
  [ARTICLE_STATUS.DRAFT]: '草稿',
  [ARTICLE_STATUS.PUBLISHED]: '已发布',
  [ARTICLE_STATUS.OFFLINE]: '已下架',
}

export const USER_STATUS_TEXT: Record<number, string> = {
  [USER_STATUS.DISABLED]: '禁用',
  [USER_STATUS.ACTIVE]: '启用',
}

export const ROLE_TEXT: Record<number, string> = {
  [ROLE.USER]: '普通用户',
  [ROLE.ADMIN]: '管理员',
}

export const VISIBILITY_TEXT: Record<number, string> = {
  [VISIBILITY.PRIVATE]: '私密',
  [VISIBILITY.PUBLIC]: '公开',
}

export const NOTIFICATION_TYPE_TEXT: Record<string, string> = {
  [NOTIFICATION_TYPE.FOLLOW]: '关注了你',
  [NOTIFICATION_TYPE.LIKE]: '点赞了你的文章',
  [NOTIFICATION_TYPE.COMMENT]: '评论了你的文章',
  [NOTIFICATION_TYPE.SYSTEM]: '系统通知',
}

export const ERROR_MESSAGES: Record<number, string> = {
  0: 'ok',
  40000: '请求参数错误',
  40100: '未认证或登录已过期',
  40300: '无权限执行该操作',
  40400: '资源不存在',
  40900: '资源冲突',
  42200: '参数校验失败',
  42900: '请求过于频繁，请稍后再试',
  50000: '服务器内部错误',
}
