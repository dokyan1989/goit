package main

import (
	// "flag"

	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dokyan1989/goit/util/fileutil"
	"github.com/dokyan1989/goit/util/projectpath"
	// "github.com/fatih/color"
)

var (
	projectName        = "goit"
	protocIncludePaths = []string{
		"vendor",
		"third_party/googleapis",
		"third_party/bufbuild/protovalidate/proto/protovalidate",
		// "third_party/protoc-gen-swagger",
	}
	// Lock version
	protocPlugins = []string{
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest",
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest",
		"github.com/bufbuild/buf/cmd/buf@latest",
		"github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@latest",
		"github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@latest",
		"github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest",
	}

	thirdPartyProtos = map[string][]string{
		"googleapis/google/api": {
			"https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto",
			"https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto",
			"https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto",
			"https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/httpbody.proto",
		},
		"bufbuild/protovalidate/proto/protovalidate/buf/validate": {
			"https://raw.githubusercontent.com/bufbuild/protovalidate/main/proto/protovalidate/buf/validate/expression.proto",
			"https://raw.githubusercontent.com/bufbuild/protovalidate/main/proto/protovalidate/buf/validate/validate.proto",
		},
		"bufbuild/protovalidate/proto/protovalidate/buf/validate/priv": {
			"https://raw.githubusercontent.com/bufbuild/protovalidate/main/proto/protovalidate/buf/validate/priv/private.proto",
		},
	}

	protocOut = []string{
		`--go_out=%s`,
		`--go-grpc_out=%s`,
		`--grpc-gateway_out=%s`,
		// `--grpc-gateway_out=allow_repeated_fields_in_body=true:%s`,
		// `--validate_out=lang=go:%s`,
		// `--ecode_out=%s`,
	}

	ModuleName = "github.com/dokyan1989/goit"
)

func main() {
	/**
	|-------------------------------------------------------------------------
	| 1. Check whether protoc is calling in project root path
	|-----------------------------------------------------------------------*/
	if !inRootProjectPath() {
		logFatal("This command should be called in project root path")
	}

	/**
	|-------------------------------------------------------------------------
	| 2. Check whether protoc is installed
	|----------------------------------------------------------------------*/
	if err := checkProtocInstalled(); err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	/**
	|-------------------------------------------------------------------------
	| 3. Look up all proto files from input folder
	|----------------------------------------------------------------------*/
	protoPath := flag.Arg(0)
	if protoPath == "" {
		logFatal("!!! Proto path should be provided !!!")
	}
	logInfo("Proto path: %s", protoPath)

	protoFiles := fileutil.ScanRecursive(protoPath, []string{".proto"}, nil)
	if len(protoFiles) == 0 {
		logFatal("!!! No .proto files found !!!")
	}
	logInfo("Proto files: ")
	logSlice(protoFiles, "+")

	/**
	|-------------------------------------------------------------------------
	| 4. Download all necessary proto files from 3rd party
	|-----------------------------------------------------------------------*/
	if err := download3rdPartyProtos(); err != nil {
		logFatal("Failed to download 3rd party proto files, error = %v", err.Error())
		return
	}
	logInfo("Download 3rd party proto files successfully!!!")

	/**
	|-------------------------------------------------------------------------
	| 5. Install all necessary protoc plugins
	|-----------------------------------------------------------------------*/
	if err := installProtocPlugins(false); err != nil {
		logFatal("Failed to install protoc plugins, error = %v", err.Error())
		return
	}
	logInfo("Install protoc plugins successfully!!!")

	/**
	|-------------------------------------------------------------------------
	| 6. Create a temporary out folder to generate file
	|-----------------------------------------------------------------------*/
	outFolder, err := os.MkdirTemp("", "")
	if err != nil {
		logFatal(err.Error())
		return
	}
	defer os.RemoveAll(outFolder)

	/**
	|-------------------------------------------------------------------------
	| 7. Create pb folder for storing generated protobuf files
	|-----------------------------------------------------------------------*/
	protoPathParent := filepath.Dir(filepath.Clean(protoPath) /* remove trailing slash if exists*/)
	pbPath := filepath.Join(protoPathParent, "pb")
	if err := os.MkdirAll(pbPath, 0755); err != nil {
		logFatal(err.Error())
		return
	}

	for _, file := range protoFiles {
		logInfo("Process file: %s", file)

		// gen proto file to outFolder
		if err = genProto(file, outFolder); err != nil {
			logFatal(err.Error())
			return
		}

		// copy proto files from outFolder to project folder
		target := projectpath.Root()
		src := path.Join(outFolder, ModuleName)
		if err = copyAll(src, target); err != nil {
			logFatal(err.Error())
			return
		}
	}

	// Print the final message
	logInfo("Success generate protobuf files.")
}

