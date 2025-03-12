I want to build a feature called api-key-env-setup.
Feature Description:
I need  a function to write .env file in openai.go

Here is the summary of the code:

File: cmd/commandbuilder.go (package cmd)
  [Package: cmd] Function: init

File: cmd/docbuilder.go (package cmd)
  [Package: cmd] Function: init

File: cmd/grab.go (package cmd)
  [Package: cmd] Function: init

File: cmd/grabpublic.go (package cmd)
  [Package: cmd] Function: init

File: cmd/implement.go (package cmd)
  [Package: cmd] Function: init

File: cmd/prompt.go (package cmd)
  [Package: cmd] Function: init

File: cmd/registered_actions.go (package cmd)
  [Package: cmd] Function: init

File: cmd/root.go (package cmd)
  [Package: cmd] Function: Execute

File: cmd/smartgrab.go (package cmd)
  [Package: cmd] Function: init

File: cmd/summary.go (package cmd)
  [Package: cmd] Function: init

File: cmd/tree.go (package cmd)
  [Package: cmd] Function: init

File: cmd/treefunc.go (package cmd)
  [Package: cmd] Function: init

File: internal/commandbuilder/commandbuilder.go (package commandbuilder)
  [Package: commandbuilder] Function: RegisterActions
  [Package: commandbuilder] Function: escapeString

File: internal/docbuilder/docbuilder.go (package docbuilder)
  [Package: docbuilder] Function: BuildReadme

File: internal/grab/grab.go (package grab)
  [Package: grab] Function: Grab
  [Package: grab] Function: isProtectedWorkspace
  [Package: grab] Function: findFileByName
  [Package: grab] Function: GrabCode
  [Package: grab] Function: GrabCodesProject
  [Package: grab] Function: getCodesProjectContent
  [Package: grab] Function: countFiles
  [Package: grab] Function: confirmAction
  [Package: grab] Function: GrabFiles
  [Package: grab] Function: GrabMultipleFolders

File: internal/grab/grab_public.go (package grab)
  [Package: grab] Function: GrabPublicFuncsWithDescriptions
  [Package: grab] Function: extractPublicFuncsWithDescriptions
  [Package: grab] Function: PrintPublicFunctions

File: internal/grab/grab_summary.go (package grab)
  [Package: grab] Function: exprToString
  [Package: grab] Function: buildSummary
  [Package: grab] Function: GrabSummary

File: internal/grab/smart_grab.go (package grab)
  [Package: grab] Function: SmartGrab
  [Package: grab] Struct: Output
  [Package: grab] Function: getFeatureBranch
  [Package: grab] Function: buildPrompt

File: internal/implement/implement.go (package implement)
  [Package: implement] Function: CreateGitBranch
  [Package: implement] Function: MergeBranch
  [Package: implement] Function: PrepareImplementPrompt
  [Package: implement] Function: Implement
  [Package: implement] Interface: ImplementationManager

File: internal/prompt/input.go (package prompt)
  [Package: prompt] Function: OpenInputInNeovim

File: internal/prompt/openai.go (package prompt)
  [Package: prompt] Function: init
  [Package: prompt] Function: SaveOutputToFile
  [Package: prompt] Struct: CodeResponse
  [Package: prompt] Function: GenerateSchema
  [Package: prompt] Function: PromptOpenai
  [Package: prompt] Function: PromptFromNeovim
  [Package: prompt] Struct: FileResponse
  [Package: prompt] Function: PromptOpenaiFiles

File: internal/prompt/output.go (package prompt)
  [Package: prompt] Function: ProcessScriptsFromOutputFile

File: internal/replbuilder/repl_builder.go (package repl)
  [Package: repl] Function: grabFunctions

File: internal/tree/print.go (package tree)
  [Package: tree] Function: shouldIgnoreFile
  [Package: tree] Function: PrintTree
  [Package: tree] Function: PrintTreeWithFunctions
  [Package: tree] Function: extractFunctions
  [Package: tree] Function: GenerateTreeString
  [Package: tree] Function: GenerateTreeWithFunctionsString
  [Package: tree] Function: CopyTreeToClipboard
  [Package: tree] Function: CopyTreeWithFunctionsToClipboard

File: internal/version/version_control.go (package version)
  [Package: version] Function: WriteReadme

File: main.go (package main)
  [Package: main] Function: main


Give me a list of files needed in order to build this feature.