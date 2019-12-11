Task InstallPackages {
    Exec { go get "github.com/urfave/cli" }
    Exec { go get "github.com/moby/buildkit/frontend/dockerfile/parser" }
    Exec { go get "github.com/docker/distribution/reference" }
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
 