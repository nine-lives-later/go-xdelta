param()

$ErrorPreference = "Stop"

$xdeltaVersion = (Select-String -Path .\src\xdelta3\xdelta3-main.h -Pattern 'Xdelta version ([0-9.]+),').Matches.Groups[1].Value

Write-Host "Xdelta version: $xdeltaVersion"

$utf8 = New-Object System.Text.UTF8Encoding $False

[System.IO.File]::WriteAllLines("version.go", "package lib`n`nconst Version = """ + $xdeltaVersion + """", $utf8)
