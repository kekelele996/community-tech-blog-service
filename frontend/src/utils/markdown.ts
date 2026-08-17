// Markdown 渲染：marked + highlight.js 代码高亮（marked v12 通过后处理高亮）
import { marked } from 'marked'
import hljs from 'highlight.js'

marked.setOptions({ gfm: true, breaks: true })

function decodeEntities(s: string): string {
  return s
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&amp;/g, '&')
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
}

function highlightBlock(code: string, lang: string): string {
  const text = decodeEntities(code)
  const language = lang.toLowerCase()
  const highlighted =
    language && hljs.getLanguage(language)
      ? hljs.highlight(text, { language }).value
      : hljs.highlightAuto(text).value
  return `<pre><code class="hljs language-${language || 'plaintext'}">${highlighted}</code></pre>`
}

export function renderMarkdown(content: string): string {
  const html = marked.parse(content || '') as string
  return html
    .replace(
      /<pre><code class="language-([^"]+)">([\s\S]*?)<\/code><\/pre>/g,
      (_m, lang: string, code: string) => highlightBlock(code, lang),
    )
    .replace(/<pre><code>([\s\S]*?)<\/code><\/pre>/g, (_m, code: string) => highlightBlock(code, ''))
}
