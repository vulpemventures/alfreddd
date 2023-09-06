package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	rw = 0766
	rx = 0755
)

func makeScaffolding(folder, ghProject, module, app string, noApiSpec bool) error {
	// list of subdirs to be created and whether they are empty
	dirs := map[string]bool{
		fmt.Sprintf("cmd/%s", app):             false,
		"pkg":                                  true,
		"scripts":                              false,
		"internal/test":                        true,
		"internal/config":                      true,
		"internal/infrastructure":              true,
		"internal/interface/grpc/handlers":     true,
		"internal/interface/grpc/interceptors": true,
		"internal/interface/grpc/permissions":  true,
		"internal/core/domain":                 true,
		"internal/core/application":            true,
		"internal/core/ports":                  true,
		".github/workflows":                    false,
	}
	if !noApiSpec {
		dirs[fmt.Sprintf("api-spec/protobuf/%s/v1", app)] = false
	}

	for dir, isEmpty := range dirs {
		path := filepath.Join(folder, dir)
		if err := makeDirectoryIfNotExists(path); err != nil {
			return err
		}
		if isEmpty {
			file := filepath.Join(path, ".gitkeep")
			if err := os.WriteFile(file, nil, rw); err != nil {
				return err
			}
		}
	}

	license := filepath.Join(folder, "LICENSE")
	licenseContent := makeLicense(ghProject)
	if err := os.WriteFile(license, licenseContent, rw); err != nil {
		return err
	}

	readme := filepath.Join(folder, "README.md")
	readmeContent := makeReadme(ghProject)
	if err := os.WriteFile(readme, readmeContent, rw); err != nil {
		return err
	}

	makefile := filepath.Join(folder, "Makefile")
	makefileContent := makeMakefile(app, noApiSpec)
	if err := os.WriteFile(makefile, makefileContent, rw); err != nil {
		return err
	}

	gitignore := filepath.Join(folder, ".gitignore")
	gitignoreContent := makeGitIngnore()
	if err := os.WriteFile(gitignore, gitignoreContent, rw); err != nil {
		return err
	}

	dockerignore := filepath.Join(folder, ".dockerignore")
	dockerignoreContent := makeDockerIngnore()
	if err := os.WriteFile(dockerignore, dockerignoreContent, rw); err != nil {
		return err
	}

	dockerfile := filepath.Join(folder, "Dockerfile")
	dockerfileContent := makeDockerfile(ghProject, app)
	if err := os.WriteFile(dockerfile, dockerfileContent, rw); err != nil {
		return err
	}

	releaseDockerfile := filepath.Join(folder, "goreleaser.Dockerfile")
	releaseDockerfileContent := makeReleaseDockerfile(ghProject, app)
	if err := os.WriteFile(releaseDockerfile, releaseDockerfileContent, rw); err != nil {
		return err
	}

	script := filepath.Join(folder, "scripts/build")
	scriptContent := makeScriptBuild(app)
	if err := os.WriteFile(script, scriptContent, rx); err != nil {
		return err
	}

	releaser := filepath.Join(folder, ".goreleaser.yaml")
	releaserContent := makeGoreleaserYaml(ghProject, app)
	if err := os.WriteFile(releaser, releaserContent, rw); err != nil {
		return err
	}

	unitAction := filepath.Join(folder, ".github/workflows/ci.unit.yaml")
	unitActionContent := makeUnitTestActionYaml()
	if err := os.WriteFile(unitAction, unitActionContent, rw); err != nil {
		return err
	}

	intergationAction := filepath.Join(folder, ".github/workflows/ci.intergation.yaml")
	intergationActionContent := makeIntegrationTestActionYaml()
	if err := os.WriteFile(intergationAction, intergationActionContent, rw); err != nil {
		return err
	}

	releaseAction := filepath.Join(folder, ".github/workflows/release.yaml")
	releaseActionContent := makeReleaseActionYaml()
	if err := os.WriteFile(releaseAction, releaseActionContent, rw); err != nil {
		return err
	}

	main := filepath.Join(folder, fmt.Sprintf("cmd/%s/main.go", app))
	mainContent := makeMainPlaceholder(module)
	if err := os.WriteFile(main, mainContent, rw); err != nil {
		return err
	}

	iface := filepath.Join(folder, "internal/interface/service.go")
	ifaceContent := makeInterfacePlaceholder(module)
	if err := os.WriteFile(iface, ifaceContent, rw); err != nil {
		return err
	}

	service := filepath.Join(folder, "internal/interface/grpc/service.go")
	serviceContent := makeServicePlaceholder()
	if err := os.WriteFile(service, serviceContent, rw); err != nil {
		return err
	}

	if !noApiSpec {
		buf := filepath.Join(folder, "api-spec/protobuf/buf.yaml")
		bufContent := makeBufYaml(ghProject)
		if err := os.WriteFile(buf, bufContent, rw); err != nil {
			return err
		}

		bufWork := filepath.Join(folder, "buf.work.yaml")
		bufWorkContent := makeBufWorkYaml()
		if err := os.WriteFile(bufWork, bufWorkContent, rw); err != nil {
			return err
		}

		bufGen := filepath.Join(folder, "buf.gen.yaml")
		bufGenContent := makeBufGenYaml(module)
		if err := os.WriteFile(bufGen, bufGenContent, rw); err != nil {
			return err
		}

		proto := filepath.Join(folder, fmt.Sprintf("api-spec/protobuf/%s/v1/service.proto", app))
		protoContent := makeProtoPlaceholder(app)
		if err := os.WriteFile(proto, protoContent, rw); err != nil {
			return err
		}
	}

	return nil
}
