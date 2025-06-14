# Vigil Rebranding Scripts

This directory contains automated scripts to rebrand the codebase from "Vigil" to "Vigil" throughout the entire project.

## Files

- `rebrand_to_vigil.ps1` - Main PowerShell script that performs the rebranding
- `rebrand_to_vigil.bat` - Windows batch file launcher for easier execution
- `REBRANDING_README.md` - This documentation file

## What Gets Changed

The rebranding script systematically replaces the following:

### Copyright Notices
- `Copyright (c) YYYY The Vigil developers` → `Copyright (c) YYYY The Vigil developers`

### Import Paths
- `github.com/vigilnetwork/vgl/` → `github.com/vigil/vgld/`
- `vigil.network/vgld/` → `vigil.org/vgld/`
- `github.com/Vigil/vglwallet/` → `github.com/vigil/vglwallet/`
- `vigil.network/vgl/wallet/` → `vigil.org/vglwallet/`
- `github.com/Vigil/VGLutil` → `github.com/vigil/vglutil`
- `github.com/Vigil/VGLjson` → `github.com/vigil/vgljson`
- `github.com/vigilnetwork/vgl/slog` → `github.com/vigil/slog`
- `github.com/vigilnetwork/vgl/vspd` → `github.com/vigil/vspd`

### Module Names (in go.mod files)
- `module github.com/vigilnetwork/vgl` → `module github.com/vigil/vgld`
- `module vigil.network/vgld` → `module vigil.org/vgld`
- `module github.com/Vigil/vglwallet` → `module github.com/vigil/vglwallet`
- `module vigil.network/vgl/wallet` → `module vigil.org/vglwallet`

### Binary and Executable Names
- `vgld` → `vgld`
- `vglwallet` → `vglwallet`
- `vglctl` → `vglctl`
- `vgldata` → `vgldata`
- `VGLpool` → `vglpool`

### Network and Protocol References
- `Vigil` → `Vigil`
- `VGL` → `VGL`
- `VGL` → `vgl`
- `VGLwire` → `vglwire`
- `VGLnet` → `vglnet`

### Text in Comments and Documentation
- `Vigil signed message` → `Vigil Signed Message`
- `amount of Vigil` → `amount of vigil`
- `Vigil network` → `Vigil network`
- `Vigil blockchain` → `Vigil blockchain`
- `Vigil protocol` → `Vigil protocol`
- `Vigil wallet` → `Vigil wallet`
- `Vigil node` → `Vigil node`
- `Vigil mining` → `Vigil mining`
- `Vigil staking` → `Vigil staking`

### URLs and Domains
- `vigil.network` → `vigil.org`
- `data.vigil.network` → `vgldata.vigil.org`

## File Types Processed

The script processes the following file types:
- Go source files (*.go)
- Go module files (*.mod)
- Go sum files (*.sum)
- Markdown files (*.md)
- Text files (*.txt)
- YAML files (*.yml, *.yaml)
- JSON files (*.json)
- Shell scripts (*.sh)
- Batch files (*.bat)
- PowerShell scripts (*.ps1)

## Directories and Files Skipped

### Skipped Directories
- `.git`
- `vendor`
- `node_modules`
- `.vscode`
- `.idea`
- `bin`
- `build`
- `dist`

### Skipped Files
- `go.sum` (checksums should be regenerated)
- `rebrand_to_vigil.ps1` (the script itself)

## Usage

### Option 1: Using the Batch File (Recommended for Windows)

1. Double-click `rebrand_to_vigil.bat`
2. Choose from the menu:
   - **Option 1**: Dry run (preview changes without modifying files)
   - **Option 2**: Full rebranding (modify files)
   - **Option 3**: Verbose output options
   - **Option 4**: Exit

### Option 2: Using PowerShell Directly

#### Dry Run (Preview Changes)
```powershell
.\rebrand_to_vigil.ps1 -DryRun -Verbose
```

#### Full Rebranding
```powershell
.\rebrand_to_vigil.ps1
```

#### With Verbose Output
```powershell
.\rebrand_to_vigil.ps1 -Verbose
```

#### Custom Root Path
```powershell
.\rebrand_to_vigil.ps1 -RootPath "C:\path\to\your\project"
```

## Parameters

- `-RootPath`: Specify the root directory to process (default: current script directory)
- `-DryRun`: Preview changes without modifying files
- `-Verbose`: Show detailed output including files processed and patterns matched

## Safety Features

1. **Dry Run Mode**: Always test with `-DryRun` first to see what changes will be made
2. **Confirmation Prompt**: The script asks for confirmation before making changes
3. **Error Handling**: Gracefully handles file access errors and continues processing
4. **Backup Recommendation**: Always commit your changes to git before running
5. **Skip Lists**: Automatically skips sensitive directories and files

## Recommended Workflow

1. **Backup**: Ensure all changes are committed to git
   ```bash
   git add .
   git commit -m "Pre-rebranding backup"
   ```

2. **Dry Run**: Test the script first
   ```powershell
   .\rebrand_to_vigil.ps1 -DryRun -Verbose
   ```

3. **Review**: Check the dry run output to ensure expected changes

4. **Execute**: Run the full rebranding
   ```powershell
   .\rebrand_to_vigil.ps1
   ```

5. **Verify**: Review the changes
   ```bash
   git diff
   ```

6. **Test**: Build and test the project
   ```bash
   go mod tidy
   go build ./...
   go test ./...
   ```

7. **Commit**: Save the rebranded code
   ```bash
   git add .
   git commit -m "Complete rebranding from Vigil to Vigil"
   ```

## Post-Rebranding Tasks

After running the script, you may need to manually update:

1. **Go Module Dependencies**: Run `go mod tidy` to update dependencies
2. **Build Scripts**: Update any custom build scripts or Makefiles
3. **Documentation**: Review and update README files, documentation
4. **Configuration Files**: Update any configuration files not covered by the script
5. **External References**: Update any external references, URLs, or links
6. **Binary Names**: Rename actual binary files if they exist
7. **Database Schemas**: Update any database schemas or migration scripts
8. **API Documentation**: Update API documentation and examples

## Troubleshooting

### PowerShell Execution Policy Error
If you get an execution policy error, run:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### File Access Errors
- Ensure no files are open in editors or IDEs
- Run as administrator if needed
- Check file permissions

### Unexpected Changes
- Always use dry run first
- Review the replacement patterns in the script
- Use git to revert if needed: `git checkout .`

## Script Statistics

The script will report:
- Total files processed
- Total files modified
- Total number of replacements made
- Execution time

## Support

If you encounter issues:
1. Check the verbose output for specific error messages
2. Ensure all files are closed and not locked
3. Verify PowerShell version compatibility
4. Review the file and directory skip lists

---

**Note**: This script is designed specifically for the Vigil project rebranding. Always review and test thoroughly before applying to production code.
