Import-Module -Name Mdbc

Connect-Mdbc -ConnectionString 'mongodb://root:Xq5xrtzEKj44ueyd@localhost:27017'

$db = Get-MdbcDatabase -Name 'SettingsService'

$collection = Get-MdbcCollection -Name 'DefaultSettings' -Database $db

Import-MdbcData -Path SampleDefaultSettings.json | Set-MdbcData -Collection $collection -Add

$collection = Get-MdbcCollection -Name 'UserSettings' -Database $db

Import-MdbcData -Path SampleUserSettings.json | Set-MdbcData -Collection $collection -Add
