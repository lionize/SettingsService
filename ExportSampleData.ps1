Connect-Mdbc -ConnectionString 'mongodb://root:9fP30ErG0fBv5R@localhost:52540'

$db = Get-MdbcDatabase -Name 'Settings'

$collection = Get-MdbcCollection -Name 'Defaults' -Database $db

Get-MdbcData | Export-MdbcData -Path SampleDefaults.json