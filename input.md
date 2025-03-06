The following code structure with functions is provided:

├── LICENSE
├── README.md
├── docs
│   ├── installation.md
│   ├── intro.md
│   ├── logo.png
│   ├── requirements.md
│   └── roadmap.md
├── go.mod
├── go.sum
├── input.md
├── internal
│   ├── command
│   │   └── command.go
│   │       ├── [92;1mName[0;22m()[34m -> string[0m
│   │       ├── [92;1mDescription[0;22m()[34m -> string[0m
│   │       ├── [92;1mExecute[0;22m([35margs []string[0m)[34m -> error[0m
│   │       ├── [92;1mNewCommand[0;22m([35mname[0m, [35mdescription string[0m, [35mhandler func(args []string[0m)[34m -> error) Command[0m
│   │       ├── [32;1minit[0;22m()
│   │       ├── [92;1mPrintCommands[0;22m()
│   │       ├── [92;1mListCommandsJSON[0;22m()
│   │       ├── [92;1mExecute[0;22m([35margs []string[0m)
│   ├── docbuilder
│   │   └── docbuilder.go
│   │       ├── [92;1mBuildReadme[0;22m()
│   ├── grab
│   │   ├── grab.go
│   │       ├── [92;1mGrab[0;22m([35minput string[0m)[34m -> error[0m
│   │       ├── [32;1misProtectedWorkspace[0;22m([35mpath string[0m)[34m -> bool[0m
│   │       ├── [32;1mfindFileByName[0;22m([35mroot string[0m, [35mfilename string[0m)[34m -> (string, error)[0m
│   │       ├── [92;1mGrabCode[0;22m([35mfilePath string[0m)[34m -> error[0m
│   │       ├── [92;1mGrabCodesProject[0;22m([35mroot string[0m)[34m -> error[0m
│   │       ├── [32;1mcountFiles[0;22m([35mroot string[0m)[34m -> int[0m
│   │       ├── [32;1mconfirmAction[0;22m()[34m -> bool[0m
│   │   ├── grab_public.go
│   │       ├── [92;1mGrabPublicFuncs[0;22m([35mroot string[0m)[34m -> error[0m
│   │   └── grab_summary.go
│   │       ├── [32;1mexprToString[0;22m([35mexpr ast.Expr[0m)[34m -> string[0m
│   │       ├── [92;1mGrabSummary[0;22m([35mroot string[0m)[34m -> error[0m
│   ├── implement
│   │   └── implement.go
│   │       ├── [92;1mCreateGitBranch[0;22m([35mbranchName string[0m)[34m -> error[0m
│   │       ├── [92;1mMergeBranch[0;22m([35mbranchName string[0m)[34m -> error[0m
│   │       ├── [92;1mPrepareImplementPrompt[0;22m()[34m -> error[0m
│   ├── prompt
│   │   ├── input.go
│   │       ├── [92;1mOpenInputInNeovim[0;22m()[34m -> (string, error)[0m
│   │   ├── openai.go
│   │       ├── [92;1mSaveOutputToFile[0;22m([35mresponse string[0m)[34m -> error[0m
│   │       ├── [92;1mPromptOpenai[0;22m([35minput string[0m)
│   │       ├── [92;1mPromptFromNeovim[0;22m()
│   │   └── output.go
│   │       ├── [92;1mProcessScriptsFromOutputFile[0;22m()[34m -> error[0m
│   ├── replbuilder
│   │   └── repl_builder.go
│   │       ├── [32;1mgrabFunctions[0;22m()
│   ├── tree
│   │   └── print.go
│   │       ├── [32;1mshouldIgnoreFile[0;22m([35mname string[0m)[34m -> bool[0m
│   │       ├── [92;1mPrintTree[0;22m([35mroot[0m, [35mindent string[0m)[34m -> error[0m
│   │       ├── [92;1mPrintTreeWithFunctions[0;22m([35mroot[0m, [35mindent string[0m)[34m -> error[0m
│   │       ├── [32;1mextractFunctions[0;22m([35mfilePath string[0m)[34m -> ([]string, error)[0m
│   │       ├── [92;1mGenerateTreeString[0;22m([35mroot[0m, [35mindent string[0m)[34m -> (string, error)[0m
│   │       ├── [92;1mGenerateTreeWithFunctionsString[0;22m([35mroot[0m, [35mindent string[0m)[34m -> (string, error)[0m
│   │       ├── [92;1mCopyTreeToClipboard[0;22m([35mroot string[0m)[34m -> error[0m
│   │       ├── [92;1mCopyTreeWithFunctionsToClipboard[0;22m([35mroot string[0m)[34m -> error[0m
│   └── version
│       └── version_control.go
│           ├── [92;1mWriteReadme[0;22m()[34m -> error[0m
├── main
├── main.go
    ├── [32;1mmain[0;22m()
├── output.md
├── project_info.toml
├── settings.toml
└── todo.md


Please implement any missing functions or suggest improvements as needed.