# sourcegraph.tar.gz 

Export SourceGraph search query results to a tar.gz file.

## Installation

```bash
git clone github.com/bthuilot/sourcegraph.tar.gz
cd sourcegraph.tar.gz
make
```

## Usage

```bash
export SOURCEGRAPH_TOKEN="your-sourcegraph-token"

sg-tar --query "content:New[A-Z][a-z]+Client\(" --output sourcegraph.tar.gz --compress
```