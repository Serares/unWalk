### unWalk is a cli utility that is able to unarchive a bunch of gzip files

Example running the tool:
`./unWalk -root /tmp/archive_tmp -dest /tmp/unarchived_files`

The above command will recursevly look for files inside `/tmp/archive_tmp` and unarchive all the files keeping the directory structure existent in the root

The above command was run on the example directory:

```
/tmp/archive_tmp
├── actions.go.gz
├── actions_test.go.gz
├── helper.go.gz
├── main.go.gz
├── main_test.go.gz
└── testdata
    └── dir.log.gz
```

And extracted the following files:
```
/tmp/unarchived_files
├── actions.go
├── actions_test.go
├── helper.go
├── main.go
├── main_test.go
└── testdata
    └── dir.log
```