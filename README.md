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

```sh
mdesales@ubuntu [10/26/2015 22:18:38] ~/go/src/github.com/marcellodesales/go-nexus-artifacts-downloader (master) $ go run main.go 
2015/10/26 22:18:50 Loading services from Nexus at https://repo1.maven.org/maven2
2015/10/26 22:18:50 Loading the package list https://repo1.maven.org/maven2/org/springframework/cloud
2015/10/26 22:18:52 Received url https://repo1.maven.org/maven2/org/springframework/cloud/spring-cloud-config-server/maven-metadata.xml
2015/10/26 22:18:52 Received url https://repo1.maven.org/maven2/org/springframework/cloud/spring-cloud-netflix-eureka-server/maven-metadata.xml
2015/10/26 22:18:52 Received url https://repo1.maven.org/maven2/org/springframework/cloud/spring-cloud-starter-eureka-server/maven-metadata.xml
2015/10/26 22:18:52 The maven list of artifacts is map[string]*nexus.Artifact{"spring-cloud-config-server":(*nexus.Artifact)(0xc2080965f0), "spring-cloud-netflix-eureka-server":(*nexus.Artifact)(0xc208096e60), "spring-cloud-starter-eureka-server":(*nexus.Artifact)(0xc208097730)} 
2015/10/26 22:18:52 DOWNLOADING NOW... 

2015/10/26 22:18:52 Processing https://repo1.maven.org/maven2/org/springframework/cloud/spring-cloud-starter-eureka-server/1.0.3.RELEASE/spring-cloud-starter-eureka-server-1.0.3.RELEASE.jar to spring-cloud-starter-eureka-server-1.0.3.RELEASE.jar
2015/10/26 22:18:52 Processing https://repo1.maven.org/maven2/org/springframework/cloud/spring-cloud-config-server/1.0.2.RELEASE/spring-cloud-config-server-1.0.2.RELEASE.jar to spring-cloud-config-server-1.0.2.RELEASE.jar
2015/10/26 22:18:52 Processing https://repo1.maven.org/maven2/org/springframework/cloud/spring-cloud-netflix-eureka-server/1.0.3.RELEASE/spring-cloud-netflix-eureka-server-1.0.3.RELEASE.jar to spring-cloud-netflix-eureka-server-1.0.3.RELEASE.jar
spring-cloud-starter-eureka-server-1.0.3.RELEASE.jar2.26 KB / 2.26 KB [===========================================================] 100.00 % 5.59 KB/s 0
spring-cloud-config-server-1.0.2.RELEASE.jar69.76 KB / 69.76 KB [===============================================================] 100.00 % 69.00 KB/s 1s
spring-cloud-netflix-eureka-server-1.0.3.RELEASE.jar399.58 KB / 457.02 KB [===============================================>------] 87.43 % 152.35 KB/s 0
[2311 bytes downloaded and saved as spring-cloud-starter-eureka-server-1.0.3.RELEASE.jar in 254.771316ms 71431 bytes downloaded and saved as spring-cloud-config-server-1.0.2.RELEASE.jar in 998.652254ms 467989 bytes downloaded and saved as spring-cloud-netflix-eureka-server-1.0.3.RELEASE.jar in 2.704306592s]
3.646426643s

```