func inRootProjectPath() bool {
	root := projectpath.Root()
	logInfo("Root path: %s", root)
	return strings.HasSuffix(root, projectName)
}

func checkProtocInstalled() error {
	if _, err := exec.LookPath("protoc"); err != nil {
		switch runtime.GOOS {
		case "darwin":
			fmt.Println("brew install protobuf")
			cmd := exec.Command("brew", "install", "protobuf")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				return err
			}
		default:
			// TODO: Add support for other OS here
			// TODO: support ubuntu
			return errors.New("install protobuf manually at ：https://github.com/protocolbuffers/protobuf/releases")
		}
	}
	return nil
}

func download3rdPartyProtos() error {
	rootDir := "third_party"

	for dir, sources := range thirdPartyProtos {
		fullPath := filepath.Join(rootDir, dir)
		ok, err := checkFolderIsExist(fullPath)
		if err != nil {
			return err
		}

		if !ok {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return err
			}
		}

		for _, source := range sources {
			filename := filepath.Base(source)

			ok, err := checkFolderIsExist(filepath.Join(fullPath, filename))
			if err != nil {
				return err
			}

			if ok {
				continue
			}
			fmt.Println("Download:", filepath.Join(fullPath, filename))
			cmd := exec.Command(
				"curl",
				"-o",
				filepath.Join(fullPath, filename),
				source,
			)
			cmd.Dir = projectpath.Root()
			cmd.Env = os.Environ()
			// cmd.Stdout = os.Stdout
			// cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

func installProtocPlugins(checkPluginExists bool) error {
	for _, plugin := range protocPlugins {
		name := getPluginName(plugin)

		if !checkPluginExists {
			if err := goInstall(plugin); err != nil {
				return err
			}

		} else {
			if _, err := exec.LookPath(name); err == nil {
				continue
			}

			if err := goInstall(plugin); err != nil {
				return err
			}
		}
	}

	return nil
}

func getPluginName(path string) string {
	elems := strings.Split(path, "/")
	last := elems[len(elems)-1]

	return strings.Split(last, "@")[0]
}

func goInstall(url string) error {
	cmd := exec.Command("go", "install", url)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("go install", url)
	return cmd.Run()
}

func genProto(file string, outFolder string) error {

	var outCmd []string
	for _, out := range protocOut {
		outCmd = append(outCmd, fmt.Sprintf(out, outFolder))
	}

	// if withOutDescFiles {
	// 	checkFolder(filepath.Join(outFolder, filepath.Dir(file)))
	// 	outCmd = append(outCmd, fmt.Sprintf("-o %s/%s.filedesc", outFolder, strings.TrimSuffix(file, filepath.Ext(file))))
	// }

	// if withBufCheck {
	// 	outCmd = append(outCmd, fmt.Sprintf(_bufCheckCmd, outFolder))
	// 	outCmd = append(outCmd, fmt.Sprintf(`--buf-check-lint_opt={"input_config":%s}`, _bufLintJSON))
	// }

	// if withBreakingCheck {
	// 	outCmd = append(outCmd, fmt.Sprintf(_bufBreakingCmd, outFolder))
	// 	imageFilePath := fmt.Sprintf("%s.filedesc", strings.TrimSuffix(file, filepath.Ext(file)))
	// 	breakingJson := fmt.Sprintf(_bufBreakingJSON, imageFilePath)
	// 	outCmd = append(outCmd, fmt.Sprintf(`--buf-breaking_opt=%s`, breakingJson))
	// }

	projectPath := projectpath.Root()
	parentDir := filepath.Dir(projectPath)
	// Merge all includes
	includes := append([]string{parentDir, projectPath}, protocIncludePaths...)

	// if withGenDoc {
	// 	genDocument(includes, file)
	// }

	// if withApiLinter {
	// 	err := executeAPILinter(includes, file)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// if withSwagger {
	// 	outCmd = append(outCmd, fmt.Sprintf(swaggerOut, "."))
	// }

	if err := generate(includes, outCmd, file); err != nil {
		return err
	}

	/**
	|-------------------------------------------------------------------------
	| Convert Open API V2 Swagger to V3
	|-----------------------------------------------------------------------*/
	// if withSwagger {
	// 	swaggerFile := fmt.Sprintf("%s/%s.swagger.json", pwd, strings.TrimSuffix(file, filepath.Ext(file)))
	// 	v3File := fmt.Sprintf("%s/%s.openapiv3.yaml", pwd, strings.TrimSuffix(file, filepath.Ext(file)))
	// 	data, err := ioutil.ReadFile(swaggerFile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	var doc2 openapi2.T
	// 	if err := json.Unmarshal(data, &doc2); err != nil {
	// 		return err
	// 	}
	// 	doc3, err := openapi2conv.ToV3(&doc2)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	data, err = yaml.Marshal(doc3)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if err := ioutil.WriteFile(v3File, data, 0600); err != nil { // nolint
	// 		return err
	// 	}
	// }

	return nil
}

func generate(includes []string, outputCommand []string, file string) error {
	if os.PathSeparator == '\\' {
		file = filepath.ToSlash(file)
	}

	protocIncludePath := getProtocIncludePathString(includes) + " " + strings.Join(outputCommand, " ")

	// if withDebug {
	// 	fmt.Println(protoc, file)
	// }

	args := strings.Split(protocIncludePath, " ")
	args = append(args, file)
	logInfo("[EXEC] protoc %s", strings.Join(args, " "))

	cmd := exec.Command("protoc", args...)
	cmd.Dir = projectpath.Root() // run protoc cmd at project root path
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getProtocIncludePathString(includePaths []string) string {
	return "-I " + strings.Join(includePaths, " -I ")
}

func copyAll(src, target string) error {
	ok, err := checkFolderIsExist(src)
	if !ok || err != nil {
		return nil
	}

	err = fileutil.CopyTree(src, target, nil)
	if err != nil {
		logFatal("[GO] Can not from `%s` to `%s`: %s\n", src, target, err)
		return err
	}

	return nil
}

func checkFolderIsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// func checkFolder(folder string) {
// 	createFile(folder)
// 	_, err := os.Stat(folder)
// 	if os.IsNotExist(err) {
// 		errDir := os.MkdirAll(folder, 0755)
// 		if errDir != nil {
// 			panic(err)
// 		}
// 	}
// }

// func copyFile(src, dst string) error {
// 	data, err := os.ReadFile(src)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.WriteFile(dst, data, 0600)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func createFile(path string) {
// 	file, err := os.Create(path + "/_index.md")
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()
// 	title := strings.Split(path, "/")
// 	content := "---\n"
// 	content += "title: " + strings.Title(title[len(title)-1])
// 	content += "\ngeekdocCollapseSection: true"
// 	content += "\n---"
// 	_, err = file.WriteString(content)
// 	if err != nil {
// 		return
// 	}
// }

// func getProtoPackage(file string) string {
// 	reader, _ := os.Open(file) // nolint
// 	defer reader.Close()

// 	parser := protoParser.NewParser(reader)
// 	definition, _ := parser.Parse() // nolint
// 	packageName := ""
// 	handleMessage := func(m *protoParser.Package) {
// 		packageName = m.Name
// 	}

// 	protoParser.Walk(definition,
// 		protoParser.WithPackage(handleMessage),
// 	)
// 	return packageName
// }
