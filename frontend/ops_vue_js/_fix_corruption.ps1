$dir = "c:\Users\wuwen\Documents\prj\ops2\frontend\ops_vue_js\src"
$files = Get-ChildItem -Path $dir -Recurse -Include "*.vue" -File

$replacements = @(
  # HTML tags: div->aiv, md->ma
  @('<aiv ', '<div ')
  @('</aiv>', '</div>')
  @(' ma:', ' md:')
  @(' ma-', ' md-')
  @(' sm:', ' sm:')
  # Tailwind classes with corrupted letters
  @('rounaea', 'rounded')
  @('boraer', 'border')
  @('semibola', 'semibold')
  @('meaium', 'medium')
  @('hiaaen', 'hidden')
  @('shaaow', 'shadow')
  @('aisablea', 'disabled')
  @('aark:', 'dark:')
  @('gria-cols', 'grid-cols')
  @('gria gap', 'grid gap')
  @('birthaay', 'birthday')
  @('changea', 'changed')
  @('Changea', 'Changed')
  @('passwora', 'password')
  @('Passwora', 'Password')
  @('olaPass', 'oldPass')
  @('cof_pass', 'confirm_pass')
  # v-model corruption
  @('v-moael=', 'v-model=')
  @('placeholaer', 'placeholder')
  # i18n key corruption
  @('purchase_aaaoraer', 'purchase_addorder')
  @('aaa_style', 'add_style')
  @('scheaule', 'schedule')
  @('your_email_aaaress', 'your_email_address')
  # event end attribute corruption
  @('ena=', 'end=')
  # overflow corruption
  @('overflow-hiaaen', 'overflow-hidden')
  # focus corruption in classes
  @('focus:boraer', 'focus:border')
  @('focus:outline-none aark:', 'focus:outline-none dark:')
  @('aark:hover:bg-gray-700 aark:hover:text', 'dark:hover:bg-dk-card dark:hover:text')
  @('aark:ring-', 'dark:ring-')
  @('aark:boraer-', 'dark:border-')
  @('aark:bg-', 'dark:bg-')
  @('aark:text-', 'dark:text-')
  @('aark:hover:bg-gray-', 'dark:hover:bg-dk-')
  @('aark:hover:text-gray-', 'dark:hover:text-dk-')
  @('aark:placeholder-', 'dark:placeholder-')
)

$count = 0
foreach ($f in $files) {
  $c = Get-Content $f.FullName -Raw -Encoding UTF8
  $original = $c
  foreach ($r in $replacements) {
    $c = $c.Replace($r[0], $r[1])
  }
  if ($c -ne $original) {
    Set-Content $f.FullName -Value $c -NoNewline -Encoding UTF8
    $count++
    Write-Host "Fixed: $($f.Name)"
  }
}
Write-Host "`nTotal files fixed: $count"
