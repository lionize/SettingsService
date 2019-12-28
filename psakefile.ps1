Task Publish -Depends Pack {
    Exec { docker login docker.io --username=tiksn }
     foreach ($VersionTag in $script:VersionTags) {
         $localTag = ($script:imageName + ":" + $VersionTag)
         $remoteTag = ("docker.io/" + $localTag)
         Exec { docker tag $localTag $remoteTag }
         Exec { docker push $remoteTag }
 
         try {
             Exec { keybase chat send --nonblock --private lionize "BUILD: Published $remoteTag" }
         }
         catch {
             Write-Warning "Failed to send notification"
         }
     }
 }
 
 Task Pack -Depends Build, EstimateVersions {
    $src = (Resolve-Path ".\src\").Path
    $tagsArguments = @()
     foreach ($VersionTag in $script:VersionTags) {
         $tagsArguments += "-t"
         $tagsArguments += ($script:imageName + ":" + $VersionTag)
     }
 
     Exec { docker build -f Dockerfile $src $tagsArguments }
 }
 
 Task EstimateVersions {
    $script:VersionTags = @()
 
    if ($Latest) {
        $script:VersionTags += 'latest'
    }
 
    if (!!($Version)) {
        $Version = [Version]$Version
 
        Assert ($Version.Revision -eq -1) "Version should be formatted as Major.Minor.Patch like 1.2.3"
        Assert ($Version.Build -ne -1) "Version should be formatted as Major.Minor.Patch like 1.2.3"
 
        $Version = $Version.ToString()
        $script:VersionTags += $Version
    }
 
    Assert $script:VersionTags "No version parameter (latest or specific version) is passed."
 }
 
Task Build -Depends BuildWinx64, BuildWinx86, BuildLinux64

Task BuildWinx64 -Depends PreBuild {
    $script:publishWinx64Folder = Join-Path -Path $script:publishFolder -ChildPath "winx64"
    $outputFile = Join-Path -Path $script:publishWinx64Folder -ChildPath "SettingsService.exe"

    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    Exec { go build -o $outputFile } -workingDirectory $script:srcFolder
}

Task BuildWinx86 -Depends PreBuild {
    $script:publishWinx86Folder = Join-Path -Path $script:publishFolder -ChildPath "winx86"
    $outputFile = Join-Path -Path $script:publishWinx86Folder -ChildPath "SettingsService.exe"
    
    $env:GOOS = "windows"
    $env:GOARCH = "386"
    Exec { go build -o $outputFile } -workingDirectory $script:srcFolder
}

Task BuildLinux64 -Depends PreBuild {
    $script:publishLinux64Folder = Join-Path -Path $script:publishFolder -ChildPath "linux64"
    $outputFile = Join-Path -Path $script:publishLinux64Folder -ChildPath "SettingsService"

    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    Exec { go build -o $outputFile } -workingDirectory $script:srcFolder
}

Task PreBuild -Depends Init, Clean, Format, InstallPackages {
    $script:publishFolder = Join-Path -Path $script:trashFolder -ChildPath "bin"
    
    Exec { swag init } -workingDirectory $script:srcFolder
}

Task InstallPackages {
    Exec { go get "go.mongodb.org/mongo-driver/mongo" }
    Exec { go get "go.mongodb.org/mongo-driver/bson" }
    Exec { go get "go.mongodb.org/mongo-driver/mongo/options" }
    Exec { go get -u "github.com/swaggo/swag/cmd/swag" }
    Exec { go get "github.com/iris-contrib/swagger" }
}

Task Format -Depends Clean {
    Exec { go fmt } -workingDirectory $script:srcFolder
}

Task Clean -Depends Init {
    Exec { go clean } -workingDirectory $script:srcFolder
}

Task Init {
    $date = Get-Date
    $ticks = $date.Ticks
    $script:imageName = "tiksn/lionize-settings-service"
    $trashFolder = Join-Path -Path . -ChildPath ".trash"
    $script:trashFolder = Join-Path -Path $trashFolder -ChildPath $ticks.ToString("D19")
    New-Item -Path $script:trashFolder -ItemType Directory | Out-Null
    $script:trashFolder = Resolve-Path -Path $script:trashFolder
    $script:srcFolder = Resolve-Path -Path ".\src\" -Relative

    if(-not $env:GOPATH){
        $env:GOPATH= Join-Path -Path "$HOME" -ChildPath 'go'
    }
    
    # $env:GOROOT="/usr/lib/go"
    # $env:PATH="$env:PATH:$env:GOROOT/bin:$env:GOPATH/bin"
}
 