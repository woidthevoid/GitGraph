# GitGraph
## What is it?
GitGraph is a CLI(Command Line Interface) tool for showing git contributions locally, in a graph similar to GitHub.
It is meant as a fun little learning project for Go(lang).
## Technologies
GitGraph solely uses Go as it packages the code contained here, into a single binary, so no installing of Go is necessary to run it. Just fetch the executeable and run it!
## How to install: 
GitGraph can be built by yourself or you can download the binary from the latest release.  
### Building yourself:
- Be sure to have at least Go v1.25 installed.
- Clone the repository.
- From root, run `go mod download` and `go mod verify`to verify that the packages match their hashes.
- From root, run `go build -o your-preferred-output-name` and build into a binary. 
- Run the binary from your console with `./binary-name`.
- Youre good to go!
### Downloading the provided binary:
- Go to Releases and find the binary version that matches your OS and architecture.
- In your terminal, run `curl -L -O https://github.com/OWNER/REPO/releases/download/v1.0.0/binary-name`or just `wget https://github.com/OWNER/REPO/releases/download/v1.0.0/binary-name`.
- Unpack the tar.gz with `tar -xzf name-of-file.tar.gz`or if its a windows file, just unzip it.
- The binary should now be showing in the folder you downloaded the zipped repo to. Run it with `./name-of-binary`.
- Youre good to go!
## How to use
GitGraph has two commands, scan and stats so its pretty straighforward. 
Before showing your contribution graph, GitGraph needs to scan some of your folders and see if they contain .git.  
Run `gitgraph scan -f path/to/your/folder` and it will scan that repo, along with all the sub repos for .git, therefore it is not recommended to just run it to scan your whole home repo. 
A .gitlocalstats file will be made in your home folder, that contains the paths to where .git was found.  
Next, run `gitgraph stats -e example@email.com` where you pass the email your contributions where made through. 
This should give you your contribution graph from the last six months, based on the repos that you scanned in the earlier step. 
That is all!
## Examples
