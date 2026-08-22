$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $PSScriptRoot
$files = Get-ChildItem -Path $root -Recurse -Filter *.go | Where-Object {
  $_.FullName -notmatch '\\tests\\' -and
  $_.FullName -notmatch '\\vendor\\' -and
  $_.FullName -notmatch '\\migrations\\' -and
  $_.Name -notlike '*_test.go'
} | Sort-Object FullName

$lineCount = 0
foreach ($file in $files) {
  $count = (Get-Content $file.FullName | Measure-Object -Line).Lines
  $lineCount += $count
  '{0,6}  {1}' -f $count, ($file.FullName.Substring($root.Length + 1))
}

""
"FILES=$($files.Count)"
"LINES=$lineCount"

