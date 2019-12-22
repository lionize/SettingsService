Import-Module -Name Mdbc

Connect-Mdbc -ConnectionString 'mongodb://root:9fP30ErG0fBv5R@localhost:52540'

$db = Get-MdbcDatabase -Name 'Settings'

$collection = Get-MdbcCollection -Name 'DefaultSettings' -Database $db

Import-MdbcData -Path SampleDefaultSettings.json | Set-MdbcData -Collection $collection -Add

$collection = Get-MdbcCollection -Name 'UserSettings' -Database $db

Import-MdbcData -Path SampleUserSettings.json | Set-MdbcData -Collection $collection -Add
