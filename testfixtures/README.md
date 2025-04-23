# Test Fixtures for GitIgnoreTree

This directory contains a hierarchical structure with multiple `.gitignore` files at different levels to test the `GitIgnoreTree` functionality.

## Structure

```
testfixtures/
├── .gitignore                 # Root level gitignore
├── important.txt              # Not ignored (negated by !important.txt rule)
├── ignored.txt                # Ignored by *.txt rule
├── root_file.go               # Not ignored
├── root_ignored_dir/          # Ignored directory
│   └── file_in_ignored_dir.go # Ignored (in ignored directory)
└── level1/
    ├── .gitignore             # Level 1 gitignore
    ├── normal_file.go         # Not ignored
    ├── ignored_in_level1.go   # Ignored by explicit rule
    ├── log_file.log           # Ignored by *.log rule
    ├── level1_ignored_dir/    # Ignored directory
    │   └── file_inside.go     # Ignored (in ignored directory)
    └── level2/
        ├── .gitignore         # Level 2 gitignore
        ├── normal_level2_file.go # Not ignored
        ├── ignored_at_level2.md  # Ignored by explicit rule
        ├── temp.tmp              # Ignored by *.tmp rule
        ├── level2_ignored/       # Ignored directory
        │   └── inside_ignored_dir.go # Ignored (in ignored directory)
        └── level3/
            ├── .gitignore     # Level 3 gitignore
            ├── normal_file.go # Not ignored
            ├── level3_ignored_file.go # Ignored by explicit rule
            ├── file.binary    # Ignored by *.binary rule
            └── something.generated.js # Ignored by **/*.generated.js rule
```

## Test Cases

The test files in `ignore_tree_test.go` verify:

1. That the GitIgnoreTree loads all `.gitignore` files correctly
2. That files are properly ignored/not ignored based on rules
3. That the RegularFilesSeq method only returns non-ignored files 