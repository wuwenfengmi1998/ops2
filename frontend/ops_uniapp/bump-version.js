/**
 * 打包前自动递增版本号
 * 用法: node bump-version.js
 *
 * - versionCode: 每次 +1（必须递增，否则 Android 无法覆盖安装）
 * - versionName: 语义化版本自动递增 patch 位（1.0.0 → 1.0.1）
 */

const fs = require('fs')
const path = require('path')

const manifestPath = path.join(__dirname, 'manifest.json')

if (!fs.existsSync(manifestPath)) {
  console.error('未找到 manifest.json，请在 ops2_uniapp 目录下运行此脚本')
  process.exit(1)
}

// manifest.json 可能包含 JS 风格注释，先去掉再解析
const raw = fs.readFileSync(manifestPath, 'utf8')
const cleaned = raw.replace(/\/\*[\s\S]*?\*\/|\/\/.*$/gm, '')
const manifest = JSON.parse(cleaned)

// 递增 versionCode
const oldCode = parseInt(manifest.versionCode) || 0
const newCode = oldCode + 1
manifest.versionCode = newCode.toString()

// 递增 versionName 的 patch 位
const parts = (manifest.versionName || '1.0.0').split('.')
if (parts.length >= 3) {
  parts[2] = (parseInt(parts[2]) + 1).toString()
} else {
  parts.push('1')
}
manifest.versionName = parts.join('.')

// 写回
fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2) + '\n')
console.log(`版本已更新: ${manifest.versionName} (${manifest.versionCode})`)
