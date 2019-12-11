Properties {
    $version = "0.0.1"
}

Task Build -Depends BuildWinx64, BuildWinx86, BuildLinux64

Task BuildWinx64 -Depends PreBuild {
    $script:publishWinx64Folder = Join-Path -Path $script:publishFolder -ChildPath "winx64"
    $outputFile = Join-Path -Path $script:publishWinx64Folder -ChildPath "SettingsService.exe"

    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    Exec { go build -o $outputFile $script:srcFolder }
}

Task BuildWinx86 -Depends PreBuild {
    $script:publishWinx86Folder = Join-Path -Path $script:publishFolder -ChildPath "winx86"
    $outputFile = Join-Path -Path $script:publishWinx86Folder -ChildPath "SettingsService.exe"
    
    $env:GOOS = "windows"
    $env:GOARCH = "386"
    Exec { go build -o $outputFile $script:srcFolder }
}

Task BuildLinux64 -Depends PreBuild {
    $script:publishLinux64Folder = Join-Path -Path $script:publishFolder -ChildPath "linux64"
    $outputFile = Join-Path -Path $script:publishLinux64Folder -ChildPath "SettingsService"

    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    Exec { go build -o $outputFile $script:srcFolder }
}

Task PreBuild -Depends Init, Clean, Format, InstallPackages {
}

Task InstallPackages {
    Exec { go get "go.mongodb.org/mongo-driver" }
}

Task Format -Depends Clean {
    Exec { go fmt $script:srcFolder }
}

Task Clean -Depends Init {
    Exec { go clean $script:srcFolder }
}

Task Init {
    $date = Get-Date
    $ticks = $date.Ticks
    $trashFolder = Join-Path -Path . -ChildPath ".trash"
    $script:trashFolder = Join-Path -Path $trashFolder -ChildPath $ticks.ToString("D19")
    New-Item -Path $script:trashFolder -ItemType Directory | Out-Null
    $script:trashFolder = Resolve-Path -Path $script:trashFolder
    $script:srcFolder = Resolve-Path -Path ".\src\" -Relative
}
 