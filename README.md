# go-nexus-artifacts-downloader

A Downloader of Nexus Artifacts while learning Go, with same Docker Download status bar CLI.

The code is fully commented if you desire to learn Go by examples.

# Sample

The main.go includes the following constants:

```go
NEXUS_SERVER*   = "https://repo1.maven.org"
NEXUS_REPO*     = "maven2"
NEXUS_GROUP_ID* = "org.springframework.cloud"
```

In addition, it includes a closure for filtering what files to download. For instance, from https://repo1.maven.org/maven2/org/springframework/cloud/ it will download only the artifacts that has a "server" in the name.

```go
serverFilter := func(artifactId string) bool {
  // runnable servers from spring https://repo1.maven.org/maven2/org/springframework/cloud/
  return strings.Contains(artifactId, "server")
}

```

Running the command...

[![asciicast](https://asciinema.org/a/d2jk0du8r003g1q70bu0a2ki7.png)](https://asciinema.org/a/d2jk0du8r003g1q70bu0a2ki7)

As a result, the files listed are downloaded:

```sh
mdesales@ubuntu [10/26/2015 22:36:34] ~/go/src/github.com/marcellodesales/go-nexus-downloader (master) $ git status
On branch master
Untracked files:
  (use "git add <file>..." to include in what will be committed)

	spring-cloud-config-server-1.0.2.RELEASE.jar
	spring-cloud-netflix-eureka-server-1.0.3.RELEASE.jar
	spring-cloud-starter-eureka-server-1.0.3.RELEASE.jar
```
